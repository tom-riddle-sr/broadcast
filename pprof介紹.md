# pprof 介紹

**pprof** 是 Go 內建的效能分析（profiling）工具，用於診斷程式的 CPU 用量、記憶體洩漏、goroutine 阻塞等問題。

## 基本概念

pprof 可以收集程式執行時的統計資料（採樣），生成配置檔案（profile），然後用視覺化或文字工具分析。

## 常見用途

| 功能 | 說明 |
|------|------|
| **CPU Profiling** | 找出最耗 CPU 的函數 |
| **Memory Profiling** | 找出記憶體用量最多的位置及洩漏 |
| **Goroutine Profiling** | 查看 goroutine 數量及阻塞狀況 |
| **Block Profiling** | 分析同步原始物件造成的等待 |
| **Mutex Profiling** | 分析互斥鎖的競爭 |

## 啟用方式

### 1. 匯入 `net/http/pprof`（已在你的程式中）
```go
import _ "net/http/pprof"
```
然後在另一個埠啟動 HTTP 伺服器即可公開 pprof 端點。

### 2. 收集與分析
```bash
# 查看即時效能
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 分析記憶體
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看 goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### 3. 視覺化
```bash
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap
```
會在瀏覽器開啟互動式火焰圖。

## 互動模式常用命令

進入 pprof 後可輸入以下命令：

| 命令 | 說明 |
|------|------|
| `top [N]` | 顯示前 N 個最耗資源的函數（預設10個） |
| `list <func>` | 逐行查看特定函數的資源消耗 |
| `web` | 生成圖表（需要 Graphviz） |
| `pdf` | 生成 PDF 報告至 /tmp |
| `quit` | 退出互動模式 |

**範例：**
```
(pprof) top 20          # 看前20個函數
(pprof) list main.main  # 查看 main 函數詳情
(pprof) quit
```

## 快速使用建議

**找 CPU 瓶頸（最常用）：**
```bash
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30
```
直接開瀏覽器火焰圖，不需進互動模式。

**找記憶體洩漏：**
```bash
go tool pprof http://localhost:6060/debug/pprof/heap
(pprof) top    # 看當前占用最多的位置
(pprof) quit
```

**查看 Goroutine 狀況：**
```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
(pprof) top
(pprof) quit
```

## 優勢

- 零侵入：引入匯入後自動可用
- 低負荷：採樣而非全記錄，不大幅拖累效能
- 多維度：涵蓋 CPU、記憶體、並行等多面向
- 視覺化：火焰圖直觀看出熱點函數

你的程式已經在 port 6060 啟動 pprof，可以直接用指令分析！
