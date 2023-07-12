package main

import (
	"flag"
	"fmt"
	"strings"
)

type Users []string

func (u *Users) Set(val string) error {
	*u = strings.Split(val, ",")
	return nil
}

func (u *Users) String() string {
	str := "["
	for _, v := range *u {
		str += v
	}
	return str + "]"
}

func main() {
	var u Users
	flag.Var(&u, "u", "用户列表")
	flag.Parse()

	for _, v := range u {
		fmt.Printf(v)
	}
}
