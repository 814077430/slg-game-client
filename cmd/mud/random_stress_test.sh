#!/bin/bash

# SLG MUD 客户端随机消息压测脚本
# 启动 100 个客户端，随机发送各种游戏消息

CLIENT_COUNT=${1:-100}
SERVER_ADDR=${2:-"localhost:8080"}
MUD_CLIENT="./mud-client"
TEST_DURATION=${3:-30}  # 测试时长（秒）

echo "╔════════════════════════════════════════════════════════╗"
echo "║     SLG MUD Client Random Message Stress Test          ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo "配置:"
echo "  客户端数量：$CLIENT_COUNT"
echo "  服务器地址：$SERVER_ADDR"
echo "  测试时长：$TEST_DURATION 秒"
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

# 命令池
COMMANDS=(
    "status"
    "look"
    "map"
    "info"
    "who"
    "n"
    "s"
    "e"
    "w"
    "work"
    "rest"
)

# 记录开始时间
START_TIME=$(date +%s.%N)

# 启动客户端
echo "=== 启动 $CLIENT_COUNT 个客户端（随机消息测试）==="
for i in $(seq 1 $CLIENT_COUNT); do
    username="random_test$i"
    password="123456"
    
    # 后台启动客户端，随机发送消息
    (
        sleep 0.1
        echo "register $username $password"
        sleep 0.5
        echo "login $username $password"
        sleep 0.5
        
        # 随机发送消息 TEST_DURATION 秒
        END_TIME=$(($(date +%s) + $TEST_DURATION))
        while [ $(date +%s) -lt $END_TIME ]; do
            # 随机选择一个命令
            cmd=${COMMANDS[$RANDOM % ${#COMMANDS[@]}]}
            echo "$cmd"
            
            # 随机等待 0.1-0.5 秒
            sleep 0.$((RANDOM % 4 + 1))
        done
        
        echo "/quit"
    ) | timeout $((TEST_DURATION + 10)) $MUD_CLIENT > /tmp/random_client_$i.log 2>&1 &
    
    # 每启动 20 个客户端输出一次进度
    if [ $((i % 20)) -eq 0 ]; then
        echo "  已启动 $i/$CLIENT_COUNT 个客户端"
    fi
done

echo ""
echo "=== 等待测试完成（约 $TEST_DURATION 秒）==="

# 等待所有后台进程完成
wait

# 记录结束时间
END_TIME=$(date +%s.%N)
DURATION=$(echo "$END_TIME - $START_TIME" | bc)
TOTAL_SECONDS=$(echo "$DURATION / 1" | bc)

echo ""
echo "=== 压测结果 ==="
echo "总耗时：${DURATION}秒"
echo "客户端数：$CLIENT_COUNT"
echo "测试时长：$TEST_DURATION 秒"
echo "平均每个客户端：$(echo "scale=3; $DURATION / $CLIENT_COUNT" | bc)秒"

# 统计
SUCCESS=0
FAILED=0
TOTAL_COMMANDS=0

for i in $(seq 1 $CLIENT_COUNT); do
    log_file="/tmp/random_client_$i.log"
    
    if grep -q "连接成功\|Connected" $log_file 2>/dev/null; then
        SUCCESS=$((SUCCESS + 1))
        
        # 统计发送的命令数
        cmd_count=$(grep -c "^>" $log_file 2>/dev/null || echo 0)
        TOTAL_COMMANDS=$((TOTAL_COMMANDS + cmd_count))
    else
        FAILED=$((FAILED + 1))
    fi
done

# 计算 TPS
if [ $TOTAL_SECONDS -gt 0 ]; then
    TPS=$(echo "scale=2; $TOTAL_COMMANDS / $TOTAL_SECONDS" | bc)
else
    TPS=0
fi

echo ""
echo "连接统计:"
echo "  成功：$SUCCESS / $CLIENT_COUNT"
echo "  失败：$FAILED / $CLIENT_COUNT"
if [ $CLIENT_COUNT -gt 0 ]; then
    echo "  成功率：$(echo "scale=2; $SUCCESS * 100 / $CLIENT_COUNT" | bc)%"
fi

echo ""
echo "消息统计:"
echo "  总消息数：$TOTAL_COMMANDS"
echo "  平均 TPS: $TPS"
echo "  平均每客户端消息：$(echo "scale=1; $TOTAL_COMMANDS / $SUCCESS" | bc 2>/dev/null || echo 0)"

# 显示部分日志样本
echo ""
echo "=== 日志样本（前 3 个客户端）==="
for i in 1 2 3; do
    echo "--- 客户端 $i ---"
    head -30 /tmp/random_client_$i.log 2>/dev/null | grep -E "(连接 | 登录|>|status|look|map|who)" | head -15
done

# 清理日志
echo ""
echo "清理临时日志..."
rm -f /tmp/random_client_*.log

echo ""
echo "=== 随机消息压测完成 ==="
