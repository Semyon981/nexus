package models

import "time"

type Message struct {
	Id_Messages int64
	Id_from     int64
	Id_to       int64
	Msg         string
	Time        time.Time
}
