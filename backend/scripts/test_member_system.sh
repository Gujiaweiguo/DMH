#!/bin/bash

# 会员系统测试脚本

echo "=== 会员系统数据库迁移测试 ==="

# 数据库配置
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"root"}
DB_PASS=${DB_PASS:-""}
DB_NAME=${DB_NAME:-"dmh"}

# 执行迁移脚本
echo "执行数据库迁移..."
mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASS} ${DB_NAME} < create_member_system_tables.sql

if [ $? -eq 0 ]; then
    echo "✓ 数据库迁移成功"
else
    echo "✗ 数据库迁移失败"
    exit 1
fi

# 验证表是否创建成功
echo ""
echo "验证表结构..."

tables=(
    "members"
    "member_profiles"
    "member_tags"
    "member_tag_links"
    "member_brand_links"
    "member_merge_requests"
    "export_requests"
)

for table in "${tables[@]}"; do
    result=$(mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASS} ${DB_NAME} -e "SHOW TABLES LIKE '${table}';" 2>/dev/null | grep ${table})
    if [ -n "$result" ]; then
        echo "✓ 表 ${table} 创建成功"
    else
        echo "✗ 表 ${table} 不存在"
    fi
done

# 检查 orders 表是否添加了新字段
echo ""
echo "检查 orders 表字段..."
result=$(mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASS} ${DB_NAME} -e "SHOW COLUMNS FROM orders LIKE 'member_id';" 2>/dev/null | grep member_id)
if [ -n "$result" ]; then
    echo "✓ orders.member_id 字段添加成功"
else
    echo "✗ orders.member_id 字段不存在"
fi

result=$(mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASS} ${DB_NAME} -e "SHOW COLUMNS FROM orders LIKE 'unionid';" 2>/dev/null | grep unionid)
if [ -n "$result" ]; then
    echo "✓ orders.unionid 字段添加成功"
else
    echo "✗ orders.unionid 字段不存在"
fi

# 检查 rewards 表是否添加了新字段
result=$(mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASS} ${DB_NAME} -e "SHOW COLUMNS FROM rewards LIKE 'member_id';" 2>/dev/null | grep member_id)
if [ -n "$result" ]; then
    echo "✓ rewards.member_id 字段添加成功"
else
    echo "✗ rewards.member_id 字段不存在"
fi

echo ""
echo "=== 测试完成 ==="
