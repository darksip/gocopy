// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Charger et valider la configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Erreur de configuration: %v", err)
	}

	// Initialiser le logger
	logger, err := InitLogger("copy.log")
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du logger: %v", err)
	}

	// Lire la liste des fichiers à copier
	files, err := ReadFilesList(config.FilesListPath)
	if err != nil {
		logger.Fatalf("Erreur lors de la lecture de la liste des fichiers: %v", err)
	}

	// Contexte pour la gestion des interruptions
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture des signaux d'interruption
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nInterruption détectée, arrêt du programme...")
		cancel()
	}()

	// Lancer la copie des fichiers
	startTime := time.Now()
	err = CopyFiles(ctx, config, files, logger)
	if err != nil {
		logger.Fatalf("Erreur lors de la copie des fichiers: %v", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Copie terminée en %v.\n", duration)
}
