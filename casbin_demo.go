package main

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start1() {
	//初始化数据库，policy表
	//InitDb()

	e := initCasbin()

	r := gin.Default()

	gp := r.Group("public")
	{
		gp.POST("/login", func(ctx *gin.Context) {
			//ctx.Request.ParseForm()
			name := ctx.Request.PostFormValue("name")
			//tokenString, _ := GenToken(name)
			ctx.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "success",
				"data": gin.H{"name": name},
			})
		})
	}

	r.Use(NewAuthorizer(e))
	r.GET("menu", func(ctx *gin.Context) {
		username, _, _ := ctx.Request.BasicAuth()
		roles, _ := e.GetImplicitRolesForUser(username)

		var menu [][]string
		for _, v := range roles {
			filteredNamedPolicy := e.GetFilteredNamedPolicy("p", 0, v)
			menu = append(menu, filteredNamedPolicy...)
		}
		filteredNamedPolicy := e.GetFilteredNamedPolicy("p", 0, username)
		menu = append(menu, filteredNamedPolicy...)

		ctx.JSON(200, menu)
	})

	g1 := r.Group("v1/")
	{
		g1.GET("h1", getGinHandler("h1", "get"))
		g1.POST("h1", getGinHandler("h1", "post"))
		g1.GET("h2", getGinHandler("h2", "get"))
	}
	g2 := r.Group("v2/")
	{
		g2.GET("h1", getGinHandler("h1", "get"))
	}
	r.GET("/v3/aa", getGinHandler("/v3/aa", "get"))
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

type ginHandler = func(ctx *gin.Context)

func initCasbin() *casbin.Enforcer {
	e, err := casbin.NewEnforcer("./rbac_model.conf", "./test.csv") //Adapter
	if err != nil {
		panic(err.Error())
	}
	//rules := [][]string{
	//	[]string{"manager", "/v1/*", "*"},
	//	[]string{"it", "/v2/*", "*"},
	//}
	//
	//_, err = e.AddNamedPolicies("p", rules)
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//}
	//
	//_, err = e.AddRoleForUser("alice", "manager")
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//}
	//_, err = e.AddRoleForUser("bob", "it")
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//}
	//e.SavePolicy()
	return e
}
