package common

const (
	//系统常量
	Success         = "ok"
	InvalidParams   = "请求参数错误"
	ParamEmptyError = "参数[%s]不能为空"
	GoroutineError  = "Goroutine内部函数异常"

	//系统配置文件
	ParseSysCfgError = "系统配置文件解析异常"
	GetSysCfgError   = "系统配置文件读取异常"

	//Redis
	RedisServerError  = "Redis服务异常"
	RedisSetError     = "写入Redis异常"
	GetRedisLockError = "获取Redis锁失败"
	GetRedisLockOK    = "获取Redis锁成功"

	//MySQL
	SearchMySQLError     = "查询MySQL数据库异常"
	VisitUserCenterError = "访问用户中心微服务异常"

	NotAuth = "没有权限" //用户不具备相应的操作权限，例如非本部门的部门管理员，修改了部门用户列表

	JsonParseError     = "JSON解析异常"
	ParseRespBodyError = "解析Body异常"

)
