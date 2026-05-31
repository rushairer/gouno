package task_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rushairer/gouno/task"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskPipeline(t *testing.T) {
	pipeline := task.NewTaskPipeline(100, 10, time.Second)
	assert.NotNil(t, pipeline)
}

func TestTaskPipelineExecuteSingleTask(t *testing.T) {
	pipeline := task.NewTaskPipeline(100, 1, 100*time.Millisecond)

	var executed atomic.Bool
	tk := &mockTask{
		runFunc: func(ctx context.Context) error {
			executed.Store(true)
			return nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		_ = pipeline.Run(ctx, 1)
	}()

	pipeline.DataChan() <- tk
	time.Sleep(300 * time.Millisecond)

	assert.True(t, executed.Load(), "task should have been executed")
}

func TestTaskPipelineExecuteMultipleTasks(t *testing.T) {
	pipeline := task.NewTaskPipeline(100, 3, 100*time.Millisecond)

	var count atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		_ = pipeline.Run(ctx, 1)
	}()

	for i := 0; i < 3; i++ {
		tk := &mockTask{
			runFunc: func(ctx context.Context) error {
				count.Add(1)
				return nil
			},
		}
		pipeline.DataChan() <- tk
	}

	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, int32(3), count.Load(), "all 3 tasks should have been executed")
}

func TestTaskPipelineErrorCollection(t *testing.T) {
	pipeline := task.NewTaskPipeline(100, 2, 100*time.Millisecond)

	var errCount atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		_ = pipeline.Run(ctx, 1)
	}()

	for i := 0; i < 2; i++ {
		tk := &mockTask{
			runFunc: func(ctx context.Context) error {
				errCount.Add(1)
				return errors.New("task error")
			},
		}
		pipeline.DataChan() <- tk
	}

	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, int32(2), errCount.Load(), "both tasks should execute even with errors")
}
