package middlewares

import (
	"encoding/json"
	"fmt"
	"gin-demo/infra/utils"
	jwtToken "gin-demo/infra/utils/jwt_token"
	"gin-demo/infra/utils/log"
	"gin-demo/infra/utils/redis"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token      = c.Request.Header.Get("token")
			noAuthUrls = utils.GetSysCfg().NoAuthUrl
			userId     string
			//user       *m.User
			err error
		)
		if v, ok := noAuthUrls[c.Request.URL.String()]; ok && c.Request.Method == v {
			c.Next()
			return
		}
		if userId, err = getUserId(token); err != nil {
			goto Abort
		}
		fmt.Println(userId)
		//if user, err = service.NewUserService().GetUser(userId); err == nil {
		//	params := map[string]float64{
		//		user.UserId: float64(time.Now().Unix()),
		//	}
		//	_ = redis.ZAdd(common.UserLatestLogin, params)
		//	c.Set("userInfo", user)
		//	c.Next()
		//	return
		//}
		log.Logger.Error(err)
	Abort:
		c.JSON(401, gin.H{
			"status": 401,
			"msg":    "用户认证失败",
			"data":   "failure",
		})
		c.Abort()
	}
}

func getUserId(token string) (userId string, err error) {
	var (
		userIdStr string
		authInfo  *jwtToken.AuthCodeJwtClaims
	)
	if userIdStr, err = redis.Get(token); err == nil && userIdStr != "" {
		err = json.Unmarshal([]byte(userIdStr), &userId)
		return
	}
	if authInfo, err = jwtToken.ParseToken(token); err == nil {
		return authInfo.UserId, err
	}
	return
}
