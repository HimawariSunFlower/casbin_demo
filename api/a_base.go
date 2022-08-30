package api

import (
	"casbin_demo/message"
	"casbin_demo/model"
	"casbin_demo/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BaseApi struct{}
type BaseApiService struct{}

var (
	BaseApiServiceApp = new(BaseApiService)
)

// Login
// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body message.LoginReq true "用户名, 密码, 验证码"
// @Success 200 {object} message.LoginResp  "返回包括用户信息,token,过期时间"
// @Router /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l message.LoginReq
	_ = c.ShouldBindJSON(&l)

	u := &model.User{UserName: l.Username, Password: l.Password}
	if user, err := BaseApiServiceApp.Login(u); err != nil {
		util.Logger.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		util.FailWithMessage("用户名不存在或者密码错误", c)
	} else {
		if user.Status != 1 {
			util.Logger.Error("登陆失败! 用户被禁止登录!")
			util.FailWithMessage("用户被禁止登录", c)
			return
		}
		b.TokenNext(c, *user)
	}

}

// TokenNext 登录以后签发jwt
func (b *BaseApi) TokenNext(c *gin.Context, user model.User) {
	j := &util.JWT{SigningKey: []byte(util.SigningKey)} // 唯一签名
	claims := j.CreateClaims(util.BaseClaims{
		Uid:      user.ID,
		Username: user.UserName,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		util.Logger.Error("获取token失败!", zap.Error(err))
		util.FailWithMessage("获取token失败", c)
		return
	}
	//if !global.GVA_CONFIG.System.UseMultipoint {
	util.OkWithDetailed(message.LoginResp{
		User:      user.UserName,
		Token:     token,
		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix(),
	}, "登录成功", c)
	return
	//}
}

func (s *BaseApiService) Login(u *model.User) (userInter *model.User, err error) {
	u.Status = 1
	return u, nil

	if nil == model.Db {
		return nil, fmt.Errorf("db not init")
	}

	var user model.User
	err = model.Db.Where("user_name = ?", u.UserName).First(&user).Error
	if err == nil {
		//todo 密码加密 校验
		//if ok := util.BcryptCheck(u.Password, user.Password); !ok {
		//	return nil, errors.New("密码错误")
		//}

		//var SysAuthorityMenus []system.SysAuthorityMenu
		//err =  model.Db.Where("sys_authority_authority_id = ?", user.AuthorityId).Find(&SysAuthorityMenus).Error
		//if err != nil {
		//	return
		//}
		//
		//var MenuIds []string
		//
		//for i := range SysAuthorityMenus {
		//	MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
		//}
		//
		//var am system.SysBaseMenu
		//ferr :=  model.Db.First(&am, "name = ? and id in (?)", user.Authority.DefaultRouter, MenuIds).Error
		//if errors.Is(ferr, gorm.ErrRecordNotFound) {
		//	user.Authority.DefaultRouter = "404"
		//}
	}

	return &user, nil
}
