package persistence

import (
	"context"
	"errors"
	"sync"
	"wz-backend-go/internal/domain/render/entity"
)

// TemplateRepositoryImpl 模板仓储实现
type TemplateRepositoryImpl struct {
	// 使用内存存储，在实际应用中这里应该是数据库等持久化存储
	store       map[string]*entity.Template
	siteMap     map[string][]string // 站点ID到模板ID的映射
	nameMap     map[string]string   // 名称到ID的映射（同一站点内唯一）
	versionMap  map[string][]string // 版本号到ID的映射
	mutex       sync.RWMutex
}

// NewTemplateRepository 创建一个新的模板仓储
func NewTemplateRepository() *TemplateRepositoryImpl {
	return &TemplateRepositoryImpl{
		store:      make(map[string]*entity.Template),
		siteMap:    make(map[string][]string),
		nameMap:    make(map[string]string),
		versionMap: make(map[string][]string),
		mutex:      sync.RWMutex{},
	}
}

// Save 保存模板
func (r *TemplateRepositoryImpl) Save(ctx context.Context, template *entity.Template) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := template.ID()
	r.store[id] = template

	// 保存站点映射
	siteID := template.SiteID()
	if _, exists := r.siteMap[siteID]; !exists {
		r.siteMap[siteID] = make([]string, 0)
	}
	
	// 检查是否已经存在，避免重复添加
	exists := false
	for _, existingID := range r.siteMap[siteID] {
		if existingID == id {
			exists = true
			break
		}
	}
	
	if !exists {
		r.siteMap[siteID] = append(r.siteMap[siteID], id)
	}

	// 保存名称映射（同一站点内唯一）
	nameKey := siteID + ":" + template.Name()
	r.nameMap[nameKey] = id

	// 保存版本映射
	version := template.Version()
	if _, exists := r.versionMap[version]; !exists {
		r.versionMap[version] = make([]string, 0)
	}
	
	// 检查是否已经存在，避免重复添加
	exists = false
	for _, existingID := range r.versionMap[version] {
		if existingID == id {
			exists = true
			break
		}
	}
	
	if !exists {
		r.versionMap[version] = append(r.versionMap[version], id)
	}

	return nil
}

// FindByID 根据ID查找模板
func (r *TemplateRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Template, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if template, exists := r.store[id]; exists {
		return template, nil
	}

	return nil, errors.New("模板不存在")
}

// FindBySiteID 查找站点的所有模板
func (r *TemplateRepositoryImpl) FindBySiteID(ctx context.Context, siteID string) ([]*entity.Template, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var templates []*entity.Template

	if ids, exists := r.siteMap[siteID]; exists {
		for _, id := range ids {
			if template, exists := r.store[id]; exists {
				templates = append(templates, template)
			}
		}
	}

	if len(templates) == 0 {
		return nil, errors.New("未找到该站点的模板")
	}

	return templates, nil
}

// FindByName 根据名称查找模板（同一站点内唯一）
func (r *TemplateRepositoryImpl) FindByName(ctx context.Context, siteID string, name string) (*entity.Template, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	nameKey := siteID + ":" + name
	if id, exists := r.nameMap[nameKey]; exists {
		if template, exists := r.store[id]; exists {
			return template, nil
		}
	}

	return nil, errors.New("模板不存在")
}

// FindByVersion 查找特定版本的所有模板
func (r *TemplateRepositoryImpl) FindByVersion(ctx context.Context, version string) ([]*entity.Template, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var templates []*entity.Template

	if ids, exists := r.versionMap[version]; exists {
		for _, id := range ids {
			if template, exists := r.store[id]; exists {
				templates = append(templates, template)
			}
		}
	}

	if len(templates) == 0 {
		return nil, errors.New("未找到该版本的模板")
	}

	return templates, nil
}

// Delete 删除模板
func (r *TemplateRepositoryImpl) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	template, exists := r.store[id]
	if !exists {
		return errors.New("模板不存在")
	}

	// 删除站点映射
	siteID := template.SiteID()
	if ids, exists := r.siteMap[siteID]; exists {
		newIDs := make([]string, 0)
		for _, existingID := range ids {
			if existingID != id {
				newIDs = append(newIDs, existingID)
			}
		}
		r.siteMap[siteID] = newIDs
	}

	// 删除名称映射
	nameKey := siteID + ":" + template.Name()
	delete(r.nameMap, nameKey)

	// 删除版本映射
	version := template.Version()
	if ids, exists := r.versionMap[version]; exists {
		newIDs := make([]string, 0)
		for _, existingID := range ids {
			if existingID != id {
				newIDs = append(newIDs, existingID)
			}
		}
		r.versionMap[version] = newIDs
	}

	// 删除存储
	delete(r.store, id)

	return nil
}

// DeleteBySiteID 删除站点的所有模板
func (r *TemplateRepositoryImpl) DeleteBySiteID(ctx context.Context, siteID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if ids, exists := r.siteMap[siteID]; exists {
		for _, id := range ids {
			if template, exists := r.store[id]; exists {
				// 删除名称映射
				nameKey := siteID + ":" + template.Name()
				delete(r.nameMap, nameKey)

				// 删除版本映射
				version := template.Version()
				if versionIDs, exists := r.versionMap[version]; exists {
					newIDs := make([]string, 0)
					for _, existingID := range versionIDs {
						if existingID != id {
							newIDs = append(newIDs, existingID)
						}
					}
					r.versionMap[version] = newIDs
				}

				// 删除存储
				delete(r.store, id)
			}
		}

		// 清空站点映射
		delete(r.siteMap, siteID)
	}

	return nil
}
