package mud

import (
	"fmt"
	"strings"
)

// CommandAlias 命令别名
type CommandAlias struct {
	Name    string
	Command string
	Desc    string
}

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
func (h *CommandHandler) ProcessCommand(input string) string {
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
		return h.processSpecialCommand(input)
	}
	
	// 发送游戏命令
	h.client.Send(input)
	return input
}

// processSpecialCommand 处理特殊命令（以/开头）
func (h *CommandHandler) processSpecialCommand(input string) string {
	parts := strings.Fields(input)
	cmd := strings.ToLower(parts[0])
	
	switch cmd {
	case "/quit", "/exit":
		h.client.Close()
		return "已断开连接"
		
	case "/clear":
		return "\033[2J\033[H" // 清屏
		
	case "/history", "/hist":
		return h.getHistory()
		
	case "/alias":
		return h.listAliases()
		
	case "/help":
		return h.getHelp()
		
	default:
		return fmt.Sprintf("未知命令：%s", cmd)
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

// getHelp 获取帮助
func (h *CommandHandler) getHelp() string {
	return `
╔════════════════════════════════════════════════════════╗
║                    SLG MUD 帮助                        ║
╠════════════════════════════════════════════════════════╣
║  游戏命令：                                             ║
║    login <用户名> <密码>    登录                        ║
║    register <用户名> <密码> 注册                        ║
║    look (l)               查看周围                      ║
║    go <方向>              移动 (north/south/east/west)  ║
║    status (st)            查看状态                      ║
║    inventory (i)          查看背包                      ║
║    build <建筑>           建造建筑                      ║
║    work                   工作                          ║
║    rest                   休息                          ║
║    say <消息>             说话                          ║
║    who                    在线玩家                      ║
║    help                   帮助                          ║
╠════════════════════════════════════════════════════════╣
║  特殊命令（以/开头）：                                   ║
║    /quit, /exit           退出游戏                      ║
║    /clear                 清屏                          ║
║    /history               命令历史                      ║
║    /alias                 别名列表                      ║
║    /help                  本帮助                        ║
╚════════════════════════════════════════════════════════╝
`
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
