package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"miai-open-the-door/actions"
	"os"
)

type DoorEnv struct {
	Key             string                  `json:"key"`               // 用户的设备key
	Topic           string                  `json:"topic"`             // 用户的主题名
	OpenDoorRequest actions.OpenDoorRequest `json:"open_door_request"` // 抓包小程序来的请求参数
}

func GetDoorEnv() (*DoorEnv, error) {
	str := os.Getenv("DOOR_ENV")
	if str == "" {
		return nil, errors.New("未获取到 DOOR_ENV 环境变量")
	}
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.New("base64解码失败,请确认base64是否完整")
	}
	fmt.Println("配置文件内容", string(b))
	var doorEnv DoorEnv
	err = json.Unmarshal(b, &doorEnv)
	return &doorEnv, err
}
