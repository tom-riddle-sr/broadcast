package main

import (
	"broadcast/internal/hub"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/gofiber/fiber/v2"
)

var globalHub *hub.ShardedHub = hub.NewShardedHub()

func main() {
	// 啟動 pprof 服務器，監聽在 localhost:6060
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	app := fiber.New()

	app.Post("/connect/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		client := hub.NewClient(id)

		globalHub.AddClient(client)

		// 重點：啟動一個管理生命週期的 Goroutine
		go func() {
			// 只要 Listen 結束（可能是因為 cancel 被呼叫），
			// 就確保從 Hub 移除，徹底回收資源。
			defer globalHub.RemoveClient(client)
			client.Listen()
		}()

		return c.SendString("User " + id + " connected and listening")
	})

	app.Post("/broadcast", func(c *fiber.Ctx) error {
		message := c.FormValue("message")
		globalHub.Broadcast(message)
		return c.SendString("Message broadcasted")
	})

	log.Fatal(app.Listen(":8080"))

}
