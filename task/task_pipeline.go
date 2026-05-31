package task

import (
	"context"
	"errors"
	"time"

	gopipeline "github.com/rushairer/go-pipeline/v2"
)

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
