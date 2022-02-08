package locality

import (
	"Sellers/internal/domain"
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("locality not found")
)

type Service interface{
	GetAll(ctx context.Context) ([]domain.Locality, error)
	Save(ctx context.Context, lo domain.Locality) (domain.Locality, error)
	Exists(ctx context.Context, zipC string) bool

}

type service struct{
	repo Repository
}


func NewService(l Repository) Service {
	return &service {repo: l,
	}
}


func (l *service) GetAll(ctx context.Context) ([]domain.Locality, error){
	ls, err := l.repo.GetAll(ctx)
	if err != nil{
		return nil, err
	}

	return ls, nil
}


func(l *service) Save(ctx context.Context, lo domain.Locality) (domain.Locality, error){

	exist := l.repo.Exists(ctx, lo.ZipCode)
	if exist{
		return domain.Locality{}, ErrNotFound
	}

	p, err := l.repo.Save(ctx, lo)

	if err != nil{
		return domain.Locality{}, err
	}

	lo.ID = p

	return lo, nil

}


func(l *service) Exists(ctx context.Context, zipC string) bool{
	return l.repo.Exists(ctx, zipC)
}