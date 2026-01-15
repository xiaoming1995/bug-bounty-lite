package domain

// SeverityStatistics 漏洞严重等级统计
type SeverityStatistics struct {
	Critical int64 `json:"critical"` // 严重
	High     int64 `json:"high"`     // 高危
	Medium   int64 `json:"medium"`   // 中危
	Low      int64 `json:"low"`      // 低危
	None     int64 `json:"none"`     // 无危害
	Total    int64 `json:"total"`    // 总数
}

// TrendItem 趋势数据项
type TrendItem struct {
	Label string `json:"label"` // 标签（如"1月"、"2024"等）
	Value int64  `json:"value"` // 数值
}

// DashboardRepository 仪表盘数据仓库接口
type DashboardRepository interface {
	// CountBySeverity 按危害等级统计已审核漏洞数量
	// authorID: nil=查询所有, 非nil=只查询指定用户的
	CountBySeverity(authorID *uint) (*SeverityStatistics, error)

	// GetTrend 获取漏洞趋势数据
	// period: "day" (当月每天), "month" (当年每月), "year" (近5年)
	// authorID: nil=查询所有, 非nil=只查询指定用户的
	GetTrend(period string, authorID *uint) ([]TrendItem, error)

	// ListByStatus 按状态获取漏洞列表
	// isPending: true=待审核, false=已审核
	// authorID: nil=查询所有, 非nil=只查询指定用户的
	ListByStatus(isPending bool, limit int, authorID *uint) ([]Report, int64, error)
}

// DashboardService 仪表盘业务逻辑接口
type DashboardService interface {
	GetStatistics(userID uint, userRole string) (*SeverityStatistics, error)
	GetTrend(period string, userID uint, userRole string) ([]TrendItem, error)
	GetReportsByType(reportType string, limit int, userID uint, userRole string) ([]Report, int64, error)
}
