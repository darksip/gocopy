// copy_test.go
package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestCopyFile_Success(t *testing.T) {
	sourceFile, err := os.CreateTemp("", "source")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier source temporaire: %v", err)
	}
	defer os.Remove(sourceFile.Name())

	destFile, err := os.CreateTemp("", "dest")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier destination temporaire: %v", err)
	}
	defer os.Remove(destFile.Name())

	// Écriture de données dans le fichier source
	sourceContent := "Contenu du fichier source"
	if _, err := sourceFile.WriteString(sourceContent); err != nil {
		t.Fatalf("Erreur lors de l'écriture du fichier source: %v", err)
	}

	logger := InitTestLogger()

	err = copyFile(sourceFile.Name(), 1, destFile.Name(), logger)
	if err != nil {
		t.Errorf("Erreur inattendue lors de la copie: %v", err)
	}

	// Vérification du contenu du fichier destination
	copiedContent, err := os.ReadFile(destFile.Name())
	if err != nil {
		t.Fatalf("Erreur lors de la lecture du fichier destination: %v", err)
	}

	if string(copiedContent) != sourceContent {
		t.Errorf("Contenu du fichier destination incorrect, attendu %q, obtenu %q", sourceContent, string(copiedContent))
	}
}

func TestCopyFile_SourceNotExist(t *testing.T) {
	destFile, err := os.CreateTemp("", "dest")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier destination temporaire: %v", err)
	}
	defer os.Remove(destFile.Name())

	logger := InitTestLogger()

	err = copyFile("fichier_inexistant.txt", 1, destFile.Name(), logger)
	if err == nil {
		t.Errorf("Une erreur était attendue pour un fichier source inexistant")
	}
}

func TestCopyFile_CopyIgnored(t *testing.T) {
	sourceFile, err := os.CreateTemp("", "source")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier source temporaire: %v", err)
	}
	defer os.Remove(sourceFile.Name())

	destFile, err := os.CreateTemp("", "dest")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier destination temporaire: %v", err)
	}
	defer os.Remove(destFile.Name())

	// Écriture de données identiques dans les deux fichiers
	content := "Contenu identique"
	if _, err := sourceFile.WriteString(content); err != nil {
		t.Fatalf("Erreur lors de l'écriture du fichier source: %v", err)
	}
	if _, err := destFile.WriteString(content); err != nil {
		t.Fatalf("Erreur lors de l'écriture du fichier destination: %v", err)
	}

	logger := InitTestLogger()

	err = copyFile(sourceFile.Name(), 1, destFile.Name(), logger)
	if err != ErrCopyIgnored {
		t.Errorf("Erreur attendue ErrCopyIgnored, obtenue: %v", err)
	}
}

func TestFilesAreEqual(t *testing.T) {
	// Création de deux fichiers identiques
	file1, err := os.CreateTemp("", "file1")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier temporaire: %v", err)
	}
	defer os.Remove(file1.Name())

	file2, err := os.CreateTemp("", "file2")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier temporaire: %v", err)
	}
	defer os.Remove(file2.Name())

	content := "Contenu identique"
	if _, err := file1.WriteString(content); err != nil {
		t.Fatalf("Erreur lors de l'écriture du fichier 1: %v", err)
	}
	if _, err := file2.WriteString(content); err != nil {
		t.Fatalf("Erreur lors de l'écriture du fichier 2: %v", err)
	}

	equal, err := filesAreEqual(file1.Name(), file2.Name())
	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if !equal {
		t.Errorf("Les fichiers devraient être égaux")
	}

	// Modifier le contenu d'un fichier
	if _, err := file2.WriteString("Changement"); err != nil {
		t.Fatalf("Erreur lors de la modification du fichier 2: %v", err)
	}

	equal, err = filesAreEqual(file1.Name(), file2.Name())
	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if equal {
		t.Errorf("Les fichiers ne devraient pas être égaux")
	}
}

func TestFileHash(t *testing.T) {
	file, err := os.CreateTemp("", "file")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier temporaire: %v", err)
	}
	defer os.Remove(file.Name())

	content := "Contenu pour le hash"
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("Erreur lors de l'écriture du fichier: %v", err)
	}

	hash, err := fileHash(file.Name())
	if err != nil {
		t.Fatalf("Erreur lors du calcul du hash: %v", err)
	}

	// Calculer le hash attendu
	expectedHash := computeMD5(content)
	if hash != expectedHash {
		t.Errorf("Hash incorrect, attendu %s, obtenu %s", expectedHash, hash)
	}
}

// Fonction auxiliaire pour initialiser un logger pour les tests
func InitTestLogger() *log.Logger {
	return log.New(io.Discard, "", log.LstdFlags)
}

// Fonction auxiliaire pour calculer le MD5 d'une chaîne
func computeMD5(content string) string {
	hasher := md5.New()
	hasher.Write([]byte(content))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
