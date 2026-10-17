package httpUtil

import (
	"sync"
	"testing"

	bunker "github.com/olympus-forge/bunker/errors"
)

func TestNewAppError_DoesNotMutateGlobalSingleton(t *testing.T) {
	origMessage := bunker.Forbidden.Message

	customMessage := "pesan error kustom untuk pengujian"
	appErr := NewAppError(bunker.StatusForbidden, customMessage)

	if appErr.Message != customMessage {
		t.Fatalf("diharapkan pesan %q, didapat %q", customMessage, appErr.Message)
	}

	if bunker.Forbidden.Message != origMessage {
		t.Fatalf("bunker.Forbidden global termutasi! semula %q, sekarang %q", origMessage, bunker.Forbidden.Message)
	}
}

func TestNewAppError_ConcurrentSafe(t *testing.T) {
	const workers = 50
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := NewAppError(bunker.StatusForbidden, "pesan dari goroutine")
			if err.Message != "pesan dari goroutine" {
				t.Errorf("pesan tidak sesuai")
			}
		}(i)
	}

	wg.Wait()
}
