package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// 服务配置
type ServiceConfig struct {
	Name    string   // 服务名称
	CmdPath string   // 命令路径
	Args    []string // 启动参数
}

// 微服务列表
var services = []ServiceConfig{
	{
		Name:    "API网关服务",
		CmdPath: "./cmd/gateway/main.go",
		Args:    []string{"run", "./cmd/gateway/main.go"},
	},
	{
		Name:    "用户服务",
		CmdPath: "./cmd/rpc/user/main.go",
		Args:    []string{"run", "./cmd/rpc/user/main.go"},
	},
	{
		Name:    "内容服务",
		CmdPath: "./cmd/rpc/content/main.go",
		Args:    []string{"run", "./cmd/rpc/content/main.go"},
	},
	{
		Name:    "文件服务",
		CmdPath: "./cmd/rpc/file/main.go",
		Args:    []string{"run", "./cmd/rpc/file/main.go"},
	},
	{
		Name:    "交互服务",
		CmdPath: "./cmd/rpc/interaction/main.go",
		Args:    []string{"run", "./cmd/rpc/interaction/main.go"},
	},
	// 可根据需要添加更多服务
}

func main() {
	fmt.Println("=== WZ Backend 开发环境启动器 ===")
	fmt.Println("此工具用于在开发环境中同时启动多个微服务")
	fmt.Println("注意: 此入口仅供开发使用，生产环境应独立部署各微服务")
	fmt.Println("==============================")

	// 命令行参数
	runAll := flag.Bool("all", false, "运行所有服务")
	runGateway := flag.Bool("gateway", false, "运行API网关服务")
	runUser := flag.Bool("user", false, "运行用户服务")
	runContent := flag.Bool("content", false, "运行内容服务")
	flag.Parse()

	var servicesToRun []ServiceConfig

	// 根据命令行参数选择要运行的服务
	if *runAll {
		servicesToRun = services
	} else {
		if *runGateway {
			servicesToRun = append(servicesToRun, services[0])
		}
		if *runUser {
			servicesToRun = append(servicesToRun, services[1])
		}
		if *runContent {
			servicesToRun = append(servicesToRun, services[2])
		}
		// 如果没有指定任何服务，则默认运行网关和用户服务
		if len(servicesToRun) == 0 {
			servicesToRun = append(servicesToRun, services[0], services[1])
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// 启动选定的服务
	for _, svc := range servicesToRun {
		wg.Add(1)
		go func(svc ServiceConfig) {
			defer wg.Done()
			runService(ctx, svc)
		}(svc)
	}

	// 处理终止信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Printf("接收到信号: %v，准备优雅关闭所有服务...\n", sig)
		cancel()
	}()

	wg.Wait()
	fmt.Println("所有服务已关闭")
}

// 运行单个服务
func runService(ctx context.Context, svc ServiceConfig) {
	fmt.Printf("正在启动 %s...\n", svc.Name)

	// 检查文件是否存在
	if _, err := os.Stat(svc.CmdPath); os.IsNotExist(err) {
		log.Printf("错误: 服务文件不存在: %s\n", svc.CmdPath)
		return
	}

	// 准备命令
	cmd := exec.Command("go", svc.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Printf("启动 %s 失败: %v\n", svc.Name, err)
		return
	}

	fmt.Printf("%s 已启动 [PID: %d]\n", svc.Name, cmd.Process.Pid)

	// 等待上下文取消或命令完成
	go func() {
		<-ctx.Done()
		fmt.Printf("正在停止 %s...\n", svc.Name)

		// 给进程一个优雅关闭的机会
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("发送SIGTERM到 %s 失败: %v\n", svc.Name, err)
			// 如果无法发送SIGTERM，则强制杀死进程
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("强制终止 %s 失败: %v\n", svc.Name, err)
			}
		}

		// 等待进程结束，但最多等待5秒
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			if err != nil {
				log.Printf("%s 已退出，状态: %v\n", svc.Name, err)
			} else {
				fmt.Printf("%s 已优雅退出\n", svc.Name)
			}
		case <-time.After(5 * time.Second):
			// 超时后强制杀死进程
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("强制终止 %s 失败: %v\n", svc.Name, err)
			} else {
				log.Printf("%s 已被强制终止\n", svc.Name)
			}
		}
	}()

	// 等待命令完成
	err := cmd.Wait()
	if ctx.Err() == nil { // 如果不是因为上下文取消
		if err != nil {
			log.Printf("%s 异常退出: %v\n", svc.Name, err)
		} else {
			log.Printf("%s 已正常退出\n", svc.Name)
		}
	}
}
