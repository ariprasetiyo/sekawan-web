package repository

import "context"

type Database interface {
	GetCount(ctx context.Context, merchantId string) string
}
