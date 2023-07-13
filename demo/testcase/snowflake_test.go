package testcase

import (
	"fmt"
	"testing"
	"time"

	"github.com/sony/sonyflake"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) error {
	// 时间因子
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}

	// 毫米
	snowflake.Epoch = st.UnixMilli()
	node, err = snowflake.NewNode(machineID)
	return nil
}

func GenID() int64 {
	return node.Generate().Int64()
}

func GenIDString() string {
	return node.Generate().String()
}

func TestGenID(t *testing.T) {
	if err := Init("2020-07-01", 1); err != nil {
		fmt.Printf("init failed, err:%v\n", err)
		return
	}

	id := GenID()
	fmt.Println(id)
}

// Sony
var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

// 需传入当前的机器ID
func InitSonyFlake(startTime string, machineID uint16) error {
	// 时间因子
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}

	settings := sonyflake.Settings{
		StartTime: st,
		MachineID: getMachineID,
	}

	sonyFlake = sonyflake.NewSonyflake(settings)
	return nil
}

// GenSonyID 生成 ID
func GenSonyID() (uint64, error) {
	if sonyFlake == nil {
		return 0, fmt.Errorf("sony flake not inited")
	}

	return sonyFlake.NextID()
}
