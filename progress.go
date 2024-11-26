// progress.go
package main

import (
	"fmt"
	"log"
	"time"
)

func trackProgress(totalFiles int, progressCh <-chan int, logger *log.Logger) {
	startTime := time.Now()
	copiedFiles := 0
	for range progressCh {
		copiedFiles++
		duration := time.Since(startTime)
		remaining := time.Duration(float64(duration) / float64(copiedFiles) * float64(totalFiles-copiedFiles))
		fmt.Printf("Progression: %d/%d, Temps restant estimé: %v\n", copiedFiles, totalFiles, remaining)
	}
}
