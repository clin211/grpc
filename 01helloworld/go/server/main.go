// 包名为main，表示这是一个可执行程序
package main

// 导入必要的包
import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/clin211/grpc/rpc"
	"google.golang.org/grpc"
)

// 定义一个变量port，用于存储服务器监听的端口号
// flag.Int函数用于解析命令行参数，返回一个指向int值的指针
var (
	port = flag.Int("port", 50052, "The server port")
)

// 定义一个结构体server，实现了pb.GreeterServer接口
type server struct {
	// pb.UnimplementedGreeterServer是protobuf生成的代码，实现了GreeterServer接口的默认实现
	pb.UnimplementedGreeterServer
}

// 实现SayHello方法，处理客户端的SayHello请求
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// 输出日志，记录接收到的请求
	log.Printf("Received: %v", in.GetName())
	// 返回一个HelloReply消息，包含了一个问候语
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// main函数是程序的入口
func main() {
	// 解析命令行参数
	flag.Parse()

	// 创建一个TCP监听器，监听指定端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	// 如果创建监听器失败，输出错误日志并退出
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 创建一个gRPC服务器
	s := grpc.NewServer()

	// 注册GreeterServer服务
	pb.RegisterGreeterServer(s, &server{})

	// 输出日志，记录服务器监听地址
	log.Printf("server listening at %v", lis.Addr())

	// 启动服务器，开始监听客户端请求
	if err := s.Serve(lis); err != nil {
		// 如果启动服务器失败，输出错误日志并退出
		log.Fatalf("failed to serve: %v", err)
	}
}
