package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wz-project/wz-backend-go/services/community-service/models"
	"github.com/wz-project/wz-backend-go/services/community-service/service"
)

// SimilarApplicationHandler 同乡申请处理器
type SimilarApplicationHandler struct {
	Service *service.SimilarApplicationService
}

// NewSimilarApplicationHandler 创建新的同乡申请处理器
func NewSimilarApplicationHandler(service *service.SimilarApplicationService) *SimilarApplicationHandler {
	return &SimilarApplicationHandler{
		Service: service,
	}
}

// CreateApplication 创建申请
func (h *SimilarApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var application models.SimilarApplication
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, "请求体解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 从请求中获取用户ID（假设已经通过认证中间件设置）
	userID := r.Context().Value("user_id")
	if userID != nil {
		application.UserID = userID.(string)
	}

	err = h.Service.CreateApplication(&application)
	if err != nil {
		http.Error(w, "创建申请失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "申请创建成功",
		"data":    application,
	})
}

// GetApplication 获取申请详情
func (h *SimilarApplicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	application, err := h.Service.GetApplicationByID(id)
	if err != nil {
		http.Error(w, "获取申请失败: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(application)
}

// ListApplications 列出申请
func (h *SimilarApplicationHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	// 获取分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 构建过滤条件
	filters := make(map[string]interface{})
	if appType := r.URL.Query().Get("type"); appType != "" {
		filters["application_type"] = appType
	}
	
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}

	applications, total, err := h.Service.ListApplications(page, pageSize, filters)
	if err != nil {
		http.Error(w, "获取申请列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       applications,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// UpdateApplication 更新申请
func (h *SimilarApplicationHandler) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updates map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "请求体解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateApplication(id, updates)
	if err != nil {
		http.Error(w, "更新申请失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "申请更新成功",
	})
}

// UpdateApplicationStatus 更新申请状态
func (h *SimilarApplicationHandler) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var statusUpdate struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, "请求体解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateApplicationStatus(id, statusUpdate.Status)
	if err != nil {
		http.Error(w, "更新状态失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "申请状态更新成功",
	})
}

// DeleteApplication 删除申请
func (h *SimilarApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.Service.DeleteApplication(id)
	if err != nil {
		http.Error(w, "删除申请失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "申请删除成功",
	})
}

// RegisterRoutes 注册路由
func (h *SimilarApplicationHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/similar-applications", h.CreateApplication).Methods("POST")
	router.HandleFunc("/similar-applications", h.ListApplications).Methods("GET")
	router.HandleFunc("/similar-applications/{id}", h.GetApplication).Methods("GET")
	router.HandleFunc("/similar-applications/{id}", h.UpdateApplication).Methods("PUT")
	router.HandleFunc("/similar-applications/{id}/status", h.UpdateApplicationStatus).Methods("PATCH")
	router.HandleFunc("/similar-applications/{id}", h.DeleteApplication).Methods("DELETE")
}
