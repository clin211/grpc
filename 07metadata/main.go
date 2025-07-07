package main

import (
	"fmt"

	"google.golang.org/grpc/metadata"
)

func main() {
	md := metadata.New(nil)

	md.Set("user-role", "user")
	md.Set("user-role", "admin")
	fmt.Println(md.Get("user-role"))

	md.Append("user-permission", "read")
	md.Append("user-permission", "write", "delete")
	fmt.Println(md.Get("user-permission"))
}
