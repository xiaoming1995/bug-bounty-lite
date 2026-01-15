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
	listPending := flag.Bool("list", false, "列出所有待审核的漏洞报告")
	listAudited := flag.Bool("audited", false, "列出所有已审核的漏洞报告")
	listAll := flag.Bool("all", false, "列出所有漏洞报告")
	approveID := flag.Int("approve", 0, "审核通过指定ID的报告 (需要 -severity 参数)")
	rejectID := flag.Int("reject", 0, "驳回指定ID的报告")
	severity := flag.String("severity", "", "设置危害等级: Critical, High, Medium, Low, None")
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
		listPendingReports()
	case *listAudited:
		listAuditedReports()
	case *listAll:
		listAllReports()
	case *approveID > 0:
		if *severity == "" {
			fmt.Printf("%s❌ 审核通过时必须提供危害等级 (-severity)%s\n", colorRed, colorReset)
			fmt.Println("可用等级: Critical, High, Medium, Low, None")
			os.Exit(1)
		}
		approveReport(uint(*approveID), *severity)
	case *rejectID > 0:
		rejectReport(uint(*rejectID))
	default:
		printHelp()
	}
}

func printBanner() {
	fmt.Printf("%s%s", colorCyan, colorBold)
	fmt.Println("╔═══════════════════════════════════════════╗")
	fmt.Println("║       🔒 漏洞审核管理工具 v1.0            ║")
	fmt.Println("╚═══════════════════════════════════════════╝")
	fmt.Printf("%s\n", colorReset)
}

func printHelp() {
	fmt.Println("用法:")
	fmt.Printf("  %s-list%s              列出所有待审核的漏洞报告\n", colorGreen, colorReset)
	fmt.Printf("  %s-audited%s           列出所有已审核的漏洞报告\n", colorGreen, colorReset)
	fmt.Printf("  %s-all%s               列出所有漏洞报告\n", colorGreen, colorReset)
	fmt.Printf("  %s-approve <ID> -severity <等级>%s  审核通过指定ID的报告\n", colorGreen, colorReset)
	fmt.Printf("  %s-reject <ID>%s       驳回指定ID的报告\n", colorGreen, colorReset)
	fmt.Printf("  %s-i%s                 交互式审核模式\n", colorGreen, colorReset)
	fmt.Println()
	fmt.Println("危害等级说明:")
	fmt.Printf("  %sCritical%s - 严重  %sHigh%s - 高危  %sMedium%s - 中危  %sLow%s - 低危  %sNone%s - 无危害\n",
		colorRed, colorReset, colorYellow, colorReset, colorCyan, colorReset, colorGreen, colorReset, colorGreen, colorReset)
	fmt.Println()
	fmt.Println("Makefile 命令示例:")
	fmt.Printf("  make vuln-list                              # 查看待审核列表\n")
	fmt.Printf("  make vuln-approve ID=5 SEVERITY=High        # 通过ID=5的报告，评为高危\n")
	fmt.Printf("  make vuln-reject ID=5                       # 驳回ID=5的报告\n")
	fmt.Printf("  make vuln-interactive                       # 交互式模式\n")
}

func listPendingReports() {
	var reports []domain.Report
	if err := db.Where("status = ?", "Pending").Order("created_at DESC").Find(&reports).Error; err != nil {
		fmt.Printf("%s❌ 查询失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	if len(reports) == 0 {
		fmt.Printf("%s✅ 暂无待审核的漏洞报告%s\n", colorGreen, colorReset)
		return
	}

	fmt.Printf("\n%s📋 待审核漏洞报告列表 (共 %d 条)%s\n", colorBold, len(reports), colorReset)
	printReportTable(reports)
}

func listAuditedReports() {
	var reports []domain.Report
	if err := db.Where("status = ?", "Audited").Order("created_at DESC").Find(&reports).Error; err != nil {
		fmt.Printf("%s❌ 查询失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	if len(reports) == 0 {
		fmt.Printf("%s暂无已审核的漏洞报告%s\n", colorYellow, colorReset)
		return
	}

	fmt.Printf("\n%s✅ 已审核漏洞报告列表 (共 %d 条)%s\n", colorBold, len(reports), colorReset)
	printReportTable(reports)
}

func listAllReports() {
	var reports []domain.Report
	if err := db.Order("created_at DESC").Find(&reports).Error; err != nil {
		fmt.Printf("%s❌ 查询失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	if len(reports) == 0 {
		fmt.Printf("%s暂无漏洞报告%s\n", colorYellow, colorReset)
		return
	}

	fmt.Printf("\n%s📊 所有漏洞报告列表 (共 %d 条)%s\n", colorBold, len(reports), colorReset)
	printReportTable(reports)
}

func printReportTable(reports []domain.Report) {
	fmt.Println(strings.Repeat("─", 100))
	fmt.Printf("%-6s %-35s %-10s %-10s %-10s %s\n", "ID", "漏洞名称", "自评等级", "通过等级", "状态", "提交时间")
	fmt.Println(strings.Repeat("─", 100))

	for _, r := range reports {
		name := truncate(r.VulnerabilityName, 33)
		selfSev := getSelfSeverity(r)
		status := getStatusDisplay(r.Status)
		severity := getSeverityDisplay(r.Severity)
		fmt.Printf("%-6d %-35s %-10s %-10s %-10s %s\n",
			r.ID, name, selfSev, severity, status, r.CreatedAt.Time().Format("2006-01-02"))
	}
	fmt.Println(strings.Repeat("─", 100))
}

func getSelfSeverity(r domain.Report) string {
	if r.SelfAssessmentID == nil {
		return "-"
	}
	var config domain.SystemConfig
	if err := db.First(&config, *r.SelfAssessmentID).Error; err != nil {
		return "-"
	}
	return config.ConfigValue
}

func getStatusDisplay(status string) string {
	switch status {
	case "Pending":
		return colorYellow + "待审核" + colorReset
	case "Audited":
		return colorGreen + "已审核" + colorReset
	case "Rejected":
		return colorRed + "已驳回" + colorReset
	default:
		return status
	}
}

func getSeverityDisplay(severity string) string {
	switch severity {
	case "Critical":
		return colorRed + colorBold + "严重" + colorReset
	case "High":
		return colorRed + "高危" + colorReset
	case "Medium":
		return colorYellow + "中危" + colorReset
	case "Low":
		return colorGreen + "低危" + colorReset
	case "None":
		return colorGreen + "无危害" + colorReset
	default:
		return "-"
	}
}

func approveReport(id uint, severity string) {
	// 验证等级
	validSeverities := []string{"Critical", "High", "Medium", "Low", "None"}
	valid := false
	for _, s := range validSeverities {
		if strings.EqualFold(severity, s) {
			severity = s
			valid = true
			break
		}
	}
	if !valid {
		fmt.Printf("%s❌ 无效的危害等级: %s%s\n", colorRed, severity, colorReset)
		fmt.Println("可用等级: Critical, High, Medium, Low, None")
		return
	}

	var report domain.Report
	if err := db.First(&report, id).Error; err != nil {
		fmt.Printf("%s❌ 报告不存在 (ID: %d)%s\n", colorRed, id, colorReset)
		return
	}

	if report.Status != "Pending" {
		fmt.Printf("%s⚠️ 报告已经审核过了 (当前状态: %s)%s\n", colorYellow, report.Status, colorReset)
		return
	}

	report.Status = "Audited"
	report.Severity = severity
	if err := db.Save(&report).Error; err != nil {
		fmt.Printf("%s❌ 审核失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	fmt.Printf("%s✅ 报告审核通过！%s\n", colorGreen, colorReset)
	fmt.Printf("   ID: %d\n", report.ID)
	fmt.Printf("   漏洞名称: %s\n", report.VulnerabilityName)
	fmt.Printf("   危害等级: %s\n", getSeverityDisplay(severity))
}

func rejectReport(id uint) {
	var report domain.Report
	if err := db.First(&report, id).Error; err != nil {
		fmt.Printf("%s❌ 报告不存在 (ID: %d)%s\n", colorRed, id, colorReset)
		return
	}

	if report.Status != "Pending" {
		fmt.Printf("%s⚠️ 报告已经审核过了 (当前状态: %s)%s\n", colorYellow, report.Status, colorReset)
		return
	}

	report.Status = "Rejected"
	if err := db.Save(&report).Error; err != nil {
		fmt.Printf("%s❌ 驳回失败: %v%s\n", colorRed, err, colorReset)
		return
	}

	fmt.Printf("%s✅ 报告已驳回%s\n", colorYellow, colorReset)
	fmt.Printf("   ID: %d\n", report.ID)
	fmt.Printf("   漏洞名称: %s\n", report.VulnerabilityName)
}

func interactiveMode() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// 获取待审核报告
		var reports []domain.Report
		if err := db.Where("status = ?", "Pending").Order("created_at ASC").Find(&reports).Error; err != nil {
			fmt.Printf("%s❌ 查询失败: %v%s\n", colorRed, err, colorReset)
			return
		}

		if len(reports) == 0 {
			fmt.Printf("\n%s✅ 所有报告已审核完毕！%s\n", colorGreen, colorReset)
			return
		}

		report := reports[0]
		printReportDetail(report)

		fmt.Printf("\n%s操作选项:%s\n", colorBold, colorReset)
		fmt.Println("  [1-5] 通过并设置等级 (1=严重 2=高危 3=中危 4=低危 5=无危害)")
		fmt.Println("  [r]   驳回")
		fmt.Println("  [s]   跳过")
		fmt.Println("  [q]   退出")
		fmt.Print("\n请选择操作: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			approveReport(report.ID, "Critical")
		case "2":
			approveReport(report.ID, "High")
		case "3":
			approveReport(report.ID, "Medium")
		case "4":
			approveReport(report.ID, "Low")
		case "5":
			approveReport(report.ID, "None")
		case "r", "R":
			rejectReport(report.ID)
		case "s", "S":
			fmt.Println("已跳过")
		case "q", "Q":
			fmt.Println("退出审核")
			return
		default:
			fmt.Printf("%s无效输入，请重试%s\n", colorRed, colorReset)
		}

		fmt.Println()
		time.Sleep(500 * time.Millisecond)
	}
}

func printReportDetail(r domain.Report) {
	fmt.Printf("\n%s═══════════════════════════════════════════════════════════════%s\n", colorCyan, colorReset)
	fmt.Printf("%s📝 漏洞报告详情 [ID: %d]%s\n", colorBold, r.ID, colorReset)
	fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n", colorCyan, colorReset)

	fmt.Printf("\n%s漏洞名称:%s %s\n", colorBold, colorReset, r.VulnerabilityName)
	fmt.Printf("%s项目ID:%s   %d\n", colorBold, colorReset, r.ProjectID)
	fmt.Printf("%s提交者ID:%s %d\n", colorBold, colorReset, r.AuthorID)
	fmt.Printf("%s自评等级:%s %s\n", colorBold, colorReset, getSelfSeverity(r))
	fmt.Printf("%s提交时间:%s %s\n", colorBold, colorReset, r.CreatedAt.Time().Format("2006-01-02 15:04:05"))

	if r.VulnerabilityURL != "" {
		fmt.Printf("%s漏洞链接:%s %s\n", colorBold, colorReset, r.VulnerabilityURL)
	}

	if r.VulnerabilityImpact != "" {
		fmt.Printf("\n%s危害描述:%s\n%s\n", colorBold, colorReset, truncate(r.VulnerabilityImpact, 200))
	}

	if r.VulnerabilityDetail != "" {
		fmt.Printf("\n%s漏洞详情:%s\n%s\n", colorBold, colorReset, truncate(r.VulnerabilityDetail, 300))
	}
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}

// 将数字转换为危害等级
func numToSeverity(num int) string {
	switch num {
	case 1:
		return "Critical"
	case 2:
		return "High"
	case 3:
		return "Medium"
	case 4:
		return "Low"
	case 5:
		return "None"
	default:
		return ""
	}
}

// 解析用户输入的数字
func parseNum(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}
