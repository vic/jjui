package graph

import (
	"io"
)

func ParseRows(reader io.Reader) []Row {
	var rows []Row
	controlChan := make(chan ControlMsg)
	defer close(controlChan)
	streamerChannel, err := ParseRowsStreaming(reader, controlChan, 50)
	if err != nil {
		return nil
	}
	for {
		controlChan <- RequestMore
		chunk := <-streamerChannel
		rows = append(rows, chunk.Rows...)
		if !chunk.HasMore {
			break
		}
	}
	return rows
}
