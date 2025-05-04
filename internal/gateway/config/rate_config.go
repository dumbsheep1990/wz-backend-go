package config

// RateConfig 定义限流器配置
type RateConfig struct {
	// 限流策略，可以是 "ip", "user", "tenant", "path"
	Strategy string `yaml:"strategy" json:"strategy"`
	// 时间窗口大小，单位秒
	IntervalSecs int `yaml:"interval_secs" json:"interval_secs"`
	// 在时间窗口内允许的最大请求数
	MaxRequests int `yaml:"max_requests" json:"max_requests"`
}

// ServiceConfig 定义服务配置
type ServiceConfig struct {
	// 熔断器配置
	CircuitBreaker struct {
		// 触发熔断的最大连续失败请求数
		MaxRequests uint32 `yaml:"max_requests" json:"max_requests"`
		// 熔断器重置的时间间隔（秒）
		IntervalSecs int `yaml:"interval_secs" json:"interval_secs"`
		// 熔断器超时时间（秒）
		TimeoutSecs int `yaml:"timeout_secs" json:"timeout_secs"`
	} `yaml:"circuit_breaker" json:"circuit_breaker"`
}
