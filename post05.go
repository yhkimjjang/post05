package post05

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int
	username string
}

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

// 연결 상세
var (
	Hostname = ""
	Port     = 5432
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*sql.DB, error) {
	// 연결 문자열
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)

	// 데이터베이스 연결
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 이 함수는 사용자 이름을 받아 ID를 반환한다
// 사용자가 존재하지 않으면 -1를 반환한다.
func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := -1
	statement := fmt.Sprintf(`select "id" from "users" where username = '%s'`, username)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userID = id
	}
	defer rows.Close()
	return userID
}

// AddUser는 데이터베이스에 새로운 사용자를 추가하고
// 해당 사용자의 User ID를 반환한다
// 에러가 발생한다면 -1를 반환한다
func AddUser(d Userdata) int {
	d.Username = strings.ToLower(d.Username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", d.Username)
		return -1
	}

	insertStatement := `insert into "users" ("username") values ($1)`
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}

	insertStatement = `insert into "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, d.ID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	return userID
}

// DeleteUser는 존재하는 사용자를 지운다.
func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	statement := fmt.Sprintf(`select "username" from "users" where id = %d`, id)
	rows, err := db.Query(statement)
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	// Userdata에서 지운다.
	deleteStatement := `delete from "userdata" where userid = $1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	// users에서 지운다
	deleteStatement = `delete from "users" where id = $1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser는 존재하는 사용자를 업데이트 한다
func UpdateUser(d Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID == -1 {
		return errors.New("User does not exits")
	}

	d.ID = userID
	updateStatement := `update "userdata" set "name"=$1, "surname"=$2, "description"=$3 where userid=$4`
	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}
	return nil
}

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	selectStatement := `select "id", "username", "name", "surname", "description" from "users", "userdata" where users.id = userdata.userid`
	rows, err := db.Query(selectStatement)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var usernsme, name, surname, descrition string
		err = rows.Scan(&id, &usernsme, &name, &surname, &descrition)
		temp := Userdata{
			ID:          id,
			Username:    usernsme,
			Name:        name,
			Surname:     surname,
			Description: descrition,
		}

		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}
