// config_test.go
package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Sauvegarder les variables d'environnement actuelles
	originalEnv := os.Environ()
	defer func() {
		// Restaurer les variables d'environnement
		os.Clearenv()
		for _, e := range originalEnv {
			kv := splitEnv(e)
			os.Setenv(kv[0], kv[1])
		}
	}()

	// Cas de test : toutes les variables sont définies correctement
	os.Setenv("SOURCE_DIR", "/source")
	os.Setenv("DEST_DIR", "/dest")
	os.Setenv("FILES_LIST_PATH", "/files.txt")
	os.Setenv("THREAD_COUNT", "4")

	config, err := LoadConfig()
	if err != nil {
		t.Errorf("Erreur inattendue lors du chargement de la configuration: %v", err)
	}

	if config.SourceDir != "/source" || config.DestDir != "/dest" || config.FilesListPath != "/files.txt" || config.ThreadCount != 4 {
		t.Errorf("Configuration incorrecte: %+v", config)
	}

	// Cas de test : variable THREAD_COUNT invalide
	os.Setenv("THREAD_COUNT", "-1")
	_, err = LoadConfig()
	if err == nil {
		t.Errorf("Erreur attendue lors du chargement de la configuration avec THREAD_COUNT invalide")
	}

	// Cas de test : variables manquantes
	os.Clearenv()
	os.Setenv("SOURCE_DIR", "/source")
	_, err = LoadConfig()
	if err == nil {
		t.Errorf("Erreur attendue lors du chargement de la configuration avec des variables manquantes")
	}
}

// Fonction auxiliaire pour diviser les variables d'environnement
func splitEnv(e string) [2]string {
	for i := 0; i < len(e); i++ {
		if e[i] == '=' {
			return [2]string{e[:i], e[i+1:]}
		}
	}
	return [2]string{e, ""}
}
