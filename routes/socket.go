
package routes
import (
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/websocket/v2"
)

// Route khusus WebSocket
func WebSocketRoute(app *fiber.App) {
  // Upgrade only if it's a WebSocket request
  app.Get("/ws", websocket.New(func(c *websocket.Conn) {
    defer c.Close()
    for {
      // Baca pesan dari client
      msgType, msg, err := c.ReadMessage()
      if err != nil {
        break
      }

      // Balas pesan ke client (echo)
      if err := c.WriteMessage(msgType, msg); err != nil {
        break
      }
    }
  }))
}
