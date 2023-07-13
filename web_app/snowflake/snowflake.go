package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(startTime string, machineID int64) error {
	// 时间因子
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}

	// 毫米
	sf.Epoch = st.UnixMilli()
	node, err = sf.NewNode(machineID)
	return nil
}

func GenID() int64 {
	return node.Generate().Int64()
}
