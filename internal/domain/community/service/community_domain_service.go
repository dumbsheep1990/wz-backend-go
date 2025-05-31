package service

import (
	"errors"
	"strings"

	"wz-backend-go/internal/domain/community/entity"
	"wz-backend-go/internal/domain/community/valueobject"
)

// CommunityDomainService 社区领域服务
type CommunityDomainService struct{}

// NewCommunityDomainService 创建社区领域服务
func NewCommunityDomainService() *CommunityDomainService {
	return &CommunityDomainService{}
}

// CreateCommunity 创建社区（包含业务规则验证）
func (s *CommunityDomainService) CreateCommunity(
	name string,
	description string,
	ownerID string,
	ownerName string,
	tags []string,
	location string,
) (*entity.Community, error) {
	// 生成社区ID
	communityID := valueobject.GenerateCommunityID()
	
	// 创建社区名称值对象
	communityName, err := valueobject.NewCommunityName(name)
	if err != nil {
		return nil, err
	}
	
	// 创建标签值对象
	communityTags, err := valueobject.NewTags(tags)
	if err != nil {
		return nil, err
	}
	
	// 验证业务规则
	if err := s.validateCommunityCreation(communityName, description, ownerID, communityTags, location); err != nil {
		return nil, err
	}
	
	// 创建社区实体
	community, err := entity.NewCommunity(
		communityID,
		communityName,
		description,
		ownerID,
		ownerName,
		communityTags,
		location,
	)
	if err != nil {
		return nil, err
	}
	
	return community, nil
}

// ValidateCommunityUpdate 验证社区更新操作
func (s *CommunityDomainService) ValidateCommunityUpdate(
	community *entity.Community,
	operatorID string,
	updates map[string]interface{},
) error {
	// 检查操作权限
	if !s.canModifyCommunity(community, operatorID) {
		return errors.New("没有权限修改此社区")
	}
	
	// 检查社区状态
	if !community.Status().IsOperational() && !community.Status().IsInactive() {
		return errors.New("当前社区状态不允许修改")
	}
	
	// 验证具体的更新字段
	for field, value := range updates {
		switch field {
		case "name":
			if name, ok := value.(string); ok {
				if _, err := valueobject.NewCommunityName(name); err != nil {
					return err
				}
			}
		case "tags":
			if tags, ok := value.([]string); ok {
				if _, err := valueobject.NewTags(tags); err != nil {
					return err
				}
			}
		case "description":
			if desc, ok := value.(string); ok {
				if len(desc) > 1000 {
					return errors.New("描述长度不能超过1000个字符")
				}
			}
		}
	}
	
	return nil
}

// ValidateCommunityStatusChange 验证社区状态变更
func (s *CommunityDomainService) ValidateCommunityStatusChange(
	community *entity.Community,
	targetStatus valueobject.CommunityStatus,
	operatorID string,
	reason string,
) error {
	// 检查操作权限
	if !s.canChangeStatus(community, operatorID, targetStatus) {
		return errors.New("没有权限修改社区状态")
	}
	
	// 检查状态转换是否合法
	if !community.Status().CanTransitionTo(targetStatus) {
		return errors.New("不允许的状态转换: " + community.Status().String() + " -> " + targetStatus.String())
	}
	
	// 某些状态变更需要提供原因
	if s.requiresReason(targetStatus) && strings.TrimSpace(reason) == "" {
		return errors.New("此状态变更需要提供原因")
	}
	
	return nil
}

// ValidateMemberJoin 验证成员加入
func (s *CommunityDomainService) ValidateMemberJoin(
	community *entity.Community,
	memberID string,
	joinMethod string,
) error {
	// 检查社区是否可以接受新成员
	if !community.CanAcceptMembers() {
		return errors.New("社区当前不接受新成员")
	}
	
	// 检查是否已经是创建者
	if community.IsOwnedBy(memberID) {
		return errors.New("创建者无需重复加入社区")
	}
	
	// 验证加入方式
	validJoinMethods := []string{"invite", "request", "public", "admin"}
	isValidMethod := false
	for _, method := range validJoinMethods {
		if joinMethod == method {
			isValidMethod = true
			break
		}
	}
	if !isValidMethod {
		return errors.New("无效的加入方式")
	}
	
	return nil
}

// ValidateMemberLeave 验证成员离开
func (s *CommunityDomainService) ValidateMemberLeave(
	community *entity.Community,
	memberID string,
	leaveReason string,
) error {
	// 创建者不能离开自己的社区
	if community.IsOwnedBy(memberID) {
		return errors.New("社区创建者不能离开社区，如需删除社区请使用删除功能")
	}
	
	// 检查社区状态
	if community.Status().IsDeleted() {
		return errors.New("无法从已删除的社区中离开")
	}
	
	return nil
}

// RecommendTags 推荐标签（基于社区内容）
func (s *CommunityDomainService) RecommendTags(name, description, location string) []string {
	var recommendations []string
	
	nameKeywords := s.extractKeywords(name)
	descKeywords := s.extractKeywords(description)
	
	// 基于名称推荐
	for _, keyword := range nameKeywords {
		if s.isValidTag(keyword) {
			recommendations = append(recommendations, keyword)
		}
	}
	
	// 基于描述推荐
	for _, keyword := range descKeywords {
		if s.isValidTag(keyword) && !s.contains(recommendations, keyword) {
			recommendations = append(recommendations, keyword)
		}
	}
	
	// 基于地区推荐
	if location != "" {
		locationTags := s.extractLocationTags(location)
		for _, tag := range locationTags {
			if !s.contains(recommendations, tag) {
				recommendations = append(recommendations, tag)
			}
		}
	}
	
	// 限制推荐数量
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}
	
	return recommendations
}

// CheckCommunityHealth 检查社区健康度
func (s *CommunityDomainService) CheckCommunityHealth(community *entity.Community) map[string]interface{} {
	health := map[string]interface{}{
		"overall_score": 0,
		"indicators":    map[string]interface{}{},
		"suggestions":   []string{},
	}
	
	score := 0
	indicators := health["indicators"].(map[string]interface{})
	suggestions := []string{}
	
	// 成员活跃度指标
	memberScore := 0
	if community.MemberCount() >= 10 {
		memberScore = 25
	} else if community.MemberCount() >= 5 {
		memberScore = 15
	} else if community.MemberCount() >= 2 {
		memberScore = 10
	}
	indicators["member_activity"] = memberScore
	score += memberScore
	
	if memberScore < 15 {
		suggestions = append(suggestions, "建议邀请更多成员加入社区")
	}
	
	// 内容丰富度指标
	contentScore := 0
	if community.PostCount() >= 20 {
		contentScore = 25
	} else if community.PostCount() >= 10 {
		contentScore = 15
	} else if community.PostCount() >= 5 {
		contentScore = 10
	}
	indicators["content_richness"] = contentScore
	score += contentScore
	
	if contentScore < 15 {
		suggestions = append(suggestions, "建议发布更多有价值的内容")
	}
	
	// 组织结构指标
	structureScore := 0
	if community.GroupCount() >= 3 {
		structureScore = 25
	} else if community.GroupCount() >= 1 {
		structureScore = 15
	}
	indicators["organization"] = structureScore
	score += structureScore
	
	if structureScore < 15 {
		suggestions = append(suggestions, "建议创建更多主题群组")
	}
	
	// 标签完整性指标
	tagScore := 0
	if community.Tags().Count() >= 5 {
		tagScore = 25
	} else if community.Tags().Count() >= 3 {
		tagScore = 15
	} else if community.Tags().Count() >= 1 {
		tagScore = 10
	}
	indicators["tag_completeness"] = tagScore
	score += tagScore
	
	if tagScore < 15 {
		suggestions = append(suggestions, "建议添加更多相关标签")
	}
	
	health["overall_score"] = score
	health["suggestions"] = suggestions
	
	return health
}

// validateCommunityCreation 验证社区创建的业务规则
func (s *CommunityDomainService) validateCommunityCreation(
	name valueobject.CommunityName,
	description string,
	ownerID string,
	tags valueobject.Tags,
	location string,
) error {
	// 描述长度限制
	if len(description) > 1000 {
		return errors.New("社区描述不能超过1000个字符")
	}
	
	// 标签相关性检查
	if !tags.IsEmpty() {
		keywords := name.GetSearchKeywords()
		hasRelevantTag := false
		for _, tag := range tags.Values() {
			for _, keyword := range keywords {
				if strings.Contains(strings.ToLower(tag), keyword) {
					hasRelevantTag = true
					break
				}
			}
			if hasRelevantTag {
				break
			}
		}
		
		// 如果没有相关标签，给出建议
		if !hasRelevantTag {
			// 这里可以记录日志或者给出警告，不阻止创建
		}
	}
	
	return nil
}

// canModifyCommunity 检查是否可以修改社区
func (s *CommunityDomainService) canModifyCommunity(community *entity.Community, operatorID string) bool {
	// 目前只有创建者可以修改，后续可以扩展为管理员等
	return community.IsOwnedBy(operatorID)
}

// canChangeStatus 检查是否可以更改状态
func (s *CommunityDomainService) canChangeStatus(
	community *entity.Community,
	operatorID string,
	targetStatus valueobject.CommunityStatus,
) bool {
	// 创建者可以进行大部分状态变更
	if community.IsOwnedBy(operatorID) {
		return true
	}
	
	// 系统管理员可以进行所有状态变更（这里简化处理）
	// 实际项目中应该有更复杂的权限系统
	return false
}

// requiresReason 检查状态变更是否需要原因
func (s *CommunityDomainService) requiresReason(status valueobject.CommunityStatus) bool {
	return status == valueobject.CommunityStatusSuspended ||
		status == valueobject.CommunityStatusDeleted ||
		status == valueobject.CommunityStatusArchived
}

// extractKeywords 从文本中提取关键词
func (s *CommunityDomainService) extractKeywords(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	var keywords []string
	
	// 停用词
	stopWords := map[string]bool{
		"的": true, "是": true, "在": true, "有": true, "和": true,
		"与": true, "或": true, "但": true, "及": true, "等": true,
		"the": true, "is": true, "are": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
	}
	
	for _, word := range words {
		if len(word) > 2 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}
	
	return keywords
}

// extractLocationTags 从地区中提取标签
func (s *CommunityDomainService) extractLocationTags(location string) []string {
	var tags []string
	
	// 简单的地区解析
	parts := strings.Split(location, "-")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) > 1 {
			tags = append(tags, part)
		}
	}
	
	return tags
}

// isValidTag 检查是否是有效的标签
func (s *CommunityDomainService) isValidTag(tag string) bool {
	return len(tag) >= 2 && len(tag) <= 10
}

// contains 检查切片是否包含某个元素
func (s *CommunityDomainService) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
} 