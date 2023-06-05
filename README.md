# mail2dingrobot

![description](./doc/description.svg)

## 实现功能

  mail2dingrobot实现了一个email发送服务器，发送给mail2dingrobot的邮件可以根据配置的规则将转发到钉钉机器人（Markdown消息），也可以同时转发到对应的邮箱。

## 适用场景

  适用于组件告警中以下的一些场景：
 - 这些组件仅支持邮件告警，希望增加钉钉告警功能
 - 这些组件无法连接公网，需要中转

## 快速入门

### 创建config.yml文件

创建`config.yml`文件，文件内容参考👇，具体值替换为真实值：

```yaml
user: "your-email@example.com"
password: 123456
port: 1025
mailAddressDingTokenMap:
  adamslin@foxmail.com: "111111111111111111111111111111111111111111"
mailClient:
  smtpServer: "smtp.example.com"
  smtpPort: 587
  tls: false
  senderEmail: "your-email@example.com"
  senderPassword: "your-email-password"
```

配置内容见[参数说明](## 参数说明)

### 运行docker

#### Linux

```bash
# 修改{your_path_to_config}为配置文件所在文件夹
sudo chown 2001:2001 {your_path_to_config}/config.yml

docker run --name mail2dingrobot \
     --restart always \
     -p 1025:1025 \
     -v /{path_to_your_config}/config.yml:/app/config.yml \
     -d adamswanglin/mail2dingrobot:latest
```

#### MacOS
```bash
# 运行容器
docker run --name mail2dingrobot \
     --restart always \
     -p 1025:1025 \
     -d adamswanglin/mail2dingrobot:latest

# 进入容器
docker exec -it mail2dingrobot sh

# 修改config.yml内容为实际内容
vi config.yml
```

### 对应应用上配置

### 验证

```bash
# 发送邮件
$ telnet 127.0.0.1 1025 

EHLO localhost
AUTH PLAIN
123456
MAIL FROM:<your-email@example.com>
RCPT TO:<to@example.com>
DATA

hello world
.

# 观察日志有无报错
docker logs mail2dingrobot
```

## 外部应用中配置

不同应用配置方式不同，大概的配置：
- 用户名：config.yml中的user
- 密码：config.yml中的password
- smtp服务器地址：部署的机器IP
- smtp端口：1025
- 是否tls：否
- 发送者：config.yml中的user
- 接收者：config.yml中的mailAddressDingTokenMap对应
另外，如果支持配置邮件内容，优先使用Markdown格式（钉钉机器人推送用的Markdown）

## 参数说明

| 名称                    | 类型    | 是否必须 | 备注                                                       |
| :---------------------- | :------ | :------- |:---------------------------------------------------------|
| user                    | string  | 是       | 当前server的登录用户名                                           |
| password                | string  | 是       | 当前server的登录密码                                            |
| port                    | integer | 是       | 当前server的邮件服务端口                                          |
| mailAddressDingTokenMap | object  | 否       | key:value, key为收件人邮箱，value为对应的钉钉机器人token；如果不配置，则可以仅做邮件中转 |
| [mailClient](#mailClient)              | object  | 否       | 邮件中转：如果配置上，会发送邮件到对应邮箱；不配置仅转发到钉钉                          |

#### <span id="mailClient">mailClient</span>

| 名称           | 类型    | 是否必须 | 备注                 |
| :------------- | :------ | :------- | :------------------- |
| smtpServer     | string  | 是       | 邮件服务器地址       |
| smtpPort       | integer | 是       | 邮件服务器smtp端口   |
| tls            | boolean | 是       | 是否为tls            |
| senderEmail    | string  | 是       | 邮件服务器登录用户名 |

[//]: # (| senderPassword | string  | 是       | 邮件服务器登录密码   |)

