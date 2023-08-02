package post

import (
	"time"
)

type DBPost struct {
	ID          int64     `json:"id" db:"post_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	DBPost // 嵌入帖子结构体

	AuthorName    string `json:"author_name"`
	CommunityName string `json:"community_name"` // 嵌入社区信息
}
