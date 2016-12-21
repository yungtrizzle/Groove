package data

import (
        "errors"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

type PostgresConfig struct {
	Host string
	Port int
	User string
	Key  string // this field depends on postgres setup
	//insert password=%s and cfg.password into Sprintf below if necessary
	Dbname string
}

var db *sql.DB

func InitPostgres(cfg *PostgresConfig) error {

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Dbname)

	var err error

	db, err = sql.Open("postgres", dbInfo)

	if err != nil {
		return err
	}

	err = db.Ping()

	if err != nil {
		return err
	}

	log.Println("Succesfully Init'd Postgres")
	return nil

}

/*Insertion funcs here*/

/*key should be hashed before transport from client*/
func RegisterUser(username, key string) error {

	insertstr := `INSERT INTO users (username, password)
                  VALUES ($1,$2)
                `
	//schema defaults to user being unique
	stmt, err := db.Prepare(insertstr)

	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()

	_, ok := stmt.Exec(username, key)

	if ok != nil {
		log.Println(ok)
		return ok
	}

	return nil
}

func RegisterRoom(adminid int, name string, description string) error {

	insert := `INSERT INTO rooms(admin_id, name, desc_)
            VALUES($1,$2,$3)
            `
	stmt, err := db.Prepare(insert)

	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()

	_, ok := stmt.Exec(adminid, name, description)

	if ok != nil {
		log.Println(ok)
		return ok
	}

	return nil

}

func JoinRoom(roomid int, userid int) error {

	insert := `INSERT INTO group_member(group_id, user_id)
            VALUES($1,$2)
            `
	stmt, err := db.Prepare(insert)

	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()

	_, ok := stmt.Exec(roomid, userid)

	if ok != nil {
		log.Println(ok)
		return ok
	}

	return nil
}

func InsertMessage(message string, userid int, groupid int) error {

	msg := `INSERT INTO message(from_user_id,dest_id,message)
        VALUES($1,$2,$3)`

	stmt, err := db.Prepare(msg)

	if err != nil {
		log.Println(err)
	}

	defer stmt.Close()

	_, ok := stmt.Exec(userid, groupid, message)

	if ok != nil {
		log.Println(ok)
		return ok
	}

	return nil
}

/*Retrieval funcs*/

func GetAllRooms() []string {

	roomsqry := `SELECT group_id, name, desc_ FROM rooms`

	var rooms []string
	var roomdesc string
	rows, ok := db.Query(roomsqry)

	if ok != nil {
		log.Println(ok)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var desc string

		ok = rows.Scan(&id, &name, &desc)

		if ok != nil {
			log.Println(ok)
		}
		roomdesc += strconv.Itoa(id)
		roomdesc += "+"
		roomdesc += name
		roomdesc += desc

		rooms = append(rooms, roomdesc)

		roomdesc = ""
	}

	return rooms
}

func GetUserRooms(userid int) []string {

	var rooms []string
	var roomdesc string

	roomqry := `SELECT group_id, name FROM rooms
            WHERE group_id IN
            (SELECT group_id FROM group_member WHERE user_id=$1)
            `
	rows, ok := db.Query(roomqry, userid)

	if ok != nil {
		log.Println(ok)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		ok = rows.Scan(&id, &name)

		if ok != nil {
			log.Println(ok)
		}

		roomdesc += strconv.Itoa(id)
		roomdesc += "+"
		roomdesc += name

		rooms = append(rooms, roomdesc)

		roomdesc = ""
	}

	return rooms
}

func Auth(username, key string) (int, error) {

	auth := `SELECT user_id, key from users 
        WHERE username=$1`

	var keyd string
	var id int

	ok := db.QueryRow(auth, username).Scan(&id, &keyd)

	switch {

	case ok == sql.ErrNoRows:
		return 0,ok

	case ok != nil:
		log.Fatal(ok)

	case keyd == key:
		return id,nil

	default:
		return 0,ok

	}
	return 0, errors.New("DB:Unknown Failure")
}
