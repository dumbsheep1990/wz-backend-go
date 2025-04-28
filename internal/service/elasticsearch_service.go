package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// ElasticsearchService Elasticsearch搜索服务接口
type ElasticsearchService interface {
	// 创建或更新文档
	IndexDocument(ctx context.Context, indexName string, documentID string, document interface{}) error
	// 删除文档
	DeleteDocument(ctx context.Context, indexName string, documentID string) error
	// 搜索文档
	SearchDocuments(ctx context.Context, indexName string, query map[string]interface{}, from, size int) (SearchResult, error)
	// 批量索引文档
	BulkIndex(ctx context.Context, indexName string, documents []map[string]interface{}, idField string) error
	// 创建索引
	CreateIndex(ctx context.Context, indexName string, mappings map[string]interface{}) error
	// 检查索引是否存在
	IndexExists(ctx context.Context, indexName string) (bool, error)
}

// SearchResult 搜索结果
type SearchResult struct {
	Total      int                      `json:"total"`
	MaxScore   float64                  `json:"max_score"`
	Hits       []map[string]interface{} `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
	TimeTook   time.Duration            `json:"time_took"`
}

type elasticsearchService struct {
	client *elasticsearch.Client
}

// NewElasticsearchService 创建Elasticsearch服务
func NewElasticsearchService(config elasticsearch.Config) (ElasticsearchService, error) {
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}

	// 测试连接
	_, err = client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %v", err)
	}

	return &elasticsearchService{
		client: client,
	}, nil
}

// IndexDocument 创建或更新文档
func (s *elasticsearchService) IndexDocument(ctx context.Context, indexName string, documentID string, document interface{}) error {
	docJSON, err := json.Marshal(document)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true", // 立即刷新，使文档可搜索
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("indexing error: %s", res.String())
	}

	return nil
}

// DeleteDocument 删除文档
func (s *elasticsearchService) DeleteDocument(ctx context.Context, indexName string, documentID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: documentID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 { // 忽略文档不存在的情况
		return fmt.Errorf("delete error: %s", res.String())
	}

	return nil
}

// SearchDocuments 搜索文档
func (s *elasticsearchService) SearchDocuments(ctx context.Context, indexName string, query map[string]interface{}, from, size int) (SearchResult, error) {
	var result SearchResult

	// 构建查询体
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": query,
		"from":  from,
		"size":  size,
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return result, err
	}

	// 开始计时
	startTime := time.Now()

	// 执行搜索
	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(indexName),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return result, fmt.Errorf("search error: %s", res.String())
	}

	// 解析搜索结果
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return result, err
	}

	// 计算耗时
	result.TimeTook = time.Since(startTime)

	// 提取结果
	if hits, ok := searchResponse["hits"].(map[string]interface{}); ok {
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				result.Total = int(value)
			}
		}

		if maxScore, ok := hits["max_score"].(float64); ok {
			result.MaxScore = maxScore
		}

		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArray {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					result.Hits = append(result.Hits, hitMap)
				}
			}
		}
	}

	// 提取聚合结果
	if aggs, ok := searchResponse["aggregations"].(map[string]interface{}); ok && len(aggs) > 0 {
		result.Aggregations = aggs
	}

	return result, nil
}

// BulkIndex 批量索引文档
func (s *elasticsearchService) BulkIndex(ctx context.Context, indexName string, documents []map[string]interface{}, idField string) error {
	if len(documents) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for _, doc := range documents {
		// 准备metadata
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
			},
		}

		// 如果提供了ID字段，使用文档中的ID
		if idField != "" {
			if id, ok := doc[idField]; ok {
				meta["index"].(map[string]interface{})["_id"] = fmt.Sprintf("%v", id)
			}
		}

		// 添加metadata
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}

		// 添加文档
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return err
		}
	}

	// 执行批量索引
	res, err := s.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		s.client.Bulk.WithContext(ctx),
		s.client.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk indexing error: %s", res.String())
	}

	// 检查响应中是否有错误
	var bulkResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkResponse); err != nil {
		return err
	}

	if hasErrors, ok := bulkResponse["errors"].(bool); ok && hasErrors {
		return errors.New("bulk operation contains errors")
	}

	return nil
}

// CreateIndex 创建索引
func (s *elasticsearchService) CreateIndex(ctx context.Context, indexName string, mappings map[string]interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mappings); err != nil {
		return err
	}

	res, err := s.client.Indices.Create(
		indexName,
		s.client.Indices.Create.WithContext(ctx),
		s.client.Indices.Create.WithBody(&buf),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create index error: %s", res.String())
	}

	return nil
}

// IndexExists 检查索引是否存在
func (s *elasticsearchService) IndexExists(ctx context.Context, indexName string) (bool, error) {
	res, err := s.client.Indices.Exists(
		[]string{indexName},
		s.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}
