package skeleton

import (
	"bytes"
	"crypto/md5"
	//"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"

	pbd "yunjing.me/phoenix/pbd/go"

	oaccount "webapi/account"
	osession "webapi/session"
)

const (
	kMinPacketSize = 2 + 4 + 2 + md5.Size // 最小包体积
	kTokenSize     = 24                   // 访问令牌长度
	kAPISecret     = "Lo01v8!P"           //
)

type MessageHandler func(*Skeleton, *osession.Session, *oaccount.Role, proto.Message) (error, uint16, proto.Message)

type Packet struct {
	rn      uint16 // 随机数
	uid     uint32 // 角色编号
	token   []byte // 访问令牌
	pid     uint16 // 协议编号
	payload []byte // 协议正文
	sign    []byte // 签名
}

// 校验签名
func (self *Packet) CheckSign() bool {
	// 校验签名
	idx := 0
	total := 6 + 2 + len([]byte(kAPISecret))
	if self.uid != 0 && (self.token != nil && len(self.token) > 0) {
		total += len(self.token)
	}

	if self.payload != nil && len(self.payload) > 0 {
		total += len(self.payload)
	}

	raw := make([]byte, total)

	binary.LittleEndian.PutUint16(raw, self.rn)
	idx += 2

	binary.LittleEndian.PutUint32(raw[2:], self.uid)
	idx += 4

	// hash.Write(raw)
	if self.uid != 0 {
		if self.token == nil || len(self.token) == 0 {
			return false
		}

		copy(raw[idx:], self.token)
		idx += len(self.token)
	}

	// raw1 := make([]byte, 2)
	binary.LittleEndian.PutUint16(raw[idx:], self.pid)
	idx += 2
	// hash.Write(raw1)

	if self.payload != nil && len(self.payload) > 0 {
		// hash.Write(self.payload)
		copy(raw[idx:], self.payload)
		idx += len(self.payload)
	}

	// hash.Write([]byte(kAPISecret))
	copy(raw[idx:], []byte(kAPISecret))
	idx += len([]byte(kAPISecret))

	hash := md5.New()
	hash.Write(raw)
	verify := hash.Sum(nil)

	//log.Printf("%v", raw)
	//log.Printf("%v, %v", verify, self.sign)

	return bytes.Equal(verify, self.sign)
}

func recv(r *http.Request) (error, *Packet) {
	var reader io.Reader = r.Body
	maxFormSize := int64(1<<63 - 1)
	if _, ok := r.Body.(*maxBytesReader); !ok {
		maxFormSize = int64(10 << 20)
		reader = io.LimitReader(r.Body, maxFormSize+1)
	}

	b, e := ioutil.ReadAll(reader)
	if e != nil {
		log.Printf("读取字节流时出错: %v", e)
		return e, nil
	}

	l := len(b)
	if l == 0 {
		log.Printf("读取字节流时包体过小: %d", l)
		return errors.New("http trunk too short"), nil
	}

	if int64(l) > maxFormSize {
		log.Printf("读取字节流时包体过大: %d", l)
		return errors.New("http trunk too large"), nil
	}

	if int64(l) < kMinPacketSize {
		log.Printf("读取字节流时包体过小1: %d", l)
		return errors.New("http trunk too short"), nil
	}

	// log.Printf("%v", b)

	packet := new(Packet)
	packet.rn = binary.LittleEndian.Uint16(b[:])
	packet.uid = binary.LittleEndian.Uint32(b[2:])
	if packet.uid != 0 {
		if l < (kMinPacketSize + kTokenSize) {
			log.Printf("读取字节流时包体积过小2: %d", l)
			return errors.New("http trunk too short"), nil
		}

		packet.token = b[6 : 6+kTokenSize]
		packet.pid = binary.LittleEndian.Uint16(b[6+kTokenSize:])
		if (8 + kTokenSize) < (l - md5.Size) {
			packet.payload = b[8+kTokenSize : l-md5.Size]
		}
	} else {
		packet.pid = binary.LittleEndian.Uint16(b[6:])
		if 8 < (l - md5.Size) {
			packet.payload = b[8:(l - md5.Size)]
		}
	}
	packet.sign = b[l-md5.Size : l]
	return nil, packet
}

// 发送消息
func send(w http.ResponseWriter, pid uint16, payload proto.Message) {
	raw1, _ := proto.Marshal(payload)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(raw1)+2))

	var raw2 []byte = make([]byte, len(raw1)+2)
	binary.LittleEndian.PutUint16(raw2, pid)
	copy(raw2[2:], raw1)

	w.Write(raw2[:])
}

func doServerInternalErrorSend(w http.ResponseWriter, err error) {
	errorID := uint32(pbd.ECODE_SERVER_INVALID)
	payload := &pbd.Error{
		Code: &errorID,
	}
	send(w, uint16(pbd.MSG_ERROR), payload)
}

// ----------------------------------------------------------------------------
