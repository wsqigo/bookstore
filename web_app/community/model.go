package community

import "time"

type DBCommunity struct {
	ID            int64     `json:"id" db:"id"`
	CommunityID   int64     `json:"community_id" db:"community_id"`
	CommunityName string    `json:"community_name" db:"community_name"`
	Introduction  string    `json:"introduction" db:"introduction"`
	CreatTime     time.Time `json:"creat_time" db:"create_time"`
	UpdateTime    time.Time `json:"update_time" db:"update_time"`
}
