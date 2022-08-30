package dto

// ProjectListQuery 定义查询项目列表结构体
type ProjectListQuery struct {
	Page int
	Size int
	Name string
	Age  int
	Code string
}
