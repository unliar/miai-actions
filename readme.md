# 小爱同学语音联动

> 项目使用了 [巴法)(https://bemfa.com) MQTT 服务

## 联动小爱同学原理

1. 注册巴法 MQTT 设备云。

2. 新建主题获取 比如 mqtt001, 尾号 3 个数字代表设备类型 001 是插座。

3. 进入主题详情, 并且设置一下设备名称, 比如 插座 。

4. 打开米家 APP -> 我的 -> 其他平台设备 -> 添加 -> 绑定你的巴法账号 -> 同步设备。

5. 理论上, 此时你对小爱同学说 打开插座, 你的 MQTT 主题订阅能收到一条 on 的消息。

6. 收到消息... 你可以发挥想象做点什么。

## 关于本项目

1. 只是一个为了开楼下门禁的项目,抓包了楼下门禁小程序 API。

2. 纯属自娱自乐。

3. 不提供特殊服务。
