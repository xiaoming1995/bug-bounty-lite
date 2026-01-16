package domain

// RankingItem 排行榜单项
type RankingItem struct {
	Rank          int    `json:"rank"`
	UserID        uint   `json:"user_id"`
	UserName      string `json:"user_name"`
	AvatarUrl     string `json:"avatar_url"`
	Points        int    `json:"points"`
	VulnCount     int    `json:"vulns"`
	CriticalCount int    `json:"critical"`
	HighCount     int    `json:"high"`
}

// RankingStatistics 排行榜全局统计
type RankingStatistics struct {
	TotalHunters int64 `json:"total_hunters"`
	TotalVulns   int64 `json:"total_vulns"`
}

// RankingRepository 排行榜仓储接口
type RankingRepository interface {
	GetGlobalRanking(limit int) ([]RankingItem, error)
	GetStatistics() (*RankingStatistics, error)
}

// RankingService 排行榜服务接口
type RankingService interface {
	GetRanking(limit int) ([]RankingItem, *RankingStatistics, error)
}
