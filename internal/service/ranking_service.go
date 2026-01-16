package service

import "bug-bounty-lite/internal/domain"

type rankingService struct {
	repo domain.RankingRepository
}

func NewRankingService(repo domain.RankingRepository) domain.RankingService {
	return &rankingService{repo: repo}
}

func (s *rankingService) GetRanking(limit int) ([]domain.RankingItem, *domain.RankingStatistics, error) {
	if limit <= 0 {
		limit = 100 // 默认返回前100名
	}

	items, err := s.repo.GetGlobalRanking(limit)
	if err != nil {
		return nil, nil, err
	}

	stats, err := s.repo.GetStatistics()
	if err != nil {
		return nil, nil, err
	}

	return items, stats, nil
}
