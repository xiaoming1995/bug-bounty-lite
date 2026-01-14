package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/config"
	"bug-bounty-lite/pkg/database"

	"gorm.io/gorm"
)

// 颜色定义
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

var db *gorm.DB

func main() {
	// 命令行参数
	listPending := flag.Bool("list", false, "列出所有待审核的文章")
	listPublished := flag.Bool("published", false, "列出所有已发布的文章")
	approveID := flag.Int("approve", 0, "审核通过指定ID的文章")
	rejectID := flag.Int("reject", 0, "驳回指定ID的文章")
	rejectReason := flag.String("reason", "", "驳回原因")
	featuredID := flag.Int("featured", 0, "设置指定ID的文章为精选")
	unfeaturedID := flag.Int("unfeatured", 0, "取消指定ID的文章精选")
	interactive := flag.Bool("i", false, "交互式审核模式")
	flag.Parse()

	// 加载配置
	cfg := config.LoadConfig()

	// 连接数据库
	db = database.InitDB(cfg)

	printBanner()

	// 根据参数执行不同操作
	switch {
	case *interactive:
		interactiveMode()
	case *listPending:
		listPendingArticles()
	case *listPublished:
		listPublishedArticles()
	case *approveID > 0:
		approveArticle(uint(*approveID))
	case *rejectID > 0:
		if *rejectReason == "" {
			fmt.Printf("%s❌ 驳回文章时必须提供原因 (-reason)%s\n", colorRed, colorReset)
			os.Exit(1)
		}
		rejectArticle(uint(*rejectID), *rejectReason)
	case *featuredID > 0:
		setFeatured(uint(*featuredID), true)
	case *unfeaturedID > 0:
		setFeatured(uint(*unfeaturedID), false)
	default:
		printHelp()
	}
}

func printBanner() {
	fmt.Printf("%s%s", colorCyan, colorBold)
	fmt.Println("╔═══════════════════════════════════════════╗")
	fmt.Println("║       📝 文章审核管理工具 v1.0            ║")
	fmt.Println("╚═══════════════════════════════════════════╝")
	fmt.Printf("%s\n", colorReset)
}

func printHelp() {
	fmt.Println("用法:")
	fmt.Printf("  %s-list%s            列出所有待审核的文章\n", colorGreen, colorReset)
	fmt.Printf("  %s-published%s       列出所有已发布的文章\n", colorGreen, colorReset)
	fmt.Printf("  %s-approve <ID>%s    审核通过指定ID的文章\n", colorGreen, colorReset)
	fmt.Printf("  %s-reject <ID> -reason \"原因\"%s  驳回指定ID的文章\n", colorGreen, colorReset)
	fmt.Printf("  %s-featured <ID>%s   设为精选\n", colorGreen, colorReset)
	fmt.Printf("  %s-unfeatured <ID>%s 取消精选\n", colorGreen, colorReset)
	fmt.Printf("  %s-i%s               交互式审核模式\n", colorGreen, colorReset)
	fmt.Println()
	fmt.Println("Makefile 命令示例:")
	fmt.Printf("  make review-list                        # 查看待审核列表\n")
	fmt.Printf("  make review-approve ID=5                # 通过ID=5的文章\n")
	fmt.Printf("  make review-reject ID=5 REASON=\"原因\"   # 驳回ID=5的文章\n")
	fmt.Printf("  make review-featured ID=5               # 设为精选\n")
	fmt.Printf("  make review-unfeatured ID=5             # 取消精选\n")
	fmt.Printf("  make review-interactive                 # 交互式模式\n")
}

func listPendingArticles() {
	var articles []domain.Article
	if err := db.Where("status = ?", "pending").Order("created_at DESC").Find(&articles).Error; err != nil {
		fmt.Printf("%s❌ 查询失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	if len(articles) == 0 {
		fmt.Printf("%s✅ 暂无待审核的文章%s\n", colorGreen, colorReset)
		return
	}

	fmt.Printf("\n%s📋 待审核文章列表 (共 %d 篇)%s\n", colorBold, len(articles), colorReset)
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("%-6s %-40s %-12s %s\n", "ID", "标题", "作者ID", "提交时间")
	fmt.Println(strings.Repeat("─", 80))

	for _, a := range articles {
		title := truncate(a.Title, 38)
		fmt.Printf("%-6d %-40s %-12d %s\n",
			a.ID, title, a.AuthorID, a.CreatedAt.Format("2006-01-02 15:04"))
	}
	fmt.Println(strings.Repeat("─", 80))
}

func listPublishedArticles() {
	var articles []domain.Article
	if err := db.Where("status = ?", "approved").Order("created_at DESC").Find(&articles).Error; err != nil {
		fmt.Printf("%s❌ 查询失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	if len(articles) == 0 {
		fmt.Printf("%s暂无已发布的文章%s\n", colorYellow, colorReset)
		return
	}

	fmt.Printf("\n%s📚 已发布文章列表 (共 %d 篇)%s\n", colorBold, len(articles), colorReset)
	fmt.Println(strings.Repeat("─", 90))
	fmt.Printf("%-6s %-35s %-8s %-8s %-12s %s\n", "ID", "标题", "精选", "浏览量", "作者ID", "发布时间")
	fmt.Println(strings.Repeat("─", 90))

	for _, a := range articles {
		title := truncate(a.Title, 33)
		featuredMark := "  "
		if a.IsFeatured {
			featuredMark = "⭐"
		}
		fmt.Printf("%-6d %-35s %-8s %-8d %-12d %s\n",
			a.ID, title, featuredMark, a.Views, a.AuthorID, a.CreatedAt.Format("2006-01-02 15:04"))
	}
	fmt.Println(strings.Repeat("─", 90))
}

func approveArticle(id uint) {
	var article domain.Article
	if err := db.First(&article, id).Error; err != nil {
		fmt.Printf("%s❌ 文章不存在 (ID: %d)%s\n", colorRed, id, colorReset)
		return
	}

	if article.Status != "pending" {
		fmt.Printf("%s⚠️  文章状态为 [%s]，无需审核%s\n", colorYellow, article.Status, colorReset)
		return
	}

	article.Status = "approved"
	article.RejectReason = ""
	article.UpdatedAt = time.Now()

	if err := db.Save(&article).Error; err != nil {
		fmt.Printf("%s❌ 审核失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	fmt.Printf("%s✅ 文章审核通过!%s\n", colorGreen, colorReset)
	fmt.Printf("   ID: %d\n", article.ID)
	fmt.Printf("   标题: %s\n", article.Title)
}

func rejectArticle(id uint, reason string) {
	var article domain.Article
	if err := db.First(&article, id).Error; err != nil {
		fmt.Printf("%s❌ 文章不存在 (ID: %d)%s\n", colorRed, id, colorReset)
		return
	}

	if article.Status != "pending" {
		fmt.Printf("%s⚠️  文章状态为 [%s]，无需审核%s\n", colorYellow, article.Status, colorReset)
		return
	}

	article.Status = "rejected"
	article.RejectReason = reason
	article.UpdatedAt = time.Now()

	if err := db.Save(&article).Error; err != nil {
		fmt.Printf("%s❌ 驳回失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	fmt.Printf("%s❌ 文章已驳回!%s\n", colorRed, colorReset)
	fmt.Printf("   ID: %d\n", article.ID)
	fmt.Printf("   标题: %s\n", article.Title)
	fmt.Printf("   原因: %s\n", reason)
}

func interactiveMode() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\n%s请选择操作:%s\n", colorBold, colorReset)
		fmt.Println("  1. 查看待审核文章列表")
		fmt.Println("  2. 审核通过文章")
		fmt.Println("  3. 驳回文章")
		fmt.Println("  4. 退出")
		fmt.Print("\n请输入选项 (1-4): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			listPendingArticles()
		case "2":
			fmt.Print("请输入要通过的文章ID: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
				approveArticle(uint(id))
			} else {
				fmt.Printf("%s无效的ID%s\n", colorRed, colorReset)
			}
		case "3":
			fmt.Print("请输入要驳回的文章ID: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				fmt.Printf("%s无效的ID%s\n", colorRed, colorReset)
				continue
			}
			fmt.Print("请输入驳回原因: ")
			reason, _ := reader.ReadString('\n')
			reason = strings.TrimSpace(reason)
			if reason == "" {
				fmt.Printf("%s必须提供驳回原因%s\n", colorRed, colorReset)
				continue
			}
			rejectArticle(uint(id), reason)
		case "4":
			fmt.Printf("%s👋 再见!%s\n", colorCyan, colorReset)
			return
		default:
			fmt.Printf("%s无效的选项%s\n", colorRed, colorReset)
		}
	}
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}

func setFeatured(id uint, featured bool) {
	var article domain.Article
	if err := db.First(&article, id).Error; err != nil {
		fmt.Printf("%s❌ 文章不存在 (ID: %d)%s\n", colorRed, id, colorReset)
		return
	}

	if article.Status != "approved" {
		fmt.Printf("%s⚠️  只有已发布的文章才能设为精选%s\n", colorYellow, colorReset)
		return
	}

	article.IsFeatured = featured
	article.UpdatedAt = time.Now()

	if err := db.Save(&article).Error; err != nil {
		fmt.Printf("%s❌ 操作失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	if featured {
		fmt.Printf("%s⭐ 已设为精选!%s\n", colorGreen, colorReset)
	} else {
		fmt.Printf("%s✓ 已取消精选%s\n", colorGreen, colorReset)
	}
	fmt.Printf("   ID: %d\n", article.ID)
	fmt.Printf("   标题: %s\n", article.Title)
}
