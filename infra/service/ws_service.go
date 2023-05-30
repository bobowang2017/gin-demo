package service

import (
	"encoding/json"
	"gin-demo/infra/common"
	"gin-demo/infra/dto"
	"gin-demo/infra/utils/log"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

// WsClient 定义websocket客户端连接信息
type WsClient struct {
	Manager    *WsClientManager
	SocketConn *websocket.Conn
	Send       chan []byte
}

// Read 读取客户端发送过来的消息
func (ws *WsClient) Read() {
	// 出现故障后把当前客户端注销
	defer func() {
		_ = ws.SocketConn.Close()
		ws.Manager.UnRegister <- ws
	}()
	ws.SocketConn.SetReadLimit(common.MaxMessageSize)
	_ = ws.SocketConn.SetReadDeadline(time.Now().Add(common.PongWait))
	ws.SocketConn.SetPongHandler(func(string) error {
		_ = ws.SocketConn.SetReadDeadline(time.Now().Add(common.PongWait))
		return nil
	})
	for {
		_, data, err := ws.SocketConn.ReadMessage()
		if err != nil {
			log.Logger.Error(err.Error())
			break
		}
		var msg dto.WsMessage
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Logger.Error(err.Error())
			break
		}
	}
}

func (ws *WsClient) Write() {
	ticker := time.NewTicker(common.PingPeriod)
	defer func() {
		ticker.Stop()
		ws.SocketConn.Close()
	}()
	for {
		select {
		case message, ok := <-ws.Send:
			ws.SocketConn.SetWriteDeadline(time.Now().Add(common.WriteWait))
			if !ok {
				// The hub closed the channel.
				ws.SocketConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := ws.SocketConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(ws.Send)
			for i := 0; i < n; i++ {
				//w.Write(newline)
				w.Write(<-ws.Send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ws.SocketConn.SetWriteDeadline(time.Now().Add(common.WriteWait))
			if err := ws.SocketConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WsClientManager 定义websocket管理结构体，所有建立连接的websocket均由此结构管理
type WsClientManager struct {
	Clients    map[*WsClient]bool // 记录在线连接
	Broadcast  chan []byte        // 触发消息广播
	Register   chan *WsClient     // 触发建立连接
	UnRegister chan *WsClient     // 触发断开连接
}

var (
	wsClientManager     *WsClientManager
	wsClientManagerOnce sync.Once
)

func GetWsClientManager() *WsClientManager {
	wsClientManagerOnce.Do(func() {
		wsClientManager = &WsClientManager{
			Broadcast:  make(chan []byte),
			Register:   make(chan *WsClient),
			UnRegister: make(chan *WsClient),
			Clients:    make(map[*WsClient]bool),
		}
	})
	return wsClientManager
}

// Start 伴随程序启动
func (wm *WsClientManager) Start() {
	for {
		select {
		case client := <-wm.Register:
			wm.Clients[client] = true
		case client := <-wm.UnRegister:
			if _, ok := wm.Clients[client]; ok {
				delete(wm.Clients, client)
				close(client.Send)
			}
		case message := <-wm.Broadcast:
			for client := range wm.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(wm.Clients, client)
				}
			}
		}
	}
}
