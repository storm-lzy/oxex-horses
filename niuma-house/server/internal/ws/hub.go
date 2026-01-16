package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"niuma-house/internal/middleware"
	"niuma-house/internal/model"
	"niuma-house/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// Client WebSocket 客户端
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   uint
	username string
}

// Hub WebSocket 中心
type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	mutex      sync.RWMutex
}

// Message 消息结构
type Message struct {
	Type       string `json:"type"` // message, notification
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
}

var hub *Hub
var once sync.Once

// GetHub 获取 Hub 单例
func GetHub() *Hub {
	once.Do(func() {
		hub = &Hub{
			clients:    make(map[uint]*Client),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			broadcast:  make(chan *Message, 256),
		}
		go hub.run()
	})
	return hub
}

// run 运行 Hub
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client
			h.mutex.Unlock()
			log.Printf("WebSocket client registered: userID=%d", client.userID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mutex.Unlock()
			log.Printf("WebSocket client unregistered: userID=%d", client.userID)

		case message := <-h.broadcast:
			h.mutex.RLock()
			// 发送给接收者
			if client, ok := h.clients[message.ReceiverID]; ok {
				data, _ := json.Marshal(message)
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client.userID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// SendMessage 发送消息
func (h *Hub) SendMessage(msg *Message) {
	h.broadcast <- msg
}

// IsOnline 检查用户是否在线
func (h *Hub) IsOnline(userID uint) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	username := middleware.GetCurrentUsername(c)

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		hub:      GetHub(),
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID,
		username: username,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// readPump 读取消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	messageRepo := repository.NewMessageRepository()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		msg.SenderID = c.userID
		msg.Timestamp = time.Now().Unix()

		// 持久化消息
		dbMsg := &model.Message{
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			IsRead:     false,
		}
		if err := messageRepo.Create(dbMsg); err != nil {
			log.Printf("Failed to save message: %v", err)
		}

		// 发送给接收者
		c.hub.SendMessage(&msg)

		// 发送回执给发送者
		msg.Type = "sent"
		response, _ := json.Marshal(msg)
		c.send <- response
	}
}

// writePump 写入消息
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
