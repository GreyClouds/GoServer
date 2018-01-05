package main

type AI interface {
	Start(client *Client)
}

type defaultAI struct {
}

func newDefaultAI() *defaultAI {
	return &defaultAI{}
}

func (self *defaultAI) Start(client *Client) {

}
