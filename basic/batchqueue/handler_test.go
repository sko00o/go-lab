package batchqueue

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkHandler(b *testing.B) {
	tests := []struct {
		name string
		ta   testArgs
	}{
		{
			name: "32B/64B",
			ta: testArgs{
				timeout:     10 * time.Millisecond,
				size:        10,
				bytes:       64,
				contentSize: 32,
			},
		},
		{
			name: "32B/64B/chan",
			ta: testArgs{
				timeout:         10 * time.Millisecond,
				size:            10,
				bytes:           64,
				contentSize:     32,
				useChannelQueue: true,
			},
		},
		{
			name: "128B/64B",
			ta: testArgs{
				timeout:     10 * time.Millisecond,
				size:        10,
				bytes:       64,
				contentSize: 128,
			},
		},
		{
			name: "128B/64B/chan",
			ta: testArgs{
				timeout:         10 * time.Millisecond,
				size:            10,
				bytes:           64,
				contentSize:     128,
				useChannelQueue: true,
			},
		},
		{
			name: "128B/256B",
			ta: testArgs{
				timeout:     10 * time.Millisecond,
				size:        10,
				bytes:       256,
				contentSize: 128,
			},
		},
		{
			name: "128B/256B/chan",
			ta: testArgs{
				timeout:         10 * time.Millisecond,
				size:            10,
				bytes:           256,
				contentSize:     128,
				useChannelQueue: true,
			},
		},
		{
			name: "128B/1024B",
			ta: testArgs{
				timeout:     10 * time.Millisecond,
				size:        10,
				bytes:       1024,
				contentSize: 128,
			},
		},
		{
			name: "128B/1024B/chan",
			ta: testArgs{
				timeout:         10 * time.Millisecond,
				size:            10,
				bytes:           1024,
				contentSize:     128,
				useChannelQueue: true,
			},
		},
		{
			name: "250ms/128B/1024B",
			ta: testArgs{
				timeout:     250 * time.Millisecond,
				size:        10,
				bytes:       1024,
				contentSize: 128,
			},
		},
		{
			name: "250ms/128B/1024B/chan",
			ta: testArgs{
				timeout:         250 * time.Millisecond,
				size:            10,
				bytes:           1024,
				contentSize:     128,
				useChannelQueue: true,
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			commonBenchmark(b, tt.ta)
		})
	}
}

type testArgs struct {
	timeout         time.Duration
	size            int
	bytes           int64
	contentSize     int
	useChannelQueue bool
}

func commonBenchmark(b *testing.B, ta testArgs) {
	b.StopTimer()
	var (
		batchTimeout = ta.timeout
		batchSize    = ta.size
		batchBytes   = ta.bytes
		content      = make([]byte, ta.contentSize)
		use          = ta.useChannelQueue
	)
	rand.Read(content)
	var cnt = 0
	defer func() {
		if b.N != cnt {
			b.Errorf("cnt %d, expect %d", cnt, b.N)
		}
	}()
	handle := func(batch *Batch) {
		for _, item := range batch.Items {
			if bytes.Equal(item.Bytes(), content) {
				cnt++
			}
		}
	}
	r := NewHandler(batchTimeout, batchSize, batchBytes, handle, use)
	b.StartTimer()
	defer r.Close()

	for i := 0; i < b.N; i++ {
		r.Write(MessageItem{Content: content})
	}
}
