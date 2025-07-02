package graph

import (
	"bufio"
	"context"
	"errors"
	"github.com/idursun/jjui/internal/jj"
	appContext "github.com/idursun/jjui/internal/ui/context"
	"io"
	"time"
)

const DefaultBatchSize = 50

type GraphStreamer struct {
	command     *appContext.StreamingCommand
	cancel      context.CancelFunc
	controlChan chan ControlMsg
	rowsChan    <-chan RowBatch
	batchSize   int
}

func NewGraphStreamer(ctx *appContext.MainContext, revset string) (*GraphStreamer, error) {
	streamerCtx, cancel := context.WithCancel(context.Background())

	command, err := ctx.RunCommandStreaming(streamerCtx, jj.Log(revset))
	if err != nil {
		cancel()
		return nil, err
	}

	// Check stderr with timeout
	errCh := make(chan error, 1)
	go func() {
		errReader := bufio.NewReader(command.ErrPipe)
		data, err := errReader.Peek(1)
		if err == nil && len(data) > 0 {
			errorData, _ := io.ReadAll(errReader)
			errCh <- errors.New(string(errorData))
		} else {
			errCh <- nil
		}
	}()

	// Wait for stderr check with timeout
	select {
	case stderrErr := <-errCh:
		if stderrErr != nil {
			cancel()
			_ = command.Close()
			return nil, stderrErr
		}
	case <-time.After(100 * time.Millisecond):
		// Timeout, assume no error and continue
	}

	// Set up stdout processing
	controlChan := make(chan ControlMsg, 1)
	reader := bufio.NewReader(command)

	rowsChan, err := ParseRowsStreaming(reader, controlChan, DefaultBatchSize)
	if err != nil {
		cancel()
		_ = command.Close()
		return nil, err
	}

	return &GraphStreamer{
		command:     command,
		cancel:      cancel,
		controlChan: controlChan,
		rowsChan:    rowsChan,
		batchSize:   DefaultBatchSize,
	}, nil
}
func (g *GraphStreamer) RequestMore() RowBatch {
	g.controlChan <- RequestMore
	return <-g.rowsChan
}

func (g *GraphStreamer) Close() {
	if g == nil {
		return
	}

	if g.controlChan != nil {
		g.controlChan <- Close
		close(g.controlChan)
		g.controlChan = nil
	}

	if g.cancel != nil {
		g.cancel()
		_ = g.command.Close()
		g.cancel = nil
	}

	g.rowsChan = nil
	g.command = nil
}
