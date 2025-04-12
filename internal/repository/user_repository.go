package repository

import (
	"context"

	"github.com/watanabenaoki/go-clean-arch/internal/domain/model"
)

// UserRepository はユーザー関連のデータアクセスを定義するインターフェース
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
}