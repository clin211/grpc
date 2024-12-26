package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/clin211/grpc/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
)

func main() {
	// 解析命令行参数
	flag.Parse()
	// 创建一个连接到服务器的客户端
	// grpc.NewClient函数用于创建一个新的客户端，需要传入服务器地址和凭证
	// grpc.WithTransportCredentials函数用于设置传输层凭证
	// insecure.NewCredentials函数用于创建不安全的凭证
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 如果创建连接失败，输出错误日志并退出
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 延迟关闭连接
	defer conn.Close()
	// 创建一个GreeterClient实例
	c := pb.NewGreeterClient(conn)

	// 创建一个上下文，用于设置超时时间
	// context.WithTimeout函数用于创建一个新的上下文，需要传入父上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 延迟取消上下文
	defer cancel()

	// 调用SayHello方法，发送请求到服务器
	// c.SayHello函数用于调用SayHello方法，需要传入上下文和请求消息
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "clina"})

	// 如果调用SayHello方法失败，输出错误日志并退出
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// 输出日志，记录服务器返回的消息
	log.Printf("Greeting: %s", r.GetMessage())
}
