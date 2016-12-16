package app

import (
	"sync"
)

func init() {
	broadcastPool = sync.Pool{
		New: func() interface{} {
			return &broadcast{}
		},
	}
}

type Message struct {
	chatid int
	room   int
	msg    string
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
	bcast broadcast
}

func (b *BroadcastWork) Close() {

	b.bcast.Close()

	//return broadcast to pool
	broadcastPool.Put(b.bcast)

}

func broadcastwork(msg Message, id []Client) *BroadcastWork {

	var bcaster *broadcast
	bcaster = broadcastPool.Get().(*broadcast)

	bcaster.msg = msg
	bcaster.recievers = make([]Client, len(id))

	copy(bcaster.recievers, id)

	return &BroadcastWork{bcast: bcaster}
}

//build the other pool here also
