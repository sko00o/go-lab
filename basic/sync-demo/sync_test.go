package sync

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	t.Run("incrementation the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := Counter{}
		counter.Inc()
		counter.Inc()
		counter.Inc()

		assertCounter(t, counter, 3)
	})

	t.Run("it run safely concurrently", func(t *testing.T) {
		want := 1000
		counter := Counter{}

		var wg sync.WaitGroup
		wg.Add(want)

		go func(w *sync.WaitGroup) {
			counter.Inc()
			w.Done()
		}(&wg)
		wg.Wait()

		assertCounter(t, counter, want)

	})
}

func assertCounter(t *testing.T, got Counter, want int) {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got %d, want %d", got.Value(), want)
	}
}
