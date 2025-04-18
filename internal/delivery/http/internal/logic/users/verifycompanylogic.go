package users

import (
	"context"
	"errors"

	"wz-backend-go/internal/delivery/http/internal/svc"
	"wz-backend-go/internal/delivery/http/internal/types"
	"wz-backend-go/internal/delivery/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCompanyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyCompanyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCompanyLogic {
	return &VerifyCompanyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// validateCompanyFields 基于公司类型验证必要字段
func (l *VerifyCompanyLogic) validateCompanyFields(req *types.VerifyCompanyReq) error {
	// 检查基本必填字段
	if req.CompanyName == "" || req.ContactPerson == "" || req.ContactPhone == "" {
		return errors.New("公司名称、联系人和电话为必填项")
	}

	// 根据公司类型验证特定字段
	switch req.CompanyType {
	case types.CompanyTypeEnterprise: // 企业
		if req.BusinessLicense == "" {
			return errors.New("企业需提供营业执照")
		}
		if req.OrgCodeCert == "" {
			return errors.New("企业需提供组织机构代码证")
		}
	case types.CompanyTypeGroup: // 集团
		if req.BusinessLicense == "" {
			return errors.New("集团需提供营业执照")
		}
		if req.OrgStructure == "" {
			return errors.New("集团需提供组织架构说明")
		}
	case types.CompanyTypeGovernment: // 政府机构/NGO/协会
		if req.UnifiedSocialCredit == "" {
			return errors.New("政府机构/NGO/协会需提供统一社会信用代码证")
		}
	case types.CompanyTypeResearchInst: // 科研所
		if req.BusinessLicense == "" {
			return errors.New("科研所需提供营业执照")
		}
		if req.ResearchCert == "" {
			return errors.New("科研所需提供科研证明")
		}
	default:
		return errors.New("公司类型不正确")
	}

	return nil
}

func (l *VerifyCompanyLogic) VerifyCompany(req *types.VerifyCompanyReq) (resp *types.VerifyCompanyResp, err error) {
	// 验证字段基于公司类型的逻辑
	if err = l.validateCompanyFields(req); err != nil {
		return nil, err
	}
	
	// 调用RPC服务进行企业认证请求处理
	rpcReq := &user.VerifyCompanyRequest{
		UserID:      l.ctx.Value("user_id").(int64),
		CompanyType: int32(req.CompanyType),
		CompanyName: req.CompanyName,
		ContactPerson: req.ContactPerson,
		ContactPhone: req.ContactPhone,
	}
	
	// 根据公司类型设置相应的字段
	switch req.CompanyType {
	case types.CompanyTypeEnterprise: // 企业
		rpcReq.BusinessLicense = req.BusinessLicense
		rpcReq.OrgCodeCert = req.OrgCodeCert
		rpcReq.AgencyCert = req.AgencyCert
		rpcReq.CommitteeLetter = req.CommitteeLetter
	case types.CompanyTypeGroup: // 集团
		rpcReq.BusinessLicense = req.BusinessLicense
		rpcReq.OrgStructure = req.OrgStructure
		rpcReq.CommitteeLetter = req.CommitteeLetter
	case types.CompanyTypeGovernment: // 政府机构/NGO/协会
		rpcReq.UnifiedSocialCredit = req.UnifiedSocialCredit
		rpcReq.CommitteeLetter = req.CommitteeLetter
	case types.CompanyTypeResearchInst: // 科研所
		rpcReq.BusinessLicense = req.BusinessLicense
		rpcReq.ResearchCert = req.ResearchCert
		rpcReq.CommitteeLetter = req.CommitteeLetter
	}
	
	// 公共字段
	rpcReq.UploadedDocument = req.UploadedDocument
	
	result, err := l.svcCtx.UserRpc.VerifyCompany(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	
	resp = &types.VerifyCompanyResp{
		Success: result.Success,
	}
	
	return
}
