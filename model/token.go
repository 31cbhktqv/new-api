package model

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	Id             int            `json:"id"`
	UserId         int            `json:"user_id" gorm:"index"`
	Key            string         `json:"key" gorm:"type:char(48);uniqueIndex"`
	Status         int            `json:"status" gorm:"default:1"`
	Name           string         `json:"name" gorm:"index" validate:"max=30"`
	CreatedTime    int64          `json:"created_time" gorm:"bigint"`
	AccessedTime   int64          `json:"accessed_time" gorm:"bigint"`
	ExpiredTime    int64          `json:"expired_time" gorm:"bigint;default:-1"`
	RemainQuota    int64          `json:"remain_quota" gorm:"bigint;default:0"`
	UnlimitedQuota bool           `json:"unlimited_quota" gorm:"default:false"`
	UsedQuota      int64          `json:"used_quota" gorm:"bigint;default:0"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func GetAllTokens(userId int, startIdx int, num int) ([]*Token, error) {
	var tokens []*Token
	var err error
	err = DB.Where("user_id = ?", userId).Order("id desc").Limit(num).Offset(startIdx).Find(&tokens).Error
	return tokens, err
}

func GetTokenByKey(key string) (*Token, error) {
	token := Token{Key: key}
	var err error
	err = DB.Where(token).First(&token).Error
	return &token, err
}

func GetTokenById(id int) (*Token, error) {
	if id == 0 {
		return nil, ErrRecordNotFound
	}
	token := Token{Id: id}
	var err error
	err = DB.First(&token, "id = ?", id).Error
	return &token, err
}

func (token *Token) Insert() error {
	var err error
	err = DB.Create(token).Error
	return err
}

func (token *Token) Update() error {
	var err error
	err = DB.Save(token).Error
	return err
}

func (token *Token) Delete() error {
	var err error
	err = DB.Delete(token).Error
	return err
}

func (token *Token) IsExpired() bool {
	if token.ExpiredTime == -1 {
		return false
	}
	return token.ExpiredTime < time.Now().Unix()
}

func (token *Token) HasQuota(quota int64) bool {
	if token.UnlimitedQuota {
		return true
	}
	return token.RemainQuota >= quota
}
