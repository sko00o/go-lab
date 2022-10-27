package batchqueue

type Item interface {
	Size() int64
	Bytes() []byte
}

type MessageItem struct {
	Content []byte
}

func (m MessageItem) Size() int64 {
	return int64(len(m.Content))
}

func (m MessageItem) Bytes() []byte {
	return m.Content
}
