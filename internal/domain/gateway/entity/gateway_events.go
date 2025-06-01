package entity

import (
	"time"
	"wz-backend-go/internal/domain/event"
	"wz-backend-go/internal/domain/gateway/valueobject"
)

const (
	// RouteCreatedEvent is fired when a route is created
	RouteCreatedEvent string = "gateway.route.created"
	
	// RouteUpdatedEvent is fired when a route is updated
	RouteUpdatedEvent string = "gateway.route.updated"
	
	// RouteDeletedEvent is fired when a route is deleted
	RouteDeletedEvent string = "gateway.route.deleted"
	
	// RouteActivatedEvent is fired when a route is activated
	RouteActivatedEvent string = "gateway.route.activated"
	
	// RouteDeactivatedEvent is fired when a route is deactivated
	RouteDeactivatedEvent string = "gateway.route.deactivated"
	
	// ServiceCreatedEvent is fired when a service is created
	ServiceCreatedEvent string = "gateway.service.created"
	
	// ServiceUpdatedEvent is fired when a service is updated
	ServiceUpdatedEvent string = "gateway.service.updated"
	
	// ServiceDeletedEvent is fired when a service is deleted
	ServiceDeletedEvent string = "gateway.service.deleted"
	
	// ServiceActivatedEvent is fired when a service is activated
	ServiceActivatedEvent string = "gateway.service.activated"
	
	// ServiceDeactivatedEvent is fired when a service is deactivated
	ServiceDeactivatedEvent string = "gateway.service.deactivated"
	
	// ServiceHealthStatusChangedEvent is fired when a service health status changes
	ServiceHealthStatusChangedEvent string = "gateway.service.health_status_changed"
	
	// RequestRoutedEvent is fired when a request is routed through the gateway
	RequestRoutedEvent string = "gateway.request.routed"
	
	// RequestRateLimitedEvent is fired when a request is rate limited
	RequestRateLimitedEvent string = "gateway.request.rate_limited"
	
	// RequestAuthFailedEvent is fired when authentication fails for a request
	RequestAuthFailedEvent string = "gateway.request.auth_failed"
	
	// RequestErrorEvent is fired when an error occurs during request processing
	RequestErrorEvent string = "gateway.request.error"
)

// RouteCreatedEventData contains data for the RouteCreatedEvent
type RouteCreatedEventData struct {
	RouteID     string `json:"route_id"`
	Path        string `json:"path"`
	ServiceName string `json:"service_name"`
	TargetURL   string `json:"target_url"`
	AuthType    string `json:"auth_type"`
}

// RouteEvent creates a new RouteCreatedEvent
func (r *Route) RouteCreatedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		RouteCreatedEvent,
		r.id.String(),
		RouteCreatedEventData{
			RouteID:     r.id.String(),
			Path:        r.path.String(),
			ServiceName: r.serviceName.String(),
			TargetURL:   r.targetURL.String(),
			AuthType:    r.authType.String(),
		},
		time.Now(),
	)
}

// RouteUpdatedEventData contains data for the RouteUpdatedEvent
type RouteUpdatedEventData struct {
	RouteID     string `json:"route_id"`
	Path        string `json:"path"`
	ServiceName string `json:"service_name"`
	TargetURL   string `json:"target_url"`
	AuthType    string `json:"auth_type"`
}

// RouteUpdatedEvent creates a new RouteUpdatedEvent
func (r *Route) RouteUpdatedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		RouteUpdatedEvent,
		r.id.String(),
		RouteUpdatedEventData{
			RouteID:     r.id.String(),
			Path:        r.path.String(),
			ServiceName: r.serviceName.String(),
			TargetURL:   r.targetURL.String(),
			AuthType:    r.authType.String(),
		},
		time.Now(),
	)
}

// RouteDeletedEventData contains data for the RouteDeletedEvent
type RouteDeletedEventData struct {
	RouteID string `json:"route_id"`
}

// RouteDeletedEvent creates a new RouteDeletedEvent
func (r *Route) RouteDeletedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		RouteDeletedEvent,
		r.id.String(),
		RouteDeletedEventData{
			RouteID: r.id.String(),
		},
		time.Now(),
	)
}

// RouteStatusChangeEventData contains data for route status change events
type RouteStatusChangeEventData struct {
	RouteID string `json:"route_id"`
	Status  bool   `json:"status"`
}

// RouteActivatedEvent creates a new RouteActivatedEvent
func (r *Route) RouteActivatedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		RouteActivatedEvent,
		r.id.String(),
		RouteStatusChangeEventData{
			RouteID: r.id.String(),
			Status:  true,
		},
		time.Now(),
	)
}

// RouteDeactivatedEvent creates a new RouteDeactivatedEvent
func (r *Route) RouteDeactivatedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		RouteDeactivatedEvent,
		r.id.String(),
		RouteStatusChangeEventData{
			RouteID: r.id.String(),
			Status:  false,
		},
		time.Now(),
	)
}

// ServiceCreatedEventData contains data for the ServiceCreatedEvent
type ServiceCreatedEventData struct {
	ServiceName string `json:"service_name"`
	BaseURL     string `json:"base_url"`
	HealthURL   string `json:"health_url,omitempty"`
	DefaultAuth string `json:"default_auth"`
}

// ServiceCreatedEvent creates a new ServiceCreatedEvent
func (s *Service) ServiceCreatedEvent() event.DomainEvent {
	data := ServiceCreatedEventData{
		ServiceName: s.name.String(),
		BaseURL:     s.baseURL.String(),
		DefaultAuth: s.defaultAuth.String(),
	}
	
	if s.healthURL.URL() != nil {
		data.HealthURL = s.healthURL.String()
	}
	
	return event.NewDomainEvent(
		ServiceCreatedEvent,
		s.name.String(),
		data,
		time.Now(),
	)
}

// ServiceUpdatedEventData contains data for the ServiceUpdatedEvent
type ServiceUpdatedEventData struct {
	ServiceName string `json:"service_name"`
	BaseURL     string `json:"base_url"`
	HealthURL   string `json:"health_url,omitempty"`
	DefaultAuth string `json:"default_auth"`
}

// ServiceUpdatedEvent creates a new ServiceUpdatedEvent
func (s *Service) ServiceUpdatedEvent() event.DomainEvent {
	data := ServiceUpdatedEventData{
		ServiceName: s.name.String(),
		BaseURL:     s.baseURL.String(),
		DefaultAuth: s.defaultAuth.String(),
	}
	
	if s.healthURL.URL() != nil {
		data.HealthURL = s.healthURL.String()
	}
	
	return event.NewDomainEvent(
		ServiceUpdatedEvent,
		s.name.String(),
		data,
		time.Now(),
	)
}

// ServiceDeletedEventData contains data for the ServiceDeletedEvent
type ServiceDeletedEventData struct {
	ServiceName string `json:"service_name"`
}

// ServiceDeletedEvent creates a new ServiceDeletedEvent
func (s *Service) ServiceDeletedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		ServiceDeletedEvent,
		s.name.String(),
		ServiceDeletedEventData{
			ServiceName: s.name.String(),
		},
		time.Now(),
	)
}

// ServiceStatusChangeEventData contains data for service status change events
type ServiceStatusChangeEventData struct {
	ServiceName string `json:"service_name"`
	Status      bool   `json:"status"`
}

// ServiceActivatedEvent creates a new ServiceActivatedEvent
func (s *Service) ServiceActivatedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		ServiceActivatedEvent,
		s.name.String(),
		ServiceStatusChangeEventData{
			ServiceName: s.name.String(),
			Status:      true,
		},
		time.Now(),
	)
}

// ServiceDeactivatedEvent creates a new ServiceDeactivatedEvent
func (s *Service) ServiceDeactivatedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		ServiceDeactivatedEvent,
		s.name.String(),
		ServiceStatusChangeEventData{
			ServiceName: s.name.String(),
			Status:      false,
		},
		time.Now(),
	)
}

// ServiceHealthStatusEventData contains data for the ServiceHealthStatusChangedEvent
type ServiceHealthStatusEventData struct {
	ServiceName  string `json:"service_name"`
	IsHealthy    bool   `json:"is_healthy"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// ServiceHealthStatusChangedEvent creates a new ServiceHealthStatusChangedEvent
func (s *Service) ServiceHealthStatusChangedEvent() event.DomainEvent {
	return event.NewDomainEvent(
		ServiceHealthStatusChangedEvent,
		s.name.String(),
		ServiceHealthStatusEventData{
			ServiceName:  s.name.String(),
			IsHealthy:    s.isHealthy,
			ErrorMessage: s.errorMessage,
		},
		time.Now(),
	)
}

// RequestRoutedEventData contains data for the RequestRoutedEvent
type RequestRoutedEventData struct {
	RequestID   string `json:"request_id"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	RouteID     string `json:"route_id"`
	ServiceName string `json:"service_name"`
	Duration    int64  `json:"duration_ms"`
	StatusCode  int    `json:"status_code"`
}

// RequestRateLimitedEventData contains data for the RequestRateLimitedEvent
type RequestRateLimitedEventData struct {
	RequestID   string `json:"request_id"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	ClientIP    string `json:"client_ip"`
	LimitType   string `json:"limit_type"`
	LimitValue  int    `json:"limit_value"`
	CurrentRate int    `json:"current_rate"`
}

// RequestAuthFailedEventData contains data for the RequestAuthFailedEvent
type RequestAuthFailedEventData struct {
	RequestID  string `json:"request_id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	ClientIP   string `json:"client_ip"`
	AuthType   string `json:"auth_type"`
	FailReason string `json:"fail_reason"`
}

// RequestErrorEventData contains data for the RequestErrorEvent
type RequestErrorEventData struct {
	RequestID   string `json:"request_id"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	ServiceName string `json:"service_name"`
	ErrorCode   string `json:"error_code"`
	ErrorMsg    string `json:"error_message"`
}
