package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type casbinRuleRecord struct {
	Id    int64  `gorm:"column:id"`
	Ptype string `gorm:"column:ptype"`
	V0    string `gorm:"column:v0"`
	V1    string `gorm:"column:v1"`
	V2    string `gorm:"column:v2"`
}

type ListCasbinRuleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCasbinRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCasbinRuleLogic {
	return &ListCasbinRuleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCasbinRuleLogic) ListCasbinRule(in *pb.ListCasbinRuleReq) (*pb.ListCasbinRuleResp, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	db := l.svcCtx.Orm.WithContext(l.ctx).Table("casbin_rule")
	if in.RoleCode != "" {
		db = db.Where("v0 = ?", in.RoleCode)
	}
	if in.Path != "" {
		db = db.Where("v1 = ?", in.Path)
	}
	if in.Method != "" {
		db = db.Where("v2 = ?", in.Method)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("查询Casbin策略总数失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	var rules []casbinRuleRecord
	offset := int((page - 1) * pageSize)
	if err := db.Order("id DESC").Limit(int(pageSize)).Offset(offset).Find(&rules).Error; err != nil {
		l.Errorf("查询Casbin策略列表失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*pb.CasbinRuleItem, 0, len(rules))
	for _, rule := range rules {
		list = append(list, &pb.CasbinRuleItem{
			Id:       rule.Id,
			Ptype:    rule.Ptype,
			RoleCode: rule.V0,
			ApiPath:  rule.V1,
			Method:   rule.V2,
		})
	}

	return &pb.ListCasbinRuleResp{
		Total: total,
		List:  list,
	}, nil
}
