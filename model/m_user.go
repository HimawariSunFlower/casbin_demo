package model

func init() {
	register(&User{})
}

type User struct {
	BaseModel
	UserName string  `gorm:"size:64;uniqueIndex;default:'';not null;;comment:用户昵称"` // 用户名
	RealName string  `gorm:"size:64;index;default:'';;comment:真实姓名"`                // 真实姓名
	Password string  `gorm:"size:40;default:'';comment:密码"`                         // 密码
	Email    *string `gorm:"size:255;comment:邮箱"`                                   // 邮箱
	Phone    *string `gorm:"size:20;comment:手机号"`                                   // 手机号
	Status   int     `gorm:"index;default:0;comment:状态(1:启用 2:停用)"`                 // 状态(1:启用 2:停用)
	Creator  uint64  `gorm:"comment:创建者"`                                           // 创建者
}

func (User) TableName() string {
	return "user"
}
