package biz

import (
	"context"
	"valuation/pkg/errorx"

	"github.com/go-kratos/kratos/v2/log"
)

type Good struct {
	Id     uint64  ` json:"id,omitempty"`
	Name   string  ` json:"name,omitempty"`
	Type   string  ` json:"type,omitempty"`
	Unit   string  ` json:"unit,omitempty"`
	Price  float32 ` json:"price,omitempty"`
	Tariff float32 ` json:"tariff,omitempty"`
	Alias  string  ` json:"alias,omitempty"`
}
type GoodRepo interface {
	Save(context.Context, *Good) (*Good, error)
	Update(context.Context, *Good) (*Good, error)
	Delete(context.Context, *Good) (*Good, error)
	FindByID(context.Context, uint64) (*Good, error)
	ExistByName(context.Context, string) (bool, error)
	ListAll(context.Context) ([]*Good, error)
	GetGoodsByWords(context.Context, string) ([]*Good, error)
}

// GreeterUsecase is a Greeter usecase.
type GoodUsecase struct {
	repo GoodRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGoodUsecase(repo GoodRepo, logger log.Logger) *GoodUsecase {
	return &GoodUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGood  creates a Good, and returns the new Good.
func (uc *GoodUsecase) CreateGood(ctx context.Context, g *Good) (*Good, error) {
	uc.log.WithContext(ctx).Infof("CreateGood: %v ", g)
	exist, err := uc.repo.ExistByName(ctx, g.Name)
	errorx.Dangerous(err)
	if exist {
		errorx.Bomb(201, "good is already exist")
	}
	return uc.repo.Save(ctx, g)
}

func (uc *GoodUsecase) GetGoodsByWords(ctx context.Context, words string) ([]*Good, error) {
	uc.log.WithContext(ctx).Infof("words: %v\n", words)
	return uc.repo.GetGoodsByWords(ctx, words)
}
