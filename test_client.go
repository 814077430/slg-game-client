package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
	pb "slg-game/protocol"
)

const (
	HeaderSize = 12
	MagicNumber = 0x534C
	ProtocolVersion = 1
)

// Packet 网络包（与服务器一致）
type Packet struct {
	Magic   uint16
	Version uint8
	Flags   uint8
	MsgID   uint32
	Data    []byte
}

// Encode 编码
func (p *Packet) Encode() []byte {
	buf := make([]byte, HeaderSize+len(p.Data))
	binary.BigEndian.PutUint16(buf[0:2], p.Magic)
	buf[2] = p.Version
	buf[3] = p.Flags
	binary.BigEndian.PutUint32(buf[4:8], p.MsgID)
	binary.BigEndian.PutUint32(buf[8:12], uint32(len(p.Data)))
	copy(buf[12:], p.Data)
	return buf
}

// Decode 解码
func Decode(reader *bufio.Reader) (*Packet, error) {
	header := make([]byte, HeaderSize)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}

	magic := binary.BigEndian.Uint16(header[0:2])
	if magic != MagicNumber {
		return nil, fmt.Errorf("invalid magic number: %x", magic)
	}

	version := header[2]
	if version != ProtocolVersion {
		return nil, fmt.Errorf("invalid version: %d", version)
	}

	flags := header[3]
	msgID := binary.BigEndian.Uint32(header[4:8])
	dataLen := binary.BigEndian.Uint32(header[8:12])

	data := make([]byte, dataLen)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}

	return &Packet{
		Magic:   magic,
		Version: version,
		Flags:   flags,
		MsgID:   msgID,
		Data:    data,
	}, nil
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║          SLG Game Client - Test                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")

	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("❌ 连接失败：%v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("✓ 连接到服务器 localhost:8080")

	reader := bufio.NewReader(conn)

	// 测试 1: 注册
	fmt.Println("\n📝 测试注册...")
	registerReq := &pb.C2S_RegisterRequest{
		Username: "testuser",
		Password: "test123",
		Email:    "test@example.com",
	}
	sendPacket(conn, 1002, registerReq)
	
	resp := readPacket(reader)
	if resp != nil {
		loginResp := &pb.S2C_RegisterResponse{}
		proto.Unmarshal(resp.Data, loginResp)
		if loginResp.Success {
			fmt.Printf("✓ 注册成功！PlayerID: %d\n", loginResp.PlayerId)
		} else {
			fmt.Printf("⚠ 注册响应：%s\n", loginResp.Message)
		}
	}

	// 测试 2: 登录
	fmt.Println("\n🔐 测试登录...")
	loginReq := &pb.C2S_LoginRequest{
		Username: "testuser",
		Password: "test123",
	}
	sendPacket(conn, 1001, loginReq)
	
	resp = readPacket(reader)
	if resp != nil {
		loginResp := &pb.S2C_LoginResponse{}
		proto.Unmarshal(resp.Data, loginResp)
		if loginResp.Success {
			fmt.Printf("✓ 登录成功！PlayerID: %d, 用户名：%s\n", loginResp.PlayerId, loginResp.PlayerData.Username)
			fmt.Printf("  资源 - 金币：%d, 木材：%d, 食物：%d\n", 
				loginResp.PlayerData.Resources["gold"],
				loginResp.PlayerData.Resources["wood"],
				loginResp.PlayerData.Resources["food"])
		} else {
			fmt.Printf("⚠ 登录响应：%s\n", loginResp.Message)
		}
	}

	// 测试 3: 移动
	fmt.Println("\n🚶 测试移动...")
	moveReq := &pb.C2S_MoveRequest{
		X: 100,
		Y: 200,
	}
	sendPacket(conn, 1003, moveReq)
	
	resp = readPacket(reader)
	if resp != nil {
		moveResp := &pb.S2C_MoveResponse{}
		proto.Unmarshal(resp.Data, moveResp)
		if moveResp.Success {
			fmt.Printf("✓ 移动成功！位置：(%d, %d)\n", moveResp.X, moveResp.Y)
		} else {
			fmt.Printf("⚠ 移动响应：%s\n", moveResp.Message)
		}
	}

	// 测试 4: 聊天
	fmt.Println("\n💬 测试聊天...")
	chatReq := &pb.C2S_ChatRequest{
		Content: "大家好，我是测试玩家！",
		Channel: "world",
	}
	sendPacket(conn, 4010, chatReq)  // 聊天协议 MsgID: 4010
	
	resp = readPacket(reader)
	if resp != nil {
		chatResp := &pb.S2C_ChatResponse{}
		proto.Unmarshal(resp.Data, chatResp)
		if chatResp.Success {
			fmt.Printf("✓ 聊天消息发送成功！\n")
		} else {
			fmt.Printf("⚠ 聊天响应：%s\n", chatResp.Message)
		}
	}

	// 等待一下看是否有广播
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n═══════════════════════════════════════════════════════")
	fmt.Println("✅ 所有功能测试完成！")
	fmt.Println("═══════════════════════════════════════════════════════")
}

func sendPacket(conn net.Conn, msgID uint32, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	packet := &Packet{
		Magic:   MagicNumber,
		Version: ProtocolVersion,
		Flags:   0,
		MsgID:   msgID,
		Data:    data,
	}
	conn.Write(packet.Encode())
}

func readPacket(reader *bufio.Reader) *Packet {
	packet, err := Decode(reader)
	if err != nil {
		fmt.Printf("读取响应失败：%v\n", err)
		return nil
	}
	return packet
}
