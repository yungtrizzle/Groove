package app

import (
	"sync"
)

var broadcastPool sync.Pool
var wPool *WorkPool
var Bhub *Hub

var token = struct{
    sync.RWMutex
    wstokens map[string]KnownClient
}{ wstokens: make(map[string]KnownClient)}

func init() {
	broadcastPool = sync.Pool{
		New: func() interface{} {
			return &broadcast{}
		},
	}

	wPool = NewPool(10) //10 workers in pool, should be configurable
        
	Bhub = NewHub()
	go Bhub.run()
}

func StartPool() {
	wPool.Close()
	wPool.Wait()
}

func execv(work *BroadcastWork) {
	wPool.Exec(work)
}

//message type and broadcast work unit types

type Message struct {
	chatid   int
	room     int
	roomName string
	text     string
}

func NewMessage(chat int, room int, text string, grp string) Message {

	msg := Message{chatid: chat, room: room, roomName: grp, text: text}
	return msg
}

func (m *Message) Clear() {

	m = &Message{}
}

//broadcast info
type broadcast struct {
	msg       Message
	recievers []*Client
}

func (b *broadcast) Close() {

	b.msg.Clear()
	b.recievers = nil
}

//work unit
type BroadcastWork struct {
	bcast *broadcast
}

func (b *BroadcastWork) Close() {

	b.bcast.Close()
	//return broadcast to pool
	broadcastPool.Put(b.bcast)
}

func (b *BroadcastWork) Execute() {

	sender := b.bcast.msg

	for i, _ := range b.bcast.recievers {

		select {
		case b.bcast.recievers[i].send <- []byte(sender.text):

		default:
			close(b.bcast.recievers[i].send)
		}

	}

	b.Close()
}

func broadcastwork(msg Message, id []*Client) *BroadcastWork {

	var bcaster *broadcast
	bcaster = broadcastPool.Get().(*broadcast)

	bcaster.msg = msg
	bcaster.recievers = make([]*Client, len(id))

	copy(bcaster.recievers, id)

	return &BroadcastWork{bcast: bcaster}
}

type Session struct {
	msgHub *Hub
}

func NewSession(h *Hub) *Session {

	return &Session{msgHub: h}
}

func (s *Session) Onnline(cc *Client) {
	s.msgHub.onnline <- cc
}

func (s *Session) Offline(cc *Client) {
	s.msgHub.offline <- cc
}

func (s *Session) Broadcast(msg Message) {
	s.msgHub.broadcast <- msg
}

func (s *Session) Join(cc *Client) {
	s.msgHub.join <- cc
}

func (s *Session) Leave(cc *Client) {
	s.msgHub.leave <- cc
}
