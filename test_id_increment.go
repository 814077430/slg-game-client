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
	HeaderSize    = 12
	MagicNumber   = 0x534C
	ProtocolVersion = 1
)

type Packet struct {
	Magic   uint16
	Version uint8
	Flags   uint8
	MsgID   uint32
	Data    []byte
}

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
	return &Packet{Magic: magic, Version: version, Flags: header[3], MsgID: msgID, Data: data}, nil
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("❌ 连接失败：%v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for i := 1; i <= 5; i++ {
		username := fmt.Sprintf("user_%d_%d", i, time.Now().UnixNano())
		regReq := &pb.C2S_RegisterRequest{Username: username, Password: "pass123", Email: "test@example.com"}
		data, _ := proto.Marshal(regReq)
		packet := &Packet{Magic: MagicNumber, Version: ProtocolVersion, Flags: 0, MsgID: 1002, Data: data}
		conn.Write(packet.Encode())
		
		resp, _ := Decode(reader)
		regResp := &pb.S2C_RegisterResponse{}
		proto.Unmarshal(resp.Data, regResp)
		fmt.Printf("✓ 注册 %s -> PlayerID: %d\n", username, regResp.PlayerId)
		time.Sleep(50 * time.Millisecond)
	}
}
