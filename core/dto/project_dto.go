package dto

// ProjectListQuery 定义查询项目列表结构体
type ProjectListQuery struct {
	Page int    `form:"page"`
	Size int    `form:"size"`
	Name string `form:"name"`
	Age  *int   `form:"age"`
	Code string `form:"code"`
}
