// progress_test.go
package main

import (
	"testing"
	"time"
)

func TestTrackProgress(t *testing.T) {
	totalFiles := 5
	progressCh := make(chan int)

	go func() {
		for i := 0; i < totalFiles; i++ {
			progressCh <- 1
			time.Sleep(10 * time.Millisecond)
		}
		close(progressCh)
	}()

	trackProgress(totalFiles, progressCh)
	// Si la fonction se termine correctement, le test est réussi
}
