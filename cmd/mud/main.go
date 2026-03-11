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
	fmt.Println("║          SLG MUD Client - World Explorer               ║")
	fmt.Println("║          1024x1024 大世界探索                           ║")
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
	ui := mud.NewGameUI()

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

	// 启动输出处理
	go handleOutput(client, ui)

	// 显示欢迎界面
	showWelcome(ui)

	// 主输入循环
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n\033[32m>\033[0m ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ 读取错误：%v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		
		// 处理命令
		result := handler.ProcessCommand(input, ui)
		
		if result != "" {
			// 特殊命令直接输出
			if strings.HasPrefix(input, "/") {
				fmt.Println(result)
			}
		}

		// 检查是否已断开
		if !client.IsConnected() {
			fmt.Println("\n已断开连接")
			break
		}
	}

	fmt.Println("\n感谢游玩，再见！")
	time.Sleep(1 * time.Second)
}

// showWelcome 显示欢迎界面
func showWelcome(ui *mud.GameUI) {
	ui.Clear()
	
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║          欢迎来到 SLG 大世界                            ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  世界尺寸：1024x1024 (1,048,576 地块)                  ║")
	fmt.Println("║  中心坐标：(512, 512)                                  ║")
	fmt.Println("║  皇城范围：(480,480)~(544,544)                         ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  资源分布（从外到内）：                                ║")
	fmt.Println("║    边缘绝境 (0 级) → 蛮荒带 (1-2 级) → 四大州 (3-4 级)    ║")
	fmt.Println("║    → 中心安全区 (5 级) → 皇城 (6 级)                    ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  输入 'help' 查看帮助，'explore' 开始探索              ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// handleOutput 处理服务器输出
func handleOutput(client *mud.MUDClient, ui *mud.GameUI) {
	for line := range client.GetOutputChan() {
		// 解析并格式化输出
		formatted := formatOutput(line, ui)
		if formatted != "" {
			fmt.Println(formatted)
		}
	}
}

// formatOutput 格式化输出
func formatOutput(line string, ui *mud.GameUI) string {
	// 解析特殊格式
	if strings.Contains(line, "成功") {
		return "\033[32m✅ " + line + "\033[0m"
	}
	
	if strings.Contains(line, "失败") || strings.Contains(line, "错误") {
		return "\033[31m❌ " + line + "\033[0m"
	}
	
	if strings.Contains(line, "欢迎") {
		return "\033[33m🎮 " + line + "\033[0m"
	}
	
	if strings.HasPrefix(line, "位置:") {
		return "\033[36m📍 " + line + "\033[0m"
	}
	
	if strings.HasPrefix(line, "资源:") {
		return "\033[33m💰 " + line + "\033[0m"
	}
	
	if strings.HasPrefix(line, "区域:") {
		return "\033[35m🗺️ " + line + "\033[0m"
	}
	
	if strings.HasPrefix(line, "地形:") {
		return "\033[32m🌲 " + line + "\033[0m"
	}
	
	if strings.HasPrefix(line, "玩家 ID:") {
		parts := strings.Split(line, ":")
		if len(parts) > 1 {
			var id uint64
			fmt.Sscanf(parts[1], "%d", &id)
			ui.UpdatePlayer(&mud.PlayerInfo{ID: id})
		}
		return "\033[34m👤 " + line + "\033[0m"
	}
	
	return line
}
