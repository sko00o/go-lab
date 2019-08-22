package leaky_buffer

type Buffer []byte

var freeList = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)

func client() {
	for {
		var b *Buffer
		select {
		case b = <-freeList:
			// Got one, nothing to do.
		default:
			// None free, so allocate a new one.
			b = new(Buffer)
		}
		load(b)         // Read next message from the net.
		serverChan <- b // Send to server.
	}
}

func server() {
	for {
		b := <-serverChan
		process(b)
		// Reuse buffer if there's room
		select {
		case freeList <- b:
			// Buffer on free list; nothing more to do.
		default:
			// Free list full, just carry on.
		}
	}
}

func load(buffer *Buffer) {
	// loading
}

func process(buffer *Buffer) {
	// processing
}
