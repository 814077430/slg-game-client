#!/bin/bash

# SLG MUD 客户端视野同步压测脚本
# 启动 100 个客户端，测试玩家视野同步功能

CLIENT_COUNT=${1:-100}
SERVER_ADDR=${2:-"localhost:8080"}
MUD_CLIENT="./mud-client"
VIEW_RANGE=10  # 视野范围（格）

echo "╔════════════════════════════════════════════════════════╗"
echo "║     SLG MUD Client Vision Sync Stress Test             ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo "配置:"
echo "  客户端数量：$CLIENT_COUNT"
echo "  服务器地址：$SERVER_ADDR"
echo "  视野范围：$VIEW_RANGE 格"
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

# 启动客户端
echo "=== 启动 $CLIENT_COUNT 个客户端（视野同步测试）==="
for i in $(seq 1 $CLIENT_COUNT); do
    username="vision_test$i"
    password="123456"
    
    # 计算出生位置（集中在中心区域，确保在彼此视野内）
    # 中心坐标 512,512，视野范围 10 格
    # 让玩家分布在 512±50 范围内，确保大部分玩家在视野内
    offset_x=$((RANDOM % 100 - 50))
    offset_y=$((RANDOM % 100 - 50))
    start_x=$((512 + offset_x))
    start_y=$((512 + offset_y))
    
    # 后台启动客户端，模拟登录、移动、查询视野
    (
        sleep 0.05  # 错开启动时间
        echo "register $username $password"
        sleep 0.3
        echo "login $username $password"
        sleep 0.3
        echo "status"
        sleep 0.3
        echo "who"  # 查询视野内玩家 ⭐
        sleep 0.3
        echo "n"    # 向北移动（触发视野同步）
        sleep 0.3
        echo "who"  # 再次查询视野
        sleep 0.3
        echo "e"    # 向东移动
        sleep 0.3
        echo "who"
        sleep 0.3
        echo "/quit"
    ) | timeout 15 $MUD_CLIENT > /tmp/vision_client_$i.log 2>&1 &
    
    # 每启动 20 个客户端输出一次进度
    if [ $((i % 20)) -eq 0 ]; then
        echo "  已启动 $i/$CLIENT_COUNT 个客户端"
    fi
done

echo ""
echo "=== 等待客户端完成 ==="

# 等待所有后台进程完成
wait

# 记录结束时间
END_TIME=$(date +%s.%N)
DURATION=$(echo "$END_TIME - $START_TIME" | bc)

echo ""
echo "=== 压测结果 ==="
echo "总耗时：${DURATION}秒"
echo "客户端数：$CLIENT_COUNT"
echo "平均每个客户端：$(echo "scale=3; $DURATION / $CLIENT_COUNT" | bc)秒"

# 统计成功/失败
SUCCESS=0
FAILED=0
WHO_SUCCESS=0
VISION_SYNC=0

for i in $(seq 1 $CLIENT_COUNT); do
    log_file="/tmp/vision_client_$i.log"
    
    if grep -q "连接成功\|Connected" $log_file 2>/dev/null; then
        SUCCESS=$((SUCCESS + 1))
        
        # 检查 who 命令是否执行
        if grep -q "视野内玩家\|players in vision\|WhoResponse" $log_file 2>/dev/null; then
            WHO_SUCCESS=$((WHO_SUCCESS + 1))
        fi
        
        # 检查是否有视野同步通知
        if grep -q "PlayerMove\|PlayerEnter\|视野" $log_file 2>/dev/null; then
            VISION_SYNC=$((VISION_SYNC + 1))
        fi
    else
        FAILED=$((FAILED + 1))
    fi
done

echo ""
echo "连接统计:"
echo "  成功：$SUCCESS / $CLIENT_COUNT"
echo "  失败：$FAILED / $CLIENT_COUNT"
if [ $CLIENT_COUNT -gt 0 ]; then
    echo "  成功率：$(echo "scale=2; $SUCCESS * 100 / $CLIENT_COUNT" | bc)%"
fi

echo ""
echo "视野同步统计:"
echo "  Who 查询成功：$WHO_SUCCESS / $SUCCESS"
echo "  视野同步通知：$VISION_SYNC / $SUCCESS"

# 显示部分日志样本
echo ""
echo "=== 日志样本（前 5 个客户端）==="
for i in 1 2 3 4 5; do
    echo "--- 客户端 $i ---"
    head -20 /tmp/vision_client_$i.log 2>/dev/null | grep -E "(连接 | 登录|who|视野|Player)" | head -10
done

# 清理日志
echo ""
echo "清理临时日志..."
rm -f /tmp/vision_client_*.log

echo ""
echo "=== 视野同步压测完成 ==="
