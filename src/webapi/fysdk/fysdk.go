package fysdk

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strconv"
)

var (
	_appId       string
	_loginSecret string
	_paySecret   string
)

func Initialize(id, loginSecret, paySecret string) {
	_appId = id
	_loginSecret = loginSecret
	_paySecret = paySecret
}

func GetLoginSecret() string {
	return _loginSecret
}

func GetPaySecret() string {
	return _paySecret
}

func ReimpleURLParamsEncode(params url.Values, isEmptyExcept bool) string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}

		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := params[k]
		prefix := k + "="
		for _, v := range vs {
			if isEmptyExcept && len(v) == 0 {
				continue
			}

			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}

func Sign(form url.Values, key string) string {
	plain := ReimpleURLParamsEncode(form, false)

	hash := md5.New()
	hash.Write([]byte(url.QueryEscape(plain)))
	hash.Write([]byte(`&`))
	hash.Write([]byte(key))

	sign := hex.EncodeToString(hash.Sum(nil))

	return sign
}

// 校验签名
func CheckSign(form url.Values, key string) bool {

	return Sign(form, key) == form.Get("sign")
}

func parseUint32(str string) uint32 {
	iv, _ := strconv.ParseUint(str, 10, 32)
	return uint32(iv)
}

func parseInt32(str string) int32 {
	iv, _ := strconv.ParseInt(str, 10, 32)
	return int32(iv)
}

func parseInt64(str string) int64 {
	iv, _ := strconv.ParseInt(str, 10, 64)
	return iv
}
