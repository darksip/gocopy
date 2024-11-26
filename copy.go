// copy.go
package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var ErrCopyIgnored = errors.New("copie ignorée")

func copyFile(source string, id int, dest string, logger *log.Logger) error {
	// Vérifier si le fichier source existe
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Créer les répertoires parents du fichier de destination si nécessaire
	err = os.MkdirAll(filepath.Dir(dest), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Impossible de créer les répertoires de destination: %w", err)
	}

	// Vérifier si le fichier de destination existe
	_, err = os.Stat(dest)
	if err == nil {
		// Comparer les hashs des fichiers pour déterminer s'ils sont identiques
		same, err := filesAreEqual(source, dest)
		if err != nil {
			return err
		}
		if same {
			// Log et affichage en cas de copie ignorée
			msg := fmt.Sprintf("Worker %d: Copie ignorée pour %s: fichiers identiques", id, filepath.Base(source))
			logger.Println(msg)
			fmt.Printf("\033[33m⚠\033[0m %s\n", msg) // Pictogramme jaune pour signaler l'ignorance
			return ErrCopyIgnored
		}
	}

	// Ouvrir le fichier source en lecture
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Créer le fichier de destination en écriture
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Impossible de créer le fichier de destination: %w", err)
	}
	defer destFile.Close()

	// Copier le contenu du fichier source vers le fichier de destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("Erreur lors de la copie: %w", err)
	}

	// Copier les permissions du fichier source vers le fichier de destination
	err = os.Chmod(dest, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("Impossible de définir les permissions du fichier de destination: %w", err)
	}

	// Copier les dates d'accès et de modification du fichier source vers le fichier de destination
	err = os.Chtimes(dest, sourceInfo.ModTime(), sourceInfo.ModTime())
	if err != nil {
		return fmt.Errorf("Impossible de définir les dates du fichier de destination: %w", err)
	}

	return nil
}

func filesAreEqual(file1, file2 string) (bool, error) {
	hash1, err := fileHash(file1)
	if err != nil {
		return false, err
	}
	hash2, err := fileHash(file2)
	if err != nil {
		return false, err
	}
	return hash1 == hash2, nil
}

func fileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
