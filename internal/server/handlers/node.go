package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bestruirui/bestsub/internal/core/node"
	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
	"github.com/bestruirui/bestsub/internal/server/middleware"
	"github.com/bestruirui/bestsub/internal/server/resp"
	"github.com/bestruirui/bestsub/internal/server/router"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func init() {
	router.NewGroupRouter("/api/v1/node").
		Use(middleware.Auth()).
		AddRoute(
			router.NewRoute("", router.GET).
				Handle(getNodes),
		).
		AddRoute(
			router.NewRoute("/log", router.GET).
				Handle(getNodeUpdateLog),
		)
}

type nodeMeta struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// getNodes 获取节点列表
// @Summary 获取节点列表
// @Description 获取节点列表，可选 sub_id 过滤
// @Tags 节点
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sub_id query string false "订阅ID，支持逗号分隔"
// @Success 200 {object} resp.ResponseStruct{data=[]node.Response} "获取成功"
// @Failure 400 {object} resp.ResponseStruct "请求参数错误"
// @Failure 401 {object} resp.ResponseStruct "未授权"
// @Router /api/v1/node [get]
func getNodes(c *gin.Context) {
	subIDRaw := strings.TrimSpace(c.Query("sub_id"))
	var nodes []nodeModel.Data
	if subIDRaw != "" {
		ids, err := parseSubIDs(subIDRaw)
		if err != nil {
			resp.ErrorBadRequest(c)
			return
		}
		nodes = *node.GetBySubId(ids)
	} else {
		nodes = node.GetAll()
	}

	respData := make([]nodeModel.Response, 0, len(nodes))
	for _, n := range nodes {
		var meta nodeMeta
		_ = yaml.Unmarshal(n.Base.Raw, &meta)

		item := nodeModel.Response{
			SubID:       n.Base.SubId,
			UniqueKey:   n.Base.UniqueKey,
			Name:        meta.Name,
			Type:        meta.Type,
			Country:     "",
			AliveStatus: 0,
		}
		if n.Info != nil {
			item.Delay = n.Info.Delay.Average()
			item.SpeedUp = n.Info.SpeedUp.Average()
			item.SpeedDown = n.Info.SpeedDown.Average()
			item.Risk = n.Info.Risk
			item.AliveStatus = n.Info.AliveStatus
			item.Country = n.Info.Country
		}
		respData = append(respData, item)
	}

	resp.Success(c, respData)
}

func parseSubIDs(raw string) ([]uint16, error) {
	parts := strings.Split(raw, ",")
	ids := make([]uint16, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseUint(part, 10, 16)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint16(id))
	}
	if len(ids) == 0 {
		return nil, errors.New("empty sub_id")
	}
	return ids, nil
}

// getNodeUpdateLog 获取订阅节点更新日志
// @Summary 获取订阅节点更新日志
// @Description 获取订阅节点更新日志
// @Tags 节点
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sub_id query int true "订阅ID"
// @Param limit query int false "返回条数"
// @Success 200 {object} resp.ResponseStruct{data=node.UpdateLogResponse} "获取成功"
// @Failure 400 {object} resp.ResponseStruct "请求参数错误"
// @Failure 401 {object} resp.ResponseStruct "未授权"
// @Router /api/v1/node/log [get]
func getNodeUpdateLog(c *gin.Context) {
	subIDStr := strings.TrimSpace(c.Query("sub_id"))
	if subIDStr == "" {
		resp.ErrorBadRequest(c)
		return
	}
	parsedID, err := strconv.ParseUint(subIDStr, 10, 16)
	if err != nil {
		resp.ErrorBadRequest(c)
		return
	}
	limit := 5
	if limitStr := strings.TrimSpace(c.Query("limit")); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			resp.ErrorBadRequest(c)
			return
		}
		if parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	resp.Success(c, node.GetUpdateLog(uint16(parsedID), limit))
}
