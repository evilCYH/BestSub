package handlers

import (
	"net/http"
	"strconv"

	"github.com/bestruirui/bestsub/internal/core/node"
	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/gin-gonic/gin"
)

func init() {
	router.NewGroupRouter("/api/v1/node/log/detail").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(GetNodeTestLogs),
		)
}

// GetNodeTestLogs godoc
// @Summary 获取节点详细测试日志
// @Description 获取订阅下各节点的详细测试日志，支持分页、级别筛选和关键词搜索
// @Tags nodes
// @Param sub_id query int true "订阅ID"
// @Param level query string false "级别筛选: info/warn/error"
// @Param keyword query string false "关键词搜索（节点名或日志内容）"
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页数量，默认50，最大100" default(50)
// @Success 200 {object} map[string]any{code=int,data=node.NodeTestLogResponse}
// @Router /api/v1/node/log/detail [get]
func GetNodeTestLogs(c *gin.Context) {
	var query nodeModel.NodeTestLogQuery

	subID, err := strconv.Atoi(c.Query("sub_id"))
	if err != nil || subID <= 0 {
		resp.Error(c, http.StatusBadRequest, "invalid sub_id")
		return
	}
	query.SubID = uint16(subID)

	query.Level = c.DefaultQuery("level", "")
	query.Keyword = c.DefaultQuery("keyword", "")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	query.Page = page

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}
	query.PageSize = pageSize

	logs, total := node.QueryNodeLogs(query)

	resp.Success(c, nodeModel.NodeTestLogResponse{
		Total: total,
		List:  logs,
	})
}
