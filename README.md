## MQ消费者部署


### 先安装RabbitMq 

> 安装支持延时消息的mq插件 rabbitmq_delayed_message_exchange

地址 `https://github.com/rabbitmq/rabbitmq-delayed-message-exchange`

启动rabbitmq_delayed_message_exchange插件

`rabbitmq-plugins enable rabbitmq_delayed_message_exchange`
 
### 安装 golang:v1.11.x

> 在linux编译

`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags -s -a -installsuffix cgo amqpserver.go`

> 在windows 上交叉编译

`set GOOS=linux`

`set GOARCH=amd64`

`set CGO_ENABLED=0`

`go build -ldflags -s -a -installsuffix cgo amqpserver.go `

### 配置说明

```json

{
  "Log4g": {
    "Path": "logs",
    "Stdout": true
  },
  "RabbitMq": {
    "DataSource": "amqp://guest:guest@127.0.0.1:5672",// MQ 协议地址
    "QueueName": "go_queue", //监听队列名
    "Consumer": "go_queue_consumer",//给消费都取个名字
    "Exchange": "go_exchange",// 消息交换机名
    "Durable": false //是否持久队列
  }
}

```

> 拉取文件
 
`git clone https://g.digi800.com/kyang1/amqp-consumer.git`


### 直接运行

> 1.运行程序 

`cd amqp-consumer` 

`chmod 755 amqpserver`

`./amqpserver -c config.json`

> 2.查看日志

`tail -f access.log|error.log`




### 基于docker运行

> 1.docker运行

`cd amqp-consumer` 

`docker build amqp-consumer .`

`docker run -i -t -d -v /root/myhome/logs:/tmp/apps/logs amqpserver`

> 2.查看日志

`tail -f /root/myhome/logs/access.log|error.log`


### 基于 docker-composer 运行

> 1.编辑docker-composer.yml

`vim docker-composer.yml`

```docker
services:
  amqpserver_1:
    build: ./amqp-consumer # 要构建的目录
    image: amqpserver_1:latest
    container_name: amqpserver_1
    volumes:
      - /root/docker/amqpserver/logs:/tmp/apps/logs
```


> 2.查看日志

`tail -f /root/docker/amqpserver/logs/access.log|error.log`



### 发送消息协议示例

```json

{
  "data": "interface", // 数据内容
  "url": "http://www.host.com", // 回调地址
  "delay": 3000, // 延迟时间以 毫秒为单位，为 `0` 时 不延时
  "retryTime": 6000 // 失败重试时间以 毫秒为单位，为 `0` 时 不进行失败重试，
}

```


### 可运行文件仓库地址

https://g.digi800.com/kyang1/amqp-consumer



