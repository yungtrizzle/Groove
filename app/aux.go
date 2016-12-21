package app

//auxiliary types and helpers to web handlers


import(
        "errors"
        "time"
        "strconv"
)


type KnownClient struct {
	User string    `json:user`
	Chat int       `json:id`
	Ticket string  `json:ticket`
}

type AuthS struct {
	User string `json:username`
	Key  string `json:key`
}




func Ticket(user string, id int) string{
    
    tik := user+strconv.FormatInt(time.Now().Unix(),10)+strconv.Itoa(id)
    return tik
}

func CacheTicket(tik string, user KnownClient){
    
    token.Lock()
    token.wstokens[tik]=user
    token.Unlock()
    
}

func ClientTicket(tik string) (KnownClient,error){
    
    token.RLock()
    kclt,ok := token.wstokens[tik]
    token.RUnlock()
    
    if !ok{
        return kclt, errors.New("No Token Found")
    }
    
    return kclt,nil
}
