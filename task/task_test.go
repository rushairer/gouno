package task_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rushairer/gouno/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTask struct {
	runFunc func(ctx context.Context) error
}

func (m *mockTask) Run(ctx context.Context) error {
	return m.runFunc(ctx)
}

func TestTaskInterface(t *testing.T) {
	var _ task.Task = &mockTask{}
}

func TestTaskRunWithSuccess(t *testing.T) {
	executed := false
	tk := &mockTask{
		runFunc: func(ctx context.Context) error {
			executed = true
			return nil
		},
	}

	err := tk.Run(context.Background())
	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestTaskRunWithError(t *testing.T) {
	expectedErr := errors.New("task failed")
	tk := &mockTask{
		runFunc: func(ctx context.Context) error {
			return expectedErr
		},
	}

	err := tk.Run(context.Background())
	assert.ErrorIs(t, err, expectedErr)
}

func TestTaskContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tk := &mockTask{
		runFunc: func(ctx context.Context) error {
			return ctx.Err()
		},
	}

	err := tk.Run(ctx)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}
