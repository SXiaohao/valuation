package service

import (
	"context"

	pb "valuation/api/valuation/v1"
	"valuation/internal/biz"
	"valuation/pkg/convertx"
	"valuation/pkg/errorx"

	"github.com/go-kratos/kratos/v2/log"
)

type GoodService struct {
	uc  *biz.GoodUsecase
	log *log.Helper
	pb.UnimplementedGoodServer
}

func NewGoodService(uc *biz.GoodUsecase, logger log.Logger) *GoodService {
	return &GoodService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *GoodService) CreateGood(ctx context.Context, req *pb.CreateGoodsRequest) (*pb.CreateGoodsReply, error) {
	good, err := s.uc.CreateGood(ctx, &biz.Good{
		Name:   req.Good.Name,
		Type:   req.Good.Type,
		Unit:   req.Good.Unit,
		Price:  req.Good.Price,
		Tariff: req.Good.Tariff,
		Alias:  convertx.Chinese2Spell(req.Good.Name),
	})
	errorx.Dangerous(err)
	return &pb.CreateGoodsReply{Id: good.Id}, nil
}
func (s *GoodService) UpdateGood(ctx context.Context, req *pb.UpdateGoodsRequest) (*pb.UpdateGoodsReply, error) {
	return &pb.UpdateGoodsReply{}, nil
}
func (s *GoodService) DeleteGood(ctx context.Context, req *pb.DeleteGoodsRequest) (*pb.DeleteGoodsReply, error) {
	return &pb.DeleteGoodsReply{}, nil
}
func (s *GoodService) GetGood(ctx context.Context, req *pb.GetGoodsRequest) (*pb.GetGoodsReply, error) {
	return &pb.GetGoodsReply{}, nil
}
func (s *GoodService) ListGoods(ctx context.Context, req *pb.ListGoodsRequest) (*pb.ListGoodsReply, error) {
	return &pb.ListGoodsReply{}, nil
}
func (s *GoodService) ListGoodsByWords(ctx context.Context, req *pb.ListGoodsByWordsRequest) (*pb.ListGoodsReply, error) {
	goods, err := s.uc.GetGoodsByWords(ctx, req.Words)
	errorx.Dangerous(err)
	var goodsInfo []*pb.GoodInfo
	for _, good := range goods {
		goodsInfo = append(goodsInfo, &pb.GoodInfo{
			Id:     good.Id,
			Name:   good.Name,
			Type:   good.Type,
			Unit:   good.Unit,
			Price:  good.Price,
			Tariff: good.Tariff,
			Alias:  good.Alias,
		})
	}
	return &pb.ListGoodsReply{Goods: goodsInfo}, nil
}
