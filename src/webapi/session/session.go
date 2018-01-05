package session

type Session struct {
	SN  uint16
	IP  string
	Uid uint32
}

func NewSession(sn uint16, uid uint32, ip string) *Session {
	self := &Session{
		SN:  sn,
		Uid: uid,
		IP:  ip,
	}

	return self
}
