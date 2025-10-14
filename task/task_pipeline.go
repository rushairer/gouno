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
		func(ctx context.Context, batchData []Task) (err error) {
			for _, task := range batchData {
				if errInner := task.Run(ctx); errInner != nil {
					return errInner
				} else {
					err = errors.Join(err, errInner)
				}
			}
			return err
		},
	)
}
