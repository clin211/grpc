package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	pb "github.com/clin211/grpc/service-types/go/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 建立连接
	conn, err := grpc.NewClient("localhost:6002",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建股票服务客户端
	client := pb.NewStockServiceClient(conn)

	// 准备订阅请求
	subscribeReq := &pb.StockSubscribeRequest{
		Symbols:  []string{"AAPL", "GOOGL", "TSLA"}, // 订阅苹果、谷歌、特斯拉的股票
		ClientId: "client_001",
	}

	log.Printf("Subscribing to stocks: %v", subscribeReq.Symbols)

	// 发起服务端流式 RPC 调用
	ctx := context.Background()
	stream, err := client.SubscribeStockPrice(ctx, subscribeReq)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// 接收价格更新流
	fmt.Println("\n📈 Stock Price Updates (Press Ctrl+C to exit):")
	fmt.Println(strings.Repeat("=", 70))

	updateCount := 0
	for {
		// 从流中接收价格更新
		update, err := stream.Recv()
		if err == io.EOF {
			// 服务端关闭了流
			log.Println("Stream ended by server")
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive update: %v", err)
		}

		// 处理和显示价格更新
		updateCount++
		displayPriceUpdate(update, updateCount)
	}

	log.Println("Stock price subscription ended")
}

// displayPriceUpdate 格式化显示股票价格更新
func displayPriceUpdate(update *pb.StockPriceUpdate, count int) {
	// 格式化时间戳
	timestamp := time.Unix(update.Timestamp, 0).Format("15:04:05")

	// 格式化价格变化显示
	var changeSymbol string
	var changeColor string
	if update.ChangePercent > 0 {
		changeSymbol = "📈"
		changeColor = "+" // 涨
	} else if update.ChangePercent < 0 {
		changeSymbol = "📉"
		changeColor = "" // 跌（负号已包含）
	} else {
		changeSymbol = "➡️"
		changeColor = " " // 平
	}

	// 格式化成交量（以K为单位显示）
	volumeK := float64(update.Volume) / 1000

	// 计算股票代码的填充空格（保证对齐）
	padding := ""
	symbolLen := len(update.Symbol)
	if symbolLen < 6 {
		padding = strings.Repeat(" ", 6-symbolLen)
	}

	// 显示格式化的价格信息
	fmt.Printf("[%d] %s %s%s | $%.2f | %s%.2f%% ($%.2f) | Vol: %.0fK | %s\n",
		count,
		changeSymbol,
		update.Symbol,
		padding, // 对齐
		update.CurrentPrice,
		changeColor,
		update.ChangePercent,
		update.ChangeAmount,
		volumeK,
		timestamp,
	)
}
