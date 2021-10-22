package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"miai-open-the-door/actions"
	"os"
	"time"
)

const server = "bemfa.com"
const port = 9501

func main() {
	// 获取配置
	env, err := GetDoorEnv()
	if err != nil {
		panic(err)
	}
	// 配置项
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", server, port))
	// 用户私钥
	opts.SetClientID(env.Key)
	// 超时时间配置
	opts.ConnectTimeout = 5 * time.Second
	opts.AutoReconnect = true
	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Printf("重连mqtt服务 \n")
	}
	// 连接成功
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Printf("连接成功 \n")
	}
	// 连接丢失
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("连接丢失: %v \n", err)
	}
	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", message.Topic())
		fmt.Printf("MSG: %s\n", message.Payload())
	})
	// 客户端
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// 订阅
	client.Subscribe(env.Topic, 0, func(client mqtt.Client, message mqtt.Message) {
		payload := string(message.Payload())
		fmt.Printf("From - Topic: %s\n", message.Topic())
		fmt.Printf("MSG - Payload: %s\n", payload)
		if payload == "on" {
			fmt.Println("开门指令触发")
			actions.OpenTheDoor(env.OpenDoorRequest)
			return
		}
		if payload == "off" {
			fmt.Println("关门指令触发")
			return
		}
		fmt.Printf("收到其他指令消息:" + payload)
	})

	c := make(chan os.Signal)
	s := <-c
	fmt.Println("退出", s)
}
