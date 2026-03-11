package client

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	HeaderSize   = 12
	MaxMsgSize   = 1024 * 1024
	MagicNumber  = 0x534C
	ProtocolVer  = 1
)

// 消息 ID 常量
const (
	MsgID_C2S_LoginRequest     = 1001
	MsgID_C2S_RegisterRequest  = 1002
	MsgID_C2S_MoveRequest      = 1003
	MsgID_C2S_BuildRequest     = 1004
	MsgID_S2C_LoginResponse    = 2001
	MsgID_S2C_RegisterResponse = 2002
	MsgID_S2C_MoveResponse     = 2003
	MsgID_S2C_BuildResponse    = 2004
)

// Client SLG 游戏客户端
type Client struct {
	serverAddr  string
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	sendChan    chan *Packet
	recvChan    chan *Packet
	closeChan   chan struct{}
	isClosed    bool
	mutex       sync.Mutex
	playerID    uint64
	username    string
	isLoggedIn  bool
}

// Packet 网络包
type Packet struct {
	MsgID uint32
	Data  []byte
}

// Encode 编码
func (p *Packet) Encode() []byte {
	buf := make([]byte, HeaderSize+len(p.Data))
	binary.BigEndian.PutUint32(buf[0:4], p.MsgID)
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(p.Data)))
	copy(buf[8:], p.Data)
	return buf
}

// Decode 解码
func Decode(reader io.Reader) (*Packet, error) {
	header := make([]byte, 8)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}

	msgID := binary.BigEndian.Uint32(header[0:4])
	msgLen := binary.BigEndian.Uint32(header[4:8])

	if msgLen > MaxMsgSize {
		return nil, errors.New("message too large")
	}

	data := make([]byte, msgLen)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}

	return &Packet{
		MsgID: msgID,
		Data:  data,
	}, nil
}

// NewClient 创建新客户端
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
		sendChan:   make(chan *Packet, 100),
		recvChan:   make(chan *Packet, 100),
		closeChan:  make(chan struct{}),
	}
}

// Connect 连接服务器
func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.serverAddr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)

	go c.sendLoop()
	go c.recvLoop()

	return nil
}

// sendLoop 发送循环
func (c *Client) sendLoop() {
	for {
		select {
		case packet := <-c.sendChan:
			if err := c.writePacket(packet); err != nil {
				c.Close()
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

// recvLoop 接收循环
func (c *Client) recvLoop() {
	for {
		packet, err := Decode(c.reader)
		if err != nil {
			c.Close()
			return
		}

		select {
		case c.recvChan <- packet:
		case <-c.closeChan:
			return
		}
	}
}

// writePacket 写入包
func (c *Client) writePacket(packet *Packet) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isClosed {
		return errors.New("connection closed")
	}

	data := packet.Encode()
	_, err := c.writer.Write(data)
	if err != nil {
		return err
	}

	return c.writer.Flush()
}

// Send 发送消息
func (c *Client) Send(msgID uint32, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	packet := &Packet{
		MsgID: msgID,
		Data:  data,
	}

	select {
	case c.sendChan <- packet:
		return nil
	default:
		return errors.New("send queue full")
	}
}

// Recv 接收消息
func (c *Client) Recv(timeout time.Duration) (*Packet, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case packet := <-c.recvChan:
		return packet, nil
	case <-timer.C:
		return nil, errors.New("timeout")
	case <-c.closeChan:
		return nil, errors.New("connection closed")
	}
}

// Close 关闭连接
func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isClosed {
		return
	}

	c.isClosed = true
	close(c.closeChan)

	if c.conn != nil {
		c.conn.Close()
	}
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return !c.isClosed && c.conn != nil
}

// GetPlayerID 获取玩家 ID
func (c *Client) GetPlayerID() uint64 {
	return c.playerID
}

// GetUsername 获取用户名
func (c *Client) GetUsername() string {
	return c.username
}

// IsLoggedIn 检查登录状态
func (c *Client) IsLoggedIn() bool {
	return c.isLoggedIn
}
