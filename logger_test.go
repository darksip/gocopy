// logger_test.go
package main

import (
	"os"
	"testing"
)

func TestInitLogger(t *testing.T) {
	logger, err := InitLogger("test.log")
	if err != nil {
		t.Errorf("Erreur inattendue lors de l'initialisation du logger: %v", err)
	}
	defer os.Remove("test.log")

	if logger == nil {
		t.Errorf("Le logger ne doit pas être nul")
	}

	// Tester l'écriture dans le logger
	logger.Println("Message de test")
}
