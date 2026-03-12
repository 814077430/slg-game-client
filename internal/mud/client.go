package mud

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

// MUDClient MUD 客户端
type MUDClient struct {
	serverAddr  string
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	isConnected bool
	isLoggedIn  bool
	username    string
	mutex       sync.Mutex
	outputChan  chan string
	inputChan   chan string
	closeChan   chan struct{}
	
	// 游戏状态
	playerID    uint64
	location    string
	stats       map[string]int
	inventory   []string
}

// NewMUDClient 创建新 MUD 客户端
func NewMUDClient(serverAddr string) *MUDClient {
	return &MUDClient{
		serverAddr: serverAddr,
		outputChan: make(chan string, 100),
		inputChan:  make(chan string, 100),
		closeChan:  make(chan struct{}),
		stats:      make(map[string]int),
		inventory:  make([]string, 0),
	}
}

// Connect 连接服务器
func (m *MUDClient) Connect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	conn, err := net.DialTimeout("tcp", m.serverAddr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("连接失败：%v", err)
	}

	m.conn = conn
	m.reader = bufio.NewReader(conn)
	m.writer = bufio.NewWriter(conn)
	m.isConnected = true

	go m.readLoop()
	go m.writeLoop()

	return nil
}

// readLoop 读取循环
func (m *MUDClient) readLoop() {
	for {
		select {
		case <-m.closeChan:
			return
		default:
			line, err := m.reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					m.outputChan <- fmt.Sprintf("❌ 读取错误：%v", err)
				}
				m.Close()
				return
			}

			line = strings.TrimSpace(line)
			if line != "" {
				m.outputChan <- line
				m.parseOutput(line)
			}
		}
	}
}

// writeLoop 写入循环
func (m *MUDClient) writeLoop() {
	for {
		select {
		case input := <-m.inputChan:
			m.mutex.Lock()
			if m.isConnected && m.writer != nil {
				m.writer.WriteString(input + "\n")
				m.writer.Flush()
			}
			m.mutex.Unlock()
		case <-m.closeChan:
			return
		}
	}
}

// parseOutput 解析服务器输出
func (m *MUDClient) parseOutput(line string) {
	// 解析游戏状态信息
	if strings.Contains(line, "Player ID:") {
		fmt.Sscanf(line, "Player ID: %d", &m.playerID)
	}
	
	if strings.Contains(line, "位置:") {
		parts := strings.Split(line, "位置:")
		if len(parts) > 1 {
			m.location = strings.TrimSpace(parts[1])
		}
	}
}

// Send 发送命令
func (m *MUDClient) Send(command string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if m.isConnected {
		m.inputChan <- command
	}
}

// SendChat 发送聊天消息（简化 Protobuf 格式）
func (m *MUDClient) SendChat(content, channel string) {
	if channel == "" {
		channel = "world"
	}
	
	// 简化格式：channel + " " + content
	// 服务器会解析为 Protobuf 格式
	text := fmt.Sprintf("%s %s", channel, content)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if m.isConnected && m.writer != nil {
		// 发送文本，服务器 router.go 会处理
		m.writer.WriteString(text)
		m.writer.WriteByte('\n')
		m.writer.Flush()
	}
}

// Close 关闭连接
func (m *MUDClient) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.isConnected {
		return
	}

	m.isConnected = false
	close(m.closeChan)

	if m.conn != nil {
		m.conn.Close()
	}
}

// IsConnected 检查连接状态
func (m *MUDClient) IsConnected() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.isConnected
}

// GetOutputChan 获取输出通道
func (m *MUDClient) GetOutputChan() <-chan string {
	return m.outputChan
}

// Login 登录
func (m *MUDClient) Login(username, password string) {
	m.username = username
	m.Send(fmt.Sprintf("login %s %s", username, password))
}

// Register 注册
func (m *MUDClient) Register(username, password string) {
	m.Send(fmt.Sprintf("register %s %s", username, password))
}

// Move 移动
func (m *MUDClient) Move(direction string) {
	m.Send(fmt.Sprintf("go %s", direction))
}

// Look 查看
func (m *MUDClient) Look() {
	m.Send("look")
}

// Status 查看状态
func (m *MUDClient) Status() {
	m.Send("status")
}

// Inventory 查看背包
func (m *MUDClient) Inventory() {
	m.Send("inventory")
}

// Build 建造
func (m *MUDClient) Build(buildingType string) {
	m.Send(fmt.Sprintf("build %s", buildingType))
}

// Work 工作
func (m *MUDClient) Work() {
	m.Send("work")
}

// Rest 休息
func (m *MUDClient) Rest() {
	m.Send("rest")
}

// Help 帮助
func (m *MUDClient) Help() {
	m.Send("help")
}

// Who 在线玩家
func (m *MUDClient) Who() {
	m.Send("who")
}

// Say 说话
func (m *MUDClient) Say(message string) {
	m.Send(fmt.Sprintf("say %s", message))
}

// GetPlayerID 获取玩家 ID
func (m *MUDClient) GetPlayerID() uint64 {
	return m.playerID
}

// GetLocation 获取位置
func (m *MUDClient) GetLocation() string {
	return m.location
}

// GetStats 获取属性
func (m *MUDClient) GetStats() map[string]int {
	return m.stats
}

// GetInventory 获取背包
func (m *MUDClient) GetInventory() []string {
	return m.inventory
}
