package model

import (
	"time"
)

// SearchLog 搜索记录模型
type SearchLog struct {
	ID         int64      `json:"id" db:"id"`
	UserID     int64      `json:"user_id" db:"user_id"`
	Keyword    string     `json:"keyword" db:"keyword"`
	Category   string     `json:"category" db:"category"`
	Tags       string     `json:"tags" db:"tags"`
	Location   string     `json:"location" db:"location"`
	Longitude  float64    `json:"longitude" db:"longitude"`
	Latitude   float64    `json:"latitude" db:"latitude"`
	ResultNum  int        `json:"result_num" db:"result_num"`
	DeviceID   string     `json:"device_id" db:"device_id"`
	IP         string     `json:"ip" db:"ip"`
	UserAgent  string     `json:"user_agent" db:"user_agent"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// SearchSuggestion 搜索建议模型
type SearchSuggestion struct {
	ID         int64      `json:"id" db:"id"`
	Keyword    string     `json:"keyword" db:"keyword"`
	Category   string     `json:"category" db:"category"`
	Frequency  int        `json:"frequency" db:"frequency"`
	IsEnabled  bool       `json:"is_enabled" db:"is_enabled"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// HotSearch 热搜词模型
type HotSearch struct {
	ID         int64      `json:"id" db:"id"`
	Keyword    string     `json:"keyword" db:"keyword"`
	Category   string     `json:"category" db:"category"`
	Count      int        `json:"count" db:"count"`
	Trend      int8       `json:"trend" db:"trend"` // 1:上升，0:持平，-1:下降
	IsPromoted bool       `json:"is_promoted" db:"is_promoted"`
	SortOrder  int        `json:"sort_order" db:"sort_order"`
	OperatorID int64      `json:"operator_id" db:"operator_id"`
	StartTime  *time.Time `json:"start_time" db:"start_time"`
	EndTime    *time.Time `json:"end_time" db:"end_time"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// SearchStatistic 搜索统计模型
type SearchStatistic struct {
	ID          int64     `json:"id" db:"id"`
	StatDate    time.Time `json:"stat_date" db:"stat_date"`
	Keyword     string    `json:"keyword" db:"keyword"`
	Category    string    `json:"category" db:"category"`
	SearchCount int       `json:"search_count" db:"search_count"`
	UserCount   int       `json:"user_count" db:"user_count"`
	ResultCount int       `json:"result_count" db:"result_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SearchSynonym 搜索同义词模型
type SearchSynonym struct {
	ID         int64     `json:"id" db:"id"`
	Word       string    `json:"word" db:"word"`
	Synonyms   string    `json:"synonyms" db:"synonyms"` // 逗号分隔的同义词列表
	IsEnabled  bool      `json:"is_enabled" db:"is_enabled"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// SearchStopword 搜索停用词模型
type SearchStopword struct {
	ID         int64     `json:"id" db:"id"`
	Word       string    `json:"word" db:"word"`
	IsEnabled  bool      `json:"is_enabled" db:"is_enabled"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// SearchIndexConfig 搜索索引配置模型
type SearchIndexConfig struct {
	ID           int64     `json:"id" db:"id"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	IndexName    string    `json:"index_name" db:"index_name"`
	Fields       string    `json:"fields" db:"fields"`             // JSON格式
	BoostFields  string    `json:"boost_fields" db:"boost_fields"` // JSON格式
	FilterFields string    `json:"filter_fields" db:"filter_fields"` // JSON格式
	IsEnabled    bool      `json:"is_enabled" db:"is_enabled"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SearchSortRule 搜索排序规则模型
type SearchSortRule struct {
	ID         int64     `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	SortFields string    `json:"sort_fields" db:"sort_fields"` // JSON格式
	IsDefault  bool      `json:"is_default" db:"is_default"`
	IsEnabled  bool      `json:"is_enabled" db:"is_enabled"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// SearchParams 搜索参数模型
type SearchParams struct {
	Keyword     string   `json:"keyword"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Location    string   `json:"location"`
	Longitude   float64  `json:"longitude"`
	Latitude    float64  `json:"latitude"`
	Radius      float64  `json:"radius"` // 搜索半径，单位：公里
	SortBy      string   `json:"sort_by"` // 排序字段
	SortOrder   string   `json:"sort_order"` // 排序方向：asc, desc
	Page        int      `json:"page"`
	PageSize    int      `json:"page_size"`
	Filters     map[string]interface{} `json:"filters"` // 其他过滤条件
	UserID      int64    `json:"user_id"`
	DeviceID    string   `json:"device_id"`
	IP          string   `json:"ip"`
	UserAgent   string   `json:"user_agent"`
}

// SearchResult 搜索结果模型
type SearchResult struct {
	Total       int         `json:"total"`
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
	TotalPages  int         `json:"total_pages"`
	HasNext     bool        `json:"has_next"`
	Keyword     string      `json:"keyword"`
	Category    string      `json:"category"`
	Results     interface{} `json:"results"`      // 搜索结果列表，根据具体搜索内容类型变化
	Suggestions []string    `json:"suggestions"`  // 搜索建议
	Duration    int64       `json:"duration"`     // 搜索耗时，单位：毫秒
}

// SearchService 搜索服务接口
type SearchService interface {
	// 通用搜索接口
	Search(params SearchParams) (*SearchResult, error)
	
	// 搜索建议
	GetSuggestions(keyword string, category string, limit int) ([]string, error)
	
	// 热搜词
	GetHotSearches(category string, limit int) ([]HotSearch, error)
	AddHotSearch(hotSearch *HotSearch) error
	UpdateHotSearch(hotSearch *HotSearch) error
	DeleteHotSearch(id int64) error
	
	// 搜索日志
	LogSearch(searchLog *SearchLog) error
	GetSearchStatistics(startDate, endDate time.Time, category string) ([]SearchStatistic, error)
	
	// 搜索配置
	GetSynonyms(word string) ([]string, error)
	UpdateSynonyms(word string, synonyms []string) error
	GetStopwords() ([]string, error)
	UpdateStopword(word string, isEnabled bool) error
	
	// 索引管理
	UpdateIndexConfig(config *SearchIndexConfig) error
	GetIndexConfig(resourceType string) (*SearchIndexConfig, error)
	
	// 排序规则管理
	GetSortRules() ([]SearchSortRule, error)
	UpdateSortRule(rule *SearchSortRule) error
}

// SearchRepository 搜索仓储接口
type SearchRepository interface {
	// 搜索日志
	SaveSearchLog(log *SearchLog) error
	GetSearchLogs(userID int64, keyword string, category string, startTime, endTime time.Time, page, pageSize int) ([]*SearchLog, int, error)
	
	// 搜索建议
	GetSuggestions(keyword string, category string, limit int) ([]*SearchSuggestion, error)
	UpdateSuggestionFrequency(keyword string, category string) error
	SaveSuggestion(suggestion *SearchSuggestion) error
	
	// 热搜词
	GetHotSearches(category string, limit int) ([]*HotSearch, error)
	SaveHotSearch(hotSearch *HotSearch) error
	UpdateHotSearch(hotSearch *HotSearch) error
	DeleteHotSearch(id int64) error
	
	// 搜索统计
	SaveSearchStatistic(statistic *SearchStatistic) error
	UpdateSearchStatistic(statistic *SearchStatistic) error
	GetSearchStatistics(startDate, endDate time.Time, category string) ([]*SearchStatistic, error)
	
	// 搜索配置
	GetSynonyms(word string) (*SearchSynonym, error)
	SaveSynonym(synonym *SearchSynonym) error
	UpdateSynonym(synonym *SearchSynonym) error
	GetAllEnabledSynonyms() ([]*SearchSynonym, error)
	
	GetStopwords() ([]*SearchStopword, error)
	SaveStopword(stopword *SearchStopword) error
	UpdateStopword(stopword *SearchStopword) error
	
	// 索引管理
	GetIndexConfig(resourceType string) (*SearchIndexConfig, error)
	SaveIndexConfig(config *SearchIndexConfig) error
	UpdateIndexConfig(config *SearchIndexConfig) error
	
	// 排序规则
	GetSortRules() ([]*SearchSortRule, error)
	GetDefaultSortRule() (*SearchSortRule, error)
	SaveSortRule(rule *SearchSortRule) error
	UpdateSortRule(rule *SearchSortRule) error
}
