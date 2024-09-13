package organization

import "context"

type Usecase interface {
	IsResponsible(ctx context.Context, userId int) error
	Update(ctx context.Context, userId int) error
}

type UsecaseLayer struct {
}

func NewUsecaseLayer() *UsecaseLayer {
	return &UsecaseLayer{}
}
