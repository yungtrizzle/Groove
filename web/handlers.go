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
        
         CheckOrigin: func(r *http.Request) bool {
        return true
    },
}


var addr = ":8080"


func serveWs(hub *app.Hub, w http.ResponseWriter, r *http.Request) {

        tik :=  r.URL.Query().Get("ticket")
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
        
        app.DeleteTicket(tik)
        
	sess := app.NewSession(hub)

	client := app.NewClient(clt.User,
		"",
		clt.Chat,
		0,
		sess,
		conn)

	client.Session.Onnline(client)
	go client.WritePump()
	go client.ReadPump()
}

func ActivateHandler(w http.ResponseWriter, r *http.Request){
     tik :=  r.URL.Query().Get("user")
     tok := r.URL.Query().Get("token")
         
     ok:=data.RetrieveToken(tok,tik)
     
     if ok != nil{
          http.Error(w, "Token Doesn't Exists or Has Expired", 404)
          log.Println(ok)
          return
    }
    data.ActivateUser(tik)
    
    w.Write([]byte("Account Activated"))
    
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
	
	if !data.IsActivated(v.User) {
         
            http.Error(w, "Unactivated User", 400)
            return
        }
        
        v.Key=app.BakeKey(v.Key)
        
	id, ok := data.Auth(v.User, v.Key)

	if ok != nil {
		http.Error(w, "Auth Failure", 400)
                log.Println("Auth Failure: ", ok)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
        
        t:=app.Ticket(v.User,id)
	res := app.KnownClient{User: v.User, Chat: id, Ticket:t}
        app.CacheTicket(t, res)
       
	json.NewEncoder(w).Encode(res)
}


func RegisterHandler(w http.ResponseWriter, r *http.Request){
    
    var v app.Reg

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
	
        v.Key = app.BakeKey(v.Key)
        
        ok:=data.RegisterUser(v.User,v.Key)
        
        if ok != nil{
            http.Error(w, "User Exists", 400)
            log.Println(ok)
            return
        }
        
        token := app.ActivateToken(v.Email,v.User)
        log.Println(token)
        
        rok := data.CacheEmailActivation(token, v.User)
        
        if err != nil{
            log.Println("Caching Error:",rok)
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(token)
}


func Router() {

	http.HandleFunc("/api/login", LoginHandler)
        http.HandleFunc("/api/register", RegisterHandler)
        http.HandleFunc("/api/activate", ActivateHandler)
	http.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(app.Bhub, w, r)
	})

          log.Println("Listening on ", addr)
	err := http.ListenAndServe(addr, nil)
      

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
