package repositories

import (
	"testing"
	"time"

	"wz-backend-go/internal/domain"
	"wz-backend-go/internal/repository/sql"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 测试友情链接仓储实现
func TestLinkRepository(t *testing.T) {
	// 创建数据库连接
	conn := sqlx.NewMysql("root:password@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=true")
	repo := sql.NewLinkRepository(conn)

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
	t.Run("Create", func(t *testing.T) {
		id, err := repo.Create(testLink)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
		testLink.ID = id
	})

	// 测试获取详情
	t.Run("GetByID", func(t *testing.T) {
		link, err := repo.GetByID(testLink.ID)
		assert.NoError(t, err)
		assert.Equal(t, testLink.Name, link.Name)
		assert.Equal(t, testLink.URL, link.URL)
		assert.Equal(t, testLink.TenantID, link.TenantID)
		assert.WithinDuration(t, time.Now(), link.CreatedAt, 1*time.Minute)
	})

	// 测试更新
	t.Run("Update", func(t *testing.T) {
		testLink.Name = "更新后的链接"
		testLink.Sort = 200
		err := repo.Update(testLink)
		assert.NoError(t, err)

		// 验证更新结果
		link, err := repo.GetByID(testLink.ID)
		assert.NoError(t, err)
		assert.Equal(t, "更新后的链接", link.Name)
		assert.Equal(t, 200, link.Sort)
	})

	// 测试列表查询
	t.Run("List", func(t *testing.T) {
		// 创建测试数据
		for i := 0; i < 5; i++ {
			link := &domain.Link{
				Name:        "链接" + string(i+1),
				URL:         "https://example" + string(i+1) + ".com",
				Sort:        i + 1,
				Status:      1,
				Description: "描述" + string(i+1),
				TenantID:    1,
			}
			_, err := repo.Create(link)
			assert.NoError(t, err)
		}

		// 测试分页查询
		links, total, err := repo.List(1, 10, map[string]interface{}{
			"tenant_id": int64(1),
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(5))
		assert.GreaterOrEqual(t, len(links), 5)

		// 测试条件查询
		links, total, err = repo.List(1, 10, map[string]interface{}{
			"name": "更新",
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		found := false
		for _, link := range links {
			if link.ID == testLink.ID {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	// 测试删除
	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(testLink.ID)
		assert.NoError(t, err)

		// 验证删除结果
		_, err = repo.GetByID(testLink.ID)
		assert.Error(t, err)
	})
}
