package config

import "fmt"

func GenServiceConfig(httpName string, grpcName string, httpPort string, grpcPort string, enableSsl bool, domain string) string {
	return fmt.Sprintf(`
system:
  ssl: %t # 是否启用ssl true的话不会使用port
  gateway_address: "%s"
  gateway_port: 8080

log:
  stdout: true
  level: -1
  file: "logs/app.log"
  format: "console" # console、json

email:
  enable: false
  smtp_server: "smtp.gmail.com"
  port: 25
  username: ""
  password: ""

livekit:
  address: ""
  port: 7880
  url: "wss://"
  api_key: ""
  secret_key: ""
  timeout: "1m"

cache:
  enable: true

mysql:
  address: "mysql"
  port: 3306
  username: "root"
  password: "Hitosea@123.."
  database: "coss"
  opts:
   allowNativePasswords: "true"
   timeout: "800ms"
   readTimeout: "200ms"
   writeTimeout: "800ms"
   parseTime: "true"
   loc: "Local"
   charset: "utf8mb4"

redis:
  proto: "tcp"
  address: "redis"
  port: 6379
  password: "Hitosea@123.."
#  protocol: 3

dtm:
  name: "dtm"
  address: "dtm"
  port: 36790

http:
  name: "0.0.0.0"
  port: %s

grpc:
  name: "0.0.0.0"
  port: %s

# 注册本服务
register:
  # 注册中心地址
  address: "consul"
  # 注册中心端口
  port: 8500
  tags: ["%s", "service", "%s service"]

discovers:
  user:
    name: "user_service"
    address: "user"
    port: 10002
    direct: true
  relation:
    name: "relation_service"
    address: "relation"
    port: 10001
    direct: true
  storage:
    name: "storage_service"
    address: "storage"
    port: 10006
    direct: true
  msg:
    name: "msg_service"
    address: "msg"
    port: 10000
    direct: true
  group:
    name: "group_service"
    address: "group"
    port: 10005
    direct: true
  push:
    name: "push_service"
    address: "push"
    port: 10007
    direct: true

message_queue:
  name: "rabbitmq"
  username: "root"
  password: "Hitosea@123.."
  address: "rabbitmq"
  port: 5672

encryption:
  enabled: false
  name: coss-im
  email: max.mustermann@example.com
  passphrase: LongSecret
  rsaBits: 2048

multiple_device_limit:
  enable: false
  max: 1

oss:
  name: "minio"
  address: "minio"
  port: 9000
  accessKey: "root"
  secretKey: "Hitosea@123.."
  ssl: false
  presignedExpires: ""
  dial: "3000ms"
  timeout: "5000ms"
`, enableSsl, domain, httpPort, grpcPort, grpcName, grpcName)
}

func GenConsulServiceConfig(httpName, grpcName, httpPort, grpcPort string, enableSsl bool, domain string) string {
	return fmt.Sprintf(`
system:
  ssl: %t # 是否启用ssl true的话不会使用port
  gateway_address: "%s"
  gateway_port: 8080
  jwt_secret: "secret"

log:
  stdout: true
  level: -1
  file: "logs/app.log"
  format: "console" # console、json

email:
  enable: false
  smtp_server: "smtp.gmail.com"
  port: 25
  username: ""
  password: ""

livekit:
 address: http://103.63.139.136
 port: 7880
 url: wss://tuo.gezi.vip
 api_key: APIbsEc4M9ceob3
 secret_key: Op5frnZoFRUlG0lnCUhlh12I1XfdrB90ZEji07fXQZbB
 timeout: 2m

cache:
  enable: true

http:
  name: "%s"
  port: %s

grpc:
  name: "%s"
  port: %s

register:
  address: "consul"
  port: 8500
  tags: ["%s", "service", "%s"]

discovers:
  msg:
    name: "msg_service"
    port: 10000
    direct: false
  user:
    name: "user_service"
    port: 10002
    direct: false
  group:
    name: "group_service"
    port: 10005
    direct: false
  relation:
    name: "relation_service"
    port: 10001
    direct: false
  storage:
    name: "storage_service"
    port: 10006
    direct: false
  push:
    name: "push_service"
    port: 10007
    direct: false

encryption:
  enabled: false
  name: coss-im
  email: max.mustermann@example.com
  passphrase: LongSecret
  rsaBits: 2048

multiple_device_limit:
  enable: false
  max: 1
`, enableSsl, domain, httpName, httpPort, grpcName, grpcPort, grpcName, grpcName)
}
