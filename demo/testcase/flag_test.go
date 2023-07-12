package testcase

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestArgs(t *testing.T) {
	// os.Args 是一个 []string
	if len(os.Args) > 0 {
		for index, arg := range os.Args {
			fmt.Printf("args[%d]=%v\n", index, arg)
		}
	}
}

var (
	host     string
	dbName   string
	port     int
	dbUser   string
	password string
)

func TestFlagType(t *testing.T) {
	flag.StringVar(&host, "host", "", "数据库地址")
	flag.StringVar(&dbName, "db_name", "", "数据库名称")
	flag.StringVar(&dbUser, "user", "", "数据库用户")
	flag.StringVar(&password, "password", "", "数据库密码")
	flag.IntVar(&port, "port", 3306, "数据库端口")

	flag.Parse()

	fmt.Println("数据库地址: ", host)
	fmt.Println("数据库名称: ", dbName)
	fmt.Println("数据库用户: ", dbUser)
	fmt.Println("数据库密码: ", password)
	fmt.Println("数据库端口: ", port)
}
