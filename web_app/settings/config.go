package settings

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
	viper.SetConfigType("yaml")   // 指定配置文件类型
	fmt.Println(os.Getwd())
	viper.AddConfigPath(".")    // 指定查找配置文件的路径
	err := viper.ReadInConfig() // 读取配置文件信息
	if err != nil {
		panic("read config file failed. " + err.Error())
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file was changed....")
	})
}
