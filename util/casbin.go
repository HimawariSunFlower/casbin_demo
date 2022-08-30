package util

import (
	"casbin_demo/model"
	"github.com/casbin/casbin/v2"
)

var Enf *casbin.Enforcer

func InitCasbin() {
	var err error
	Enf, err = casbin.NewEnforcer("./rbac_model.conf", model.Adapter) //Adapter
	if err != nil {
		panic(err)
	}
	Enf.EnableAutoSave(true)
}
