package main

import (
	"bufio" // Import manquant pour lire la liste des fichiers
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

const maxRetries = 3  // Nombre maximum de tentatives en cas d'échec de copie
const workerCount = 3 // Nombre de workers pour la copie parallèle

func main() {
	// Chargement des variables d'environnement depuis le fichier .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env: %v", err)
	}

	// Récupération des variables d'environnement
	sourceDir := os.Getenv("SOURCE_DIR")
	destDir := os.Getenv("DEST_DIR")
	filesListPath := os.Getenv("FILES_LIST_PATH")

	// Vérification que les variables d'environnement sont définies
	if sourceDir == "" || destDir == "" || filesListPath == "" {
		log.Fatal("Les variables d'environnement SOURCE_DIR, DEST_DIR et FILES_LIST_PATH doivent être définies.")
	}

	// Lecture de la liste des fichiers à copier
	files, err := readFilesList(filesListPath)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de la liste des fichiers: %v", err)
	}

	// Initialisation du fichier de log
	logFile, err := os.OpenFile("copy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier de log: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	startTime := time.Now()

	var wg sync.WaitGroup
	fileCh := make(chan string, len(files))
	progressCh := make(chan int, len(files))
	errorCh := make(chan error) // Canal pour collecter les erreurs non fatales

	// Remplir le canal avec les fichiers à copier
	for _, file := range files {
		fileCh <- file
	}
	close(fileCh)

	// Lancer les workers pour la copie parallèle
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, sourceDir, destDir, fileCh, progressCh, errorCh, &wg, logger)
	}

	// Goroutine pour suivre la progression de la copie
	go func() {
		totalFiles := len(files)
		copiedFiles := 0
		for range progressCh {
			copiedFiles++
			duration := time.Since(startTime)
			remaining := time.Duration(float64(duration) / float64(copiedFiles) * float64(totalFiles-copiedFiles))
			fmt.Printf("Progression: %d/%d, Temps restant estimé: %v\n", copiedFiles, totalFiles, remaining)
		}
	}()

	// Goroutine pour gérer les erreurs non fatales
	go func() {
		for err := range errorCh {
			fmt.Printf("Erreur : %v\n", err)
		}
	}()

	// Attendre que tous les workers aient terminé
	wg.Wait()
	close(progressCh)
	close(errorCh)

	fmt.Println("Copie terminée.")
}

// Fonction worker pour copier les fichiers en parallèle
func worker(id int, sourceDir, destDir string, fileCh <-chan string, progressCh chan<- int, errorCh chan<- error, wg *sync.WaitGroup, logger *log.Logger) {
	defer wg.Done()
	for file := range fileCh {
		sourcePath := filepath.Join(sourceDir, file)
		destPath := filepath.Join(destDir, file)

		retries := 0
		for {
			err := copyFile(sourcePath, destPath)
			if err == nil {
				logger.Printf("Worker %d: Copie réussie de %s vers %s\n", id, sourcePath, destPath)
				progressCh <- 1
				break
			}

			retries++
			if retries >= maxRetries {
				errMsg := fmt.Errorf("Worker %d: Échec de la copie de %s après %d tentatives: %v", id, sourcePath, retries, err)
				logger.Println(errMsg)
				errorCh <- errMsg
				break
			}

			logger.Printf("Worker %d: Erreur lors de la copie de %s, nouvelle tentative (%d/%d)\n", id, sourcePath, retries, maxRetries)
			time.Sleep(2 * time.Second) // Attente avant de retenter
		}
	}
}

// Fonction pour copier un fichier du chemin source vers le chemin destination
func copyFile(source, dest string) error {
	// Créer les répertoires parents du fichier de destination si nécessaire
	err := os.MkdirAll(filepath.Dir(dest), os.ModePerm)
	if err != nil {
		return fmt.Errorf("impossible de créer les répertoires de destination: %w", err)
	}

	// Ouvrir le fichier source en lecture
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("impossible d'ouvrir le fichier source: %w", err)
	}
	defer sourceFile.Close()

	// Créer le fichier de destination en écriture
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier de destination: %w", err)
	}
	defer destFile.Close()

	// Copier le contenu du fichier source vers le fichier de destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("erreur lors de la copie: %w", err)
	}

	// Copier les permissions du fichier source vers le fichier de destination
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("impossible d'obtenir les informations du fichier source: %w", err)
	}
	err = os.Chmod(dest, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("impossible de définir les permissions du fichier de destination: %w", err)
	}

	return nil
}

// Fonction pour lire la liste des fichiers à copier depuis un fichier texte
func readFilesList(filePath string) ([]string, error) {
	// Ouvrir le fichier contenant la liste des fichiers
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir la liste des fichiers: %w", err)
	}
	defer file.Close()

	var files []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la liste des fichiers: %w", err)
	}

	return files, nil
}
