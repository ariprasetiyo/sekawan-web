package repository

import (
	"context"
	baseModel "sekawan-web/app/main/server"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	queryGetMerchantIdByUserIdChatat string = "select merchant_id  from merchant_users_chatat where user_id_chatat = ?"
)

func NewDatabase(db *gorm.DB, base *baseModel.PostgreSQLClientRepository) Database {
	return &databaseImpl{db: db, base: base}
}

type databaseImpl struct {
	db   *gorm.DB
	base *baseModel.PostgreSQLClientRepository
}

func (di databaseImpl) GetCount(ctx context.Context, userId string) string {
	var merchantId string
	err := di.db.WithContext(ctx).Raw(queryGetMerchantIdByUserIdChatat, userId).Scan(&merchantId)
	if err != nil {
		logrus.Errorln("error", err.Error)
	}
	return merchantId
}
