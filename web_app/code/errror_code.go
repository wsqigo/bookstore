package code

import "errors"

var (
	ErrorUserExist    = errors.New("username already exist")
	ErrorUserNotExist = errors.New("username already exist")
	ErrorUserNotLogin = errors.New("user is not login")
	ErrorInvalidID    = errors.New("id is not valid")
)
