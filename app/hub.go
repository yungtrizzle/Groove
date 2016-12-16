package app

import (
	"github.com/yungtrizzle/groove/data"
	"log"
)

/*
 * hub has a single broadcast channel
 * each message is picked up by a goroutine
 * and sent.
 */

type hub struct {
	join      chan *Client
	leave     chan *Client
	broadcast chan *Message
	online    map[int]Client
}

func NewHub() *hub {

	hb := hub{
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan *Message),
		online:    make(map[int]Client),
	}
	return &hb
}

func (h *hub) run() {

	for {
		select {

		case clientn := <-h.join:

			data.EnterRoom(clientn.Room, clientn.User, clientn.Chatid) //ignoring errors here
			h.online[clientn.Chatid] = clientn

		case clientL := <-h.leave:

			data.Leave(clientL.Room, clientL.User, clientL.Chatid)
			if _, ok := h.online[clientL.Chatid]; ok {
				delete(h.online, clientL)
				close(clientL.send)
			}

		case msg := <-h.broadcast:

			mems, ok := data.RoomMembers(msg.Room) //slice of room member id's:integer
			var clist []Client

			if ok != nil {
				log.Println(ok)
				break
			}
			//find the relevant clients and build a work unit
			for _, clyents := range mems {

				con := h.online[clyents]
				clist = append(clist, con)
			}

			work := broadcastwork(msg, clist)
			clist = nil
		}
	}

}
