package user

import "context"

type Usecase interface {
	GetData(ctx context.Context, userId int) error
	Update(ctx context.Context, userId int) error
}

type UsecaseLayer struct {
}

func NewUsecaseLayer() *UsecaseLayer {
	return &UsecaseLayer{}
}
