package main

import (
	"fmt"
	"miai-open-the-door/actions"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const server = "bemfa.com"
const port = 9501

// 订阅消息
func sub(client mqtt.Client, env *DoorEnv) {
	if token := client.Subscribe(env.Topic, 0, func(client mqtt.Client, message mqtt.Message) {
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
	}); token.Wait() && token.Error() != nil {
		fmt.Println("订阅失败")
		panic(token.Error())
	}
	fmt.Println("订阅成功")
}
func checkConnection(client mqtt.Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !client.IsConnected() {
				fmt.Println("MQTT连接已断开，尝试重新连接...")
				if token := client.Connect(); token.Wait() && token.Error() != nil {
					fmt.Printf("重新连接失败: %v\n", token.Error())
				} else {
					fmt.Println("重新连接成功")
				}
			} else {
				fmt.Println("MQTT连接正常")
			}
		}
	}
}
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
	opts.KeepAlive = 60
	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Printf("mqtt服务 重连... \n")
		sub(client, env)
	}
	// 连接成功
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Printf("mqtt服务 连接成功 \n")
		sub(client, env)
	}
	// 连接丢失
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("mqtt服务 连接丢失: %v \n", err)
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
	// 启动检查连接的goroutine
	go checkConnection(client, 30*time.Second)

	// 监听OS信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-c
		fmt.Println("退出", s)
		client.Disconnect(1000)

		fmt.Println("已断开连接")
	}

}
