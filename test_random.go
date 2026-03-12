package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
	pb "slg-game/protocol"
)

const (
	HeaderSize    = 12
	MagicNumber   = 0x534C
	ProtocolVersion = 1
)

// Packet 网络包
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
		return nil, fmt.Errorf("invalid magic number")
	}

	version := header[2]
	if version != ProtocolVersion {
		return nil, fmt.Errorf("invalid version")
	}

	msgID := binary.BigEndian.Uint32(header[4:8])
	dataLen := binary.BigEndian.Uint32(header[8:12])

	data := make([]byte, dataLen)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}

	return &Packet{
		Magic:   magic,
		Version: version,
		Flags:   header[3],
		MsgID:   msgID,
		Data:    data,
	}, nil
}

func sendPacket(conn net.Conn, msgID uint32, msg proto.Message) error {
	data, _ := proto.Marshal(msg)
	packet := &Packet{
		Magic:   MagicNumber,
		Version: ProtocolVersion,
		Flags:   0,
		MsgID:   msgID,
		Data:    data,
	}
	_, err := conn.Write(packet.Encode())
	return err
}

func readPacket(reader *bufio.Reader) (*Packet, error) {
	return Decode(reader)
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║       SLG Game - Random Command Test                   ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")

	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("❌ 连接失败：%v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("✓ 连接到服务器 localhost:8080")

	reader := bufio.NewReader(conn)
	rand.Seed(time.Now().UnixNano())

	// 生成随机用户名
	username := fmt.Sprintf("test_user_%d", rand.Intn(10000))
	password := "password123"

	// 测试 1: 注册
	fmt.Println("\n📝 测试 1: 注册")
	regReq := &pb.C2S_RegisterRequest{
		Username: username,
		Password: password,
		Email:    fmt.Sprintf("%s@test.com", username),
	}
	if err := sendPacket(conn, 1002, regReq); err != nil {
		fmt.Printf("❌ 发送注册请求失败：%v\n", err)
		return
	}
	resp, _ := readPacket(reader)
	if resp != nil {
		regResp := &pb.S2C_RegisterResponse{}
		proto.Unmarshal(resp.Data, regResp)
		if regResp.Success {
			fmt.Printf("✓ 注册成功！PlayerID: %d\n", regResp.PlayerId)
		} else {
			fmt.Printf("⚠ 注册响应：%s\n", regResp.Message)
		}
	}
	time.Sleep(100 * time.Millisecond)

	// 测试 2: 登录
	fmt.Println("\n🔐 测试 2: 登录")
	loginReq := &pb.C2S_LoginRequest{
		Username: username,
		Password: password,
	}
	if err := sendPacket(conn, 1001, loginReq); err != nil {
		fmt.Printf("❌ 发送登录请求失败：%v\n", err)
		return
	}
	resp, _ = readPacket(reader)
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
	time.Sleep(100 * time.Millisecond)

	// 测试 3-7: 随机指令测试（5 次）
	commands := []string{"move", "build", "chat", "move", "build"}
	for i, cmd := range commands {
		fmt.Printf("\n🎲 测试 %d: 随机指令 [%s]\n", i+3, cmd)

		switch cmd {
		case "move":
			// 移动
			x := int32(rand.Intn(1000))
			y := int32(rand.Intn(1000))
			moveReq := &pb.C2S_MoveRequest{
				X: x,
				Y: y,
			}
			fmt.Printf("  发送移动指令：位置 (%d, %d)\n", x, y)
			if err := sendPacket(conn, 1003, moveReq); err != nil {
				fmt.Printf("❌ 发送移动请求失败：%v\n", err)
				continue
			}
			resp, _ = readPacket(reader)
			if resp != nil {
				moveResp := &pb.S2C_MoveResponse{}
				proto.Unmarshal(resp.Data, moveResp)
				if moveResp.Success {
					fmt.Printf("✓ 移动成功！位置：(%d, %d)\n", moveResp.X, moveResp.Y)
				} else {
					fmt.Printf("⚠ 移动响应：%s\n", moveResp.Message)
				}
			}

		case "build":
			// 建造
			buildTypes := []string{"farm", "lumber_mill", "mine", "barracks"}
			buildType := buildTypes[rand.Intn(len(buildTypes))]
			x := int32(rand.Intn(100))
			y := int32(rand.Intn(100))
			buildReq := &pb.C2S_BuildRequest{
				BuildingType: buildType,
				X:            x,
				Y:            y,
			}
			fmt.Printf("  发送建造指令：建筑 [%s] 位置 (%d, %d)\n", buildType, x, y)
			if err := sendPacket(conn, 1004, buildReq); err != nil {
				fmt.Printf("❌ 发送建造请求失败：%v\n", err)
				continue
			}
			resp, _ = readPacket(reader)
			if resp != nil {
				buildResp := &pb.S2C_BuildResponse{}
				proto.Unmarshal(resp.Data, buildResp)
				if buildResp.Success {
					fmt.Printf("✓ 建造成功！建筑：%s, 等级：%d\n", buildResp.Building.BuildingType, buildResp.Building.Level)
				} else {
					fmt.Printf("⚠ 建造响应：%s\n", buildResp.Message)
				}
			}

		case "chat":
			// 聊天
			messages := []string{
				"大家好！",
				"有人吗？",
				"测试聊天功能",
				"SLG 游戏真好玩",
				"今天天气不错",
			}
			msg := messages[rand.Intn(len(messages))]
			chatReq := &pb.C2S_ChatRequest{
				Content: msg,
				Channel: "world",
			}
			fmt.Printf("  发送聊天指令：消息 [%s]\n", msg)
			if err := sendPacket(conn, 4010, chatReq); err != nil {
				fmt.Printf("❌ 发送聊天请求失败：%v\n", err)
				continue
			}
			resp, _ = readPacket(reader)
			if resp != nil {
				chatResp := &pb.S2C_ChatResponse{}
				proto.Unmarshal(resp.Data, chatResp)
				if chatResp.Success {
					fmt.Printf("✓ 聊天成功！时间戳：%d\n", chatResp.Timestamp)
				} else {
					fmt.Printf("⚠ 聊天响应：%s\n", chatResp.Message)
				}
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	// 测试 8: 查询视野内玩家
	fmt.Println("\n👥 测试 8: 查询视野内玩家")
	whoReq := &pb.C2S_WhoRequest{}
	if err := sendPacket(conn, 1005, whoReq); err != nil {
		fmt.Printf("❌ 发送查询请求失败：%v\n", err)
		return
	}
	resp, _ = readPacket(reader)
	if resp != nil {
		whoResp := &pb.S2C_WhoResponse{}
		proto.Unmarshal(resp.Data, whoResp)
		if whoResp.Success {
			fmt.Printf("✓ 查询成功！视野内玩家数：%d\n", len(whoResp.Players))
			for i, p := range whoResp.Players {
				if i < 5 { // 只显示前 5 个
					fmt.Printf("  - 玩家 %d: ID=%d, 名字=%s, 位置=(%d,%d)\n",
						i+1, p.PlayerId, p.Username, p.X, p.Y)
				}
			}
			if len(whoResp.Players) > 5 {
				fmt.Printf("  ... 还有 %d 个玩家\n", len(whoResp.Players)-5)
			}
		} else {
			fmt.Printf("⚠ 查询响应：%s\n", whoResp.Message)
		}
	}

	fmt.Println("\n═══════════════════════════════════════════════════════")
	fmt.Println("✅ 所有随机指令测试完成！")
	fmt.Println("═══════════════════════════════════════════════════════")
}
