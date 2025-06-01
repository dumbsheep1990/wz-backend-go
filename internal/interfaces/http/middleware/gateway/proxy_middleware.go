package gateway

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wz-backend-go/internal/application/gateway/dto"
)

// ProxyMiddleware handles proxying requests to backend services
type ProxyMiddleware struct {
	client *http.Client
}

// NewProxyMiddleware creates a new ProxyMiddleware
func NewProxyMiddleware() *ProxyMiddleware {
	return &ProxyMiddleware{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// HandleProxy returns a handler function to proxy requests
func (m *ProxyMiddleware) HandleProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get target URL from context (set by route finder middleware)
		targetURLInterface, exists := c.Get("target_url")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "PROXY_ERROR",
				Message: "No target URL found",
			})
			return
		}

		targetURLStr, ok := targetURLInterface.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "PROXY_ERROR",
				Message: "Invalid target URL format",
			})
			return
		}

		// Parse target URL
		targetURL, err := url.Parse(targetURLStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "PROXY_ERROR",
				Message: "Failed to parse target URL",
				Details: err.Error(),
			})
			return
		}

		// Create proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Save original director for error handling
		originalDirector := proxy.Director

		// Customize director
		proxy.Director = func(req *http.Request) {
			// Call original director
			originalDirector(req)

			// Update request URL
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = targetURL.Path

			// Update query params if any
			targetQuery := targetURL.RawQuery
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}

			// Add X-Forwarded headers
			if clientIP := c.ClientIP(); clientIP != "" {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
			req.Header.Set("X-Forwarded-Host", c.Request.Host)
			req.Header.Set("X-Forwarded-Proto", c.Request.URL.Scheme)

			// Add gateway headers
			req.Header.Set("X-Gateway-Version", "1.0.0")
		}

		// Set error handler
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			c.AbortWithStatusJSON(http.StatusBadGateway, dto.ErrorResponse{
				Code:    "PROXY_ERROR",
				Message: "Failed to proxy request to backend service",
				Details: err.Error(),
			})
		}

		// Set ModifyResponse handler to capture and modify response
		proxy.ModifyResponse = func(resp *http.Response) error {
			// Add X-Gateway headers to response
			resp.Header.Set("X-Gateway-Time", time.Now().Format(time.RFC3339))
			
			// Read the response body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}
			resp.Body.Close()
			
			// Restore the response body
			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			
			// Check for 5xx errors from backend
			if resp.StatusCode >= 500 {
				// Log the error (in a real implementation this would go to a logger)
				fmt.Printf("Backend service error: %s %d %s\n", targetURLStr, resp.StatusCode, string(bodyBytes))
			}
			
			return nil
		}

		// Record start time for metrics
		startTime := time.Now()

		// Create a custom response writer to capture the status code
		responseWriter := &responseWriter{ResponseWriter: c.Writer}
		c.Writer = responseWriter

		// Serve the request
		proxy.ServeHTTP(c.Writer, c.Request)

		// Record metrics after serving the request
		duration := time.Since(startTime)
		statusCode := responseWriter.statusCode
		
		// In a real implementation, these metrics would be sent to a metrics system
		fmt.Printf("Request to %s completed in %v with status code %d\n", 
			targetURLStr, duration, statusCode)
		
		// Stop further processing since we've already sent the response
		c.Abort()
	}
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	gin.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and passes it to the underlying ResponseWriter
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write writes the response body and sets the status code to 200 if not already set
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}
