package llm

import (
	"bufio"
	"encoding/json"
	"go.qingyu31.com/gtl"
	"io"
	"time"
)

type Result[T any] struct {
	Response *T
	Reason   string
	Metric   Metric
}

type StreamResult[T any] struct {
	Iterator StreamIterator[T]
	Reason   string
	Metric   Metric
}

type Metric struct {
	InputTokens  int
	OutputTokens int
	Latency      time.Duration
}

type StreamIterator[T any] interface {
	HasNext() bool
	Next() (*T, error)
	Close() error
}

type ItemIterator[T any] struct {
	queue gtl.List[*T]
}

func NewItemIterator[T any]() *ItemIterator[T] {
	r := new(ItemIterator[T])
	r.queue = gtl.NewLinkedList[*T]()
	return r
}

func (i *ItemIterator[T]) Write(data *T) {
	i.queue.PushFront(data)
}

func (i *ItemIterator[T]) HasNext() bool {
	return i.queue.Len() > 0
}

func (i *ItemIterator[T]) Next() (*T, error) {
	if !i.HasNext() {
		return nil, io.EOF
	}
	return i.queue.PopBack().Value(), nil
}

func (i *ItemIterator[T]) Close() error {
	return nil
}

type BytesReader[T any] struct {
	reader  io.ReadCloser // Required for Closing
	scanner *bufio.Scanner
	value   *T
}

func NewBytesReader[T any](r io.ReadCloser) *BytesReader[T] {
	return &BytesReader[T]{reader: r, scanner: bufio.NewScanner(r)}
}

func (sr *BytesReader[T]) HasNext() bool {
	return sr.scanner.Scan()
}

func (sr *BytesReader[T]) Next() (*T, error) {
	for sr.scanner.Scan() {
		line := sr.scanner.Text()

		if line == "" || line[0] == ':' {
			continue
		}
		data := new(T)
		err := json.Unmarshal([]byte(line), data)
		return data, err
	}
	return nil, sr.scanner.Err()
}

func (sr *BytesReader[T]) Close() error {
	return sr.reader.Close()
}
