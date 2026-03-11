package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"slg-game-client/internal/mud"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║          SLG MUD Client                                ║")
	fmt.Println("║          Text-Based SLG Game Client                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 获取服务器地址
	serverAddr := "localhost:8080"
	if len(os.Args) > 1 {
		serverAddr = os.Args[1]
	}

	// 创建 MUD 客户端
	client := mud.NewMUDClient(serverAddr)
	handler := mud.NewCommandHandler(client)

	// 连接服务器
	fmt.Printf("正在连接到 %s...\n", serverAddr)
	if err := client.Connect(); err != nil {
		fmt.Printf("❌ 连接失败：%v\n", err)
		fmt.Println("请确保服务器正在运行 (slg-server)")
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("✅ 连接成功！")
	fmt.Println()
	fmt.Println("输入 'help' 查看帮助，'/quit' 退出游戏")
	fmt.Println()

	// 启动输出处理
	go handleOutput(client)

	// 主输入循环
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ 读取错误：%v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		
		// 处理命令
		result := handler.ProcessCommand(input)
		
		if result != "" {
			// 特殊命令直接输出
			if strings.HasPrefix(input, "/") {
				fmt.Println(result)
			}
		}

		// 检查是否已断开
		if !client.IsConnected() {
			fmt.Println("已断开连接")
			break
		}
	}

	fmt.Println("感谢游玩，再见！")
	time.Sleep(1 * time.Second)
}

// handleOutput 处理服务器输出
func handleOutput(client *mud.MUDClient) {
	for line := range client.GetOutputChan() {
		// 解析并格式化输出
		formatted := formatOutput(line)
		if formatted != "" {
			fmt.Println(formatted)
		}
	}
}

// formatOutput 格式化输出
func formatOutput(line string) string {
	// 解析特殊格式
	if strings.Contains(line, "成功") {
		return "✅ " + line
	}
	
	if strings.Contains(line, "失败") || strings.Contains(line, "错误") {
		return "❌ " + line
	}
	
	if strings.Contains(line, "欢迎") {
		return "🎮 " + line
	}
	
	if strings.HasPrefix(line, "位置:") {
		return "📍 " + line
	}
	
	if strings.HasPrefix(line, "资源:") {
		return "💰 " + line
	}
	
	return line
}
