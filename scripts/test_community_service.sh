#!/bin/bash

# 社区服务测试脚本
# 本脚本用于测试社区服务的基本功能，包括创建社区、创建群组、发布帖子等

# 颜色输出函数
print_green() {
  echo -e "\033[0;32m$1\033[0m"
}

print_yellow() {
  echo -e "\033[0;33m$1\033[0m"
}

print_red() {
  echo -e "\033[0;31m$1\033[0m"
}

# 设置基本变量
API_URL="http://localhost:8084/api/v1"
TOKEN=""
COMMUNITY_ID=""
GROUP_ID=""
POST_ID=""

# 简单检查社区服务是否在运行
check_service() {
  print_yellow "检查社区服务状态..."
  
  response=$(curl -s -o /dev/null -w "%{http_code}" $API_URL/health)
  
  if [ "$response" == "200" ]; then
    print_green "社区服务正常运行"
    return 0
  else
    print_red "社区服务未运行或异常，状态码: $response"
    print_yellow "请先启动社区服务: cd /Users/wxn/Desktop/wz-project/wz-backend-go/services/community-service && go run main.go"
    return 1
  fi
}

# 用户登录
login() {
  print_yellow "执行用户登录..."
  
  response=$(curl -s -X POST $API_URL/auth/login \
    -H "Content-Type: application/json" \
    -d '{
      "username": "admin",
      "password": "admin123"
    }')
  
  TOKEN=$(echo $response | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
  
  if [ -n "$TOKEN" ]; then
    print_green "登录成功，获取到Token"
  else
    print_red "登录失败，无法获取Token"
    print_yellow "请检查社区服务是否正确运行，以及用户数据是否已初始化"
    exit 1
  fi
}

# 创建社区
create_community() {
  print_yellow "创建测试社区..."
  
  response=$(curl -s -X POST $API_URL/communities \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "name": "测试社区",
      "description": "这是一个用于测试的社区",
      "tags": ["测试", "技术交流"],
      "location": "江苏省-南京市"
    }')
  
  COMMUNITY_ID=$(echo $response | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
  
  if [ -n "$COMMUNITY_ID" ]; then
    print_green "社区创建成功，ID: $COMMUNITY_ID"
  else
    print_red "社区创建失败"
    print_yellow "响应内容: $response"
    exit 1
  fi
}

# 创建群组
create_group() {
  print_yellow "在社区中创建测试群组..."
  
  response=$(curl -s -X POST $API_URL/groups \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "name": "测试群组",
      "description": "这是一个用于测试的群组",
      "community_id": "'$COMMUNITY_ID'",
      "tags": ["测试", "讨论"]
    }')
  
  GROUP_ID=$(echo $response | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
  
  if [ -n "$GROUP_ID" ]; then
    print_green "群组创建成功，ID: $GROUP_ID"
  else
    print_red "群组创建失败"
    print_yellow "响应内容: $response"
  fi
}

# 发布帖子
create_post() {
  print_yellow "在群组中发布测试帖子..."
  
  response=$(curl -s -X POST $API_URL/posts \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "title": "测试帖子标题",
      "content": "这是一个测试帖子的内容，用于验证社区服务的功能是否正常。",
      "community_id": "'$COMMUNITY_ID'",
      "group_id": "'$GROUP_ID'",
      "tags": ["测试"]
    }')
  
  POST_ID=$(echo $response | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
  
  if [ -n "$POST_ID" ]; then
    print_green "帖子发布成功，ID: $POST_ID"
  else
    print_red "帖子发布失败"
    print_yellow "响应内容: $response"
  fi
}

# 添加评论
create_comment() {
  print_yellow "添加测试评论..."
  
  response=$(curl -s -X POST $API_URL/comments \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "content": "这是一条测试评论",
      "post_id": "'$POST_ID'"
    }')
  
  COMMENT_ID=$(echo $response | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
  
  if [ -n "$COMMENT_ID" ]; then
    print_green "评论添加成功，ID: $COMMENT_ID"
  else
    print_red "评论添加失败"
    print_yellow "响应内容: $response"
  fi
}

# 查询社区列表
list_communities() {
  print_yellow "查询社区列表..."
  
  response=$(curl -s -X GET $API_URL/communities)
  
  print_green "社区列表查询结果："
  echo $response | python -m json.tool
}

# 查询特定社区的群组
list_groups() {
  print_yellow "查询社区的群组列表..."
  
  response=$(curl -s -X GET "$API_URL/groups?community_id=$COMMUNITY_ID")
  
  print_green "群组列表查询结果："
  echo $response | python -m json.tool
}

# 查询特定群组的帖子
list_posts() {
  print_yellow "查询群组的帖子列表..."
  
  response=$(curl -s -X GET "$API_URL/posts?group_id=$GROUP_ID")
  
  print_green "帖子列表查询结果："
  echo $response | python -m json.tool
}

# 执行测试流程
main() {
  print_yellow "开始执行社区服务测试..."
  
  # 检查服务状态
  check_service || exit 1
  
  # 执行测试步骤
  login
  create_community
  create_group
  create_post
  create_comment
  
  # 查询数据
  list_communities
  list_groups
  list_posts
  
  print_green "社区服务测试完成!"
}

# 运行主函数
main
