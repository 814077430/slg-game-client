package mud

import (
	"fmt"
	"strings"
	"sync"
)

// WorldInfo 世界信息
type WorldInfo struct {
	Size        int
	CenterStart int
	CenterEnd   int
	CastleStart int
	CastleEnd   int
}

// PlayerInfo 玩家信息
type PlayerInfo struct {
	ID        uint64
	Username  string
	X         int32
	Y         int32
	Zone      string
	Gold      int64
	Wood      int64
	Food      int64
	Stone     int64
	Level     int32
}

// TileInfo 地块信息
type TileInfo struct {
	X         int32
	Y         int32
	TileType  string
	Zone      string
	OwnerID   uint64
	Resource  map[string]int32
	Passable  bool
	CityType  string
}

// GameUI 游戏界面
type GameUI struct {
	client    *MUDClient
	player    *PlayerInfo
	world     *WorldInfo
	visible   []*TileInfo
	mutex     sync.Mutex
	width     int
	height    int
}

// NewGameUI 创建游戏界面
func NewGameUI(client *MUDClient) *GameUI {
	return &GameUI{
		client: client,
		player: &PlayerInfo{},
		world: &WorldInfo{
			Size:        1024,
			CenterStart: 384,
			CenterEnd:   640,
			CastleStart: 480,
			CastleEnd:   544,
		},
		visible: make([]*TileInfo, 0),
		width:   80,
		height:  24,
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
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	
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
				tile := ui.getTile(x, y)
				if tile != nil {
					fmt.Print(ui.getTileSymbol(tile))
				} else {
					fmt.Print("?")
				}
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
func (ui *GameUI) getTileSymbol(tile *TileInfo) string {
	if tile.CityType != "" {
		return "#"
	}
	
	switch tile.TileType {
	case "plain":
		return "."
	case "mountain":
		return "^"
	case "river":
		return "~"
	case "forest":
		return "*"
	case "hill":
		return "+"
	case "desert":
		return "="
	case "snow":
		return "x"
	default:
		return "."
	}
}

// getTile 获取地块（模拟）
func (ui *GameUI) getTile(x, y int32) *TileInfo {
	// 实际应该从服务器获取
	return &TileInfo{
		X:        x,
		Y:        y,
		TileType: "plain",
		Zone:     ui.getZoneName(x, y),
		Passable: true,
	}
}

// getZoneName 获取区域名称
func (ui *GameUI) getZoneName(x, y int32) string {
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

// ShowAreaInfo 显示区域信息
func (ui *GameUI) ShowAreaInfo() {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║  当前区域信息                                          ║")
	fmt.Printf("║  区域：%s                                      ║\n", ui.getZoneName(ui.player.X, ui.player.Y))
	fmt.Printf("║  世界尺寸：%dx%d                                ║\n", ui.world.Size, ui.world.Size)
	fmt.Printf("║  皇城坐标：(%d,%d)~(%d,%d)                    ║\n",
		ui.world.CastleStart, ui.world.CastleStart,
		ui.world.CastleEnd, ui.world.CastleEnd)
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
	fmt.Println("║  游戏命令：                                             ║")
	fmt.Println("║    build <建筑>  建造建筑                              ║")
	fmt.Println("║    work          工作                                  ║")
	fmt.Println("║    rest          休息                                  ║")
	fmt.Println("║    say <消息>    说话                                  ║")
	fmt.Println("║    who           在线玩家                              ║")
	fmt.Println("╠════════════════════════════════════════════════════════╣")
	fmt.Println("║  特殊命令（以/开头）：                                  ║")
	fmt.Println("║    /quit         退出游戏                              ║")
	fmt.Println("║    /clear        清屏                                  ║")
	fmt.Println("║    /help         本帮助                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
}

// UpdatePlayer 更新玩家信息
func (ui *GameUI) UpdatePlayer(info *PlayerInfo) {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	ui.player = info
}

// Move 移动
func (ui *GameUI) Move(direction string) {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	
	dx, dy := ui.getDirectionDelta(direction)
	ui.player.X += dx
	ui.player.Y += dy
	
	// 边界检查
	if ui.player.X < 0 {
		ui.player.X = 0
	}
	if ui.player.X >= int32(ui.world.Size) {
		ui.player.X = int32(ui.world.Size) - 1
	}
	if ui.player.Y < 0 {
		ui.player.Y = 0
	}
	if ui.player.Y >= int32(ui.world.Size) {
		ui.player.Y = int32(ui.world.Size) - 1
	}
	
	ui.client.Send(fmt.Sprintf("move %d %d", ui.player.X, ui.player.Y))
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

// Refresh 刷新界面
func (ui *GameUI) Refresh() {
	ui.Clear()
	ui.ShowHeader()
	ui.ShowPlayerInfo()
	ui.ShowMiniMap()
	ui.ShowAreaInfo()
}

// GetPlayerX 获取玩家 X 坐标
func (ui *GameUI) GetPlayerX() int32 {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	return ui.player.X
}

// GetPlayerY 获取玩家 Y 坐标
func (ui *GameUI) GetPlayerY() int32 {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()
	return ui.player.Y
}
