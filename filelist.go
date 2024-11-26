// filelist.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFilesList(filePath string) ([]string, error) {
	// Ouvrir le fichier contenant la liste des fichiers
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Impossible d'ouvrir la liste des fichiers: %w", err)
	}
	defer file.Close()

	var files []string
	// Scanner chaque ligne du fichier pour obtenir les chemins des fichiers
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			files = append(files, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture de la liste des fichiers: %w", err)
	}

	return files, nil
}
