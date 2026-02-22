# Shutdown Tool (Windows Service Edition)

这是一个轻量级的远程控制工具，使用 Go 语言编写。
**零配置，开箱即用。** 它不仅可以直接运行，还可以**注册为 Windows 服务**，实现开机自启和后台静默运行。

## 📥 下载

请前往 GitHub Actions 页面下载最新的 `shutdown-tool-windows` 构建产物。

## 🚀 使用方法

### 1. 方式一：直接运行 (调试用)

双击 `shutdown-tool.exe` 或在命令行运行。
此时会有黑框窗口，关闭窗口程序就会停止。

### 2. 方式二：安装为服务 (推荐)

以**管理员身份**打开 CMD 或 PowerShell，进入程序所在目录：

```bash
# 安装服务
shutdown-tool.exe install

# 启动服务
shutdown-tool.exe start
```

安装后，程序将在后台静默运行，开机自启，无黑框干扰。

**其他命令：**
- 停止服务：`shutdown-tool.exe stop`
- 卸载服务：`shutdown-tool.exe uninstall`

### 3. 手机端控制

确保手机和电脑在同一 Wi-Fi，且防火墙允许 **8080** 端口。

-   **关机**：`http://<电脑IP>:8080/execute/shutdown`
-   **睡眠**：`http://<电脑IP>:8080/execute/sleep`
-   **锁屏**：`http://<电脑IP>:8080/execute/lock`
-   **取消关机**：`http://<电脑IP>:8080/execute/abort`

## 许可证

MIT
