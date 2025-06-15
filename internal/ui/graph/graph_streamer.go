package graph

import (
	"bufio"
	"context"
	"errors"
	"github.com/idursun/jjui/internal/jj"
	appContext "github.com/idursun/jjui/internal/ui/context"
	"io"
)

const DefaultBatchSize = 50

type GraphStreamer struct {
	command     *appContext.StreamingCommand
	cancel      context.CancelFunc
	controlChan chan ControlMsg
	rowsChan    <-chan RowBatch
	batchSize   int
}

func NewGraphStreamer(ctx appContext.AppContext, revset string) (*GraphStreamer, error) {
	streamerCtx, cancel := context.WithCancel(context.Background())

	command, err := ctx.RunCommandStreaming(streamerCtx, jj.Log(revset))
	if err != nil {
		cancel()
		return nil, err
	}

	controlChan := make(chan ControlMsg, 1)

	stdoutChan := make(chan struct{})
	stderrChan := make(chan error)
	reader := bufio.NewReader(command)

	go func() {
		if _, err := reader.Peek(1); err == nil {
			stdoutChan <- struct{}{}
		}
	}()

	// Check stderr for data
	go func() {
		errReader := bufio.NewReader(command.ErrPipe)
		if _, err := errReader.Peek(1); err == nil {
			// There's error data available
			errorOutput, _ := io.ReadAll(errReader)
			stderrChan <- errors.New(string(errorOutput))
		}
	}()

	// Wait for either stdout or stderr to have data
	select {
	case <-stdoutChan:
		// Data is available on stdout, proceed with parsing
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

	case err := <-stderrChan:
		// An actual error occurred
		cancel()
		_ = command.Close()
		return nil, err
	}
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
