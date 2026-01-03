#!/bin/bash

# 测试后端 API 接口

BASE_URL="http://localhost:8888"

echo "=== Day 2 Backend API 测试 ==="
echo ""

# 测试1: 创建活动
echo "测试1: 创建活动"
curl -X POST "${BASE_URL}/api/v1/campaigns" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "春季招生活动",
    "description": "推荐即送现金奖励",
    "formFields": ["姓名", "手机号", "意向课程"],
    "rewardRule": 15.00,
    "startTime": "2026-01-01T00:00:00Z",
    "endTime": "2026-03-31T23:59:59Z"
  }'
echo -e "\n"

# 测试2: 获取活动列表
echo "测试2: 获取活动列表"
curl "${BASE_URL}/api/v1/campaigns?page=1&pageSize=10"
echo -e "\n"

# 测试3: 获取活动详情
echo "测试3: 获取活动详情 (ID=1)"
curl "${BASE_URL}/api/v1/campaigns/detail?id=1"
echo -e "\n"

# 测试4: 更新活动
echo "测试4: 更新活动 (ID=1)"
curl -X PUT "${BASE_URL}/api/v1/campaigns/update?id=1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "春季招生活动(已更新)",
    "description": "推荐即送现金奖励 - 活动延期",
    "formFields": ["姓名", "手机号", "意向课程", "备注"],
    "rewardRule": 20.00,
    "startTime": "2026-01-01T00:00:00Z",
    "endTime": "2026-04-30T23:59:59Z",
    "status": "active"
  }'
echo -e "\n"

# 测试5: 搜索活动
echo "测试5: 搜索活动 (关键词: 招生)"
curl "${BASE_URL}/api/v1/campaigns?keyword=招生"
echo -e "\n"

echo "=== 测试完成 ==="
