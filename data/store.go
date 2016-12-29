//redis online store

package data

import (
	"errors"
	"github.com/mediocregopher/radix.v2/pool"
	"log"
	"strconv"
	"strings"
)

type RedisConfig struct {
	Protocl  string
	Addr     string
	Port     int
	Poolsize int
}

var redispool *pool.Pool

func InitRedis(cfg *RedisConfig) error {

	addr := cfg.Addr + ":" + strconv.Itoa(cfg.Port)
	var ok error

	redispool, ok = pool.New(cfg.Protocl, addr, cfg.Poolsize)

	if ok != nil {
		return ok
	}

	log.Println("Succesfully Init'd Redis")
	return nil
}

func command(cmd string, args ...interface{}) error {

	resp := redispool.Cmd(cmd, args)

	if resp.Err != nil {
		return resp.Err
	}

	return nil
}

func CacheEmailActivation(usser string, activation string) error {
	return command("SETEX", usser, "86400", activation)
}

func CacheWSTicket(usr string, tik string) error {
	return command("SETEX", usr, "5184000", tik)
}

func RetrieveToken(user string, tok string) error {

	resp := redispool.Cmd("GET", user)

	if resp.Err != nil {
		return resp.Err
	}

	code, ok := resp.Str()

	if ok != nil {
		log.Println("Retrieval Failure: ", ok)
		return ok
	}

	if code != tok {
		return errors.New("No Token Found")
	}

	res := redispool.Cmd("DEL", user)

	if res.Err != nil {
		log.Println(res.Err)
	}

	return nil

}

func RetrieveTicket(us string, tok string) error {

	resp := redispool.Cmd("GET", us)

	if resp.Err != nil {
		return resp.Err
	}

	code, ok := resp.Str()

	if ok != nil {
		log.Println("Retrieval Failure: ", ok)
		return ok
	}

	if code != tok {
		return errors.New("No Token Found")
	}

	return nil
}

/*client entering and leaving rooms state changes*/

func EnterRoom(room, chatid string, id int) error {

	chatid += ":"
	chatid += strconv.Itoa(id)
	return command("SADD", room, chatid)
}

func LeaveRoom(room, chatid string, id int) error {
	chatid += ":"
	chatid += strconv.Itoa(id)
	return command("SREM", room, chatid)
}

func RoomMembers(room string) ([]int, error) {

	var members []int
	resp := redispool.Cmd("SMEMBERS", room)

	mem, ok := resp.List()

	if ok != nil {
		return members, ok
	}

	for _, sid := range mem {

		split := strings.Split(sid, ":")

		id, _ := strconv.Atoi(split[1])

		members = append(members, id)
	}

	return members, nil
}

/*
 * messages are written to permanent storage directly
 * and then stored in redis temporarily
 */
func Enqueue(room int, chatid int, message string) error {

	var smsg string
	smsg += strconv.Itoa(room)
	smsg += "+"
	smsg += strconv.Itoa(chatid)
	smsg += "+"
	smsg += message

	return command("LPUSH", "newMsgs", smsg)
}
