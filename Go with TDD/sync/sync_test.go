package sync

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {

	t.Run("it also runs safely in a testing enviornment ", func(t *testing.T) {

		wantedCount := 1000
		counter := Counter{}

		var wg sync.WaitGroup
		wg.Add(wantedCount)

		// new way to write loops in go basically, elimitte this bs - 	for i := 0; i < wantedCount; i++ {
		for range wantedCount {

			go func() {
				counter.Inc()
				wg.Done()
			}()
		}
		wg.Wait() // this one waits out till it can to make sure all processes end before the loop retrns anything

		assertCounter(t, &counter, wantedCount)
	})
}

func assertCounter(t testing.TB, got *Counter, want int) {
	t.Helper()

	if got.Value() != want {
		t.Errorf("got %d, want %d", got.Value(), 3)
	}

}
