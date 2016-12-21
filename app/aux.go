package app

//auxiliary types and helpers to web handlers


import(
        "errors"
        "time"
        "strconv"
        "crypto/sha512"
        "encoding/hex"
)


type KnownClient struct {
	User string    `json:"user"`
	Chat int       `json: id`
	Ticket string  `json:"ticket"`
}

type AuthS struct {
	User string `json:"username"`
	Key  string `json:"key"`
}

type Reg struct{
    Email string `json:"email"`
    User string `json:"username"`
    Key  string `json:"key"`
    
}


func BakeKey(key string) string{
    h := sha512.New()
    h.Write([]byte(key))
    bake := hex.EncodeToString(h.Sum(nil))
    return bake
    
}

func ActivateToken(email, user string) string{
    return user+email+strconv.FormatInt(time.Now().Unix(),10)
}

func Ticket(user string, id int) string{
    
    tik := user+strconv.FormatInt(time.Now().Unix(),10)+strconv.Itoa(id)
    return tik
}


func DeleteTicket(tik string){
    token.Lock()
     if _,ok:=token.wstokens[tik]; ok{
        delete(token.wstokens,tik)
     }
    token.Unlock()
}

func CacheTicket(tik string, user KnownClient){
    
    token.Lock()
     if _,ok:=token.wstokens[tik]; !ok{
    token.wstokens[tik]=user
     }
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
