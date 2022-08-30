package main

import (
	_ "casbin_demo/docs"
	"casbin_demo/model"
	"casbin_demo/util"
	"github.com/gin-gonic/gin"
)

var (
	//Enf *casbin.Enforcer
	//Logger *zap.SugaredLogger
	Router *gin.Engine
	//Jwt    *util.JWT
	//SingleFlight
)

// @title           Swagger Example API
// @version         2.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath  /
func main() {
	model.InitDb()

	util.InitCasbin()
	util.InitJWT()
	util.InitLogger()

	Router = gin.New()
	router(Router)
	Router.Run(":8080")
}

// 获取系统全部路由
func GetAllRoutes() []map[string]string {
	routers := []map[string]string{}

	if Router == nil {
		return routers
	}
	appRouters := Router.Routes()
	for _, route := range appRouters {
		// fmt.Printf("Method: %s, Path: %s \n", route.Method, route.Path)
		routers = append(routers, map[string]string{
			"url":    route.Path,
			"name":   route.Path,
			"method": route.Method,
		})
	}
	return routers
}
