package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/clin211/grpc/service-types/go/rpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 建立连接
	conn, err := grpc.NewClient("localhost:6004",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建聊天服务客户端
	client := pb.NewChatServiceClient(conn)

	// 获取用户名
	fmt.Print("请输入您的用户名: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	if username == "" {
		username = "匿名用户"
	}

	// 生成用户ID
	userID := uuid.New().String()

	log.Printf("正在以用户名 '%s' 加入聊天室...", username)

	// 创建双向流
	ctx := context.Background()
	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("Failed to create chat stream: %v", err)
	}

	// 发送加入消息
	joinMsg := &pb.ChatMessage{
		UserId:   userID,
		Username: username,
		Content:  "", // 加入消息内容为空
		Type:     pb.MessageType_USER_JOIN,
		RoomId:   "general",
	}

	if err := stream.Send(joinMsg); err != nil {
		log.Fatalf("Failed to send join message: %v", err)
	}

	// 启动消息接收处理器
	go handleReceiveMessages(stream)

	// 显示聊天界面
	displayChatInterface(username)

	// 处理用户输入
	handleUserInput(stream, userID, username, scanner)
}

// handleReceiveMessages 处理接收消息
func handleReceiveMessages(stream pb.ChatService_ChatClient) {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("\n💔 与服务器的连接已断开")
			os.Exit(0)
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}

		displayMessage(msg)
	}
}

// handleUserInput 处理用户输入
func handleUserInput(stream pb.ChatService_ChatClient, userID, username string, scanner *bufio.Scanner) {
	for {
		fmt.Print("💬 ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// 检查退出命令
		if input == "/quit" || input == "/exit" {
			fmt.Println("👋 再见！")
			break
		}

		// 检查帮助命令
		if input == "/help" {
			displayHelp()
			continue
		}

		// 跳过空消息
		if input == "" {
			continue
		}

		// 发送文本消息
		msg := &pb.ChatMessage{
			UserId:   userID,
			Username: username,
			Content:  input,
			Type:     pb.MessageType_TEXT,
			RoomId:   "general",
		}

		if err := stream.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
			break
		}
	}
}

// displayMessage 格式化显示消息
func displayMessage(msg *pb.ChatMessage) {
	timestamp := time.Unix(msg.Timestamp, 0).Format("15:04:05")

	switch msg.Type {
	case pb.MessageType_TEXT:
		fmt.Printf("\r[%s] %s: %s\n💬 ", timestamp, msg.Username, msg.Content)
	case pb.MessageType_USER_JOIN:
		fmt.Printf("\r🎉 [%s] %s\n💬 ", timestamp, msg.Content)
	case pb.MessageType_USER_LEAVE:
		fmt.Printf("\r👋 [%s] %s\n💬 ", timestamp, msg.Content)
	case pb.MessageType_SYSTEM:
		fmt.Printf("\r📢 [%s] %s\n💬 ", timestamp, msg.Content)
	default:
		fmt.Printf("\r❓ [%s] %s: %s\n💬 ", timestamp, msg.Username, msg.Content)
	}
}

// displayChatInterface 显示聊天界面
func displayChatInterface(username string) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("💬 欢迎来到 gRPC 聊天室，%s！\n", username)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📝 输入消息后按回车发送")
	fmt.Println("💡 输入 /help 查看帮助，输入 /quit 退出聊天室")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
}

// displayHelp 显示帮助信息
func displayHelp() {
	fmt.Println("\r" + strings.Repeat("=", 50))
	fmt.Println("📖 聊天室帮助")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("💬 直接输入文字发送消息")
	fmt.Println("❓ /help  - 显示此帮助信息")
	fmt.Println("👋 /quit  - 退出聊天室")
	fmt.Println("👋 /exit  - 退出聊天室")
	fmt.Println(strings.Repeat("=", 50))
}
