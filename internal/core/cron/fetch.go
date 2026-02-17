package cron

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/bestruirui/bestsub/internal/core/fetch"
	"github.com/bestruirui/bestsub/internal/database/op"
	subModel "github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/utils/generic"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/robfig/cron/v3"
)

var fetchFunc = generic.MapOf[uint16, cronFunc]{}
var fetchScheduled = generic.MapOf[uint16, cron.EntryID]{}
var fetchRunning = generic.MapOf[uint16, context.CancelFunc]{}

func FetchLoad() {
	subData, err := op.GetSubList(context.Background())
	if err != nil {
		log.Errorf("failed to load sub data: %v", err)
		return
	}
	for i := range subData {
		FetchAdd(&subData[i])
	}
}

func FetchAdd(data *subModel.Data) error {
	subID := data.ID
	cronExpr := data.CronExpr
	config := data.Config
	enable := data.Enable

	fetchFunc.Store(subID, cronFunc{
		fn: func() {
			fetchCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			fetchRunning.Store(subID, cancel)
			defer func() {
				cancel()
				fetchRunning.Delete(subID)
			}()
			result := fetch.Do(fetchCtx, subID, config)
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 2*time.Second)
			if err := op.UpdateSubResult(updateCtx, subID, result); err != nil {
				log.Warnf("failed to update sub result: %v", err)
			}
			updateCancel()
			sub, err := op.GetSubByID(context.Background(), subID)
			if err != nil {
				log.Warnf("failed to get sub by id: %v", err)
				return
			}
			if !sub.Enable {
				FetchDisable(subID)
				log.Infof("fetch task %d auto disable", subID)
			}
		},
		cronExpr: cronExpr,
	})
	if enable {
		FetchEnable(subID)
	}
	return nil
}

func FetchRun(subID uint16) subModel.Result {
	if ft, ok := fetchFunc.Load(subID); ok {
		ft.fn()
	} else {
		log.Warnf("fetch task %d not found", subID)
		return subModel.Result{
			Msg:     "fetch task not found",
			LastRun: time.Now(),
		}
	}
	sub, err := op.GetSubByID(context.Background(), subID)
	if err != nil {
		log.Warnf("failed to get sub by id: %v", err)
		return subModel.Result{
			Msg:     "fetch task not found",
			LastRun: time.Now(),
		}
	}
	var result subModel.Result
	if err := json.Unmarshal([]byte(sub.Result), &result); err != nil {
		log.Warnf("failed to unmarshal sub result: %v", err)
		return subModel.Result{
			Msg:     "failed to unmarshal sub result",
			LastRun: time.Now(),
		}
	}
	return result
}

func FetchEnable(subID uint16) error {
	if _, ok := fetchScheduled.Load(subID); ok {
		log.Warnf("fetch task %d already scheduled", subID)
		return nil
	}
	if ft, ok := fetchFunc.Load(subID); ok {
		entryID, err := scheduler.AddFunc(ft.cronExpr,
			func() {
				time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
				ft.fn()
			})
		if err != nil {
			log.Errorf("failed to add task: %v", err)
			return err
		}
		fetchScheduled.Store(subID, entryID)
	}
	return nil
}
func FetchDisable(subID uint16) error {
	if entryID, ok := fetchScheduled.Load(subID); ok {
		scheduler.Remove(entryID)
		fetchScheduled.Delete(subID)
		if cancel, ok := fetchRunning.Load(subID); ok {
			cancel()
			fetchRunning.Delete(subID)
		}
	}
	return nil
}

func FetchRemove(subID uint16) error {
	if entryID, ok := fetchScheduled.Load(subID); ok {
		scheduler.Remove(entryID)
		fetchScheduled.Delete(subID)
		fetchFunc.Delete(subID)
		if cancel, ok := fetchRunning.Load(subID); ok {
			cancel()
			fetchRunning.Delete(subID)
		}
	}
	return nil
}
func FetchStop(subID uint16) error {
	if cancel, ok := fetchRunning.Load(subID); ok {
		cancel()
		fetchRunning.Delete(subID)
	}
	return nil
}
func FetchUpdate(data *subModel.Data) error {
	FetchRemove(data.ID)
	FetchAdd(data)
	return nil
}
func FetchStatus(subID uint16) string {
	if _, ok := fetchRunning.Load(subID); ok {
		return RunningStatus
	}
	if _, ok := fetchScheduled.Load(subID); ok {
		return ScheduledStatus
	}
	if _, ok := fetchFunc.Load(subID); ok {
		return PendingStatus
	}
	return DisabledStatus
}
