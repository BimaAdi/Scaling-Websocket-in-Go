package examples

import (
	"log"

	server "github.com/BimaAdi/use-websocket-scaler/fiber-websocket-server"
	adapter "github.com/BimaAdi/use-websocket-scaler/memoryadapter"
	wsc "github.com/BimaAdi/use-websocket-scaler/websocketscaler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
)

func FiberMemoryServe() {
	engine := html.New("./templates", ".html")
	app := fiber.New(
		fiber.Config{
			Views: engine,
		},
	)

	ws := server.NewFiberWebsocketServer()
	ma := adapter.NewMemoryAdapter()
	scaler := wsc.NewScaler(ws, ma)
	chatNamespace := wsc.InitNamespace()
	chatNamespace.On("/message", func(ws wsc.Adapter, s string) {
		ws.Broadcast(s)
	})
	scaler.Of("/chat", chatNamespace)
	server := ws.Serve(scaler)

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"host": c.Hostname(),
		})
	})

	app.Get("/ws/:namespace", server)

	log.Fatal(app.Listen(":3000"))
}
