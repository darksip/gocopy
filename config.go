// config.go
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SourceDir     string
	DestDir       string
	FilesListPath string
	ThreadCount   int
}

func LoadConfig() (*Config, error) {
	// Charger les variables d'environnement depuis le fichier .env
	if err := godotenv.Load(".env"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("erreur lors du chargement du fichier .env: %v", err)
		}
	}

	// Récupérer les variables d'environnement
	sourceDir := os.Getenv("SOURCE_DIR")
	destDir := os.Getenv("DEST_DIR")
	filesListPath := os.Getenv("FILES_LIST_PATH")
	threadCountStr := os.Getenv("THREAD_COUNT")

	// Lecture du chemin du fichier de liste à partir de la ligne de commande si présent
	if len(os.Args) > 1 {
		filesListPath = os.Args[1]
	}

	// Valider les variables d'environnement
	missingVars := []string{}
	if sourceDir == "" {
		missingVars = append(missingVars, "SOURCE_DIR")
	}
	if destDir == "" {
		missingVars = append(missingVars, "DEST_DIR")
	}
	if filesListPath == "" {
		missingVars = append(missingVars, "FILES_LIST_PATH")
	}
	if threadCountStr == "" {
		missingVars = append(missingVars, "THREAD_COUNT")
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("les variables d'environnement suivantes sont manquantes: %v", missingVars)
	}

	// Conversion de THREAD_COUNT en entier
	threadCount, err := strconv.Atoi(threadCountStr)
	if err != nil || threadCount <= 0 {
		return nil, fmt.Errorf("THREAD_COUNT doit être un entier positif")
	}

	return &Config{
		SourceDir:     sourceDir,
		DestDir:       destDir,
		FilesListPath: filesListPath,
		ThreadCount:   threadCount,
	}, nil
}
