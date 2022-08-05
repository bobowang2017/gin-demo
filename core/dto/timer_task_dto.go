package dto

type TimerTaskDto struct {
	Url    string
	Args   map[string]string
	Method string
	Header map[string]string
}

type TimerTaskAddDto struct {
	Name        string            `binding:"required"`
	Params      map[string]string `binding:"required"`
	Description string            `binding:"omitempty"`
	Cron        string            `binding:"required"`
	StopAt      string            `binding:"omitempty"`
}
