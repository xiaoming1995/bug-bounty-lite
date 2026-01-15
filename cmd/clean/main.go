package main

import (
	"bufio"
	"bug-bounty-lite/internal/seeder"
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 命令行参数
	helpFlag := flag.Bool("help", false, "Show help message")
	allFlag := flag.Bool("all", false, "Clean all test data")
	usersFlag := flag.Bool("users", false, "Clean only user data (keeps admin)")
	projectsFlag := flag.Bool("projects", false, "Clean only project data")
	reportsFlag := flag.Bool("reports", false, "Clean only report data")
	articlesFlag := flag.Bool("articles", false, "Clean only article data")
	confirmFlag := flag.Bool("confirm", false, "Skip confirmation prompt")
	statsFlag := flag.Bool("stats", false, "Show data statistics only")
	flag.Parse()

	// 显示帮助
	if *helpFlag {
		printHelp()
		return
	}

	fmt.Println("Bug Bounty Lite - Test Data Cleaner")
	fmt.Println("====================================")

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db := database.InitDB(cfg)

	// 创建清理器
	cleaner := seeder.NewCleaner(db)

	// 只显示统计信息
	if *statsFlag {
		cleaner.PrintStatistics()
		return
	}

	// 检查是否指定了清理类型
	if !*allFlag && !*usersFlag && !*projectsFlag && !*reportsFlag && !*articlesFlag {
		fmt.Println("[ERROR] Please specify what to clean. Use -help for options.")
		fmt.Println("")
		printHelp()
		os.Exit(1)
	}

	// 确认操作
	if !*confirmFlag {
		fmt.Println("\n[WARNING] This will permanently delete data!")
		fmt.Print("Type 'yes' to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input != "yes" {
			fmt.Println("[ABORT] Operation cancelled.")
			os.Exit(0)
		}
	}

	// 显示清理前的统计
	fmt.Println("\n>>> Before cleaning:")
	cleaner.PrintStatistics()

	// 执行清理
	var err error
	if *allFlag {
		err = cleaner.CleanAll()
	} else {
		if *reportsFlag {
			if e := cleaner.CleanReportComments(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanReports(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
		}
		if *articlesFlag {
			if e := cleaner.CleanArticleComments(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanArticleLikes(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanArticleViews(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanArticles(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
		}
		if *projectsFlag {
			if e := cleaner.CleanProjectAttachments(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanProjectTasks(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanProjectAssignments(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanProjects(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
		}
		if *usersFlag {
			if e := cleaner.CleanUserUpdateLogs(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanUserInfoChanges(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
			if e := cleaner.CleanUsers(); e != nil {
				fmt.Printf("[ERROR] %v\n", e)
			}
		}
	}

	if err != nil {
		fmt.Printf("\n[ERROR] Clean failed: %v\n", err)
		os.Exit(1)
	}

	// 显示清理后的统计
	fmt.Println("\n>>> After cleaning:")
	cleaner.PrintStatistics()

	fmt.Println("\n[OK] Clean operation completed!")
}

func printHelp() {
	fmt.Println("Bug Bounty Lite - Test Data Cleaner")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/clean/main.go [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -all       Clean all test data (keeps admin users and system configs)")
	fmt.Println("  -users     Clean only user data (keeps admin users)")
	fmt.Println("  -projects  Clean only project data")
	fmt.Println("  -reports   Clean only report data")
	fmt.Println("  -articles  Clean only article data")
	fmt.Println("  -stats     Show data statistics only (no cleaning)")
	fmt.Println("  -confirm   Skip confirmation prompt (for automation)")
	fmt.Println("  -help      Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/clean/main.go -stats           # View current data stats")
	fmt.Println("  go run cmd/clean/main.go -all             # Clean all data (interactive)")
	fmt.Println("  go run cmd/clean/main.go -all -confirm    # Clean all data (no prompt)")
	fmt.Println("  go run cmd/clean/main.go -users -confirm  # Clean only users")
	fmt.Println("")
	fmt.Println("Or use batch script:")
	fmt.Println("  run.bat clean          # Interactive clean all")
	fmt.Println("  run.bat clean-force    # Clean all without prompt")
	fmt.Println("  run.bat clean-stats    # View statistics")
}
