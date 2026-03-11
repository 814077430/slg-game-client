# SLG MUD Client

SLG 游戏 MUD（文字冒险）客户端 - 基于文本的交互式游戏客户端

## 🎮 特性

- ✅ 纯文字界面
- ✅ 命令别名系统
- ✅ 命令历史记录
- ✅ 实时输出格式化
- ✅ 支持 TCP 长连接
- ✅ 与 slg-game 服务器兼容

## 🚀 快速开始

### 编译

```bash
cd cmd/mud
go build -o mud-client .
```

### 运行

```bash
# 连接本地服务器
./mud-client

# 连接远程服务器
./mud-client server:port
```

## 📖 游戏命令

### 基础命令

| 命令 | 别名 | 说明 |
|------|------|------|
| `login <用户名> <密码>` | - | 登录 |
| `register <用户名> <密码>` | - | 注册 |
| `look` | `l` | 查看周围 |
| `go <方向>` | `n/s/e/w` | 移动 |
| `status` | `st` | 查看状态 |
| `inventory` | `i` | 查看背包 |
| `build <建筑>` | - | 建造 |
| `work` | - | 工作 |
| `rest` | - | 休息 |
| `say <消息>` | - | 说话 |
| `who` | - | 在线玩家 |
| `help` | - | 帮助 |

### 特殊命令（以/开头）

| 命令 | 说明 |
|------|------|
| `/quit` | 退出游戏 |
| `/clear` | 清屏 |
| `/history` | 命令历史 |
| `/alias` | 别名列表 |
| `/help` | 帮助 |

## 🎯 游戏流程

```
1. 启动客户端
   $ ./mud-client

2. 注册账号
   > register player1 password123

3. 登录
   > login player1 password123

4. 查看周围
   > look

5. 移动探索
   > go north
   > go east

6. 建造建筑
   > build farm

7. 工作获取资源
   > work

8. 查看状态
   > status
```

## 📁 项目结构

```
slg-game-client/
├── cmd/
│   ├── main.go           # 原客户端入口
│   └── mud/
│       └── main.go       # MUD 客户端入口
├── internal/
│   ├── client/           # 协议客户端
│   └── mud/              # MUD 客户端核心
│       ├── client.go     # MUD 连接
│       └── handler.go    # 命令处理
├── go.mod
└── README.md
```

## 🔧 自定义别名

在 `handler.go` 中添加：

```go
h.RegisterAlias("缩写", "完整命令", "描述")
```

## 📊 输出格式

客户端自动格式化服务器输出：

- ✅ 成功消息
- ❌ 失败消息
- 🎮 欢迎消息
- 📍 位置信息
- 💰 资源信息

## 🛠️ 开发

### 添加新命令

1. 在 `handler.go` 的 `ProcessCommand` 中添加命令解析
2. 在 `client.go` 中添加对应的发送方法
3. 更新帮助文档

### 添加新别名

```go
h.RegisterAlias("x", "examine", "检查物品")
```

## 📝 许可证

MIT License
