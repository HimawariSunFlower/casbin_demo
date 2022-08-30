package main

import (
	api2 "casbin_demo/api"
	"casbin_demo/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func router(r *gin.Engine) {
	CasbinApiApp := new(api2.CasbinApi)
	BaseApi := new(api2.BaseApi)
	//MenuApiApp := api2.MenuApi{}

	public := r.Group("")
	private := r.Group("")
	public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := public.Group("base")
	{
		api.POST("/login", BaseApi.Login)
		api.GET("/logout", getGinHandler("logout", "GET"))
	}

	private.Use(util.AuthRequired()).Use(NewAuthorizer(util.Enf))
	{
		private.GET("user", getGinHandler("show", "GET"))
		private.PUT("user", getGinHandler("save", "PUT"))
		private.POST("user", getGinHandler("add", "POST"))
		private.DELETE("user", getGinHandler("delete", "DELETE"))
	}

	{
		private.GET("role", func(context *gin.Context) {
			util.Enf.SavePolicy()
		})
		private.PUT("role", getGinHandler("save", "PUT"))
		private.POST("role", getGinHandler("add", "POST"))
		private.GET("role/tree", getGinHandler("tree", "GET"))
		role := private.Group("role")
		role.GET("permission", CasbinApiApp.GetPolicyPathByRole)
		role.PUT("permission", CasbinApiApp.UpdateCasbinByRole)
		role.POST("permission", CasbinApiApp.AddPolicyPathByRole)
		role.DELETE("permission", CasbinApiApp.RemoveRole)
		role.POST("permission-extends", CasbinApiApp.AddRolePolicyExtends)
		role.DELETE("permission-extends", CasbinApiApp.RemoveRolePolicyExtends)
	}

	{
		private.GET("menu", getGinHandler("show", "GET"))
		private.PUT("menu", getGinHandler("save", "PUT"))
		private.POST("menu", getGinHandler("add", "POST"))
		private.DELETE("menu", getGinHandler("delete", "DELETE"))
		private.GET("menu/tree", getGinHandler("tree", "GET"))
	}

	{
		private.GET("work1", getGinHandler("show", "GET"))
		private.PUT("work1", getGinHandler("save", "PUT"))
		private.POST("work1", getGinHandler("add", "POST"))
		private.DELETE("work1", getGinHandler("delete", "DELETE"))
	}
}

func getGinHandler(str, method string) ginHandler {
	return func(ctx *gin.Context) {
		name, _, _ := ctx.Request.BasicAuth()
		ctx.JSONP(0, fmt.Sprintf("[%s] ,%s 访问 %s", method, name, str))
	}
}
