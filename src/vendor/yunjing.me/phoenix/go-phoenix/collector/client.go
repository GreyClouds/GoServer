package collector

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang/snappy"
	"github.com/pkg/errors"
	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
	"golang.org/x/crypto/pbkdf2"
)

const (
	SALT = "Phoenix-BI-2017"
)

type compStream struct {
	conn net.Conn
	w    *snappy.Writer
	r    *snappy.Reader
}

func (c *compStream) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

func (c *compStream) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	err = c.w.Flush()
	return n, err
}

func (c *compStream) Close() error {
	return c.conn.Close()
}

func newCompStream(conn net.Conn) *compStream {
	c := new(compStream)
	c.conn = conn
	c.w = snappy.NewBufferedWriter(conn)
	c.r = snappy.NewReader(conn)
	return c
}

// 数据收集客户端
type Client struct {
	queue chan []byte
}

func newClient(addr string) *Client {
	self := &Client{
		queue: make(chan []byte, 1024),
	}

	go self.start(addr)

	return self
}

func (self *Client) start(addr string) {
	config := Config{}
	config.RemoteAddr = addr
	config.Key = "it's a secrect"
	config.Crypt = "aes"
	config.Mode = "normal"
	config.Conn = 1
	config.AutoExpire = 0
	config.ScavengeTTL = 600
	config.MTU = 1350
	config.SndWnd = 128
	config.RcvWnd = 512
	config.DataShard = 10
	config.ParityShard = 3
	config.DSCP = 0
	config.NoComp = false
	config.AckNodelay = false
	config.NoDelay = 0
	config.Interval = 50
	config.Resend = 0
	config.NoCongestion = 0
	config.SockBuf = 4194304
	config.KeepAlive = 10
	config.Log = ""
	config.SnmpLog = ""
	config.SnmpPeriod = 60

	switch config.Mode {
	case "normal":
		config.NoDelay, config.Interval, config.Resend, config.NoCongestion = 0, 40, 2, 1
	case "fast":
		config.NoDelay, config.Interval, config.Resend, config.NoCongestion = 0, 30, 2, 1
	case "fast2":
		config.NoDelay, config.Interval, config.Resend, config.NoCongestion = 1, 20, 2, 1
	case "fast3":
		config.NoDelay, config.Interval, config.Resend, config.NoCongestion = 1, 10, 2, 1
	}

	pass := pbkdf2.Key([]byte(config.Key), []byte(SALT), 4096, 32, sha1.New)
	var block kcp.BlockCrypt
	switch config.Crypt {
	case "tea":
		block, _ = kcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		block, _ = kcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		block, _ = kcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		block, _ = kcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		block, _ = kcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		block, _ = kcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		block, _ = kcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		block, _ = kcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		block, _ = kcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		block, _ = kcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		block, _ = kcp.NewSalsa20BlockCrypt(pass)
	default:
		config.Crypt = "aes"
		block, _ = kcp.NewAESBlockCrypt(pass)
	}

	log.Println("encryption:", config.Crypt)
	log.Println("nodelay parameters:", config.NoDelay, config.Interval, config.Resend, config.NoCongestion)
	log.Println("remote address:", config.RemoteAddr)
	log.Println("sndwnd:", config.SndWnd, "rcvwnd:", config.RcvWnd)
	log.Println("compression:", !config.NoComp)
	log.Println("mtu:", config.MTU)
	log.Println("datashard:", config.DataShard, "parityshard:", config.ParityShard)
	log.Println("acknodelay:", config.AckNodelay)
	log.Println("dscp:", config.DSCP)
	log.Println("sockbuf:", config.SockBuf)
	log.Println("keepalive:", config.KeepAlive)
	log.Println("conn:", config.Conn)
	log.Println("autoexpire:", config.AutoExpire)
	log.Println("scavengettl:", config.ScavengeTTL)
	log.Println("snmplog:", config.SnmpLog)
	log.Println("snmpperiod:", config.SnmpPeriod)

	smuxConfig := smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = config.SockBuf
	smuxConfig.KeepAliveInterval = time.Duration(config.KeepAlive) * time.Second

	createConn := func() (*smux.Session, error) {
		kcpconn, err := kcp.DialWithOptions(config.RemoteAddr, block, config.DataShard, config.ParityShard)
		if err != nil {
			return nil, errors.Wrap(err, "createConn()")
		}
		kcpconn.SetStreamMode(true)
		kcpconn.SetWriteDelay(true)
		kcpconn.SetNoDelay(config.NoDelay, config.Interval, config.Resend, config.NoCongestion)
		kcpconn.SetWindowSize(config.SndWnd, config.RcvWnd)
		kcpconn.SetMtu(config.MTU)
		kcpconn.SetACKNoDelay(config.AckNodelay)

		// if err := kcpconn.SetDSCP(config.DSCP); err != nil {
		// 	log.Println("SetDSCP:", err)
		// }
		if err := kcpconn.SetReadBuffer(config.SockBuf); err != nil {
			log.Println("SetReadBuffer:", err)
		}
		if err := kcpconn.SetWriteBuffer(config.SockBuf); err != nil {
			log.Println("SetWriteBuffer:", err)
		}

		// stream multiplex
		var session *smux.Session
		if config.NoComp {
			session, err = smux.Client(kcpconn, smuxConfig)
		} else {
			session, err = smux.Client(newCompStream(kcpconn), smuxConfig)
		}
		if err != nil {
			return nil, errors.Wrap(err, "createConn()")
		}
		log.Println("connection:", kcpconn.LocalAddr(), "->", kcpconn.RemoteAddr())
		return session, nil
	}

	// wait until a connection is ready
	waitConn := func() *smux.Session {
		for {
			if session, err := createConn(); err == nil {
				return session
			} else {
				log.Println("re-connecting:", err)
				time.Sleep(time.Second)
			}
		}
	}

	numconn := uint16(config.Conn)
	muxes := make([]struct {
		session *smux.Session
		ttl     time.Time
	}, numconn)

	for k := range muxes {
		muxes[k].session = waitConn()
		muxes[k].ttl = time.Now().Add(time.Duration(config.AutoExpire) * time.Second)
	}

	chScavenger := make(chan *smux.Session, 128)
	go scavenger(chScavenger, config.ScavengeTTL)
	go snmpLogger(config.SnmpLog, config.SnmpPeriod)

	rr := uint16(0)
	for {
		select {
		case p := <-self.queue:
			idx := rr % numconn

			// do auto expiration && reconnection
			if muxes[idx].session.IsClosed() || (config.AutoExpire > 0 && time.Now().After(muxes[idx].ttl)) {
				chScavenger <- muxes[idx].session
				muxes[idx].session = waitConn()
				muxes[idx].ttl = time.Now().Add(time.Duration(config.AutoExpire) * time.Second)
			}

			handleClient(muxes[idx].session, p)
			rr++
		}
	}
}

func handleClient(sess *smux.Session, raw []byte) {
	p2, err := sess.OpenStream()
	if err != nil {
		log.Printf("[数据采集]打开连接时出错: %v", err)
		return
	}
	defer p2.Close()

	// log.Printf("%v", raw)

	if _, err := p2.Write(raw); err != nil {
		log.Printf("写入数据流时出错: %v", err)
	}
}

func (self *Client) Send(raw []byte) {
	go func() {

		self.queue <- raw
	}()
}

type scavengeSession struct {
	session *smux.Session
	ts      time.Time
}

func scavenger(ch chan *smux.Session, ttl int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var sessionList []scavengeSession
	for {
		select {
		case sess := <-ch:
			sessionList = append(sessionList, scavengeSession{sess, time.Now()})
			// log.Println("session marked as expired")
		case <-ticker.C:
			var newList []scavengeSession
			for k := range sessionList {
				s := sessionList[k]
				if s.session.NumStreams() == 0 || s.session.IsClosed() {
					// log.Println("session normally closed")
					s.session.Close()
				} else if ttl >= 0 && time.Since(s.ts) >= time.Duration(ttl)*time.Second {
					log.Println("session reached scavenge ttl")
					s.session.Close()
				} else {
					newList = append(newList, sessionList[k])
				}
			}
			sessionList = newList
		}
	}
}

func snmpLogger(path string, interval int) {
	if path == "" || interval == 0 {
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			f, err := os.OpenFile(time.Now().Format(path), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Println(err)
				return
			}
			w := csv.NewWriter(f)
			// write header in empty file
			if stat, err := f.Stat(); err == nil && stat.Size() == 0 {
				if err := w.Write(append([]string{"Unix"}, kcp.DefaultSnmp.Header()...)); err != nil {
					log.Println(err)
				}
			}
			if err := w.Write(append([]string{fmt.Sprint(time.Now().Unix())}, kcp.DefaultSnmp.ToSlice()...)); err != nil {
				log.Println(err)
			}
			kcp.DefaultSnmp.Reset()
			w.Flush()
			f.Close()
		}
	}
}
