package web

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/yungtrizzle/groove/app"
	"github.com/yungtrizzle/groove/data"
	"log"
	"net/http"
)



var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


var addr = ":8080"


func serveWs(hub *app.Hub, w http.ResponseWriter, r *http.Request) {

        tik :=  r.URL.Query().Get("id")
        clt,ok := app.ClientTicket(tik)
        
        if ok != nil{
            log.Println(ok)
            return
        }
    
	conn, err := upgrader.Upgrade(w, r, nil)
        
	if err != nil {
		log.Println(err)
		return
	}

	sess := app.NewSession(hub)

	client := app.NewClient(clt.User,
		"",
		clt.Chat,
		0,
		sess,
		conn)

	client.Session.Onnline(client)
	go client.WritePump()
	client.ReadPump()
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var v app.AuthS

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
                log.Println("No Request Body")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&v)

	if err != nil {
		http.Error(w, err.Error(), 400)
                log.Println(err)
		return
	}

	id, ok := data.Auth(v.User, v.Key)

	if ok != nil {
		http.Error(w, "Auth Failure", 400)
                log.Println("Auth Failure")
		return
	}

	w.Header().Set("Content-Type", "application/json")
        
        t:=app.Ticket(v.User,id)
	res := app.KnownClient{User: v.User, Chat: id, Ticket:t}
        app.CacheTicket(t, res)
       
	json.NewEncoder(w).Encode(res)
}


func RegisterHandler(w http.ResponseWriter, r *http.Request){
    
}


func Router() {

	http.HandleFunc("/api/login", LoginHandler)
	http.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(app.Bhub, w, r)
	})

          log.Println("Listening on ", addr)
	err := http.ListenAndServe(addr, nil)
      

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
