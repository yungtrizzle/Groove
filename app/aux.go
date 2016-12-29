package app

//auxiliary types and helpers to web handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/yungtrizzle/groove/data"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type KnownClient struct {
	User   string `json:"user"`
	Chat   int    `json: id`
	Ticket string `json:"ticket"`
}

type AuthS struct {
	User string `json:"username"`
	Key  string `json:"key"`
}

type Reg struct {
	Email string `json:"email"`
	User  string `json:"username"`
	Key   string `json:"key"`
}

func BakeKey(key string) string {
	h := sha512.New()
	h.Write([]byte(key))
	bake := hex.EncodeToString(h.Sum(nil))
	return bake

}

func ActivateToken(email, user string) string {
	return user + email + strconv.FormatInt(time.Now().Unix(), 10)
}

func Ticket(user string, id int) string {

	tik := user + strconv.FormatInt(time.Now().Unix(), 10) + "_" + strconv.Itoa(id)
	return tik
}

func CacheTicket(tik string, user KnownClient) {

	data.CacheWSTicket(user.User, tik)
}

func ClientTicket(tik string) (KnownClient, error) {

	re := regexp.MustCompile("^[A-Za-z]*")
	us := re.FindString(tik)

	err := data.RetrieveTicket(us, tik)

	idstr := strings.Split(tik, "_")

	cid, _ := strconv.Atoi(idstr[len(idstr)-1])
	kc := KnownClient{User: us, Chat: cid, Ticket: tik}

	return kc, err
}
