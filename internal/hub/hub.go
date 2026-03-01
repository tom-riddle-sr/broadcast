package hub

import (
	"sync"
	"sync/atomic"
)

type UserID string

type Hub struct {
	Clients     map[UserID]*Client // 用戶ID對應的Client
	mu          sync.RWMutex       // 用於保護Clients的讀寫操作
	DroppedMsgs int64              // 紀錄因為Buffer滿了而被丟棄的訊息數
}

// 添加用戶到Hub
func (h *Hub) AddClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Clients[UserID(client.ID)] = client
}

func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.Clients[UserID(client.ID)]; !exists {
		return //用戶不存在，無需移除
	}

	delete(h.Clients, UserID(client.ID))
	client.Cancel()
}

func (h *Hub) Broadcast(message string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.Clients {
		select {
		case client.MsgChan <- message:
			//成功發送消息
		default:
			atomic.AddInt64(&h.DroppedMsgs, 1)
		}
	}
}
