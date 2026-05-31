package task

import (
	"context"
	"errors"
	"time"

	gopipeline "github.com/rushairer/go-pipeline/v2"
)

// NewTaskPipeline creates a buffered pipeline that batches and flushes Task items.
// Tasks within a batch are executed sequentially; errors from all tasks are joined
// and returned together rather than stopping at the first failure.
//
// Parameters:
//   - buffSize: internal channel buffer capacity
//   - flushSize: number of items per batch before flushing
//   - flushInterval: maximum time to wait before flushing an incomplete batch
func NewTaskPipeline(buffSize uint32, flushSize uint32, flushInterval time.Duration) *gopipeline.StandardPipeline[Task] {
	return gopipeline.NewStandardPipeline(
		gopipeline.PipelineConfig{
			BufferSize:               buffSize,
			FlushSize:                flushSize,
			FlushInterval:            flushInterval,
			DrainOnCancel:            true,
			DrainGracePeriod:         time.Second,
			MaxConcurrentFlushes:     8,
			FinalFlushOnCloseTimeout: 10 * time.Second,
		},
		// 执行批次中的所有 task，收集全部错误后统一返回，不在首个错误时中断
		func(ctx context.Context, batchData []Task) (err error) {
			for _, task := range batchData {
				err = errors.Join(err, task.Run(ctx))
			}
			return err
		},
	)
}
