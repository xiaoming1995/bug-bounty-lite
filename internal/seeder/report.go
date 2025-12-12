package seeder

import (
	"bug-bounty-lite/internal/domain"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// ReportSeeder 报告测试数据填充器
type ReportSeeder struct {
	db *gorm.DB
}

func NewReportSeeder(db *gorm.DB) *ReportSeeder {
	return &ReportSeeder{db: db}
}

// Seed 填充报告测试数据（追加模式）
func (s *ReportSeeder) Seed(force bool) error {
	// 1. 获取白帽子用户
	fmt.Println("\n[Step 1] Loading whitehat users...")
	var whitehats []domain.User
	if err := s.db.Where("role = ?", "whitehat").Find(&whitehats).Error; err != nil {
		return fmt.Errorf("failed to get whitehat users: %w", err)
	}
	if len(whitehats) == 0 {
		return fmt.Errorf("no whitehat users found, please run seed-users first")
	}
	fmt.Printf("[INFO] Found %d whitehat users\n", len(whitehats))

	// 2. 获取项目列表
	fmt.Println("\n[Step 2] Loading projects...")
	var projects []domain.Project
	if err := s.db.Where("status = ?", "active").Find(&projects).Error; err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}
	if len(projects) == 0 {
		return fmt.Errorf("no active projects found, please run seed-projects first")
	}
	fmt.Printf("[INFO] Found %d active projects\n", len(projects))

	// 3. 获取漏洞类型配置
	fmt.Println("\n[Step 3] Loading vulnerability types...")
	var vulnTypes []domain.SystemConfig
	if err := s.db.Where("config_type = ? AND status = ?", "vulnerability_type", "active").Find(&vulnTypes).Error; err != nil {
		return fmt.Errorf("failed to get vulnerability types: %w", err)
	}
	if len(vulnTypes) == 0 {
		return fmt.Errorf("no vulnerability types found, please run 'make migrate' first")
	}
	fmt.Printf("[INFO] Found %d vulnerability types\n", len(vulnTypes))

	// 4. 获取危害等级配置
	fmt.Println("\n[Step 4] Loading severity levels...")
	var severityLevels []domain.SystemConfig
	if err := s.db.Where("config_type = ? AND status = ?", "severity_level", "active").Find(&severityLevels).Error; err != nil {
		return fmt.Errorf("failed to get severity levels: %w", err)
	}
	if len(severityLevels) == 0 {
		// 如果没有危害等级配置，也不报错，只是不使用自评ID
		fmt.Println("[WARN] No severity levels found, self_assessment_id will be null")
	} else {
		fmt.Printf("[INFO] Found %d severity levels\n", len(severityLevels))
	}

	// 5. 生成测试报告数据（每次都生成新的）
	fmt.Println("\n[Step 5] Generating new test reports...")
	return s.generateReports(whitehats, projects, vulnTypes, severityLevels, force)
}

// generateReports 生成测试报告
func (s *ReportSeeder) generateReports(whitehats []domain.User, projects []domain.Project, vulnTypes []domain.SystemConfig, severityLevels []domain.SystemConfig, force bool) error {
	rand.Seed(time.Now().UnixNano())

	// 定义测试数据模板
	reportTemplates := []struct {
		VulnerabilityName   string
		VulnerabilityImpact string
		VulnerabilityDetail string
		VulnerabilityURL    string
		Severity            string
		Status              string
		VulnTypeKey         string // 用于匹配特定漏洞类型
	}{
		{
			VulnerabilityName:   "登录页面SQL注入漏洞",
			VulnerabilityImpact: "攻击者可通过SQL注入获取数据库敏感信息，包括用户账号密码，可能导致整个数据库被拖库",
			VulnerabilityDetail: "在登录页面的用户名输入框中输入 ' OR '1'='1 可绕过身份验证。\n\n复现步骤：\n1. 访问登录页面\n2. 用户名输入：admin' OR '1'='1\n3. 密码任意输入\n4. 点击登录，成功绕过验证\n\n建议修复方案：使用参数化查询或预编译语句",
			VulnerabilityURL:    "https://example.com/login",
			Severity:            "Critical",
			Status:              "Pending",
			VulnTypeKey:         "SQL_INJECTION",
		},
		{
			VulnerabilityName:   "用户资料页面存储型XSS",
			VulnerabilityImpact: "攻击者可注入恶意脚本，窃取其他用户的Cookie，可被用于会话劫持",
			VulnerabilityDetail: "在个人简介字段中输入 <script>alert(document.cookie)</script> 后，其他用户访问该页面时会执行恶意脚本。\n\n影响范围：所有访问该用户资料页面的用户\n建议修复：对用户输入进行HTML实体编码",
			VulnerabilityURL:    "https://example.com/profile/edit",
			Severity:            "High",
			Status:              "Triaged",
			VulnTypeKey:         "XSS",
		},
		{
			VulnerabilityName:   "文件上传CSRF漏洞",
			VulnerabilityImpact: "攻击者可诱导用户上传恶意文件到服务器",
			VulnerabilityDetail: "文件上传接口未验证CSRF Token，攻击者可构造恶意HTML页面诱导用户上传文件。\n\nPOC代码已附在附件中。",
			VulnerabilityURL:    "https://example.com/upload",
			Severity:            "Medium",
			Status:              "Resolved",
			VulnTypeKey:         "CSRF",
		},
		{
			VulnerabilityName:   "订单API接口越权访问",
			VulnerabilityImpact: "普通用户可访问其他用户订单信息，存在严重的数据泄露风险",
			VulnerabilityDetail: "通过修改请求中的order_id参数，可以访问其他用户的订单信息。\n\n复现步骤：\n1. 登录普通用户账号A\n2. 访问 /api/orders/1001 获取自己的订单\n3. 修改订单ID为 /api/orders/1002\n4. 成功获取其他用户订单信息\n\n这是一个典型的水平越权漏洞",
			VulnerabilityURL:    "https://example.com/api/orders/1001",
			Severity:            "High",
			Status:              "Closed",
			VulnTypeKey:         "BROKEN_ACCESS_CONTROL",
		},
		{
			VulnerabilityName:   "用户列表API敏感信息泄露",
			VulnerabilityImpact: "攻击者可获取用户手机号、邮箱、身份证号等敏感信息",
			VulnerabilityDetail: "用户列表API返回了用户的完整手机号和邮箱，未做脱敏处理。\n\nAPI响应示例：\n{\n  \"users\": [\n    {\"name\": \"张三\", \"phone\": \"13800138000\", \"email\": \"zhangsan@example.com\"}\n  ]\n}\n\n建议：敏感字段需要进行脱敏处理，如：138****8000",
			VulnerabilityURL:    "https://example.com/api/users?page=1",
			Severity:            "Medium",
			Status:              "Pending",
			VulnTypeKey:         "SENSITIVE_DATA_EXPOSURE",
		},
		{
			VulnerabilityName:   "头像上传任意文件上传漏洞",
			VulnerabilityImpact: "攻击者可上传WebShell获取服务器控制权限",
			VulnerabilityDetail: "头像上传功能仅在前端验证文件扩展名，通过Burp Suite拦截修改请求后可上传.php文件。\n\n复现步骤：\n1. 准备恶意PHP文件并修改扩展名为.jpg\n2. 使用Burp拦截上传请求\n3. 修改文件名为 shell.php\n4. 上传成功并可访问执行\n\n建议：后端需要验证文件类型、扩展名，最好使用白名单机制",
			VulnerabilityURL:    "https://example.com/avatar/upload",
			Severity:            "Critical",
			Status:              "Triaged",
			VulnTypeKey:         "FILE_UPLOAD",
		},
		{
			VulnerabilityName:   "注册功能弱密码策略",
			VulnerabilityImpact: "用户账号容易被暴力破解，存在撞库攻击风险",
			VulnerabilityDetail: "系统允许设置6位纯数字密码（如：123456），且无登录失败锁定机制。\n\n测试结果：\n- 允许密码：123456（通过）\n- 允许密码：111111（通过）\n- 无密码复杂度要求\n- 无登录失败次数限制\n\n建议：要求密码至少8位，包含大小写字母和数字，并实现登录失败锁定",
			VulnerabilityURL:    "https://example.com/register",
			Severity:            "Low",
			Status:              "Pending",
			VulnTypeKey:         "SECURITY_MISCONFIGURATION",
		},
		{
			VulnerabilityName:   "搜索功能反射型XSS漏洞",
			VulnerabilityImpact: "攻击者可构造恶意链接诱导用户点击，窃取用户凭证",
			VulnerabilityDetail: "搜索功能的关键词参数未做过滤和转义。\n\nPOC URL：\nhttps://example.com/search?q=<script>alert(document.cookie)</script>\n\n该漏洞可被用于钓鱼攻击",
			VulnerabilityURL:    "https://example.com/search?q=test",
			Severity:            "Medium",
			Status:              "Resolved",
			VulnTypeKey:         "XSS",
		},
		{
			VulnerabilityName:   "文件下载目录遍历漏洞",
			VulnerabilityImpact: "攻击者可读取服务器任意文件，包括配置文件和源代码",
			VulnerabilityDetail: "通过修改下载接口的文件路径参数，可读取系统敏感文件。\n\nPOC：\nGET /download?file=../../../etc/passwd HTTP/1.1\n\n成功读取到 /etc/passwd 文件内容",
			VulnerabilityURL:    "https://example.com/download?file=report.pdf",
			Severity:            "High",
			Status:              "Pending",
			VulnTypeKey:         "PATH_TRAVERSAL",
		},
		{
			VulnerabilityName:   "登录重定向开放重定向漏洞",
			VulnerabilityImpact: "攻击者可利用此漏洞进行钓鱼攻击，诱导用户访问恶意网站",
			VulnerabilityDetail: "登录成功后的重定向URL参数未验证，可被利用跳转到任意外部网站。\n\n恶意链接示例：\nhttps://example.com/login?redirect=https://evil.com/fake-login\n\n用户登录后会被重定向到钓鱼页面",
			VulnerabilityURL:    "https://example.com/login?redirect=https://example.com/dashboard",
			Severity:            "Low",
			Status:              "Closed",
			VulnTypeKey:         "OPEN_REDIRECT",
		},
		{
			VulnerabilityName:   "图片预览SSRF漏洞",
			VulnerabilityImpact: "攻击者可探测内网服务、读取云服务元数据，可能导致内网渗透",
			VulnerabilityDetail: "图片预览功能接受任意URL参数，服务端会请求该URL获取图片。\n\n测试案例：\n1. 探测内网：/preview?url=http://192.168.1.1:8080\n2. 读取AWS元数据：/preview?url=http://169.254.169.254/latest/meta-data/\n\n均可成功请求并返回内容",
			VulnerabilityURL:    "https://example.com/preview?url=https://example.com/image.jpg",
			Severity:            "High",
			Status:              "Triaged",
			VulnTypeKey:         "SSRF",
		},
		{
			VulnerabilityName:   "JWT Token签名验证缺陷",
			VulnerabilityImpact: "攻击者可伪造Token获取任意用户权限，包括管理员权限",
			VulnerabilityDetail: "将JWT的签名算法修改为none后，服务端仍接受该Token。\n\n复现步骤：\n1. 获取正常JWT Token\n2. 解码Token，修改alg字段为'none'\n3. 修改payload中的用户ID为管理员ID\n4. 移除签名部分\n5. 使用修改后的Token请求API\n6. 成功获取管理员权限\n\n这是一个严重的认证绕过漏洞",
			VulnerabilityURL:    "https://example.com/api/admin/users",
			Severity:            "Critical",
			Status:              "Pending",
			VulnTypeKey:         "AUTHENTICATION_BYPASS",
		},
		{
			VulnerabilityName:   "评论功能DOM型XSS",
			VulnerabilityImpact: "可在用户浏览器中执行恶意JavaScript代码",
			VulnerabilityDetail: "评论内容通过innerHTML直接插入DOM，未做转义处理。\n\n恶意评论内容：\n<img src=x onerror=alert('XSS')>\n\n所有查看该评论的用户浏览器都会执行恶意代码",
			VulnerabilityURL:    "https://example.com/post/123/comments",
			Severity:            "Medium",
			Status:              "Pending",
			VulnTypeKey:         "XSS",
		},
		{
			VulnerabilityName:   "支付金额篡改漏洞",
			VulnerabilityImpact: "攻击者可修改订单支付金额，造成经济损失",
			VulnerabilityDetail: "支付接口的金额参数未在后端验证，可通过拦截修改支付金额。\n\n复现步骤：\n1. 添加商品到购物车，总价100元\n2. 点击支付，拦截请求\n3. 修改amount参数为0.01\n4. 支付成功，实际只扣款0.01元\n\n建议：后端需要验证订单金额与实际商品金额是否一致",
			VulnerabilityURL:    "https://example.com/api/pay",
			Severity:            "Critical",
			Status:              "Triaged",
			VulnTypeKey:         "BUSINESS_LOGIC",
		},
		{
			VulnerabilityName:   "验证码绕过漏洞",
			VulnerabilityImpact: "攻击者可绕过验证码进行暴力破解或批量注册",
			VulnerabilityDetail: "验证码验证存在缺陷，同一验证码可重复使用。\n\n复现步骤：\n1. 获取验证码图片和captcha_id\n2. 识别验证码内容\n3. 使用该验证码尝试登录\n4. 即使登录失败，验证码不会失效\n5. 可继续使用同一验证码尝试不同密码\n\n建议：验证码使用后应立即失效",
			VulnerabilityURL:    "https://example.com/api/captcha",
			Severity:            "Medium",
			Status:              "Resolved",
			VulnTypeKey:         "BROKEN_ACCESS_CONTROL",
		},
		{
			VulnerabilityName:   "接口未授权访问",
			VulnerabilityImpact: "未登录用户可访问需要认证的API接口",
			VulnerabilityDetail: "部分管理接口缺少认证检查，未登录即可直接访问。\n\n受影响接口：\n- GET /api/admin/stats（系统统计）\n- GET /api/admin/logs（操作日志）\n\n这些接口应该只允许管理员访问",
			VulnerabilityURL:    "https://example.com/api/admin/stats",
			Severity:            "High",
			Status:              "Pending",
			VulnTypeKey:         "BROKEN_ACCESS_CONTROL",
		},
	}

	// 创建漏洞类型映射
	vulnTypeMap := make(map[string]domain.SystemConfig)
	for _, vt := range vulnTypes {
		vulnTypeMap[vt.ConfigKey] = vt
	}

	// 每次生成 10-16 个随机报告（追加模式）
	numReports := rand.Intn(7) + 10
	fmt.Printf("[INFO] Generating %d new reports...\n", numReports)

	successCount := 0
	for i := 0; i < numReports; i++ {
		// 随机选择一个报告模板
		template := reportTemplates[rand.Intn(len(reportTemplates))]

		// 随机分配用户和项目
		author := whitehats[rand.Intn(len(whitehats))]
		project := projects[rand.Intn(len(projects))]

		// 匹配漏洞类型，如果找不到则随机选择
		var vulnType domain.SystemConfig
		if vt, ok := vulnTypeMap[template.VulnTypeKey]; ok {
			vulnType = vt
		} else {
			vulnType = vulnTypes[rand.Intn(len(vulnTypes))]
		}

		// 随机选择危害自评（70%的概率有自评）
		var selfAssessmentID *uint
		if len(severityLevels) > 0 && rand.Float32() > 0.3 {
			id := severityLevels[rand.Intn(len(severityLevels))].ID
			selfAssessmentID = &id
		}

		// 添加时间戳使漏洞名称唯一
		timestamp := time.Now().UnixNano()
		uniqueName := fmt.Sprintf("%s_%d", template.VulnerabilityName, timestamp/1000000+int64(i))

		report := domain.Report{
			ProjectID:           project.ID,
			VulnerabilityName:   uniqueName,
			VulnerabilityTypeID: vulnType.ID,
			VulnerabilityImpact: template.VulnerabilityImpact,
			SelfAssessmentID:    selfAssessmentID,
			VulnerabilityURL:    template.VulnerabilityURL,
			VulnerabilityDetail: template.VulnerabilityDetail,
			Severity:            template.Severity,
			Status:              template.Status,
			AuthorID:            author.ID,
		}

		if err := s.db.Create(&report).Error; err != nil {
			log.Printf("[WARN] Failed to create report #%d: %v", i+1, err)
		} else {
			successCount++
			fmt.Printf("[OK] #%d %s | 项目: %s | 提交者: %s (%s) | 类型: %s | 状态: %s\n",
				report.ID,
				truncateString(report.VulnerabilityName, 30),
				truncateString(project.Name, 15),
				author.Name,
				author.Username,
				vulnType.ConfigValue,
				report.Status,
			)
		}

		// 短暂延迟确保时间戳不同
		time.Sleep(time.Millisecond)
	}

	fmt.Printf("\n[INFO] Seeded %d/%d reports successfully\n", successCount, numReports)

	// 打印统计信息
	s.printStatistics()

	return nil
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-2]) + ".."
}

// printStatistics 打印统计信息
func (s *ReportSeeder) printStatistics() {
	fmt.Println("\n========== 测试数据统计 ==========")

	// 按用户统计
	type userStat struct {
		Username string
		Count    int64
	}
	var userStats []userStat
	s.db.Table("reports").
		Select("users.username, count(*) as count").
		Joins("JOIN users ON reports.author_id = users.id").
		Group("users.username").
		Scan(&userStats)

	fmt.Println("\n📊 按提交者统计:")
	for _, stat := range userStats {
		fmt.Printf("   %s: %d 条报告\n", stat.Username, stat.Count)
	}

	// 按状态统计
	type statusStat struct {
		Status string
		Count  int64
	}
	var statusStats []statusStat
	s.db.Table("reports").
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusStats)

	fmt.Println("\n📋 按状态统计:")
	statusMap := map[string]string{
		"Pending":  "待审核",
		"Triaged":  "已确认",
		"Resolved": "已修复",
		"Closed":   "已关闭",
	}
	for _, stat := range statusStats {
		name := statusMap[stat.Status]
		if name == "" {
			name = stat.Status
		}
		fmt.Printf("   %s (%s): %d 条\n", stat.Status, name, stat.Count)
	}

	// 按严重程度统计
	type severityStat struct {
		Severity string
		Count    int64
	}
	var severityStats []severityStat
	s.db.Table("reports").
		Select("severity, count(*) as count").
		Group("severity").
		Scan(&severityStats)

	fmt.Println("\n🔥 按危害等级统计:")
	severityOrder := []string{"Critical", "High", "Medium", "Low"}
	for _, sev := range severityOrder {
		for _, stat := range severityStats {
			if stat.Severity == sev {
				fmt.Printf("   %s: %d 条\n", stat.Severity, stat.Count)
				break
			}
		}
	}

	fmt.Println("\n===================================")
}
