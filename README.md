# Shutdown Tool (Go Version)

这是一个轻量级的远程控制工具（Airbridge 的极简替代品），使用 Go 语言编写。
**无需本地安装 Go 环境**，本项目已配置 GitHub Actions 自动云端编译。

## 📥 如何下载

你不需要自己编译代码！

1.  点击本项目 GitHub 页面上方的 **Actions** 标签。
2.  点击最新的一个 Workflow Run (通常是 "Build Shutdown Tool")。
3.  在页面底部的 **Artifacts** 区域，点击 `shutdown-tool-windows` 下载。
4.  解压下载的压缩包，即可得到 `shutdown-tool.exe`。

## 🚀 快速开始

### 1. 配置

在 `shutdown-tool.exe` 同级目录下创建一个 `config.yaml` 文件：

```yaml
port: "8080" # 监听端口
commands:
  shutdown: "shutdown /s /t 0" # 关机
  sleep: "rundll32.exe powrprof.dll,SetSuspendState 0,1,0" # 睡眠
  # 你可以在这里添加任意 cmd 命令
```

### 2. 运行

双击 `shutdown-tool.exe` 即可启动。
看到 `Server starting on port 8080...` 表示启动成功。

### 3. 使用方法 (手机端)

确保手机和电脑在同一 Wi-Fi 网络下。

1.  **获取电脑 IP**：在电脑终端输入 `ipconfig` 查看 IPv4 地址 (例如 `192.168.1.100`)。
2.  **发送命令**：
    -   **关机**：浏览器访问 `http://192.168.1.100:8080/execute/shutdown`
    -   **睡眠**：浏览器访问 `http://192.168.1.100:8080/execute/sleep`

推荐使用 Android App **"HTTP Shortcuts"** 创建桌面快捷方式，一键触发。

## 🛠️ 自己编译 (可选)

如果你确实想自己编译：

```bash
go build -o shutdown-tool.exe main.go
```

## 许可证

MIT
