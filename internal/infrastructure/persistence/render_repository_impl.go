package persistence

import (
	"context"
	"errors"
	"sync"
	"time"
	"wz-backend-go/internal/domain/render/entity"
	"wz-backend-go/internal/domain/render/valueobject"
)

// RenderResultRepositoryImpl 渲染结果仓储实现
type RenderResultRepositoryImpl struct {
	// 使用内存存储，在实际应用中这里应该是数据库或Redis等存储
	store     map[string]*entity.RenderResult
	cacheKeys map[string]string // 缓存键到ID的映射
	groups    map[string][]string // 组到ID的映射
	siteMap   map[string][]string // 站点ID到渲染结果ID的映射
	mutex     sync.RWMutex
}

// NewRenderResultRepository 创建一个新的渲染结果仓储
func NewRenderResultRepository() *RenderResultRepositoryImpl {
	return &RenderResultRepositoryImpl{
		store:     make(map[string]*entity.RenderResult),
		cacheKeys: make(map[string]string),
		groups:    make(map[string][]string),
		siteMap:   make(map[string][]string),
		mutex:     sync.RWMutex{},
	}
}

// Save 保存渲染结果
func (r *RenderResultRepositoryImpl) Save(ctx context.Context, result *entity.RenderResult) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := result.ID().String()
	r.store[id] = result

	// 保存缓存键映射
	if result.CacheStrategy().Enabled() {
		cacheKey := result.CacheStrategy().CacheKey()
		r.cacheKeys[cacheKey] = id

		// 保存组映射
		for _, group := range result.CacheStrategy().CacheGroups() {
			if _, exists := r.groups[group]; !exists {
				r.groups[group] = make([]string, 0)
			}
			r.groups[group] = append(r.groups[group], id)

			// 如果组是站点ID，也保存到站点映射
			if _, exists := r.siteMap[group]; !exists {
				r.siteMap[group] = make([]string, 0)
			}
			r.siteMap[group] = append(r.siteMap[group], id)
		}
	}

	return nil
}

// FindByID 根据ID查找渲染结果
func (r *RenderResultRepositoryImpl) FindByID(ctx context.Context, id valueobject.RenderID) (*entity.RenderResult, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if result, exists := r.store[id.String()]; exists {
		return result, nil
	}

	return nil, errors.New("渲染结果不存在")
}

// FindByCacheKey 根据缓存键查找渲染结果
func (r *RenderResultRepositoryImpl) FindByCacheKey(ctx context.Context, cacheKey string) (*entity.RenderResult, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if id, exists := r.cacheKeys[cacheKey]; exists {
		if result, exists := r.store[id]; exists {
			return result, nil
		}
	}

	return nil, errors.New("缓存不存在")
}

// Delete 删除渲染结果
func (r *RenderResultRepositoryImpl) Delete(ctx context.Context, id valueobject.RenderID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	idStr := id.String()
	result, exists := r.store[idStr]
	if !exists {
		return errors.New("渲染结果不存在")
	}

	// 删除缓存键映射
	if result.CacheStrategy().Enabled() {
		cacheKey := result.CacheStrategy().CacheKey()
		delete(r.cacheKeys, cacheKey)

		// 删除组映射
		for _, group := range result.CacheStrategy().CacheGroups() {
			if ids, exists := r.groups[group]; exists {
				newIDs := make([]string, 0)
				for _, existingID := range ids {
					if existingID != idStr {
						newIDs = append(newIDs, existingID)
					}
				}
				r.groups[group] = newIDs
			}

			// 如果组是站点ID，也从站点映射中删除
			if ids, exists := r.siteMap[group]; exists {
				newIDs := make([]string, 0)
				for _, existingID := range ids {
					if existingID != idStr {
						newIDs = append(newIDs, existingID)
					}
				}
				r.siteMap[group] = newIDs
			}
		}
	}

	// 删除存储
	delete(r.store, idStr)

	return nil
}

// DeleteBySiteID 删除指定站点的所有渲染结果
func (r *RenderResultRepositoryImpl) DeleteBySiteID(ctx context.Context, siteID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if ids, exists := r.siteMap[siteID]; exists {
		for _, id := range ids {
			if result, exists := r.store[id]; exists {
				// 删除缓存键映射
				if result.CacheStrategy().Enabled() {
					cacheKey := result.CacheStrategy().CacheKey()
					delete(r.cacheKeys, cacheKey)
				}

				// 删除存储
				delete(r.store, id)
			}
		}

		// 清空站点映射
		delete(r.siteMap, siteID)
	}

	// 从组映射中删除站点ID
	delete(r.groups, siteID)

	return nil
}

// DeleteByGroup 删除指定组的所有渲染结果
func (r *RenderResultRepositoryImpl) DeleteByGroup(ctx context.Context, group string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if ids, exists := r.groups[group]; exists {
		for _, id := range ids {
			if result, exists := r.store[id]; exists {
				// 删除缓存键映射
				if result.CacheStrategy().Enabled() {
					cacheKey := result.CacheStrategy().CacheKey()
					delete(r.cacheKeys, cacheKey)
				}

				// 删除存储
				delete(r.store, id)
			}
		}

		// 清空组映射
		delete(r.groups, group)
	}

	return nil
}

// DeleteExpired 删除所有过期的渲染结果
func (r *RenderResultRepositoryImpl) DeleteExpired(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	for id, result := range r.store {
		if result.CacheStrategy().Enabled() && now.After(result.ExpiresAt()) {
			// 删除缓存键映射
			cacheKey := result.CacheStrategy().CacheKey()
			delete(r.cacheKeys, cacheKey)

			// 从组映射中删除
			for _, group := range result.CacheStrategy().CacheGroups() {
				if ids, exists := r.groups[group]; exists {
					newIDs := make([]string, 0)
					for _, existingID := range ids {
						if existingID != id {
							newIDs = append(newIDs, existingID)
						}
					}
					r.groups[group] = newIDs
				}

				// 如果组是站点ID，也从站点映射中删除
				if ids, exists := r.siteMap[group]; exists {
					newIDs := make([]string, 0)
					for _, existingID := range ids {
						if existingID != id {
							newIDs = append(newIDs, existingID)
						}
					}
					r.siteMap[group] = newIDs
				}
			}

			// 删除存储
			delete(r.store, id)
		}
	}

	return nil
}

// FindExpiring 查找即将过期的渲染结果
func (r *RenderResultRepositoryImpl) FindExpiring(ctx context.Context, within time.Duration) ([]*entity.RenderResult, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	expiringResults := make([]*entity.RenderResult, 0)
	now := time.Now()
	expiryThreshold := now.Add(within)

	for _, result := range r.store {
		if result.CacheStrategy().Enabled() && result.ExpiresAt().Before(expiryThreshold) && result.ExpiresAt().After(now) {
			expiringResults = append(expiringResults, result)
		}
	}

	return expiringResults, nil
}
