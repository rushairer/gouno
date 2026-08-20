package task_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	gopipeline "github.com/rushairer/go-pipeline/v2"
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

func TestTaskPipelineErrorJoin(t *testing.T) {
	pipeline := task.NewTaskPipeline(100, 2, 100*time.Millisecond)
	// 先初始化错误通道容量，避免 Run 以容量 1 初始化后第二个错误被丢弃
	errs := pipeline.ErrorChan(2)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	runDone := make(chan struct{})
	go func() {
		defer close(runDone)
		_ = pipeline.Run(ctx, 1)
	}()

	errA := errors.New("error A")
	errB := errors.New("error B")
	pipeline.DataChan() <- &mockTask{runFunc: func(ctx context.Context) error { return errA }}
	pipeline.DataChan() <- &mockTask{runFunc: func(ctx context.Context) error { return errB }}

	select {
	case err := <-errs:
		// 同一批次的两个错误应通过 errors.Join 合并返回
		if err == nil {
			t.Fatal("expected joined error, got nil")
		}
		if !errors.Is(err, errA) || !errors.Is(err, errB) {
			t.Fatalf("expected both errors joined, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for pipeline error")
	}

	// 错误已确认，主动取消以停止管道（否则会运行至 ctx 超时）
	cancel()
	select {
	case <-runDone:
	case <-time.After(5 * time.Second):
		t.Fatal("pipeline should stop after cancel")
	}
}

func TestTaskPipelineContextCancelDrainsAndStops(t *testing.T) {
	pipeline := task.NewTaskPipeline(100, 1, 50*time.Millisecond)

	var executed atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErr := make(chan error, 1)
	go func() {
		runErr <- pipeline.Run(ctx, 1)
	}()

	pipeline.DataChan() <- &mockTask{
		runFunc: func(ctx context.Context) error {
			executed.Add(1)
			return nil
		},
	}

	// 等待任务被执行
	deadline := time.Now().Add(3 * time.Second)
	for executed.Load() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("timeout waiting for task execution")
		}
		time.Sleep(5 * time.Millisecond)
	}

	cancel()
	select {
	case err := <-runErr:
		// DrainOnCancel: 收尾后以"上下文已关闭"退出
		if !errors.Is(err, gopipeline.ErrContextIsClosed) {
			t.Fatalf("expected ErrContextIsClosed, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("pipeline should stop after context cancel")
	}
}

func TestTaskPipelineCancelEmptyPipeline(t *testing.T) {
	pipeline := task.NewTaskPipeline(100, 2, time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErr := make(chan error, 1)
	go func() {
		runErr <- pipeline.Run(ctx, 1)
	}()

	// 无任务时立即取消：应快速返回且不 panic
	cancel()
	select {
	case err := <-runErr:
		if !errors.Is(err, gopipeline.ErrContextIsClosed) {
			t.Fatalf("expected ErrContextIsClosed, got: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("cancel of empty pipeline should return promptly")
	}
}
