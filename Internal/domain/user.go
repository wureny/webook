package domain

import "time"

type User struct {
	Id       uint64
	Email    string
	Password string
	Ctime    time.Time
}

/*type Address struct {
}
*/
