package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func NewAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	a := &BasicAuthorizer{enforcer: e}

	return func(c *gin.Context) {
		if !a.CheckPermission(c.Request) {
			a.RequirePermission(c)
		}
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

// GetUserName gets the username from the request.
// Currently, only HTTP basic authentication is supported
func (a *BasicAuthorizer) GetUserName(r *http.Request) string {
	//todo 不绑定basic auth,jwt取userid再去取角色
	username, _, _ := r.BasicAuth()
	return username
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request) bool {
	user := a.GetUserName(r)
	method := r.Method
	path := r.URL.Path

	allowed, err := a.enforcer.Enforce(user, path, method)
	if err != nil {
		fmt.Errorf(err.Error())
		//todo logerror 不应当panic
		//panic(err)
	}

	return allowed
}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, "403 没有访问权限")
}
