#!/bin/bash

# 批量清理所有brand logic文件中的brand_admin权限检查

echo "开始清理brand_admin引用..."

for file in create_brand_asset_logic.go delete_brand_asset_logic.go get_brand_assets_logic.go update_brand_asset_logic.go get_brands_logic.go get_brand_stats_logic.go; do
  echo ""
  echo "=== 处理: $file ==="
  
  # 备份文件
  cp "$file" "${file}.backup"
  
  # 查找并删除包含brand_admin的行（包括userId和role声明）
  # 使用awk进行精确处理
  
  awk '
  BEGIN {
    skip = 0
    deleted = 0
  }
  
  # 跳过userId和role声明
  /userId := int64\(userInfo\["userId"\]\)\.(float64\)\)/ {
    print "  跳过userId声明"
    next
  }
  
  /role := userInfo\["role"\]\.\(string\)\)/ {
    print "  跳过role声明"
    next
  }
  
  # 跳过整个brand_admin权限检查块
  /if role == "brand_admin"/ {
    print "  开始跳过brand_admin权限检查块"
    skip = 1
    next
  }
  
  # 跳过块中的所有行
  skip == 1 && /}/ && brand_admin引用行/ {
    print "  跳过块中的行"
    next
  }
  
  # 删除else if判断（不再需要）
  /} else if role != "platform_admin"/ {
    print "  删除不需要的else判断"
    next
  }
  
  # 输出所有其他行
  {
    deleted++
    print
  }
  
  END {
    print "  共删除" deleted "行"
  }
  ' "$file" > "${file}.tmp"
  
  # 替换原文件
  mv "${file}.tmp" "$file"
  
  if [ $? -eq 0 ]; then
    echo "  ✅ 修复成功"
    rm "${file}.backup"
  else
    echo "  ❌ 修复失败，恢复备份"
    mv "${file}.backup" "$file"
  fi
done

echo ""
echo "所有文件处理完成"
echo "请验证："
echo "  grep -rn \"brand_admin\" *.go | wc -l"
