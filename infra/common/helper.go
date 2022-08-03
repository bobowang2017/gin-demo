package common

import "github.com/gin-gonic/gin"

const (
	//系统编码
	SystemCode = "22"
	//系统可忽略级别
	SystemIgnoreLevel = "2"
	//系统警告级别
	SystemWarnLevel = "3"
	//系统报错级别
	SystemErrorLevel = "4"
)

/*
 * 正常成功返回结果
 */
func RespSuccessJSON(c *gin.Context, data interface{}) {

	c.JSON(200, gin.H{
		"status": "200",
		"msg":    "success",
		"data":   data,
	})
}

/*
 * 分页返回结果
 */
func RespSuccessPageJSON(c *gin.Context, data interface{}, total interface{}) {
	c.JSON(200, gin.H{
		"status": "200",
		"msg":    "success",
		"data":   data,
		"total":  total,
	})
}

func releaseCode(status string, level string) string {
	return level + "-" + SystemCode + "-" + status
}

/*
 * 异常报错-400
 */
func RespInputErrorJSON(c *gin.Context, msg string) {

	c.JSON(400, gin.H{
		"status": releaseCode("400", "4"),
		"msg":    msg,
	})

}

/*
 * 异常报错-401
 */
func RespAuthErrorJSON(c *gin.Context, msg string) {

	c.JSON(401, gin.H{
		"status": releaseCode("401", "4"),
		"msg":    msg,
	})
}

/*
 * 异常报错-404
 */
func RespNoFoundErrorJSON(c *gin.Context, msg string) {
	c.JSON(404, gin.H{
		"status": releaseCode("404", "4"),
		"msg":    msg,
	})
}

/*
 * 异常报错-404-错误信息是列表
 */
func RespNoFoundErrorsJSON(c *gin.Context, msg []string) {
	c.JSON(404, gin.H{
		"status": releaseCode("404", "4"),
		"msg":    msg,
	})
}

/*
 * 异常报错-500
 */
func RespInternalErrorJSON(c *gin.Context, msg string) {
	c.JSON(500, gin.H{
		"status": releaseCode("500", "4"),
		"msg":    msg,
	})
}

/*
 * 异常报错-502
 */
func RespGatewayErrorJSON(c *gin.Context, msg string) {
	c.JSON(502, gin.H{
		"status": releaseCode("502", "4"),
		"msg":    msg,
	})
}

/*
 * 异常报错-501
 */
func RespFuncErrorJSON(c *gin.Context, msg string) {
	c.JSON(501, gin.H{
		"status": releaseCode("501", "4"),
		"msg":    msg,
	})
}
