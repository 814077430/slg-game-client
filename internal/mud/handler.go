package mud

import (
	"fmt"
	"strings"
	"time"
)

// CommandHandler 命令处理器
type CommandHandler struct {
	client  *MUDClient
	aliases map[string]string
	history []string
	histIdx int
}

// NewCommandHandler 创建命令处理器
func NewCommandHandler(client *MUDClient) *CommandHandler {
	h := &CommandHandler{
		client:  client,
		aliases: make(map[string]string),
		history: make([]string, 0),
		histIdx: -1,
	}
	
	// 注册默认别名
	h.RegisterAlias("l", "look", "查看周围")
	h.RegisterAlias("n", "go north", "向北")
	h.RegisterAlias("s", "go south", "向南")
	h.RegisterAlias("e", "go east", "向东")
	h.RegisterAlias("w", "go west", "向西")
	h.RegisterAlias("ne", "go northeast", "向东北")
	h.RegisterAlias("nw", "go northwest", "向西北")
	h.RegisterAlias("se", "go southeast", "向东南")
	h.RegisterAlias("sw", "go southwest", "向西南")
	h.RegisterAlias("i", "inventory", "背包")
	h.RegisterAlias("st", "status", "状态")
	h.RegisterAlias("h", "help", "帮助")
	h.RegisterAlias("exp", "explore", "探索")
	h.RegisterAlias("x", "scan", "扫描周边")
	h.RegisterAlias("c", "chat", "聊天")
	h.RegisterAlias("say", "chat", "说话")
	h.RegisterAlias("w", "chat world", "全服聊天")
	h.RegisterAlias("shout", "chat world", "全服喊话")
	
	return h
}

// RegisterAlias 注册别名
func (h *CommandHandler) RegisterAlias(alias, command, desc string) {
	h.aliases[alias] = command
}

// ProcessCommand 处理命令
func (h *CommandHandler) ProcessCommand(input string, ui *GameUI) string {
	input = strings.TrimSpace(input)
	
	if input == "" {
		return ""
	}
	
	// 添加到历史记录
	h.history = append(h.history, input)
	if len(h.history) > 100 {
		h.history = h.history[1:]
	}
	h.histIdx = len(h.history)
	
	// 处理别名
	parts := strings.Fields(input)
	if len(parts) > 0 {
		if alias, ok := h.aliases[parts[0]]; ok {
			input = alias
			if len(parts) > 1 {
				input += " " + strings.Join(parts[1:], " ")
			}
		}
	}
	
	// 处理特殊命令
	if strings.HasPrefix(input, "/") {
		return h.processSpecialCommand(input, ui)
	}
	
	// 处理游戏命令
	return h.processGameCommand(input, ui)
}

// processSpecialCommand 处理特殊命令（以/开头）
func (h *CommandHandler) processSpecialCommand(input string, ui *GameUI) string {
	parts := strings.Fields(input)
	cmd := strings.ToLower(parts[0])
	
	switch cmd {
	case "/quit", "/exit":
		h.client.Close()
		return "已断开连接"
		
	case "/clear":
		ui.Clear()
		return ""
		
	case "/history", "/hist":
		return h.getHistory()
		
	case "/alias":
		return h.listAliases()
		
	case "/help":
		ui.ShowHelp()
		return ""
		
	case "/map":
		ui.ShowMiniMap()
		return ""
		
	case "/refresh":
		ui.Refresh()
		return ""
		
	default:
		return fmt.Sprintf("未知命令：%s", cmd)
	}
}

// processGameCommand 处理游戏命令
func (h *CommandHandler) processGameCommand(input string, ui *GameUI) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return ""
	}
	
	cmd := strings.ToLower(parts[0])
	
	switch cmd {
	case "login":
		if len(parts) >= 3 {
			h.client.Login(parts[1], parts[2])
			return fmt.Sprintf("正在登录：%s", parts[1])
		}
		return "用法：login <用户名> <密码>"
		
	case "register":
		if len(parts) >= 3 {
			h.client.Register(parts[1], parts[2])
			return fmt.Sprintf("正在注册：%s", parts[1])
		}
		return "用法：register <用户名> <密码>"
		
	case "look", "l":
		ui.ShowMiniMap()
		ui.ShowAreaInfo()
		h.client.Send("look")
		return ""
		
	case "go":
		if len(parts) >= 2 {
			ui.Move(parts[1])
			return fmt.Sprintf("正在向 %s 移动...", parts[1])
		}
		return "用法：go <方向> (n/s/e/w/ne/nw/se/sw)"
		
	case "n", "north":
		ui.Move("north")
		return "正在向北移动..."
		
	case "s", "south":
		ui.Move("south")
		return "正在向南移动..."
		
	case "e", "east":
		ui.Move("east")
		return "正在向东移动..."
		
	case "w", "west":
		ui.Move("west")
		return "正在向西移动..."
		
	case "status", "st":
		ui.ShowPlayerInfo()
		h.client.Send("status")
		return ""
		
	case "map":
		ui.ShowMiniMap()
		return ""
		
	case "info":
		ui.ShowAreaInfo()
		return ""
		
	case "build":
		if len(parts) >= 2 {
			x := ui.GetPlayerX()
			y := ui.GetPlayerY()
			h.client.Send(fmt.Sprintf("build %s %d %d", parts[1], x, y))
			return fmt.Sprintf("正在建造 %s...", parts[1])
		}
		return "用法：build <建筑类型>"
		
	case "work":
		h.client.Send("work")
		return "正在工作..."
		
	case "rest":
		h.client.Send("rest")
		return "正在休息..."
		
	case "say":
		if len(parts) >= 2 {
			h.client.Send(fmt.Sprintf("say %s", strings.Join(parts[1:], " ")))
			return ""
		}
		return "用法：say <消息>"
		
	case "who", "players":
		h.client.Send("who")
		return "获取视野内玩家..."
		
	case "help", "h":
		ui.ShowHelp()
		return ""
		
	case "explore", "exp":
		return h.explore(ui)
		
	case "scan", "x":
		return h.scan(ui)
		
	case "travel":
		if len(parts) >= 3 {
			var x, y int32
			fmt.Sscanf(parts[1], "%d", &x)
			fmt.Sscanf(parts[2], "%d", &y)
			return h.travel(ui, x, y)
		}
		return "用法：travel <x> <y>"
		
	case "chat", "c":
		if len(parts) >= 2 {
			// 默认全服频道
			content := strings.Join(parts[1:], " ")
			return h.sendChat(ui, content, "world")
		}
		return "用法：chat <消息> 或 c <消息>"
		
	case "chatworld", "cw", "shout":
		if len(parts) >= 2 {
			content := strings.Join(parts[1:], " ")
			return h.sendChat(ui, content, "world")
		}
		return "用法：chatworld <消息> 或 cw <消息>"
		
	case "chataliance", "ca":
		if len(parts) >= 2 {
			content := strings.Join(parts[1:], " ")
			return h.sendChat(ui, content, "alliance")
		}
		return "用法：chataliance <消息> 或 ca <消息>"
		
	case "chathistory", "ch":
		return h.showChatHistory(ui)
		
	case "clearchat", "cc":
		return h.clearChatHistory(ui)
		
	default:
		// 未知命令，直接发送
		h.client.Send(input)
		return fmt.Sprintf("发送命令：%s", input)
	}
}

// sendChat 发送聊天消息
func (h *CommandHandler) sendChat(ui *GameUI, content, channel string) string {
	if len(content) == 0 {
		return "❌ 消息不能为空"
	}
	if len(content) > 500 {
		return "❌ 消息长度不能超过 500 字符"
	}
	
	// 发送聊天请求
	h.client.SendChat(content, channel)
	
	channelName := map[string]string{
		"world":     "全服",
		"alliance":  "联盟",
	}[channel]
	
	return fmt.Sprintf("✅ 已发送到 [%s] 频道", channelName)
}

// showChatHistory 显示聊天历史
func (h *CommandHandler) showChatHistory(ui *GameUI) string {
	history := ui.GetChatHistory()
	
	if len(history) == 0 {
		return "📭 暂无聊天历史"
	}
	
	output := "\n╔════════════════════════════════════════════════════════╗\n"
	output += "║                  聊天历史记录                  ║\n"
	output += "╠════════════════════════════════════════════════════════╣\n"
	
	for _, msg := range history {
		timeStr := time.UnixMilli(msg.Timestamp).Format("15:04:05")
		channelIcon := map[string]string{
			"world":    "🌍",
			"alliance": "🏰",
		}[msg.Channel]
		
		output += fmt.Sprintf("║ %s [%s] %s: %s\n", timeStr, channelIcon, msg.Username, msg.Content)
	}
	
	output += "╚════════════════════════════════════════════════════════╝\n"
	
	return output
}

// clearChatHistory 清空聊天历史
func (h *CommandHandler) clearChatHistory(ui *GameUI) string {
	ui.ClearChatHistory()
	return "✅ 聊天历史已清空"
}

// explore 探索当前区域
func (h *CommandHandler) explore(ui *GameUI) string {
	x := ui.GetPlayerX()
	y := ui.GetPlayerY()
	
	// 显示当前位置信息
	output := fmt.Sprintf("\n╔════════════════════════════════════════════════════════╗\n")
	output += fmt.Sprintf("║          当前位置：(%5d, %5d)                      ║\n", x, y)
	output += fmt.Sprintf("║          区域：%s                                  ║\n", ui.GetZoneName(x, y))
	output += fmt.Sprintf("╠════════════════════════════════════════════════════════╣\n")
	output += fmt.Sprintf("║  世界信息：                                            ║\n")
	output += fmt.Sprintf("║    世界尺寸：1024x1024                                ║\n")
	output += fmt.Sprintf("║    中心坐标：(512, 512)                               ║\n")
	output += fmt.Sprintf("║    皇城范围：(480,480)~(544,544)                      ║\n")
	output += fmt.Sprintf("╠════════════════════════════════════════════════════════╣\n")
	output += fmt.Sprintf("║  资源分布（从外到内）：                                ║\n")
	output += fmt.Sprintf("║    边缘绝境 (0 级) → 蛮荒带 (1-2 级) → 四大州 (3-4 级)    ║\n")
	output += fmt.Sprintf("║    → 中心安全区 (5 级) → 皇城 (6 级)                    ║\n")
	output += fmt.Sprintf("╚════════════════════════════════════════════════════════╝\n")
	
	return output
}

// scan 扫描周边环境
func (h *CommandHandler) scan(ui *GameUI) string {
	x := ui.GetPlayerX()
	y := ui.GetPlayerY()
	
	output := fmt.Sprintf("\n╔════════════════════════════════════════════════════════╗\n")
	output += fmt.Sprintf("║          周边环境扫描 (%5d, %5d)                  ║\n", x, y)
	output += fmt.Sprintf("╠════════════════════════════════════════════════════════╣\n")
	output += fmt.Sprintf("║  小地图：                                              ║\n")
	
	// 显示周边 5x5 区域
	for dy := -2; dy <= 2; dy++ {
		output += "║    "
		for dx := -2; dx <= 2; dx++ {
			scanX := x + int32(dx)
			scanY := y + int32(dy)
			
			if dx == 0 && dy == 0 {
				output += "\033[33m@\033[0m " // 玩家位置（黄色）
			} else {
				// 简单模拟地形符号
				zone := ui.GetZoneName(scanX, scanY)
				switch zone {
				case "皇城":
					output += "\033[31m#\033[0m " // 红色
				case "安全区":
					output += "\033[32m.\033[0m " // 绿色
				case "雍州", "青州", "扬州", "荆州":
					output += "\033[36m*\033[0m " // 青色
				case "蛮荒":
					output += "\033[33m^\033[0m " // 黄色
				default:
					output += "\033[30m=\033[0m " // 灰色
				}
			}
		}
		output += "                                ║\n"
	}
	
	output += fmt.Sprintf("╠════════════════════════════════════════════════════════╣\n")
	output += fmt.Sprintf("║  图例：@=玩家  .=安全区  *=州域  #=皇城  ^=蛮荒        ║\n")
	output += fmt.Sprintf("╚════════════════════════════════════════════════════════╝\n")
	
	return output
}

// travel 移动到指定坐标
func (h *CommandHandler) travel(ui *GameUI, x, y int32) string {
	currX := ui.GetPlayerX()
	currY := ui.GetPlayerY()
	
	// 计算移动方向
	dx := x - currX
	dy := y - currY
	
	output := fmt.Sprintf("正在从 (%d,%d) 移动到 (%d,%d)...\n", currX, currY, x, y)
	
	// 模拟移动过程
	steps := 0
	for dx != 0 || dy != 0 {
		if dx > 0 {
			ui.Move("east")
			dx--
		} else if dx < 0 {
			ui.Move("west")
			dx++
		}
		
		if dy > 0 {
			ui.Move("south")
			dy--
		} else if dy < 0 {
			ui.Move("north")
			dy++
		}
		
		steps++
		if steps > 10 {
			break // 限制最大步数
		}
	}
	
	return output
}

// getHistory 获取历史记录
func (h *CommandHandler) getHistory() string {
	if len(h.history) == 0 {
		return "没有历史记录"
	}
	
	var sb strings.Builder
	sb.WriteString("=== 命令历史 ===\n")
	for i, cmd := range h.history {
		if i >= len(h.history)-20 { // 只显示最近 20 条
			sb.WriteString(fmt.Sprintf("%3d: %s\n", i+1, cmd))
		}
	}
	return sb.String()
}

// listAliases 列出别名
func (h *CommandHandler) listAliases() string {
	var sb strings.Builder
	sb.WriteString("=== 命令别名 ===\n")
	for alias, cmd := range h.aliases {
		sb.WriteString(fmt.Sprintf("  %-5s -> %s\n", alias, cmd))
	}
	return sb.String()
}

// GetPreviousCommand 获取上一条命令
func (h *CommandHandler) GetPreviousCommand() string {
	if h.histIdx > 0 {
		h.histIdx--
		return h.history[h.histIdx]
	}
	return ""
}

// GetNextCommand 获取下一条命令
func (h *CommandHandler) GetNextCommand() string {
	if h.histIdx < len(h.history)-1 {
		h.histIdx++
		return h.history[h.histIdx]
	}
	h.histIdx = len(h.history)
	return ""
}
