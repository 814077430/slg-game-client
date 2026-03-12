#!/bin/bash

# SLG MUD 客户端聊天测试脚本
# 启动 10 个客户端，互相发送聊天消息

CLIENT_COUNT=${1:-10}
SERVER_ADDR=${2:-"localhost:8080"}
MUD_CLIENT="./mud-client"

echo "╔════════════════════════════════════════════════════════╗"
echo "║          SLG MUD Client Chat Test                      ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo "配置:"
echo "  客户端数量：$CLIENT_COUNT"
echo "  服务器地址：$SERVER_ADDR"
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
START_TIME=$(date +%s.%N)

echo "=== 启动 $CLIENT_COUNT 个客户端（聊天测试）==="

# 启动客户端
for i in $(seq 1 $CLIENT_COUNT); do
    username="chat_test$i"
    password="123456"
    
    # 后台启动客户端，模拟聊天
    (
        sleep 0.2
        echo "register $username $password"
        sleep 0.5
        echo "login $username $password"
        sleep 0.5
        echo "status"
        sleep 0.5
        
        # 发送聊天消息
        echo "chat 大家好，我是玩家 $i！"
        sleep 0.5
        
        # 等待其他玩家消息
        sleep 2
        
        # 回复其他玩家
        echo "chat 玩家 $i 来报道了！"
        sleep 0.5
        
        # 全服喊话
        echo "cw 测试全服聊天 $i"
        sleep 0.5
        
        # 查看聊天历史
        echo "ch"
        sleep 0.5
        
        # 退出
        echo "/quit"
    ) | timeout 15 $MUD_CLIENT > /tmp/chat_client_$i.log 2>&1 &
    
    echo "  启动客户端 $i: $username"
done

echo ""
echo "=== 等待聊天测试完成（约 10 秒）==="
sleep 10

# 记录结束时间
END_TIME=$(date +%s.%N)
DURATION=$(echo "$END_TIME - $START_TIME" | bc)

echo ""
echo "=== 聊天测试结果 ==="
echo "总耗时：${DURATION}秒"
echo ""

# 统计
SUCCESS=0
FAILED=0
CHAT_SENT=0

for i in $(seq 1 $CLIENT_COUNT); do
    log_file="/tmp/chat_client_$i.log"
    
    if grep -q "连接成功\|Connected" $log_file 2>/dev/null; then
        SUCCESS=$((SUCCESS + 1))
        
        # 统计发送的聊天消息
        chat_count=$(grep -c "chat\|cw" $log_file 2>/dev/null || echo 0)
        CHAT_SENT=$((CHAT_SENT + chat_count))
    else
        FAILED=$((FAILED + 1))
    fi
done

echo "连接统计:"
echo "  成功：$SUCCESS / $CLIENT_COUNT"
echo "  失败：$FAILED / $CLIENT_COUNT"
if [ $CLIENT_COUNT -gt 0 ]; then
    echo "  成功率：$(echo "scale=2; $SUCCESS * 100 / $CLIENT_COUNT" | bc)%"
fi

echo ""
echo "聊天统计:"
echo "  发送消息数：$CHAT_SENT"
if [ $SUCCESS -gt 0 ]; then
    AVG=$(echo "scale=1; $CHAT_SENT / $SUCCESS" | bc 2>/dev/null || echo "0")
    echo "  平均每客户端：$AVG 条"
fi

# 显示部分日志样本
echo ""
echo "=== 聊天日志样本（前 3 个客户端）==="
for i in 1 2 3; do
    echo "--- 客户端 $i ---"
    grep -E "(chat|cw|ch|聊天|Chat|Broadcast)" /tmp/chat_client_$i.log 2>/dev/null | head -10
done

# 显示服务器日志
echo ""
echo "=== 服务器聊天日志（最近 20 条）==="
tail -200 /tmp/server.log | grep -E "(Chat|chat|Broadcast)" | tail -20

# 清理日志
echo ""
echo "清理临时日志..."
rm -f /tmp/chat_client_*.log

echo ""
echo "=== 聊天测试完成 ==="
