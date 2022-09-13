package dto

// ProjectListQuery 定义查询项目列表结构体
type ProjectListQuery struct {
	Page int    `form:"page"`
	Size int    `form:"size"`
	Name string `form:"name"`
	Age  *int   `form:"age"`
	Code string `form:"code"`
}

// 定义创建项目结构体
type ProjectCreateDto struct {
	Name string `binding:"required,min=3,max=10,not_allow_blank"`
	Age  int    `binding:"required,gt=0"`
	Code string `binding:"required,min=3,max=10"`
}
