package hub

import (
	"hash/fnv"
	"sync"
)

type UserID string

// 定義分段數量，通常是 2 的次方 (如 32, 64)
const shardCount = 32

type ShardedHub struct {
	shards []*shard
}

type shard struct {
	mu      sync.RWMutex
	clients map[UserID]*Client
}

var slicePool = sync.Pool{
	New: func() interface{} {
		// 預設一個合適的容量，例如 512
		s := make([]*Client, 0, 512)
		return &s
	},
}

func NewShardedHub() *ShardedHub {
	h := &ShardedHub{
		shards: make([]*shard, shardCount),
	}
	for i := 0; i < shardCount; i++ {
		h.shards[i] = &shard{
			clients: make(map[UserID]*Client),
		}
	}
	return h
}

// getShard 根據 UserID 計算雜湊值，分散到不同的 Shard
// 這樣可以降低簽名競爭，提高並行效能
func (h *ShardedHub) getShard(id UserID) *shard {
	hasher := fnv.New32a()
	hasher.Write([]byte(id))
	return h.shards[hasher.Sum32()%uint32(shardCount)]
}

func (h *ShardedHub) AddClient(c *Client) {
	s := h.getShard(UserID(c.ID))
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[UserID(c.ID)] = c
}

func (h *ShardedHub) RemoveClient(c *Client) {
	s := h.getShard(UserID(c.ID))
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, UserID(c.ID))
}

func (h *ShardedHub) Broadcast(msg string) {
	for _, s := range h.shards {
		s.mu.RLock()
		if len(s.clients) == 0 {
			s.mu.RUnlock()
			continue
		}

		// 從池子借一個 Slice 指標
		p := slicePool.Get().(*[]*Client)
		clients := (*p)[:0] // 重置長度，保留容量

		for _, c := range s.clients {
			clients = append(clients, c)
		}
		s.mu.RUnlock()

		for _, c := range clients {
			select {
			case c.MsgChan <- msg:
			default:
			}
		}

		// 用完放回去，這一步最重要！
		*p = clients
		slicePool.Put(p)
	}
}
