// worker.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const maxRetries = 3 // Nombre maximum de tentatives en cas d'échec de copie

func CopyFiles(ctx context.Context, config *Config, files []string, logger *log.Logger) error {
	fileCh := make(chan string)
	progressCh := make(chan int)
	errorCh := make(chan error)
	doneCh := ctx.Done()

	// Lancer les workers
	for i := 0; i < config.ThreadCount; i++ {
		go worker(i, config.SourceDir, config.DestDir, fileCh, progressCh, errorCh, doneCh, logger)
	}

	// Envoi des fichiers à copier
	go func() {
		defer close(fileCh)
		for _, file := range files {
			select {
			case <-doneCh:
				return
			case fileCh <- file:
			}
		}
	}()

	// Suivi de la progression
	go trackProgress(len(files), progressCh, logger)

	// Gestion des erreurs
	var copyErr error
	for {
		select {
		case <-doneCh:
			return fmt.Errorf("Copie interrompue")
		case err := <-errorCh:
			logger.Println(err)
			copyErr = err
		default:
			if copyErr != nil {
				return copyErr
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func worker(id int, sourceDir, destDir string, fileCh <-chan string, progressCh chan<- int, errorCh chan<- error, doneCh <-chan struct{}, logger *log.Logger) {
	for {
		select {
		case <-doneCh:
			logger.Printf("Worker %d: Arrêté suite à une interruption\n", id)
			return
		case file, ok := <-fileCh:
			if !ok {
				return
			}

			sourcePath := filepath.Join(sourceDir, file)
			destPath := filepath.Join(destDir, file)

			retries := 0
			for {
				err := copyFile(sourcePath, id, destPath, logger)
				if err == nil {
					progressCh <- 1
					break
				}
				// Gestion du cas de copie ignorée sans retry
				if err == ErrCopyIgnored {
					progressCh <- 1
					break
				}
				// Gestion de la source manquante sans retry
				if os.IsNotExist(err) {
					errMsg := fmt.Errorf("Worker %d: Fichier source manquant %s", id, sourcePath)
					errorCh <- errMsg
					break
				}
				// Gestion des tentatives en cas d'échec
				retries++
				if retries >= maxRetries {
					errMsg := fmt.Errorf("Worker %d: Échec de la copie de %s après %d tentatives: %v", id, sourcePath, retries, err)
					errorCh <- errMsg
					break
				}

				// Nouvelle tentative après attente
				logger.Printf("Worker %d: Erreur lors de la copie de %s, nouvelle tentative (%d/%d)\n", id, sourcePath, retries, maxRetries)
				time.Sleep(2 * time.Second)
			}
		}
	}
}
