package sqlx_demo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_queryRowDemo(t *testing.T) {
	err := initDB()
	assert.Nil(t, err)

	//queryRowDemo()
	//queryMultiRowDemo()
	//insertUserDemo()
	namedQuery()
}

func TestBatchInsert(t *testing.T) {
	err := initDB()
	assert.Nil(t, err)

	defer db.Close()

	u1 := user{Name: "七米", Age: 18}
	u2 := user{Name: "q1mi", Age: 28}
	u3 := user{Name: "小王子", Age: 38}

	// 方法1
	users := []any{&u1, &u2, &u3}
	err = BatchInsertUsers2(users)
	assert.Nilf(t, err, "BatchInsertUser2 failed")
}

func TestQueryByIDs(t *testing.T) {
	users, err := QueryByIDs([]int{1, 2, 3})
	assert.Nil(t, err)
	fmt.Println(users)
}
