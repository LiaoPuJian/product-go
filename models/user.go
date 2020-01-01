package models

type User struct {
	ID       int64  `json:"id" sql:"id"`
	NickName string `json:"nick_name" sql:"nick_name"`
	UserName string `json:"user_name" sql:"user_name"`
	Password string `json:"-" sql:"password"`
}
