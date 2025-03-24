package screen

import (
	"bufio"
	"bytes"
	"io"
)

func Parse(raw []byte) []Segment {
	var segments []Segment
	for segment := range ParseFromReader(bytes.NewReader(raw)) {
		segments = append(segments, *segment)
	}
	return segments
}

func ParseFromReader(r io.Reader) <-chan *Segment {
	ch := make(chan *Segment)
	go func() {
		defer close(ch)
		var buffer bytes.Buffer
		currentParams := ""
		reader := bufio.NewReader(r)

		for {
			b, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}

			if b == 0x1B {
				peekBytes, err := reader.Peek(1)
				if err != nil {
					buffer.WriteByte(b)
					break
				}

				if len(peekBytes) >= 1 && peekBytes[0] == '[' {
					_, _ = reader.Discard(1)
					if buffer.Len() > 0 {
						ch <- &Segment{
							Text:   buffer.String(),
							Params: currentParams,
						}
						buffer.Reset()
					}

					var seq bytes.Buffer
					for {
						c, err := reader.ReadByte()
						if err != nil || c == 'm' {
							break
						}
						seq.WriteByte(c)
					}

					paramStr := seq.String()
					if paramStr == "0" {
						currentParams = ""
					} else {
						currentParams = paramStr
					}
				} else {
					buffer.WriteByte(b)
					if len(peekBytes) >= 1 {
						nextByte, _ := reader.ReadByte()
						buffer.WriteByte(nextByte)
					}
				}
			} else {
				buffer.WriteByte(b)
			}
		}

		if buffer.Len() > 0 {
			ch <- &Segment{
				Text:   buffer.String(),
				Params: currentParams,
			}
		}
	}()
	return ch
}
