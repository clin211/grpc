package main

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "github.com/clin211/grpc/service-types/go/rpc"
	"google.golang.org/grpc"
)

// StockService 实现
type stockService struct {
	pb.UnimplementedStockServiceServer
	// 存储股票的基础价格，用于模拟价格变化
	stockPrices map[string]float64
	mutex       sync.RWMutex
}

// 初始化股票基础价格
func newStockService() *stockService {
	return &stockService{
		stockPrices: map[string]float64{
			"AAPL":  150.00,  // 苹果
			"GOOGL": 2800.00, // 谷歌
			"TSLA":  250.00,  // 特斯拉
			"MSFT":  300.00,  // 微软
			"AMZN":  3200.00, // 亚马逊
			"META":  280.00,  // Meta
			"NVDA":  400.00,  // 英伟达
		},
	}
}

// SubscribeStockPrice 实现服务端流式 RPC
func (s *stockService) SubscribeStockPrice(req *pb.StockSubscribeRequest, stream pb.StockService_SubscribeStockPriceServer) error {
	log.Printf("Client %s subscribed to stocks: %v", req.ClientId, req.Symbols)

	// 验证股票代码
	validSymbols := make([]string, 0)
	for _, symbol := range req.Symbols {
		s.mutex.RLock()
		if _, exists := s.stockPrices[symbol]; exists {
			validSymbols = append(validSymbols, symbol)
		} else {
			log.Printf("Invalid symbol: %s", symbol)
		}
		s.mutex.RUnlock()
	}

	if len(validSymbols) == 0 {
		log.Printf("No valid symbols for client %s", req.ClientId)
		return nil
	}

	// 为每个有效的股票代码发送初始价格
	for _, symbol := range validSymbols {
		s.mutex.RLock()
		basePrice := s.stockPrices[symbol]
		s.mutex.RUnlock()

		update := &pb.StockPriceUpdate{
			Symbol:        symbol,
			CurrentPrice:  basePrice,
			ChangeAmount:  0.0,
			ChangePercent: 0.0,
			Timestamp:     time.Now().Unix(),
			Volume:        rand.Int63n(1000000) + 100000, // 随机成交量 100K-1.1M
		}

		if err := stream.Send(update); err != nil {
			log.Printf("Failed to send initial price for %s: %v", symbol, err)
			return err
		}
	}

	// 持续推送价格更新
	ticker := time.NewTicker(2 * time.Second) // 每2秒更新一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 为每个订阅的股票生成价格更新
			for _, symbol := range validSymbols {
				update := s.generatePriceUpdate(symbol)
				if err := stream.Send(update); err != nil {
					log.Printf("Failed to send price update for %s: %v", symbol, err)
					return err
				}
			}

		case <-stream.Context().Done():
			// 客户端断开连接
			log.Printf("Client %s disconnected", req.ClientId)
			return nil
		}
	}
}

// generatePriceUpdate 生成模拟的股票价格更新
func (s *stockService) generatePriceUpdate(symbol string) *pb.StockPriceUpdate {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	basePrice := s.stockPrices[symbol]

	// 生成 -5% 到 +5% 的随机价格变化
	changePercent := (rand.Float64() - 0.5) * 10 // -5% 到 +5%
	changeAmount := basePrice * changePercent / 100
	newPrice := basePrice + changeAmount

	// 更新存储的价格（模拟真实价格变化）
	s.stockPrices[symbol] = newPrice

	return &pb.StockPriceUpdate{
		Symbol:        symbol,
		CurrentPrice:  newPrice,
		ChangeAmount:  changeAmount,
		ChangePercent: changePercent,
		Timestamp:     time.Now().Unix(),
		Volume:        rand.Int63n(1000000) + 100000, // 随机成交量
	}
}

func main() {
	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册股票服务
	stockSvc := newStockService()
	pb.RegisterStockServiceServer(server, stockSvc)

	// 监听端口
	lis, err := net.Listen("tcp", ":6002")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Stock Service Server started on :6002")
	log.Println("Available stocks: AAPL, GOOGL, TSLA, MSFT, AMZN, META, NVDA")

	// 启动服务
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
