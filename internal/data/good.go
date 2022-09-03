package data

import (
	"context"

	"valuation/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type goodRepo struct {
	data *Data
	log  *log.Helper
}

func NewGoodRepo(data *Data, logger log.Logger) biz.GoodRepo {
	return &goodRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// GetGoodsByWords select * from goods where name like '%words%' or alias like '%words%'
func (g *goodRepo) GetGoodsByWords(ctx context.Context, words string) (goods []*biz.Good, err error) {
	res := g.data.db.Where("name LIKE ? or alias LIKE ?", "%"+words+"%", "%"+words+"%").Find(&goods)
	if res.Error != nil {
		return nil, res.Error
	}
	return
}

func (g *goodRepo) Delete(ctx context.Context, good *biz.Good) (*biz.Good, error) {
	//TODO implement me
	panic("implement me")
}

func (g *goodRepo) Save(ctx context.Context, good *biz.Good) (*biz.Good, error) {
	res := g.data.db.Create(good)
	if res.Error != nil {
		return good, res.Error
	}
	return good, nil
}

func (g *goodRepo) Update(ctx context.Context, good *biz.Good) (*biz.Good, error) {
	//TODO implement me
	panic("implement me")
}

func (g *goodRepo) FindByID(ctx context.Context, u uint64) (*biz.Good, error) {
	//TODO implement me
	panic("implement me")
}

func (g *goodRepo) ExistByName(ctx context.Context, s string) (bool, error) {
	var cnt int64
	if err := g.data.db.Table("goods").Where("name = ?", s).Count(&cnt).Error; err != nil {
		return false, err
	}
	if cnt != 0 {
		return true, nil
	}
	return false, nil
}

func (g *goodRepo) ListAll(ctx context.Context) ([]*biz.Good, error) {
	//TODO implement me
	panic("implement me")
}
