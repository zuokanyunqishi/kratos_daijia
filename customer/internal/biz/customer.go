package biz

import (
	"context"
	"customer/api/verifyCode"
	"customer/internal/conf"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"time"
)

type CustomerRepo interface {
	CachePhoneCode(ctx context.Context, telephone, phoneCode string, liftTime int64) error
	MakeVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error)
	GetVerifyCode(ctx context.Context, telephone string) string
	GetCustomerByTelephone(ctx context.Context, telephone string) (Customer, error)
	QuickCreateCustomerByPhone(ctx context.Context, telephone string) (Customer, error)
	UpdateCustomerToken(ctx context.Context, c *Customer) (*Customer, error)
	GetTokenById(ctx context.Context, id int64) (string, error)
	DeleteToken(ctx context.Context, id int64) error
}

// Customer Model
type Customer struct {
	gorm.Model
	CustomerWork
	CustomerToken
}

type CustomerWork struct {
	Telephone     string `gorm:"type:varchar(15);unique" json:"telephone,omitempty"`
	TelephoneCode string `gorm:"type:varchar(15)" json:"telephone_code"`
	Name          string `gorm:"type:varchar(150)" json:"name,omitempty"`
	Email         string `gorm:"type:varchar(55)" json:"email,omitempty"`
	Wechat        string `gorm:"type:varchar(150)" json:"wechat,omitempty"`
	//CityId    uint32
}

type CustomerToken struct {
	Token          string       `gorm:"type:varchar(1000)" json:"token,omitempty"`
	TokenCreatedAt sql.NullTime `json:"token_created_at"`
}

// CustomerUsecase is a Customer usecase.
type CustomerUsecase struct {
	repo    CustomerRepo
	log     *log.Helper
	cnfAuth *conf.Auth
}

// NewCustomerUsecase NewGreeterUsecase new a Customer usecase.
func NewCustomerUsecase(config *conf.Auth, repo CustomerRepo, logger log.Logger) *CustomerUsecase {
	return &CustomerUsecase{cnfAuth: config, repo: repo, log: log.NewHelper(logger)}
}

func (u *CustomerUsecase) CachePhoneCode(ctx context.Context, phone, code string, expireTime int64) error {
	return u.repo.CachePhoneCode(ctx, phone, code, expireTime)
}

func (u *CustomerUsecase) MakeVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error) {
	return u.repo.MakeVerifyCode(ctx, length, t)
}

func (u *CustomerUsecase) GetVerifyCode(ctx context.Context, telephone string) string {
	return u.repo.GetVerifyCode(ctx, telephone)
}

func (u *CustomerUsecase) GetRepo() CustomerRepo {
	return u.repo
}

func (u *CustomerUsecase) GenerateTokenAndSave(ctx context.Context, customer *Customer, tokenLife time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		// 签发机构
		Issuer: "Daijia",
		// 说明
		Subject: "customer-authentication",
		// 签发给谁
		Audience:  []string{"customer", "other"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLife)),
		NotBefore: nil,
		IssuedAt:  nil,
		ID:        fmt.Sprintf("%d", customer.ID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.cnfAuth.ApiKey))
	if err != nil {
		return "", errors.New("biz:customer:GenerateTokenAndSave# signToken fail")
	}
	customer.Token = signedToken
	if _, err := u.GetRepo().UpdateCustomerToken(ctx, customer); err != nil {
		return "", errors.New("biz:customer:GenerateTokenAndSave# updateCustomer fail")
	}

	return signedToken, nil
}
