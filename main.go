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
		CmdPath: "./services/gateway-service/main.go",
		Args:    []string{"run", "./services/gateway-service/main.go"},
	},
	{
		Name:    "用户服务",
		CmdPath: "./cmd/user/main.go",
		Args:    []string{"run", "./cmd/user/main.go"},
	},
	{
		Name:    "内容服务",
		CmdPath: "./cmd/content/main.go",
		Args:    []string{"run", "./cmd/content/main.go"},
	},
	{
		Name:    "AI服务",
		CmdPath: "./cmd/ai/main.go",
		Args:    []string{"run", "./cmd/ai/main.go"},
	},
	{
		Name:    "学习服务",
		CmdPath: "./cmd/learn/main.go",
		Args:    []string{"run", "./cmd/learn/main.go"},
	},
	{
		Name:    "商务服务",
		CmdPath: "./cmd/commerce/main.go",
		Args:    []string{"run", "./cmd/commerce/main.go"},
	},
	{
		Name:    "交易服务",
		CmdPath: "./cmd/trade/main.go",
		Args:    []string{"run", "./cmd/trade/main.go"},
	},
	{
		Name:    "导航服务",
		CmdPath: "./cmd/navigation/main.go",
		Args:    []string{"run", "./cmd/navigation/main.go"},
	},
	{
		Name:    "后台管理服务",
		CmdPath: "./cmd/admin/main.go",
		Args:    []string{"run", "./cmd/admin/main.go"},
	},
	{
		Name:    "页面服务",
		CmdPath: "./cmd/page/main.go",
		Args:    []string{"run", "./cmd/page/main.go"},
	},
	{
		Name:    "站点服务",
		CmdPath: "./cmd/site/main.go",
		Args:    []string{"run", "./cmd/site/main.go"},
	},
	{
		Name:    "组件服务",
		CmdPath: "./cmd/component/main.go",
		Args:    []string{"run", "./cmd/component/main.go"},
	},
	{
		Name:    "渲染服务",
		CmdPath: "./cmd/render-service/main.go",
		Args:    []string{"run", "./cmd/render-service/main.go"},
	},
	{
		Name:    "主题服务",
		CmdPath: "./cmd/theme/main.go",
		Args:    []string{"run", "./cmd/theme/main.go"},
	},
	// RPC服务（可选启动）
	{
		Name:    "内容RPC服务",
		CmdPath: "./cmd/rpc/content/main.go",
		Args:    []string{"run", "./cmd/rpc/content/main.go"},
	},
	{
		Name:    "AI RPC服务",
		CmdPath: "./cmd/rpc/ai/main.go",
		Args:    []string{"run", "./cmd/rpc/ai/main.go"},
	},
	{
		Name:    "交易RPC服务",
		CmdPath: "./cmd/rpc/trade/main.go",
		Args:    []string{"run", "./cmd/rpc/trade/main.go"},
	},
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
	runAI := flag.Bool("ai", false, "运行AI服务")
	runLearn := flag.Bool("learn", false, "运行学习服务")
	runCommerce := flag.Bool("commerce", false, "运行商务服务")
	runTrade := flag.Bool("trade", false, "运行交易服务")
	runNavigation := flag.Bool("navigation", false, "运行导航服务")
	runAdmin := flag.Bool("admin", false, "运行后台管理服务")
	runPage := flag.Bool("page", false, "运行页面服务")
	runSite := flag.Bool("site", false, "运行站点服务")
	runComponent := flag.Bool("component", false, "运行组件服务")
	runRender := flag.Bool("render", false, "运行渲染服务")
	runTheme := flag.Bool("theme", false, "运行主题服务")
	runRPC := flag.Bool("rpc", false, "运行所有RPC服务")
	flag.Parse()

	var servicesToRun []ServiceConfig

	// 根据命令行参数选择要运行的服务
	if *runAll {
		servicesToRun = services[:14] // 不包含RPC服务
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
		if *runAI {
			servicesToRun = append(servicesToRun, services[3])
		}
		if *runLearn {
			servicesToRun = append(servicesToRun, services[4])
		}
		if *runCommerce {
			servicesToRun = append(servicesToRun, services[5])
		}
		if *runTrade {
			servicesToRun = append(servicesToRun, services[6])
		}
		if *runNavigation {
			servicesToRun = append(servicesToRun, services[7])
		}
		if *runAdmin {
			servicesToRun = append(servicesToRun, services[8])
		}
		if *runPage {
			servicesToRun = append(servicesToRun, services[9])
		}
		if *runSite {
			servicesToRun = append(servicesToRun, services[10])
		}
		if *runComponent {
			servicesToRun = append(servicesToRun, services[11])
		}
		if *runRender {
			servicesToRun = append(servicesToRun, services[12])
		}
		if *runTheme {
			servicesToRun = append(servicesToRun, services[13])
		}
		if *runRPC {
			// 添加所有RPC服务
			servicesToRun = append(servicesToRun, services[14:]...)
		}
		// 如果没有指定任何服务，则默认运行核心服务
		if len(servicesToRun) == 0 {
			fmt.Println("未指定服务，运行核心服务（网关、用户、内容、AI）...")
			servicesToRun = append(servicesToRun, services[0], services[1], services[2], services[3])
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
