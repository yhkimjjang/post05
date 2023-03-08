package post05

import (
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
