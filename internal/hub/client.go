package hub

import "context"

type Client struct {
	ID      string
	MsgChan chan string        //接收消息的管道
	Ctx     context.Context    //用於監聽取消事件(ex: 斷線)
	Cancel  context.CancelFunc //用於主動取消監聽事件
}

func NewClient(userID string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ID:      userID,
		MsgChan: make(chan string, 1000000),
		Ctx:     ctx,
		Cancel:  cancel,
	}
}

func (c *Client) Listen() {
	for {
		select {
		case msg := <-c.MsgChan:
			// 這裡可以處理接收到的消息，例如發送給前端或其他處理邏輯
			// 目前只是簡單地打印出來
			println("Client " + c.ID + " received message: " + msg)
		case <-c.Ctx.Done():
			println("Client " + c.ID + " is disconnecting...")
			return
		}
	}
}
