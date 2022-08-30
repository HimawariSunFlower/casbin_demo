package model

type BaseModel struct {
	ID int64 `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`

	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoCreateTime"`
}
