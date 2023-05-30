package service

import (
	"encoding/json"
	"fmt"
	"gin-demo/infra/common"
	"gin-demo/infra/dto"
	"gin-demo/infra/utils/log"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

// WsClientManager 定义websocket管理结构体，所有建立连接的websocket均由此结构管理
type WsClientManager struct {
	Clients    map[*WsClient]bool // 记录所有连接
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
			log.Logger.Info("Register +1 :", len(wm.Clients))
		case client := <-wm.UnRegister:
			if _, ok := wm.Clients[client]; ok {
				delete(wm.Clients, client)
				log.Logger.Info("UnRegister", len(wm.Clients))
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

// WsClient 定义websocket客户端连接信息
type WsClient struct {
	Manager        *WsClientManager
	SocketConn     *websocket.Conn
	ChatGptService IChatGptService
	Send           chan []byte
}

type NewWsClientOpt func(*WsClient)

func WithConnOpt(conn *websocket.Conn) NewWsClientOpt {
	return func(client *WsClient) {
		client.SocketConn = conn
	}
}

func GetWsClient(opts ...NewWsClientOpt) *WsClient {
	client := &WsClient{
		Manager:        GetWsClientManager(),
		ChatGptService: GetChatGptService(),
		Send:           make(chan []byte, 256),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// Read 读取客户端发送过来的消息
func (ws *WsClient) Read() {
	var ticker *time.Ticker
	// 出现故障后把当前客户端注销
	defer func() {
		_ = ws.SocketConn.Close()
		ws.Manager.UnRegister <- ws
	}()
	defer ticker.Stop()
	ws.SocketConn.SetReadLimit(common.MaxMessageSize)
	_ = ws.SocketConn.SetReadDeadline(time.Now().Add(common.PongWait))
	ws.SocketConn.SetPongHandler(func(string) error {
		_ = ws.SocketConn.SetReadDeadline(time.Now().Add(common.PongWait))
		return nil
	})
	ticker = time.NewTicker(3 * time.Second)
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
			continue
		}
		//ws.ChatGptService.analysis("hahahaha")
		ws.Send <- []byte("开始响应")
		cnt := 0
		for cnt < 5 {
			select {
			case <-ticker.C:
				ws.Send <- []byte(fmt.Sprintf("响应中{%d}", cnt))
				cnt += 1
			}
		}
		ws.Send <- []byte("结束响应")
	}
}

func (ws *WsClient) Write() {
	ticker := time.NewTicker(common.PingPeriod)
	defer func() {
		ticker.Stop()
		_ = ws.SocketConn.Close()
	}()
	for {
		select {
		case message, ok := <-ws.Send:
			_ = ws.SocketConn.SetWriteDeadline(time.Now().Add(common.WriteWait))
			if !ok {
				// The hub closed the channel.
				_ = ws.SocketConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := ws.SocketConn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Logger.Error(err)
			}
		case <-ticker.C:
			log.Logger.Info("Ping")
			_ = ws.SocketConn.SetWriteDeadline(time.Now().Add(common.WriteWait))
			if err := ws.SocketConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
