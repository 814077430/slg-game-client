package mud

import (
	"fmt"
	"strings"
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
		
	case "who":
		h.client.Send("who")
		return "获取在线玩家列表..."
		
	case "help", "h":
		ui.ShowHelp()
		return ""
		
	default:
		// 未知命令，直接发送
		h.client.Send(input)
		return fmt.Sprintf("发送命令：%s", input)
	}
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
