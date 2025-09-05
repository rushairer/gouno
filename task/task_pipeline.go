package task

import (
	"context"
	"errors"
	"time"

	gopipeline "github.com/rushairer/go-pipeline"
)

func NewTaskPipeline(buffSize uint32, flushSize uint32, flushInterval time.Duration) *gopipeline.Pipeline[Task] {
	return gopipeline.NewPipeline(
		gopipeline.PipelineConfig{
			BufferSize:    buffSize,
			FlushSize:     flushSize,
			FlushInterval: flushInterval,
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
