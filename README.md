# Shutdown Tool (Windows Service Edition)

这是一个轻量级的远程控制工具，使用 Go 语言编写。
**零配置，开箱即用。** 它不仅可以直接运行，还可以**注册为 Windows 服务**，实现开机自启和后台静默运行。

## 📥 下载

请前往 GitHub Actions 页面下载最新的 `shutdown-tool-windows` 构建产物。

## 🔍 服务名称

如果你安装了服务，它在 Windows 服务列表 (`services.msc`) 中的名字是：
-   **显示名称**: `Remote Shutdown Service`
-   **服务名称**: `RemoteShutdown`

你可以通过按 `Win + R` 输入 `services.msc` 找到它。

## 🚀 使用方法

### 1. 方式一：放到启动文件夹 (推荐用于睡眠/锁屏)

如果你主要使用 **睡眠 (Sleep)** 或 **锁屏 (Lock)** 功能，**不要安装为服务**。
因为 Windows 服务运行在隔离会话中，无法控制用户的桌面锁屏。

**步骤：**
1.  按 `Win + R`，输入 `shell:startup` 打开启动文件夹。
2.  将 `shutdown-tool.exe` 放入该文件夹。
3.  下次开机它会自动后台运行（无黑框）。

### 2. 方式二：安装为服务 (推荐用于关机)

如果你主要使用 **关机 (Shutdown)** 功能，服务模式最稳定。

以**管理员身份**打开 CMD 或 PowerShell，进入程序所在目录：

```bash
# 安装服务
shutdown-tool.exe install

# 启动服务
shutdown-tool.exe start
```

**日志文件**：程序会在同级目录下生成 `shutdown-tool.log`，如果有问题可以查看该文件。

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
