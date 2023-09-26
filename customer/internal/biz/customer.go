package biz

import (
	"context"
	"customer/api/verifyCode"
	"database/sql"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type CustomerRepo interface {
	CachePhoneCode(ctx context.Context, telephone, phoneCode string, liftTime int64) error
	GetVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error)
}

// Customer Model
type Customer struct {
	gorm.Model
	CustomerWork
	CustomerToken
}

type CustomerWork struct {
	ID            int    `gorm:"primaryKey" json:"id"`
	Telephone     string `gorm:"type:varchar(15);unique" json:"telephone,omitempty"`
	TelephoneCode string `gorm:"type:varchar(15)" json:"telephone_code"`
	Name          string `gorm:"type:varchar(150)" json:"name,omitempty"`
	Email         string `gorm:"type:varchar(55)" json:"email,omitempty"`
	Wechat        string `gorm:"type:varchar(150)" json:"wechat,omitempty"`
	//CityId    uint32
}

type CustomerToken struct {
	Token          string       `gorm:"type:varchar(255)" json:"token,omitempty"`
	TokenCreatedAt sql.NullTime `json:"token_created_at"`
}

// CustomerUsecase GreeterUsecase is a Customer usecase.
type CustomerUsecase struct {
	repo CustomerRepo
	log  *log.Helper
}

// NewCustomerUsecase NewGreeterUsecase new a Customer usecase.
func NewCustomerUsecase(repo CustomerRepo, logger log.Logger) *CustomerUsecase {
	return &CustomerUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (u *CustomerUsecase) SetVerifyCode(ctx context.Context, phone, code string, expireTime int64) error {
	return u.repo.CachePhoneCode(ctx, phone, code, expireTime)
}

func (u *CustomerUsecase) GetVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error) {
	return u.repo.GetVerifyCode(ctx, length, t)
}
