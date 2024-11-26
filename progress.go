// progress.go
package main

import (
	"fmt"
	"strings"
	"time"
)

func trackProgress(totalFiles int, progressCh <-chan int) {
	startTime := time.Now()
	copiedFiles := 0
	for range progressCh {
		copiedFiles++
		duration := time.Since(startTime)
		remaining := time.Duration(float64(duration) / float64(copiedFiles) * float64(totalFiles-copiedFiles))

		// Calculer le pourcentage de progression
		percent := float64(copiedFiles) / float64(totalFiles) * 100

		// Créer une barre de progression simple
		width := 50
		completed := int(float64(width) * float64(copiedFiles) / float64(totalFiles))
		bar := strings.Repeat("=", completed) + strings.Repeat("-", width-completed)

		// Afficher la barre de progression et les informations sur la même ligne
		fmt.Printf("\r[%s] %.2f%% (%d/%d) Temps restant estimé: %v",
			bar, percent, copiedFiles, totalFiles, remaining)
	}
	// Ajouter une nouvelle ligne à la fin pour ne pas écraser la dernière mise à jour
	fmt.Println()
}
