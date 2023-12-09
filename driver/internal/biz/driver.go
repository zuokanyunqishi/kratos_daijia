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
	Telephone     string         `gorm:"column:telephone;type:varchar(16);uniqueIndex;" json:"telephone"`
	Token         sql.NullString `gorm:"column:token;type:varchar(2047);" json:"token"`
	Name          sql.NullString `gorm:"column:name;type:varchar(255);index;" json:"name"`
	Status        sql.NullString `gorm:"column:status;type:enum('out','in','listen','stop');" json:"status"`
	IdNumber      sql.NullString `gorm:"column:id_number;type:varchar(18);uniqueIndex;" json:"id_number"`
	IdImageA      sql.NullString `gorm:"column:id_image_a;type:varchar(255);" json:"id_image_a"`
	LicenceImageA sql.NullString `gorm:"column:licence_image_a;type:varchar(255);" json:"licence_image_a"`
	LicenceImageB sql.NullString `gorm:"column:licence_image_b;type:varchar(255);" json:"licence_image_b"`
	DistinctCode  sql.NullString `gorm:"column:distinct_code;type:varchar(16);index;" json:"distinct_code"`
	AuditAt       sql.NullTime   `gorm:"column:audit_at;type:datetime;index;" json:"audit_at"`
	TelephoneBak  sql.NullString `gorm:"column:telephone_bak;type:varchar(16);" json:"telephone_bak"`
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
