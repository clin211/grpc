package main

import (
	"fmt"

	"google.golang.org/grpc/metadata"
)

func main() {
	basicMetadataExample()
}

func basicMetadataExample() {
	// 创建元数据
	md := metadata.Pairs(
		"user-id", "12345",
		"session-token", "abc123xyz",
		"client-version", "1.2.0",
	)

	// 查看元数据内容
	fmt.Printf("元数据内容: %v\n", md) // 元数据内容: map[client-version:[1.2.0] session-token:[abc123xyz] user-id:[12345]]

	// 获取特定键的值
	userIDs := md.Get("user-id")
	if len(userIDs) > 0 {
		fmt.Printf("用户ID: %s\n", userIDs[0]) // 用户ID: 12345
	}

	// 添加更多值到同一个键
	md.Append("user-role", "admin", "user")
	fmt.Printf("用户角色: %v\n", md.Get("user-role")) // 用户角色: [admin user]
}
