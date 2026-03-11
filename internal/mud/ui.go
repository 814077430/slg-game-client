package mud

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// ChatMessage 聊天消息
type ChatMessage struct {
	PlayerID  uint64 `json:"player_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Channel   string `json:"channel"` // "world" / "alliance"
}

// GameUI 游戏界面
type GameUI struct {
	player      *PlayerInfo
	history     []*ChatMessage
	historyLock sync.RWMutex
	maxHistory  int
	width       int
	height      int
}

// PlayerInfo 玩家信息
type PlayerInfo struct {
	ID       uint64
	Username string
	X        int32
	Y        int32
	Zone     string
	Gold     int64
	Wood     int64
	Food     int64
	Stone    int64
	Level    int32
}

// NewGameUI 创建游戏界面
func NewGameUI() *GameUI {
	return &GameUI{
		player:     &PlayerInfo{},
		history:    make([]*ChatMessage, 0),
		maxHistory: 100, // 最多保留 100 条聊天历史
		width:      80,
		height:     24,
	}
}

// Clear 清屏
func (ui *GameUI) Clear() {
	fmt.Print("\033[2J\033[H")
}

// ShowHeader 显示头部
func (ui *GameUI) ShowHeader() {
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║          SLG MUD Client - 1024x1024                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
}

// ShowPlayerInfo 显示玩家信息
func (ui *GameUI) ShowPlayerInfo() {
	ui.historyLock.RLock()
	defer ui.historyLock.RUnlock()
	
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║  玩家信息                                              ║")
	fmt.Printf("║  ID: %-10d  坐标：(%5d, %5d)  区域：%-10s  ║\n",
		ui.player.ID, ui.player.X, ui.player.Y, ui.player.Zone)
	fmt.Printf("║  等级：%-4d  金币：%-8d  木材：%-8d        ║\n",
		ui.player.Level, ui.player.Gold, ui.player.Wood)
	fmt.Printf("║  粮食：%-8d  石料：%-8d                    ║\n",
		ui.player.Food, ui.player.Stone)
	fmt.Println("╚════════════════════════════════════════════════════════╝")
}

// ShowMiniMap 显示小地图（简化版）
func (ui *GameUI) ShowMiniMap() {
	ui.historyLock.RLock()
	defer ui.historyLock.RUnlock()
	
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║  小地图 (周边 5 格)                                     ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	
	// 显示玩家周边 5x5 区域
	for dy := -2; dy <= 2; dy++ {
		fmt.Print("║  ")
		for dx := -2; dx <= 2; dx++ {
			x := ui.player.X + int32(dx)
			y := ui.player.Y + int32(dy)
			
			if dx == 0 && dy == 0 {
				fmt.Print("@") // 玩家位置
			} else {
				tile := ui.getTileSymbol(x, y)
				fmt.Print(tile)
			}
			fmt.Print(" ")
		}
		fmt.Println("                                   ║")
	}
	
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  @:玩家  .:平原  ^:山地  ~:河流  *:森林  #:城市        ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
}

// getTileSymbol 获取地块符号
func (ui *GameUI) getTileSymbol(x, y int32) string {
	// 简单模拟地形
	zone := ui.GetZoneName(x, y)
	switch zone {
	case "皇城":
		return "#"
	case "安全区":
		return "."
	case "雍州", "青州", "扬州", "荆州":
		return "*"
	case "蛮荒":
		return "^"
	default:
		return "="
	}
}

// ShowAreaInfo 显示区域信息
func (ui *GameUI) ShowAreaInfo() {
	ui.historyLock.RLock()
	defer ui.historyLock.RUnlock()
	
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║  当前区域信息                                          ║")
	fmt.Printf("║  区域：%s                                      ║\n", ui.getZoneDetailed(ui.player.X, ui.player.Y))
	fmt.Printf("║  世界尺寸：1024x1024                                ║\n")
	fmt.Printf("║  皇城坐标：(480,480)~(544,544)                    ║\n")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
}

// ShowHelp 显示帮助
func (ui *GameUI) ShowHelp() {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║                    SLG MUD 帮助                        ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  移动命令：                                             ║")
	fmt.Println("║    n/north     向北移动                                ║")
	fmt.Println("║    s/south     向南移动                                ║")
	fmt.Println("║    e/east      向东移动                                ║")
	fmt.Println("║    w/west      向西移动                                ║")
	fmt.Println("║    ne          向东北                                  ║")
	fmt.Println("║    nw          向西北                                  ║")
	fmt.Println("║    se          向东南                                  ║")
	fmt.Println("║    sw          向西南                                  ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  查看命令：                                             ║")
	fmt.Println("║    look (l)    查看周围                                ║")
	fmt.Println("║    map         显示小地图                              ║")
	fmt.Println("║    status (st) 查看状态                                ║")
	fmt.Println("║    info        区域信息                                ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  聊天命令：                                             ║")
	fmt.Println("║    chat <消息>  发送消息 (默认全服)                     ║")
	fmt.Println("║    c <消息>     聊天快捷                               ║")
	fmt.Println("║    cw <消息>    全服聊天                               ║")
	fmt.Println("║    ca <消息>    联盟聊天                               ║")
	fmt.Println("║    ch           查看聊天历史                           ║")
	fmt.Println("║    cc           清空聊天历史                           ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  游戏命令：                                             ║")
	fmt.Println("║    build <建筑>  建造建筑                              ║")
	fmt.Println("║    work          工作                                  ║")
	fmt.Println("║    rest          休息                                  ║")
	fmt.Println("║    who           视野内玩家                            ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  特殊命令（以/开头）：                                  ║")
	fmt.Println("║    /quit         退出游戏                              ║")
	fmt.Println("║    /clear        清屏                                  ║")
	fmt.Println("║    /help         本帮助                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
}

// UpdatePlayer 更新玩家信息
func (ui *GameUI) UpdatePlayer(info *PlayerInfo) {
	ui.historyLock.Lock()
	defer ui.historyLock.Unlock()
	ui.player = info
}

// Move 移动
func (ui *GameUI) Move(direction string) {
	ui.historyLock.Lock()
	defer ui.historyLock.Unlock()
	
	dx, dy := ui.getDirectionDelta(direction)
	ui.player.X += dx
	ui.player.Y += dy
	
	// 边界检查
	if ui.player.X < 0 {
		ui.player.X = 0
	}
	if ui.player.X >= 1024 {
		ui.player.X = 1023
	}
	if ui.player.Y < 0 {
		ui.player.Y = 0
	}
	if ui.player.Y >= 1024 {
		ui.player.Y = 1023
	}
}

// getDirectionDelta 获取方向增量
func (ui *GameUI) getDirectionDelta(dir string) (int32, int32) {
	switch strings.ToLower(dir) {
	case "n", "north":
		return 0, -1
	case "s", "south":
		return 0, 1
	case "e", "east":
		return 1, 0
	case "w", "west":
		return -1, 0
	case "ne", "northeast":
		return 1, -1
	case "nw", "northwest":
		return -1, -1
	case "se", "southeast":
		return 1, 1
	case "sw", "southwest":
		return -1, 1
	default:
		return 0, 0
	}
}

// GetPlayerX 获取玩家 X 坐标
func (ui *GameUI) GetPlayerX() int32 {
	ui.historyLock.RLock()
	defer ui.historyLock.RUnlock()
	return ui.player.X
}

// GetPlayerY 获取玩家 Y 坐标
func (ui *GameUI) GetPlayerY() int32 {
	ui.historyLock.RLock()
	defer ui.historyLock.RUnlock()
	return ui.player.Y
}

// GetZoneName 获取区域名称
func (ui *GameUI) GetZoneName(x, y int32) string {
	// 边缘绝境
	if x < 64 || x >= 1024-64 || y < 64 || y >= 1024-64 {
		return "绝境"
	}
	
	// 中心区域
	if x >= 384 && x < 640 && y >= 384 && y < 640 {
		if x >= 480 && x < 544 && y >= 480 && y < 544 {
			return "皇城"
		}
		return "安全区"
	}
	
	// 四大州
	mid := int32(512)
	if x < mid && y < mid {
		return "雍州"
	} else if x >= mid && y < mid {
		return "青州"
	} else if x < mid && y >= mid {
		return "扬州"
	} else {
		return "荆州"
	}
}

// getZoneDetailed 获取详细区域名称
func (ui *GameUI) getZoneDetailed(x, y int32) string {
	// 边缘绝境
	if x < 64 || x >= 1024-64 || y < 64 || y >= 1024-64 {
		return "边缘绝境"
	}
	
	// 中心区域
	if x >= 384 && x < 640 && y >= 384 && y < 640 {
		if x >= 480 && x < 544 && y >= 480 && y < 544 {
			return "皇城"
		}
		return "中心安全区"
	}
	
	// 蛮荒带
	barbarianStart := int32(640)
	barbarianEnd := int32(768)
	
	if (x >= barbarianStart && x < barbarianEnd) ||
	   (x >= 1024-barbarianEnd && x < 1024-barbarianStart) ||
	   (y >= barbarianStart && y < barbarianEnd) ||
	   (y >= 1024-barbarianEnd && y < 1024-barbarianStart) {
		return "蛮荒带"
	}
	
	// 四大州
	mid := int32(512)
	if x < mid && y < mid {
		return "雍州（西北）"
	} else if x >= mid && y < mid {
		return "青州（东北）"
	} else if x < mid && y >= mid {
		return "扬州（西南）"
	} else {
		return "荆州（东南）"
	}
}

// AddChatMessage 添加聊天消息
func (ui *GameUI) AddChatMessage(msg *ChatMessage) {
	ui.historyLock.Lock()
	defer ui.historyLock.Unlock()
	
	ui.history = append(ui.history, msg)
	
	// 限制历史记录数量
	if len(ui.history) > ui.maxHistory {
		ui.history = ui.history[1:]
	}
}

// GetChatHistory 获取聊天历史
func (ui *GameUI) GetChatHistory() []*ChatMessage {
	ui.historyLock.RLock()
	defer ui.historyLock.RUnlock()
	
	// 返回副本
	history := make([]*ChatMessage, len(ui.history))
	copy(history, ui.history)
	return history
}

// ClearChatHistory 清空聊天历史
func (ui *GameUI) ClearChatHistory() {
	ui.historyLock.Lock()
	defer ui.historyLock.Unlock()
	ui.history = make([]*ChatMessage, 0)
}

// ShowChatMessages 显示最近聊天消息
func (ui *GameUI) ShowChatMessages(count int) {
	ui.historyLock.RLock()
	defer ui.historyLock.RUnlock()
	
	if len(ui.history) == 0 {
		fmt.Println("📭 暂无聊天消息")
		return
	}
	
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║                  最近聊天                              ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	
	start := len(ui.history) - count
	if start < 0 {
		start = 0
	}
	
	for i := start; i < len(ui.history); i++ {
		msg := ui.history[i]
		timeStr := time.UnixMilli(msg.Timestamp).Format("15:04:05")
		channelIcon := map[string]string{
			"world":    "🌍",
			"alliance": "🏰",
		}[msg.Channel]
		
		// 截断过长的消息
		content := msg.Content
		if len(content) > 50 {
			content = content[:47] + "..."
		}
		
		fmt.Printf("║ %s %s [%s]: %s\n", timeStr, channelIcon, msg.Username, content)
	}
	
	fmt.Println("╚════════════════════════════════════════════════════════╝")
}

// Refresh 刷新界面
func (ui *GameUI) Refresh() {
	ui.Clear()
	ui.ShowHeader()
	ui.ShowPlayerInfo()
	ui.ShowMiniMap()
	ui.ShowAreaInfo()
	ui.ShowChatMessages(5) // 显示最近 5 条聊天
}
