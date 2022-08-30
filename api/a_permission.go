package api

import (
	"casbin_demo/message"
	"casbin_demo/model"
	"casbin_demo/util"
	"errors"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CasbinApi struct{}

type CasbinService struct{}

var (
	CasbinServiceApp = new(CasbinService)
)

//api级方法

// UpdateCasbinByRole pass
// @Tags Casbin
// @Summary 更新角色api权限
// @accept application/json
// @Produce application/json
// @Param data body message.CasbinInReceive true "角色id, 权限模型列表"
// @Success 200 {string} string "更新角色api权限"
// @Router /role/permission [put]
func (cas *CasbinApi) UpdateCasbinByRole(c *gin.Context) {
	var input message.CasbinInReceive
	_ = c.ShouldBindJSON(&input)

	if err := CasbinServiceApp.UpdateRoleCasbin(input.RoleId, input.CasbinInfos); err != nil {
		util.Logger.Error("更新权限失败!", zap.Error(err))
		util.FailWithMessage("更新权限失败", c)
	} else {
		util.OkWithMessage("更新权限成功", c)
	}
}

// GetPolicyPathByRole pass
// @Tags Casbin
// @Summary 获取权限列表
// @accept application/json
// @Produce application/json
// @Param data body message.CasbinInReceive true "角色id"
// @Success 200 {object} message.PolicyPathResponse "获取权限列表,返回包括casbin详情列表"
// @Router /role/permission [get]
func (cas *CasbinApi) GetPolicyPathByRole(c *gin.Context) {
	roleId := c.Request.FormValue("roleId")

	paths := CasbinServiceApp.GetPolicyPathByRole(roleId)
	util.OkWithDetailed(message.PolicyPathResponse{Paths: paths}, "获取成功", c)
}

// AddPolicyPathByRole pass
// @Tags Casbin
// @Summary 给角色新增权限
// @accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param data body message.CasbinInReceive true "角色id, 权限模型列表"
// @Success 200 {string} string "给角色新增权限"
// @Router /role/permission [post]
func (cas *CasbinApi) AddPolicyPathByRole(c *gin.Context) {
	var input message.CasbinInReceive
	_ = c.ShouldBindJSON(&input)

	if err := CasbinServiceApp.AddPolicyPathByRole(input.RoleId, input.CasbinInfos); err != nil {
		util.Logger.Error("新增权限失败!", zap.Error(err))
		util.FailWithMessage("新增权限失败", c)
	} else {
		util.OkWithMessage("新增权限成功", c)
	}
}

// AddRolePolicyExtends
// @Tags Casbin
// @Summary 角色继承
// @accept application/json
// @Produce application/json
// @Param data body message.CasbinInReceive true "角色id, 继承角色id"
// @Success 200 {string} string "给角色新增继承"
// @Router /role/permission-extends [post]
func (cas *CasbinApi) AddRolePolicyExtends(c *gin.Context) {
	var input message.CasbinInReceive
	_ = c.ShouldBindJSON(&input)

	if err := CasbinServiceApp.AddRoleExtends(input.RoleId, input.ExtendsRoleIds); err != nil {
		util.Logger.Error("新增角色继承失败!", zap.Error(err))
		util.FailWithMessage("新增角色继承失败", c)
	} else {
		util.OkWithMessage("新增角色继承成功", c)
	}
}

// RemoveRolePolicyExtends
// @Tags Casbin
// @Summary 删除角色继承
// @accept application/json
// @Produce application/json
// @Param data body message.CasbinInReceive true "角色id, 继承角色id"
// @Success 200 {string} string "删除角色继承"
// @Router /role/permission-extends [delete]
func (cas *CasbinApi) RemoveRolePolicyExtends(c *gin.Context) {
	var input message.CasbinInReceive
	_ = c.ShouldBindJSON(&input)

	if err := CasbinServiceApp.RemoveRoleExtends(input.RoleId, input.ExtendsRoleIds); err != nil {
		util.Logger.Error("删除角色继承失败!", zap.Error(err))
		util.FailWithMessage("删除角色继承失败", c)
	} else {
		util.OkWithMessage("删除角色继承成功", c)
	}

}

// AddRole
// AddPolicyPathByRole,用新角色名就是新增角色
// AddRolePolicyExtends，给新角色继承也是新增角色

// RemoveRole
// @Tags Casbin
// @Summary 删除角色
// @accept application/json
// @Produce application/json
// @Param data body message.CasbinInReceive true "角色id"
// @Success 200 {string} string "删除角色，改角色被人继承时不能被删除，请先取消继承关系"
// @Router /role/permission  [delete]
func (cas *CasbinApi) RemoveRole(c *gin.Context) {
	var input message.CasbinInReceive
	_ = c.ShouldBindJSON(&input)

	if err := CasbinServiceApp.RemoveRole(input.RoleId); err != nil {
		util.Logger.Error("删除角色失败!", zap.Error(err))
		util.FailWithMessage("删除角色失败", c)
	} else {
		util.OkWithMessage("删除角色成功", c)
	}
}

//Service层方法

// UpdateRoleCasbin
// 跟新casbin api
func (casbinService *CasbinService) UpdateRoleCasbin(role string, casbinInfos []*message.CasbinInfo) error {
	casbinService.ClearCasbin(0, role)
	var rules [][]string
	for _, v := range casbinInfos {
		rules = append(rules, []string{role, v.Path, v.Method})
	}
	e := util.Enf
	_, err := e.AddPolicies(rules)
	return err
}

// GetPolicyPathByRole
// 根据角色获得对应权限列表
func (casbinService *CasbinService) GetPolicyPathByRole(role string) (pathMaps []*message.CasbinInfo) {
	e := util.Enf
	list, _ := e.GetImplicitPermissionsForUser(role)
	for _, v := range list {
		pathMaps = append(pathMaps, &message.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

// AddPolicyPathByRole
// 给角色添加权限
func (casbinService *CasbinService) AddPolicyPathByRole(role string, casbinInfos []*message.CasbinInfo) error {
	var rules [][]string
	for _, v := range casbinInfos {
		rules = append(rules, []string{role, v.Path, v.Method})
	}
	_, err := util.Enf.AddPolicies(rules)
	return err
}

// AddRoleExtends
// 给角色添加继承角色
func (casbinService *CasbinService) AddRoleExtends(role string, extendIds []string) error {
	_, err := util.Enf.AddRolesForUser(role, extendIds)
	return err
}

// RemoveRoleExtends
// 删除角色继承角色
func (casbinService *CasbinService) RemoveRoleExtends(role string, extendIds []string) error {
	ret := ""
	for _, v := range extendIds {
		_, err := util.Enf.DeleteRoleForUser(role, v)
		if err != nil {
			ret += err.Error()
			ret += "  "
		}
	}
	if len(ret) > 0 {
		return errors.New(ret)
	}
	return nil
}

// RemoveRole
// !!!不能删除带有用户的角色
func (casbinService *CasbinService) RemoveRole(role string) error {
	users, err := util.Enf.GetUsersForRole(role)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("请先解除" + role + "的继承关系")
	}

	_, err = util.Enf.DeleteRole(role)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCasbinApi
// 跟新权限路径
func (casbinService *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := model.Db.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}

// ClearCasbin
// 清除匹配的权限
func (casbinService *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := util.Enf
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}
