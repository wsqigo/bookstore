package post

import (
	"bookstore/web_app/dao/mysql"
	"bookstore/web_app/snowflake"
)

var db = mysql.GetDBConn()

func GenAndInsertPost(post DBPost) error {
	// 1. 生成 post id
	post.ID = snowflake.GenID()

	// 2. 保存到数据库
	sqlStr := `insert into post post_id, title, content, author_id, community_id
    values (?, ?, ?, ?, ?)`

	_, err := db.Exec(sqlStr, post.ID, post.Title, post.Content, post.AuthorID)
	return err
}

// GetPostByID 根据帖子 id 查询帖子详情数据
func GetPostByID(pid int64) (DBPost, error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id = ?`

	var post DBPost
	err := db.Get(&post, sqlStr, pid)
	if err != nil {
		return DBPost{}, err
	}

	return post, nil
}
