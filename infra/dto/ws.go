package dto

// WsMessage 定义websocket消息结构体
type WsMessage struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}
