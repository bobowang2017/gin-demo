package dto

type BasePageQuery struct {
	Page  int `form:"page" binding:"omitempty,gt=0"`
	Limit int `form:"limit" binding:"omitempty,gt=0"`
}
