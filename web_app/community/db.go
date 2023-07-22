package community

import (
	"bookstore/web_app/code"
	"bookstore/web_app/dao/mysql"
	"database/sql"

	"go.uber.org/zap"
)

var db = mysql.GetDBConn()

func getCommunityList() ([]DBCommunity, error) {
	sqlStr := "select community_id, community_name from community"

	var res []DBCommunity
	err := db.Select(&res, sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}

func GetCommunityDetailByID(id int64) (DBCommunity, error) {
	sqlStr := `select community_id, community_name, introduction, create_time
from community
where id = ?`

	var res DBCommunity
	err := db.Get(&res, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = code.ErrorInvalidID
		}

		return DBCommunity{}, err
	}

	return res, nil
}
