鉴权：提供IM_token
用户信息管理
长连接维持
消息资源上传
消息接收、送达

protoc编译命令：
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto
