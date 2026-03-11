package client

import (
	"errors"
	"time"

	"google.golang.org/protobuf/proto"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Success  bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message  string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	PlayerId uint64 `protobuf:"varint,3,opt,name=player_id,proto3" json:"player_id,omitempty"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Email    string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Success  bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message  string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	PlayerId uint64 `protobuf:"varint,3,opt,name=player_id,proto3" json:"player_id,omitempty"`
}

// MoveRequest 移动请求
type MoveRequest struct {
	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	Y int32 `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
}

// MoveResponse 移动响应
type MoveResponse struct {
	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	X       int32  `protobuf:"varint,3,opt,name=x,proto3" json:"x,omitempty"`
	Y       int32  `protobuf:"varint,4,opt,name=y,proto3" json:"y,omitempty"`
}

// BuildRequest 建造请求
type BuildRequest struct {
	BuildingType string `protobuf:"bytes,1,opt,name=building_type,proto3" json:"building_type,omitempty"`
	X            int32  `protobuf:"varint,2,opt,name=x,proto3" json:"x,omitempty"`
	Y            int32  `protobuf:"varint,3,opt,name=y,proto3" json:"y,omitempty"`
}

// BuildResponse 建造响应
type BuildResponse struct {
	Success  bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message  string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

// Login 登录
func (c *Client) Login(username, password string) (*LoginResponse, error) {
	if !c.IsConnected() {
		return nil, errors.New("not connected")
	}

	req := &LoginRequest{
		Username: username,
		Password: password,
	}

	if err := c.Send(MsgID_C2S_LoginRequest, req); err != nil {
		return nil, err
	}

	packet, err := c.Recv(10 * time.Second)
	if err != nil {
		return nil, err
	}

	resp := &LoginResponse{}
	if err := proto.Unmarshal(packet.Data, resp); err != nil {
		return nil, err
	}

	if resp.Success {
		c.playerID = resp.PlayerId
		c.username = username
		c.isLoggedIn = true
	}

	return resp, nil
}

// Register 注册
func (c *Client) Register(username, password, email string) (*RegisterResponse, error) {
	if !c.IsConnected() {
		return nil, errors.New("not connected")
	}

	req := &RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	}

	if err := c.Send(MsgID_C2S_RegisterRequest, req); err != nil {
		return nil, err
	}

	packet, err := c.Recv(10 * time.Second)
	if err != nil {
		return nil, err
	}

	resp := &RegisterResponse{}
	if err := proto.Unmarshal(packet.Data, resp); err != nil {
		return nil, err
	}

	if resp.Success {
		c.playerID = resp.PlayerId
		c.username = username
		c.isLoggedIn = true
	}

	return resp, nil
}

// Move 移动
func (c *Client) Move(x, y int32) (*MoveResponse, error) {
	if !c.IsConnected() {
		return nil, errors.New("not connected")
	}

	if !c.isLoggedIn {
		return nil, errors.New("not logged in")
	}

	req := &MoveRequest{X: x, Y: y}

	if err := c.Send(MsgID_C2S_MoveRequest, req); err != nil {
		return nil, err
	}

	packet, err := c.Recv(10 * time.Second)
	if err != nil {
		return nil, err
	}

	resp := &MoveResponse{}
	if err := proto.Unmarshal(packet.Data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Build 建造
func (c *Client) Build(buildingType string, x, y int32) (*BuildResponse, error) {
	if !c.IsConnected() {
		return nil, errors.New("not connected")
	}

	if !c.isLoggedIn {
		return nil, errors.New("not logged in")
	}

	req := &BuildRequest{
		BuildingType: buildingType,
		X:            x,
		Y:            y,
	}

	if err := c.Send(MsgID_C2S_BuildRequest, req); err != nil {
		return nil, err
	}

	packet, err := c.Recv(10 * time.Second)
	if err != nil {
		return nil, err
	}

	resp := &BuildResponse{}
	if err := proto.Unmarshal(packet.Data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
