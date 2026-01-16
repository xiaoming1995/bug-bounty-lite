package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	Service domain.RankingService
}

func NewRankingHandler(s domain.RankingService) *RankingHandler {
	return &RankingHandler{Service: s}
}

func (h *RankingHandler) GetRanking(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	items, stats, err := h.Service.GetRanking(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取排行榜失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":       items,
		"statistics": stats,
	})
}
