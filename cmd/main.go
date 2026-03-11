package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"slg-game-client/internal/client"
)

func main() {
	serverAddr := flag.String("server", "localhost:8080", "Server address")
	username := flag.String("username", "", "Username")
	password := flag.String("password", "password123", "Password")
	flag.Parse()

	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║          SLG Game Client                               ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 创建客户端
	c := client.NewClient(*serverAddr)

	// 连接服务器
	fmt.Printf("Connecting to %s...\n", *serverAddr)
	if err := c.Connect(); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	fmt.Println("✅ Connected!")
	defer c.Close()

	// 登录
	if *username != "" {
		fmt.Printf("Logging in as %s...\n", *username)
		resp, err := c.Login(*username, *password)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}

		if resp.Success {
			fmt.Printf("✅ Login successful! Player ID: %d\n", resp.PlayerId)
		} else {
			log.Fatalf("Login failed: %s", resp.Message)
		}
	}

	// 示例操作
	fmt.Println()
	fmt.Println("=== Testing Operations ===")

	// 移动测试
	fmt.Println("Testing Move...")
	moveResp, err := c.Move(100, 200)
	if err != nil {
		fmt.Printf("❌ Move failed: %v\n", err)
	} else if moveResp.Success {
		fmt.Printf("✅ Move successful! Position: (%d, %d)\n", moveResp.X, moveResp.Y)
	} else {
		fmt.Printf("❌ Move failed: %s\n", moveResp.Message)
	}

	time.Sleep(1 * time.Second)

	// 建造测试
	fmt.Println("Testing Build...")
	buildResp, err := c.Build("farm", 50, 50)
	if err != nil {
		fmt.Printf("❌ Build failed: %v\n", err)
	} else if buildResp.Success {
		fmt.Printf("✅ Build successful!\n")
	} else {
		fmt.Printf("❌ Build failed: %s\n", buildResp.Message)
	}

	fmt.Println()
	fmt.Println("=== Client Ready ===")
	fmt.Println("Press Ctrl+C to exit...")

	// 保持连接
	select {}
}
