package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/clin211/grpc/service-types/go/rpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// chatService 实现
type chatService struct {
	pb.UnimplementedChatServiceServer

	// 存储所有连接的客户端
	clients map[string]*clientConnection
	// 用于保护clients map的互斥锁
	clientsMutex sync.RWMutex

	// 消息广播通道
	broadcast chan *pb.ChatMessage
}

// clientConnection 表示一个客户端连接
type clientConnection struct {
	userID   string
	username string
	stream   pb.ChatService_ChatServer
	send     chan *pb.ChatMessage
}

// newChatService 创建聊天服务实例
func newChatService() *chatService {
	service := &chatService{
		clients:   make(map[string]*clientConnection),
		broadcast: make(chan *pb.ChatMessage, 100),
	}

	// 启动消息广播处理器
	go service.handleBroadcast()

	return service
}

// Chat 实现双向流式 RPC
func (s *chatService) Chat(stream pb.ChatService_ChatServer) error {
	var client *clientConnection

	defer func() {
		// 清理客户端连接
		if client != nil {
			s.removeClient(client)

			// 发送用户离开通知
			leaveMsg := &pb.ChatMessage{
				MessageId: uuid.New().String(),
				UserId:    "system",
				Username:  "系统",
				Content:   fmt.Sprintf("用户 %s 离开了聊天室", client.username),
				Timestamp: time.Now().Unix(),
				Type:      pb.MessageType_USER_LEAVE,
				RoomId:    "general",
			}
			s.broadcast <- leaveMsg
		}
	}()

	// 处理客户端消息
	for {
		// 接收客户端发送的消息
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Client %s disconnected", getClientID(client))
			break
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		// 处理第一条消息（用户加入）
		if client == nil {
			client = &clientConnection{
				userID:   msg.UserId,
				username: msg.Username,
				stream:   stream,
				send:     make(chan *pb.ChatMessage, 10),
			}

			// 添加客户端到连接池
			s.addClient(client)

			// 启动消息发送处理器
			go s.handleClientSend(client)

			log.Printf("User %s (%s) joined the chat room", client.username, client.userID)

			// 发送用户加入通知
			joinMsg := &pb.ChatMessage{
				MessageId: uuid.New().String(),
				UserId:    "system",
				Username:  "系统",
				Content:   fmt.Sprintf("欢迎 %s 加入聊天室！", client.username),
				Timestamp: time.Now().Unix(),
				Type:      pb.MessageType_USER_JOIN,
				RoomId:    "general",
			}
			s.broadcast <- joinMsg

			// 如果第一条消息就是文本消息，也要处理
			if msg.Type == pb.MessageType_TEXT {
				s.handleTextMessage(msg, client)
			}
		} else {
			// 处理后续消息
			s.handleMessage(msg, client)
		}
	}

	return nil
}

// handleMessage 处理接收到的消息
func (s *chatService) handleMessage(msg *pb.ChatMessage, client *clientConnection) {
	switch msg.Type {
	case pb.MessageType_TEXT:
		s.handleTextMessage(msg, client)
	default:
		log.Printf("Unknown message type: %v", msg.Type)
	}
}

// handleTextMessage 处理文本消息
func (s *chatService) handleTextMessage(msg *pb.ChatMessage, client *clientConnection) {
	// 设置消息元数据
	msg.MessageId = uuid.New().String()
	msg.UserId = client.userID
	msg.Username = client.username
	msg.Timestamp = time.Now().Unix()
	msg.RoomId = "general"

	log.Printf("Message from %s: %s", client.username, msg.Content)

	// 广播消息给所有客户端
	s.broadcast <- msg
}

// addClient 添加客户端到连接池
func (s *chatService) addClient(client *clientConnection) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	s.clients[client.userID] = client
	log.Printf("Added client %s, total clients: %d", client.username, len(s.clients))
}

// removeClient 从连接池移除客户端
func (s *chatService) removeClient(client *clientConnection) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	if _, exists := s.clients[client.userID]; exists {
		delete(s.clients, client.userID)
		close(client.send)
		log.Printf("Removed client %s, total clients: %d", client.username, len(s.clients))
	}
}

// handleBroadcast 处理消息广播
func (s *chatService) handleBroadcast() {
	for msg := range s.broadcast {
		s.clientsMutex.RLock()
		for _, client := range s.clients {
			select {
			case client.send <- msg:
				// 消息发送成功
			default:
				// 客户端发送通道已满，跳过该客户端
				log.Printf("Client %s send channel is full, skipping message", client.username)
			}
		}
		s.clientsMutex.RUnlock()
	}
}

// handleClientSend 处理单个客户端的消息发送
func (s *chatService) handleClientSend(client *clientConnection) {
	for msg := range client.send {
		if err := client.stream.Send(msg); err != nil {
			log.Printf("Error sending message to client %s: %v", client.username, err)
			break
		}
	}
}

// getClientID 获取客户端标识（用于日志）
func getClientID(client *clientConnection) string {
	if client == nil {
		return "unknown"
	}
	return fmt.Sprintf("%s(%s)", client.username, client.userID)
}

func main() {
	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册聊天服务
	chatSvc := newChatService()
	pb.RegisterChatServiceServer(server, chatSvc)

	// 监听端口
	lis, err := net.Listen("tcp", ":6004")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("💬 Chat Room Server started on :6004")
	log.Println("Waiting for users to join the chat room...")

	// 启动服务
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
