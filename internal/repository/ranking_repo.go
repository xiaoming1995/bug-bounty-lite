package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type rankingRepo struct {
	db *gorm.DB
}

func NewRankingRepo(db *gorm.DB) domain.RankingRepository {
	return &rankingRepo{db: db}
}

func (r *rankingRepo) GetGlobalRanking(limit int) ([]domain.RankingItem, error) {
	var items []domain.RankingItem

	// 积分规则：严重40、高危30、中危20、低危10
	sql := `
		SELECT 
			u.id as user_id,
			u.name as user_name,
			a.url as avatar_url,
			COALESCE(SUM(CASE 
				WHEN LOWER(r.severity) = 'critical' THEN 40
				WHEN LOWER(r.severity) = 'high' THEN 30
				WHEN LOWER(r.severity) = 'medium' THEN 20
				WHEN LOWER(r.severity) = 'low' THEN 10
				ELSE 0 
			END), 0) as points,
			COUNT(r.id) as vuln_count,
			COALESCE(SUM(CASE WHEN LOWER(r.severity) = 'critical' THEN 1 ELSE 0 END), 0) as critical_count,
			COALESCE(SUM(CASE WHEN LOWER(r.severity) = 'high' THEN 1 ELSE 0 END), 0) as high_count
		FROM users u
		LEFT JOIN reports r ON u.id = r.author_id AND LOWER(r.status) IN ('audited', 'triaged', 'resolved', 'closed') AND r.deleted_at IS NULL
		LEFT JOIN avatars a ON u.avatar_id = a.id
		WHERE u.role = 'whitehat'
		GROUP BY u.id
		ORDER BY points DESC, vuln_count DESC
		LIMIT ?
	`

	rows, err := r.db.Raw(sql, limit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rank := 1
	for rows.Next() {
		var item domain.RankingItem
		var avatarUrl *string // 使用指针处理可能为 NULL 的情况
		err := rows.Scan(
			&item.UserID,
			&item.UserName,
			&avatarUrl,
			&item.Points,
			&item.VulnCount,
			&item.CriticalCount,
			&item.HighCount,
		)
		if err != nil {
			return nil, err
		}
		if avatarUrl != nil {
			item.AvatarUrl = *avatarUrl
		}
		item.Rank = rank
		items = append(items, item)
		rank++
	}

	return items, nil
}

func (r *rankingRepo) GetStatistics() (*domain.RankingStatistics, error) {
	stats := &domain.RankingStatistics{}

	// 统计注册白帽子总数 (角色为 whitehat 的用户)
	err := r.db.Model(&domain.User{}).Where("role = ?", "whitehat").Count(&stats.TotalHunters).Error
	if err != nil {
		return nil, err
	}

	// 统计已发现漏洞总数 (仅统计所有非 Pending 状态的有效漏洞，与计分对齐)
	err = r.db.Model(&domain.Report{}).Where("LOWER(status) IN ('audited', 'triaged', 'resolved', 'closed')").Count(&stats.TotalVulns).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}
