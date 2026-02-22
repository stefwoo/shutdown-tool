# Shutdown Tool (Python Version)

这是一个极其简单的 Python 脚本，允许你通过 HTTP 请求远程控制你的电脑（关机、睡眠、执行脚本等）。

它模仿了 **Airbridge** 的核心功能，但只用了一个 Python 文件和一个 JSON 配置文件即可运行，无需编译，依赖仅仅是 Python 标准库。

## 特性

-   🐍 **纯 Python**：只需安装 Python 即可运行，无第三方依赖。
-   ⚙️ **简单配置**：通过 `config.json` 自定义命令。
-   🚀 **极速启动**：双击脚本即可运行。

## 快速开始

### 1. 准备环境

确保电脑上安装了 Python (建议 3.6+)。

### 2. 下载与配置

1.  下载 `server.py` 和 `config.json` 到同一个文件夹。
2.  (可选) 修改 `config.json` 定义你的命令：

```json
{
  "port": 8080,
  "commands": {
    "shutdown": "shutdown /s /t 0",
    "sleep": "rundll32.exe powrprof.dll,SetSuspendState 0,1,0",
    "lock": "rundll32.exe user32.dll,LockWorkStation"
  }
}
```

### 3. 运行

-   **Windows**: 双击 `server.py` (如果安装 Python 时勾选了关联文件)，或者在命令行运行：
    ```bash
    python server.py
    ```

看到 `Serving on port 8080` 即表示启动成功。

### 4. 使用方法

确保手机和电脑在同一局域网。

1.  **获取电脑 IP**：在命令行输入 `ipconfig` 查看 IPv4 地址。
2.  **发送请求**：
    -   关机：浏览器访问 `http://<电脑IP>:8080/execute/shutdown`
    -   睡眠：浏览器访问 `http://<电脑IP>:8080/execute/sleep`
    -   查看所有命令：直接访问 `http://<电脑IP>:8080/`

### 5. Android 配合使用

推荐使用 **HTTP Shortcuts** App 创建桌面快捷方式，一键触发。

## 许可证

MIT
