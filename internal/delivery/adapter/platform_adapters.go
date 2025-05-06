package adapter

import (
	"context"
	"reflect"
	"time"

	"github.com/wxnacy/wz-backend-go/internal/pkg/tenantctx"
)

// MobileResponseAdapter adapts responses for mobile platforms
type MobileResponseAdapter struct {
	*tenantctx.BaseResponseAdapter
}

// NewMobileResponseAdapter creates a new mobile response adapter
func NewMobileResponseAdapter() *MobileResponseAdapter {
	return &MobileResponseAdapter{
		BaseResponseAdapter: tenantctx.NewBaseResponseAdapter(tenantctx.PlatformMobile),
	}
}

// AdaptResponse implements ResponseAdapter.AdaptResponse for mobile platforms
func (a *MobileResponseAdapter) AdaptResponse(ctx context.Context, response interface{}) (interface{}, bool) {
	// Get the response type
	respValue := reflect.ValueOf(response)
	
	// Handle pointer types
	if respValue.Kind() == reflect.Ptr {
		if respValue.IsNil() {
			return response, false
		}
		respValue = respValue.Elem()
	}
	
	// Only process struct types
	if respValue.Kind() != reflect.Struct {
		return response, false
	}
	
	// For mobile responses, add mobile-specific fields or transform date formats
	wasModified := false
	
	// Check if this is a standard API response structure
	if respValue.Type().Name() == "Response" || respValue.Type().Name() == "APIResponse" {
		// Process standard API response
		if dataField := respValue.FieldByName("Data"); dataField.IsValid() && dataField.CanInterface() {
			dataValue := dataField.Interface()
			
			// Process the data field
			if adaptedData, modified := a.adaptDataForMobile(dataValue); modified {
				// In a real implementation, you'd need to set the modified data back to the response
				// This example just demonstrates the concept
				wasModified = true
			}
		}
	} else {
		// Direct processing of the response object
		wasModified = a.processFields(respValue)
	}
	
	return response, wasModified
}

// adaptDataForMobile adapts data for mobile platforms
func (a *MobileResponseAdapter) adaptDataForMobile(data interface{}) (interface{}, bool) {
	if data == nil {
		return nil, false
	}
	
	dataValue := reflect.ValueOf(data)
	
	// Handle pointer types
	if dataValue.Kind() == reflect.Ptr {
		if dataValue.IsNil() {
			return data, false
		}
		dataValue = dataValue.Elem()
	}
	
	// Process based on type
	if dataValue.Kind() == reflect.Struct {
		return data, a.processFields(dataValue)
	} else if dataValue.Kind() == reflect.Slice {
		// For slices, process each element
		wasModified := false
		for i := 0; i < dataValue.Len(); i++ {
			elemValue := dataValue.Index(i)
			if elemValue.Kind() == reflect.Struct || 
			   (elemValue.Kind() == reflect.Ptr && !elemValue.IsNil() && elemValue.Elem().Kind() == reflect.Struct) {
				
				if a.processFields(elemValue) {
					wasModified = true
				}
			}
		}
		return data, wasModified
	}
	
	return data, false
}

// processFields processes struct fields for mobile platform
func (a *MobileResponseAdapter) processFields(structValue reflect.Value) bool {
	wasModified := false
	
	// Ensure we're working with a struct
	if structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			return false
		}
		structValue = structValue.Elem()
	}
	
	if structValue.Kind() != reflect.Struct {
		return false
	}
	
	// Process fields
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structValue.Type().Field(i)
		
		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}
		
		// Handle time fields - mobile often needs timestamps or formatted dates
		if field.Type() == reflect.TypeOf(time.Time{}) {
			// For mobile, we might want to convert time.Time to unix timestamp
			// In a real implementation, you would set the value back to the field
			// if field.CanSet() { ... }
			wasModified = true
		}
		
		// Handle nested structs recursively
		if field.Kind() == reflect.Struct {
			if a.processFields(field) {
				wasModified = true
			}
		} else if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
			// Handle slices of structs
			for j := 0; j < field.Len(); j++ {
				if a.processFields(field.Index(j)) {
					wasModified = true
				}
			}
		}
	}
	
	return wasModified
}

// UniAppResponseAdapter adapts responses for UniApp platforms
type UniAppResponseAdapter struct {
	*tenantctx.BaseResponseAdapter
}

// NewUniAppResponseAdapter creates a new UniApp response adapter
func NewUniAppResponseAdapter() *UniAppResponseAdapter {
	return &UniAppResponseAdapter{
		BaseResponseAdapter: tenantctx.NewBaseResponseAdapter(tenantctx.PlatformUniApp),
	}
}

// AdaptResponse implements ResponseAdapter.AdaptResponse for UniApp platforms
func (a *UniAppResponseAdapter) AdaptResponse(ctx context.Context, response interface{}) (interface{}, bool) {
	// Get the response type
	respValue := reflect.ValueOf(response)
	
	// Handle pointer types
	if respValue.Kind() == reflect.Ptr {
		if respValue.IsNil() {
			return response, false
		}
		respValue = respValue.Elem()
	}
	
	// Only process struct types
	if respValue.Kind() != reflect.Struct {
		return response, false
	}
	
	// UniApp may require specific field transformations or response structure
	// This is a simplified example that would be customized based on actual requirements
	
	// Check if this response has image URLs that need to be transformed for UniApp
	if imageField := respValue.FieldByName("ImageURL"); imageField.IsValid() && imageField.Kind() == reflect.String {
		// In a real implementation, you might transform image URLs for UniApp compatibility
		// if imageField.CanSet() { ... }
		return response, true
	}
	
	// Check if there are platform-specific fields to add/modify
	if dataField := respValue.FieldByName("Data"); dataField.IsValid() {
		// Process data field for UniApp
		// In a real implementation, this would transform the data as needed
		return response, true
	}
	
	return response, false
}
