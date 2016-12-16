package app

import (
	"sync"
)

var wPool *WorkPool

func init() {
	broadcastPool = sync.Pool{
		New: func() interface{} {
			return &broadcast{}
		},
	}

	wPool = NewPool(10) //10 workers in pool
}

func execv(work *BroadcastWork) {
	wPool.Exec(work)
}

type Message struct {
	chatid int
	room   int
	text   string
}

func NewMessage(chat int, room int, text string) *Message {

	msg := Message{chat, room, text}
	return &msg
}

func (m *Message) Clear() {

	m := Message{}
}

type broadcast struct {
	msg       Message
	recievers []Client
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

func broadcastwork(msg Message, id []Client) *BroadcastWork {

	var bcaster *broadcast
	bcaster = broadcastPool.Get().(*broadcast)

	bcaster.msg = msg
	bcaster.recievers = make([]Client, len(id))

	copy(bcaster.recievers, id)

	return &BroadcastWork{bcast: bcaster}
}
