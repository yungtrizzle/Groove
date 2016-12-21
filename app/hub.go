package app

import (
	"github.com/yungtrizzle/groove/data"
	"log"
)

/*
 * hub has a single broadcast channel
 * each message is picked up by a broadcastwork unit
 * and sent.
 */

type Hub struct {
	join      chan *Client
	leave     chan *Client
	onnline   chan *Client
	offline   chan *Client
	broadcast chan Message
	online    map[int]*Client
}

func NewHub() *Hub {

	hb := Hub{
		join:      make(chan *Client),
		leave:     make(chan *Client),
		onnline:   make(chan *Client),
		offline:   make(chan *Client),
		broadcast: make(chan Message),
		online:    make(map[int]*Client),
	}
	return &hb
}

func (h *Hub) run() {

	for {
		select {

		case clientn := <-h.join:

			data.EnterRoom(clientn.Room, clientn.User, clientn.Chatid) //ignoring errors here
		case clientL := <-h.leave:

			data.LeaveRoom(clientL.Room, clientL.User, clientL.Chatid)

		case clientno := <-h.onnline:

			if _, ok := h.online[clientno.Chatid]; !ok {
				h.online[clientno.Chatid] = clientno
			}

		case cliento := <-h.offline:

			if _, ok := h.online[cliento.Chatid]; ok {
				delete(h.online, cliento.Chatid)
				close(cliento.send)
			}

		case msg := <-h.broadcast:
			//store msg before broadcasting
			data.InsertMessage(msg.text, msg.chatid, msg.room)
			data.Enqueue(msg.room, msg.chatid, msg.text)

			mems, ok := data.RoomMembers(msg.roomName) //slice of room member id's:integer
			var clist []*Client

			if ok != nil {
				log.Println(ok)
				break //break to beginning of for
			}
			//find the relevant clients and build a work unit
			for _, clyents := range mems {

				cid := clyents
				con := h.online[cid]

				if msg.chatid != con.Chatid {
					clist = append(clist, con)
				}

			}

			work := broadcastwork(msg, clist)
			execv(work) //push work unit into pool
			clist = nil
		}
	}

}
