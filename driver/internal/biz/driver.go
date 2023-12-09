package biz

import (
	"context"
	"database/sql"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

const (
	DriverStatusOut    = "out"
	DriverStatusIn     = "in"
	DriverStatusListen = "listen"
	DriverStatusStop   = "stop"
)

type Driver struct {
	gorm.Model
	DriverWork
}

// DriverWork 司机的业务模型
type DriverWork struct {
	Telephone     string       `gorm:"column:telephone;type:varchar(16);" json:"telephone"`
	Token         string       `gorm:"column:token;type:varchar(2047);" json:"token"`
	Name          string       `gorm:"column:name;type:varchar(255);" json:"name"`
	Status        string       `gorm:"column:status;type:enum('out','in','listen','stop');" json:"status"`
	IdNumber      string       `gorm:"column:id_number;type:varchar(18);" json:"id_number"`
	IdImageA      string       `gorm:"column:id_image_a;type:varchar(255);" json:"id_image_a"`
	LicenceImageA string       `gorm:"column:licence_image_a;type:varchar(255);" json:"licence_image_a"`
	LicenceImageB string       `gorm:"column:licence_image_b;type:varchar(255);" json:"licence_image_b"`
	DistinctCode  string       `gorm:"column:distinct_code;type:varchar(16);" json:"distinct_code"`
	AuditAt       sql.NullTime `gorm:"column:audit_at;type:datetime;" json:"audit_at"`
	TelephoneBak  string       `gorm:"column:telephone_bak;type:varchar(16);" json:"telephone_bak"`
}

type VerifyCodeRepo interface {
	GetVerifyCode(ctx context.Context, phone string, service string, lifeTime int64) (string, error)
	ValidateVerifyCode(ctx context.Context, phone string, service string) error
}

type DriverBiz struct {
	vc VerifyCodeRepo
}

func NewDriverBiz(vc VerifyCodeRepo) *DriverBiz {
	return &DriverBiz{vc: vc}
}

func (d *DriverBiz) GetVerifyCode(ctx context.Context, phone string, expireTime int64) (string, error) {
	ctx, span := otel.Tracer("biz:driver").Start(ctx, "GetVerifyCode")
	defer span.End()
	return d.vc.GetVerifyCode(ctx, phone, "driver:", expireTime)
}
