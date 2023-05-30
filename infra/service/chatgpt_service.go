package service

import (
	"fmt"
	"sync"
)

var (
	chatGptServiceIns     IChatGptService
	chatGptServiceInsOnce sync.Once
)

type IChatGptService interface {
	analysis(data string) error
}

type chatGptService struct {
}

func (c *chatGptService) analysis(data string) error {
	fmt.Println("analysis")
	return nil
}

func GetChatGptService() IChatGptService {
	chatGptServiceInsOnce.Do(func() {
		chatGptServiceIns = &chatGptService{}
	})
	return chatGptServiceIns
}
