package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	pb "github.com/clin211/grpc/service-types/go/rpc"
	"google.golang.org/grpc"
)

// fileService 实现
type fileService struct {
	pb.UnimplementedFileServiceServer
	uploadDir string // 文件上传目录
}

// newFileService 创建文件服务实例
func newFileService() *fileService {
	uploadDir := "./uploads"
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Failed to create upload directory: %v", err)
	}

	return &fileService{
		uploadDir: uploadDir,
	}
}

// UploadFile 实现客户端流式 RPC
func (s *fileService) UploadFile(stream pb.FileService_UploadFileServer) error {
	startTime := time.Now()

	var fileID string
	var filename string
	var totalChunks int32
	var receivedChunks int32
	var totalSize int64
	var fileHash string

	// 用于存储接收到的文件块
	var fileData []byte

	log.Println("Starting file upload...")

	for {
		// 从流中接收文件块
		chunk, err := stream.Recv()
		if err == io.EOF {
			// 客户端已发送完所有块
			log.Printf("File upload completed for %s", filename)
			break
		}
		if err != nil {
			log.Printf("Failed to receive chunk: %v", err)
			return err
		}

		// 处理第一个块，获取文件基本信息
		if receivedChunks == 0 {
			fileID = chunk.FileId
			filename = chunk.Filename
			totalChunks = chunk.TotalChunks
			fileHash = chunk.FileHash

			log.Printf("Receiving file: %s (ID: %s, Total chunks: %d)",
				filename, fileID, totalChunks)
		}

		// 验证文件块信息
		if chunk.FileId != fileID {
			return fmt.Errorf("file ID mismatch: expected %s, got %s",
				fileID, chunk.FileId)
		}

		if chunk.ChunkNumber != receivedChunks {
			return fmt.Errorf("chunk number mismatch: expected %d, got %d",
				receivedChunks, chunk.ChunkNumber)
		}

		// 将块数据添加到文件数据中
		fileData = append(fileData, chunk.Data...)
		totalSize += int64(chunk.ChunkSize)
		receivedChunks++

		log.Printf("Received chunk %d/%d (size: %d bytes)",
			chunk.ChunkNumber+1, totalChunks, chunk.ChunkSize)

		// 如果是最后一块，验证总数
		if chunk.IsLast {
			if receivedChunks != totalChunks {
				log.Printf("Warning: received chunks (%d) != total chunks (%d)",
					receivedChunks, totalChunks)
			}
			break
		}
	}

	// 验证文件完整性（如果提供了哈希值）
	if fileHash != "" {
		actualHash := fmt.Sprintf("%x", md5.Sum(fileData))
		if actualHash != fileHash {
			return fmt.Errorf("file hash mismatch: expected %s, got %s",
				fileHash, actualHash)
		}
		log.Printf("File hash verification passed: %s", actualHash)
	}

	// 保存文件到磁盘
	filePath := filepath.Join(s.uploadDir, filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		log.Printf("Failed to save file: %v", err)
		return err
	}

	uploadDuration := time.Since(startTime)

	// 发送响应
	response := &pb.FileUploadResponse{
		Success:           true,
		Message:           fmt.Sprintf("File '%s' uploaded successfully", filename),
		FilePath:          filePath,
		FileSize:          totalSize,
		FileId:            fileID,
		ChunksReceived:    receivedChunks,
		UploadTimeSeconds: uploadDuration.Seconds(),
	}

	log.Printf("File saved to: %s (Size: %d bytes, Duration: %.2fs)",
		filePath, totalSize, uploadDuration.Seconds())

	return stream.SendAndClose(response)
}

func main() {
	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册文件服务
	fileSvc := newFileService()
	pb.RegisterFileServiceServer(server, fileSvc)

	// 监听端口
	lis, err := net.Listen("tcp", ":6003")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("File Upload Service Server started on :6003")
	log.Printf("Upload directory: %s", fileSvc.uploadDir)

	// 启动服务
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
