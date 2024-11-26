// worker_test.go
package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCopyFiles_Cancellation(t *testing.T) {
	// Configuration du test
	config := &Config{
		SourceDir:   os.TempDir(),
		DestDir:     filepath.Join(os.TempDir(), "dest"),
		ThreadCount: 2,
	}
	os.Mkdir(config.DestDir, 0755)
	defer os.RemoveAll(config.DestDir)

	// Création de fichiers temporaires à copier
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, file := range files {
		path := filepath.Join(config.SourceDir, file)
		os.WriteFile(path, []byte("Contenu"), 0644)
		defer os.Remove(path)
	}

	logger := InitTestLogger()

	// Contexte avec annulation
	ctx, cancel := context.WithCancel(context.Background())

	// Lancer la copie des fichiers
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := CopyFiles(ctx, config, files, logger)
	if err == nil {
		t.Errorf("Une erreur était attendue en raison de l'annulation du contexte")
	}
}
