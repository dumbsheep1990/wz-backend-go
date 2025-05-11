package services

import (
	"testing"

	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/service"
	"wz-backend-go/tests/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// 测试友情链接服务实现
func TestLinkService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的仓储
	mockRepo := mocks.NewMockLinkRepository(ctrl)
	linkService := service.NewLinkService(mockRepo)

	// 测试数据
	testLink := &domain.Link{
		Name:        "测试链接",
		URL:         "https://example.com",
		Logo:        "https://example.com/logo.png",
		Sort:        100,
		Status:      1,
		Description: "这是一个测试链接",
		TenantID:    1,
	}

	// 测试创建
	t.Run("CreateLink", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any()).Return(int64(1), nil)

		id, err := linkService.CreateLink(testLink)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	// 测试获取详情
	t.Run("GetLinkById", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(int64(1)).Return(testLink, nil)

		link, err := linkService.GetLinkById(1)
		assert.NoError(t, err)
		assert.Equal(t, testLink.Name, link.Name)
		assert.Equal(t, testLink.URL, link.URL)
	})

	// 测试获取列表
	t.Run("ListLinks", func(t *testing.T) {
		expectedLinks := []*domain.Link{testLink}
		expectedTotal := int64(1)
		query := map[string]interface{}{"name": "测试"}

		mockRepo.EXPECT().List(1, 10, query).Return(expectedLinks, expectedTotal, nil)

		links, total, err := linkService.ListLinks(1, 10, query)
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
		assert.Equal(t, expectedLinks, links)
	})

	// 测试更新
	t.Run("UpdateLink", func(t *testing.T) {
		testLink.ID = 1
		testLink.Name = "更新后的链接"

		mockRepo.EXPECT().GetByID(int64(1)).Return(testLink, nil)
		mockRepo.EXPECT().Update(testLink).Return(nil)

		err := linkService.UpdateLink(testLink)
		assert.NoError(t, err)
	})

	// 测试删除
	t.Run("DeleteLink", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(int64(1)).Return(testLink, nil)
		mockRepo.EXPECT().Delete(int64(1)).Return(nil)

		err := linkService.DeleteLink(1)
		assert.NoError(t, err)
	})
}
