# Go 的 `context` 介紹

`context` 是 Go 用來管理「請求生命週期」的標準工具，常用在 goroutine、HTTP、資料庫查詢、RPC 等情境。

它主要解決三件事：

- 取消通知（cancellation）
- 時間控制（timeout / deadline）
- 傳遞請求範圍資料（request-scoped values）

---

## 1. 為什麼需要 `context`？

當你同時跑多個 goroutine，常會遇到：

- 使用者已離線，但背景工作還在跑
- 查詢超時了，子流程卻沒停
- 上層任務失敗，底下任務不知道要結束

`context` 可以把「停止訊號」一路往下傳，讓所有子工作一致地停止。

---

## 2. 常用建立方式

### `context.Background()`

- 建立根 context（通常在 `main` 或初始化流程當起點）

### `context.WithCancel(parent)`

- 建立可手動取消的子 context
- 回傳 `(ctx, cancel)`

### `context.WithTimeout(parent, d)`

- 建立有超時時間的 context
- 到時間會自動取消

### `context.WithDeadline(parent, t)`

- 指定某個絕對時間點取消

### `context.WithValue(parent, key, value)`

- 傳遞「請求範圍」的小型資料（例如 trace id）

---

## 3. 你現在這行程式碼在做什麼？

```go
ctx, cancel := context.WithCancel(context.Background())
```

意思是：

- `context.Background()` 建立父 context
- `WithCancel(...)` 產生子 context
- `ctx`：給 goroutine 監聽 `<-ctx.Done()`
- `cancel`：主動發出取消訊號（例如 client 斷線）

一旦呼叫 `cancel()`，所有監聽這個 `ctx` 的工作都可以結束。

---

## 4. 最常見的監聽寫法

```go
select {
case <-ctx.Done():
	// 收到取消或超時
	return ctx.Err() // context.Canceled / context.DeadlineExceeded
case msg := <-msgCh:
	// 正常收到訊息
	_ = msg
}
```

---

## 5. 實務範例：搭配 timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // 建議一定要呼叫，避免資源沒釋放

err := doWork(ctx)
if err != nil {
	if errors.Is(err, context.DeadlineExceeded) {
		// 超時
	}
}
```

---

## 6. 實務規範（重要）

- 函式參數慣例：`func XXX(ctx context.Context, ...)`
- `ctx` 盡量當第一個參數
- 有拿到 `cancel` 通常就要 `defer cancel()`
- 不要把 `context` 當成 struct 的長期狀態儲存（除非有明確生命週期設計）
- 不要把大量資料塞進 `WithValue`

---

## 7. 常見錯誤

- 忘記呼叫 `cancel()`，造成 goroutine 或 timer 泄漏
- 不檢查 `ctx.Done()`，導致取消訊號無效
- 把 `context` 拿來傳 business object
- 用 `string` 當 `WithValue` key（容易撞 key）

建議 key 型別：

```go
type ctxKey string

const userIDKey ctxKey = "user_id"
```

---

## 8. 在你的廣播聊天室場景怎麼用

你可以這樣思考：

1. 每個 client 建立一個 `ctx` + `cancel`
2. 負責送訊息的 goroutine 都 `select` 監聽 `ctx.Done()`
3. client 離線時呼叫 `cancel()`
4. 所有相關 goroutine 收到訊號後安全退出

這樣可以避免「人已離線，但訊息處理還卡著」的問題。

---

## 9. 一句話總結

`context` 是 Go 裡「跨 goroutine 傳遞取消與時間控制」的標準機制，讓你的併發程式可以可控地啟動，也可控地結束。
