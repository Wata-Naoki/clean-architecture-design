package usecase

import (
	"context"

	"github.com/watanabenaoki/go-clean-arch/internal/domain/model"
	"github.com/watanabenaoki/go-clean-arch/internal/repository"
)


type UserUsecase interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase {
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUsecase) Create(ctx context.Context, user *model.User) error {
	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) Update(ctx context.Context, user *model.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *userUsecase) Delete(ctx context.Context, id int64) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *userUsecase) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	return u.userRepo.List(ctx, limit, offset)
}