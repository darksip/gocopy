// filelist_test.go
package main

import (
	"os"
	"reflect"
	"testing"
)

func TestReadFilesList(t *testing.T) {
	// Création d'un fichier temporaire avec du contenu
	tempFile, err := os.CreateTemp("", "filelist")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier temporaire: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := "file1.txt\nfile2.txt\n\n# Commentaire\nfile3.txt\n"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Erreur lors de l'écriture du fichier temporaire: %v", err)
	}

	files, err := ReadFilesList(tempFile.Name())
	if err != nil {
		t.Errorf("Erreur inattendue: %v", err)
	}

	expected := []string{"file1.txt", "file2.txt", "# Commentaire", "file3.txt"}
	if !reflect.DeepEqual(files, expected) {
		t.Errorf("Résultat attendu %v, obtenu %v", expected, files)
	}
}

func TestReadFilesList_FileNotFound(t *testing.T) {
	_, err := ReadFilesList("fichier_inexistant.txt")
	if err == nil {
		t.Errorf("Une erreur était attendue pour un fichier inexistant")
	}
}
