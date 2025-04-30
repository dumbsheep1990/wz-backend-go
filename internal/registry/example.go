package registry

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServiceExample 演示如何使用服务注册与发现、健康检查和服务实例管理功能
func ServiceExample() {
	// 1. 配置Nacos客户端
	config := &NacosConfig{
		ServerAddr: GetEnvOrDefault("NACOS_SERVER_ADDR", "127.0.0.1"),
		ServerPort: uint64(GetEnvAsInt("NACOS_SERVER_PORT", 8848)),
		Namespace:  GetEnvOrDefault("NACOS_NAMESPACE", "public"),
		Group:      GetEnvOrDefault("NACOS_GROUP", "DEFAULT_GROUP"),
		LogDir:     GetEnvOrDefault("NACOS_LOG_DIR", "logs/nacos"),
		CacheDir:   GetEnvOrDefault("NACOS_CACHE_DIR", "cache/nacos"),
		LogLevel:   GetEnvOrDefault("NACOS_LOG_LEVEL", "info"),
		Username:   GetEnvOrDefault("NACOS_USERNAME", ""),
		Password:   GetEnvOrDefault("NACOS_PASSWORD", ""),
	}

	// 2. 创建服务注册对象
	registry, err := NewNacosRegistry(config)
	if err != nil {
		log.Fatalf("创建Nacos注册中心失败: %v", err)
	}

	// 3. 创建服务实例管理器
	instanceManager := NewNacosInstanceManager(registry)

	// 4. 获取本地IP和端口
	ip, err := GetLocalIP()
	if err != nil {
		log.Fatalf("获取本地IP失败: %v", err)
	}

	port, err := GetFreePort()
	if err != nil {
		log.Fatalf("获取空闲端口失败: %v", err)
	}

	// 5. 设置服务信息
	serviceName := GetEnvOrDefault("SERVICE_NAME", "example-service")
	healthCheckPort := GetEnvAsInt("HEALTH_CHECK_PORT", port+1) // 健康检查使用不同的端口

	// 6. 创建并启动健康检查服务
	healthChecker := NewHealthCheckServer(healthCheckPort, 30*time.Second)
	
	// 注册健康检查
	healthChecker.RegisterCheck("self", func() error {
		// 一个简单的自检
		return nil
	})

	// 注册Nacos健康检查
	healthChecker.RegisterCheck("nacos", NacosHealthCheck(registry))

	// 启动健康检查服务
	ctx, cancel := context.WithCancel(context.Background())
	if err := healthChecker.Start(ctx); err != nil {
		log.Fatalf("启动健康检查服务失败: %v", err)
	}

	// 7. 注册服务实例
	instance := &InstanceInfo{
		ServiceName:     serviceName,
		IP:              ip,
		Port:            port,
		Status:          StatusStarting,
		Metadata:        map[string]string{"version": "1.0.0"},
		StartTime:       time.Now(),
		HealthCheckPath: fmt.Sprintf("http://%s:%d/health", ip, healthCheckPort),
	}

	if err := instanceManager.RegisterInstance(serviceName, instance); err != nil {
		log.Fatalf("注册服务实例失败: %v", err)
	}

	// 8. 启动心跳
	if err := instanceManager.StartHeartbeat(serviceName, instance.InstanceID, 15*time.Second); err != nil {
		log.Printf("启动心跳失败: %v", err)
	}

	// 9. 更新服务状态为UP
	if err := instanceManager.UpdateInstanceStatus(serviceName, instance.InstanceID, StatusUp); err != nil {
		log.Printf("更新服务状态失败: %v", err)
	}

	// 10. 启动示例HTTP服务
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "您好，这是 %s 服务!", serviceName)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP服务启动失败: %v", err)
		}
	}()

	log.Printf("服务已启动: %s (IP: %s, Port: %d, Health: %s)", 
		serviceName, ip, port, instance.HealthCheckPath)

	// 11. 订阅服务变更
	registry.Subscribe(serviceName, func(instances []ServiceInstance) {
		log.Printf("服务 %s 实例列表变更，当前实例数: %d", serviceName, len(instances))
		for i, inst := range instances {
			log.Printf("  实例 #%d: %s:%d (健康: %v)", i+1, inst.IP, inst.Port, inst.Healthy)
		}
	})

	// 12. 服务发现示例（获取其他实例）
	go func() {
		for {
			time.Sleep(30 * time.Second)
			
			instances, err := instanceManager.GetInstances(serviceName)
			if err != nil {
				log.Printf("获取服务实例失败: %v", err)
				continue
			}
			
			log.Printf("发现 %s 服务实例数: %d", serviceName, len(instances))
		}
	}()

	// 13. 设置优雅关闭
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-gracefulShutdown
	log.Println("收到关闭信号，开始优雅关闭...")

	// 14. 更新服务状态为OUT_OF_SERVICE
	if err := instanceManager.UpdateInstanceStatus(serviceName, instance.InstanceID, StatusOutOfService); err != nil {
		log.Printf("更新服务状态失败: %v", err)
	}

	// 15. 停止心跳
	if err := instanceManager.StopHeartbeat(serviceName, instance.InstanceID); err != nil {
		log.Printf("停止心跳失败: %v", err)
	}

	// 16. 取消订阅
	if err := registry.Unsubscribe(serviceName); err != nil {
		log.Printf("取消订阅失败: %v", err)
	}

	// 17. 关闭HTTP服务
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer httpCancel()
	if err := server.Shutdown(httpCtx); err != nil {
		log.Printf("HTTP服务关闭失败: %v", err)
	}

	// 18. 注销服务实例
	if err := instanceManager.DeregisterInstance(serviceName, instance.InstanceID); err != nil {
		log.Printf("注销服务实例失败: %v", err)
	}

	// 19. 停止健康检查服务
	cancel()
	if err := healthChecker.Stop(); err != nil {
		log.Printf("健康检查服务关闭失败: %v", err)
	}

	// 20. 关闭注册中心客户端
	if err := registry.Close(); err != nil {
		log.Printf("注册中心客户端关闭失败: %v", err)
	}

	log.Println("服务已安全关闭")
}
