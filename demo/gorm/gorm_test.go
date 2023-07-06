package gorm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 修改大小
// 另一种方式 Name string `gorm:"type:varchar(12)"`
type Student struct {
	ID     uint   `gorm:"size:3"`
	Name   string `gorm:"size:8"`
	Age    int    `gorm:"size:3"`
	Gender bool
	Email  *string `gorm:"size:128"` // 使用指针能存储空值
}

func TestName(t *testing.T) {
	// gorm 会创建一个表
	// create table `students` (`name` longtext,`age` bigint,`my_student` longtext)
	//db.AutoMigrate(&Student{})

	//email := "wsqigo1@gmail.com"
	//// 添加记录
	//stu := &Student{
	//	Name:   "wsqigo",
	//	Age:    28,
	//	Gender: true,
	//	Email:  &email,
	//}
	//
	//err := db.Create(&stu).Error
	//assert.Nil(t, err)

	// 批量插入
	//students := make([]Student, 0, 10)
	//for i := 0; i < 10; i++ {
	//	students = append(students, Student{
	//		Name:   fmt.Sprintf("wsqigo_%d", i),
	//		Age:    21 + i,
	//		Gender: true,
	//		Email:  &email,
	//	})
	//}
	//
	//err := db.Create(&students).Error
	//assert.Nil(t, err)

	// 单条数据查询
	var stu Student
	db = db.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})
	err := db.Take(&stu).Error
	assert.Nil(t, err)
	fmt.Println(stu)
	var stu1 Student
	err = db.First(&stu1).Error
	assert.Nil(t, err)
	fmt.Println(stu1)
	var stu2 Student
	err = db.Last(&stu2).Error
	assert.Nil(t, err)
	fmt.Println(stu2)

	var stu3 Student
	err = db.Take(&stu3, 8).Error
	assert.Nil(t, err)
	fmt.Println(stu3)

	var stu4 Student
	err = db.Take(&stu4, "name = ?", "wsqigo_7").Error
	assert.Nil(t, err)
	fmt.Println(stu4)

}
