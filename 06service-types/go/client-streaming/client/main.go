package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/clin211/grpc/service-types/go/rpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	chunkSize = 1024 * 1024 // 1MB 每块
)

func main() {
	// 建立连接
	conn, err := grpc.NewClient("localhost:6003",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建文件服务客户端
	client := pb.NewFileServiceClient(conn)

	// 模拟上传一个测试文件
	testFilename := "test_upload.txt"
	if err := createTestFile(testFilename); err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFilename) // 清理测试文件

	log.Printf("Uploading file: %s", testFilename)

	// 上传文件
	if err := uploadFile(client, testFilename); err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	log.Println("File upload completed successfully!")
}

// createTestFile 创建一个测试文件
func createTestFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入测试数据（约 3MB）
	testData := "Hello, this is a test file for gRPC client streaming upload demo.\n"
	for i := 0; i < 50000; i++ {
		file.WriteString(fmt.Sprintf("[%d] %s", i, testData))
	}

	return nil
}

// uploadFile 上传文件的主要逻辑
func uploadFile(client pb.FileServiceClient, filename string) error {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	fileSize := fileInfo.Size()
	totalChunks := int32((fileSize + chunkSize - 1) / chunkSize) // 向上取整

	log.Printf("File size: %d bytes, Chunk size: %d bytes, Total chunks: %d",
		fileSize, chunkSize, totalChunks)

	// 计算文件哈希值
	fileHash, err := calculateFileHash(filename)
	if err != nil {
		log.Printf("Failed to calculate file hash: %v", err)
		fileHash = "" // 继续上传，不验证哈希
	}

	// 生成文件ID
	fileID := uuid.New().String()

	// 创建客户端流
	ctx := context.Background()
	stream, err := client.UploadFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to create upload stream: %v", err)
	}

	// 分块发送文件
	buffer := make([]byte, chunkSize)
	chunkNumber := int32(0)

	log.Println("Starting file upload...")
	startTime := time.Now()

	for {
		// 读取文件块
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file: %v", err)
		}

		if bytesRead == 0 {
			break // 文件读取完毕
		}

		// 判断是否为最后一块
		isLast := (chunkNumber == totalChunks-1) || err == io.EOF

		// 创建文件块消息
		chunk := &pb.FileChunk{
			FileId:      fileID,
			Filename:    filepath.Base(filename),
			ChunkNumber: chunkNumber,
			TotalChunks: totalChunks,
			Data:        buffer[:bytesRead],
			ChunkSize:   int32(bytesRead),
			IsLast:      isLast,
			FileHash:    fileHash,
		}

		// 发送文件块
		if err := stream.Send(chunk); err != nil {
			return fmt.Errorf("failed to send chunk %d: %v", chunkNumber, err)
		}

		chunkNumber++
		log.Printf("Sent chunk %d/%d (size: %d bytes)", chunkNumber, totalChunks, bytesRead)

		// 如果是最后一块，结束循环
		if isLast {
			break
		}
	}

	// 关闭发送流并接收响应
	response, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to receive upload response: %v", err)
	}

	uploadDuration := time.Since(startTime)

	// 显示上传结果
	displayUploadResult(response, uploadDuration)

	return nil
}

// calculateFileHash 计算文件的MD5哈希值
func calculateFileHash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// displayUploadResult 显示上传结果
func displayUploadResult(response *pb.FileUploadResponse, uploadDuration time.Duration) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📁 File Upload Result")
	fmt.Println(strings.Repeat("=", 60))

	if response.Success {
		fmt.Printf("✅ Status: %s\n", response.Message)
		fmt.Printf("📄 File ID: %s\n", response.FileId)
		fmt.Printf("💾 Server Path: %s\n", response.FilePath)
		fmt.Printf("📊 File Size: %d bytes (%.2f MB)\n",
			response.FileSize, float64(response.FileSize)/(1024*1024))
		fmt.Printf("📦 Chunks Received: %d\n", response.ChunksReceived)
		fmt.Printf("⏱️  Server Process Time: %.2f seconds\n", response.UploadTimeSeconds)
		fmt.Printf("🚀 Client Total Time: %.2f seconds\n", uploadDuration.Seconds())

		// 计算传输速度
		speedMBps := float64(response.FileSize) / (1024 * 1024) / uploadDuration.Seconds()
		fmt.Printf("📈 Transfer Speed: %.2f MB/s\n", speedMBps)
	} else {
		fmt.Printf("❌ Upload failed: %s\n", response.Message)
	}

	fmt.Println(strings.Repeat("=", 60))
}
