# OpenCode 使用指南

## 安装

### 最简单的方式
```bash
curl -fsSL https://opencode.ai/install | bash
```

### 其他安装方式

**使用包管理器：**
- **npm：** `npm install -g opencode-ai`
- **Bun：** `bun install -g opencode-ai`
- **Homebrew (macOS/Linux)：** `brew install anomalyco/tap/opencode`
- **Windows Chocolatey：** `choco install opencode`
- **Windows Scoop：** `scoop install extras/opencode`
- **Arch Linux：** `paru -S opencode-bin`

**使用 Docker：**
```bash
docker run -it --rm ghcr.io/anomalyco/opencode
```

## 配置

### 1. 连接 LLM 提供商
```bash
opencode
/connect
```

### 2. 推荐使用 OpenCode Zen
- 选择 opencode 提供商
- 访问 opencode.ai/auth
- 注册并添加付款信息
- 复制 API 密钥并粘贴

### 3. 其他提供商
- OpenAI (GPT)
- Anthropic (Claude)
- Google (Gemini)
- 75+ 其他提供商

## 项目初始化

```bash
cd /path/to/your/project
opencode
/init
```

这会：
- 分析您的项目结构
- 创建 `AGENTS.md` 文件
- 理解项目编码模式

## 基本使用

### 启动 OpenCode
```bash
opencode
```

### 模式切换
- **Tab 键**：在计划模式和构建模式间切换
- 右下角显示当前模式

### 文件引用
使用 `@` 符号引用文件：
```
How is authentication handled in @src/auth.js
```

### 主要功能

#### 1. 提问
```
解释 @src/components/Header.tsx 的工作原理
```

#### 2. 添加功能
**步骤：**
1. 按 Tab 进入计划模式
2. 描述需求：
   ```
   当用户删除笔记时，在数据库中标记为已删除。
   创建显示最近删除笔记的页面。
   用户可以恢复或永久删除笔记。
   ```
3. 审查计划
4. 按 Tab 返回构建模式
5. 实施功能：
   ```
   Sounds good! Go ahead and make the changes
   ```

#### 3. 直接修改
```
给 /settings 路由添加身份验证，
参考 @src/notes.ts 中的实现方式
```

#### 4. 撤销/重做
```bash
/undo    # 撤销更改
/redo    # 重做更改
```

#### 5. 分享会话
```bash
/share   # 创建分享链接
```

## 高级功能

### 拖放图片
- 拖拽图片到终端
- OpenCode 会扫描并添加到对话中
- 可用作设计参考

### 多会话
- 支持并行运行多个代理
- 在同一项目上同时工作

### LSP 支持
- 自动加载适合的语言服务器
- 提供智能代码补全

## 常用命令

```bash
/help           # 获取帮助
/connect        # 连接提供商
/init           # 初始化项目
/undo           # 撤销更改
/redo           # 重做更改
/share          # 分享会话
/clear          # 清屏
/exit           # 退出
```

## 自定义配置

### 主题选择
```bash
# 在配置文件中设置
theme: "dark"   # 或 "light", "auto"
```

### 键位绑定
自定义快捷键以提高效率。

### 代码格式化
配置项目使用的格式化工具。

## 故障排除

### 常见问题
1. **终端兼容性**：使用现代终端（WezTerm、Alacritty、Kitty）
2. **API 密钥**：确保密钥正确配置
3. **网络问题**：检查防火墙设置

### 获取帮助
- Discord 社区：https://opencode.ai/discord
- GitHub Issues：https://github.com/anomalyco/opencode/issues
- 文档：https://opencode.ai/docs

## 最佳实践

1. **详细描述**：给 OpenCode 充足的上下文
2. **分步进行**：复杂功能先制定计划
3. **及时撤销**：不满意的更改立即撤销
4. **保存配置**：将 AGENTS.md 提交到 Git
5. **定期更新**：保持 OpenCode 为最新版本

---

## 快速参考

| 功能 | 命令/快捷键 |
|------|-------------|
| 启动 | `opencode` |
| 计划模式 | `Tab` |
| 引用文件 | `@文件名` |
| 撤销 | `/undo` |
| 重做 | `/redo` |
| 分享 | `/share` |
| 帮助 | `/help` |
| 退出 | `/exit` |

💡 **提示**：像和初级开发者对话一样与 OpenCode 交流，给足上下文和示例！