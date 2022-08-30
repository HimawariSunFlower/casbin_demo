package api

import (
	"casbin_demo/model"
	"casbin_demo/util"
)

type MenuApi struct{}
type MenuService struct{}

var (
	MenuServiceApp = MenuService{}
)

// GetList pass
// @Tags Menu
// @Summary 获得用户菜单
// @accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object}  message.LoginResp  ""
// @Router /menu [get]
func (cas *MenuApi) GetList() {

	model.Db.Model(&model.SysBaseMenu{}).Where(" = ?")
}

// 返回casbin的menu，需要根据router的路径再转化成菜单

func getMenuByUser(username string) [][]string {
	var menu [][]string
	e := util.Enf
	roles, _ := e.GetImplicitRolesForUser(username)
	for _, v := range roles {
		filteredNamedPolicy := e.GetFilteredNamedPolicy("p", 0, v)
		menu = append(menu, filteredNamedPolicy...)
	}
	filteredNamedPolicy := e.GetFilteredNamedPolicy("p", 0, username)
	menu = append(menu, filteredNamedPolicy...)
	return menu
}

func GetMenuByUser(username string) [][]string {
	//casbinMenus := getMenuByUser(username)

	//todo

	return nil
}
