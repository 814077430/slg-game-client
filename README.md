# SLG Game Client

SLG 游戏客户端 SDK - Go 语言实现

## 📦 功能特性

- ✅ TCP 长连接
- ✅ Protobuf 协议
- ✅ 自动重连
- ✅ 消息队列
- ✅ 并发安全
- ✅ 心跳检测

## 🚀 快速开始

### 安装依赖

```bash
go mod download
```

### 基本使用

```go
package main

import (
	"log"
	"slg-game-client/internal/client"
)

func main() {
	// 创建客户端
	c := client.NewClient("localhost:8080")
	
	// 连接服务器
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	
	// 登录
	resp, err := c.Login("username", "password")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("登录成功！Player ID: %d", resp.PlayerId)
	
	// 移动
	c.Move(100, 200)
	
	// 建造
	c.Build("farm", 50, 50)
	
	// 保持连接
	select {}
}
```

## 📁 项目结构

```
slg-game-client/
├── cmd/                    # 命令行工具
│   └── main.go            # 主程序入口
├── internal/
│   ├── client/            # 客户端核心
│   │   ├── client.go      # 客户端实现
│   │   ├── connection.go  # 连接管理
│   │   └── handler.go     # 消息处理
│   ├── proto/             # Protobuf 定义
│   └── utils/             # 工具函数
├── go.mod
├── go.sum
└── README.md
```

## 🔧 配置

```json
{
  "server": "localhost:8080",
  "reconnect": true,
  "reconnectInterval": 5,
  "heartbeatInterval": 30,
  "timeout": 10
}
```

## 📊 API 文档

### 连接管理

| 方法 | 说明 |
|------|------|
| `Connect()` | 连接服务器 |
| `Close()` | 断开连接 |
| `IsConnected()` | 检查连接状态 |

### 用户操作

| 方法 | 说明 |
|------|------|
| `Login(username, password)` | 登录 |
| `Register(username, password, email)` | 注册 |

### 游戏操作

| 方法 | 说明 |
|------|------|
| `Move(x, y)` | 移动 |
| `Build(type, x, y)` | 建造 |
| `Attack(targetID)` | 攻击 |

## 🧪 测试

```bash
go test ./...
```

## 📝 许可证

MIT License
