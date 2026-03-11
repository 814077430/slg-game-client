#!/bin/bash

# SLG MUD 客户端并发压测脚本
# 启动 100 个客户端并发连接服务器

CLIENT_COUNT=${1:-100}
SERVER_ADDR=${2:-"localhost:8080"}
MUD_CLIENT="./mud-client"

echo "╔════════════════════════════════════════════════════════╗"
echo "║          SLG MUD Client Stress Test                    ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo "配置:"
echo "  客户端数量：$CLIENT_COUNT"
echo "  服务器地址：$SERVER_ADDR"
echo "  客户端程序：$MUD_CLIENT"
echo ""

# 检查客户端是否存在
if [ ! -f "$MUD_CLIENT" ]; then
    echo "❌ 客户端不存在：$MUD_CLIENT"
    exit 1
fi

# 检查服务器是否运行
if ! pgrep -f slg-server > /dev/null; then
    echo "❌ 服务器未运行，请先启动 slg-server"
    exit 1
fi
echo "✅ 服务器运行中"
echo ""

# 记录开始时间
START_TIME=$(date +%s)

# 启动客户端
echo "=== 启动 $CLIENT_COUNT 个客户端 ==="
for i in $(seq 1 $CLIENT_COUNT); do
    username="player$i"
    password="123456"
    
    # 后台启动客户端，模拟登录和简单操作
    (
        sleep 0.1  # 错开启动时间
        echo "register $username $password"
        sleep 0.5
        echo "login $username $password"
        sleep 0.5
        echo "status"
        sleep 0.5
        echo "look"
        sleep 0.5
        echo "n"
        sleep 0.5
        echo "work"
        sleep 0.5
        echo "/quit"
    ) | timeout 10 $MUD_CLIENT > /tmp/client_$i.log 2>&1 &
    
    # 每启动 10 个客户端输出一次进度
    if [ $((i % 10)) -eq 0 ]; then
        echo "  已启动 $i/$CLIENT_COUNT 个客户端"
    fi
done

echo ""
echo "=== 等待客户端完成 ==="

# 等待所有后台进程完成
wait

# 记录结束时间
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
echo "=== 压测结果 ==="
echo "总耗时：${DURATION}秒"
echo "客户端数：$CLIENT_COUNT"
echo "平均每个客户端：$(echo "scale=2; $DURATION / $CLIENT_COUNT" | bc)秒"

# 统计成功/失败
SUCCESS=0
FAILED=0

for i in $(seq 1 $CLIENT_COUNT); do
    if grep -q "连接成功" /tmp/client_$i.log 2>/dev/null || \
       grep -q "Connected" /tmp/client_$i.log 2>/dev/null; then
        SUCCESS=$((SUCCESS + 1))
    else
        FAILED=$((FAILED + 1))
    fi
done

echo ""
echo "连接统计:"
echo "  成功：$SUCCESS / $CLIENT_COUNT"
echo "  失败：$FAILED / $CLIENT_COUNT"
echo "  成功率：$(echo "scale=2; $SUCCESS * 100 / $CLIENT_COUNT" | bc)%"

# 清理日志
echo ""
echo "清理临时日志..."
rm -f /tmp/client_*.log

echo ""
echo "=== 压测完成 ==="
