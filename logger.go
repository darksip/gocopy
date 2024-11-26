// logger.go
package main

import (
	"io"
	"log"
	"os"
)

func InitLogger(logPath string) (*log.Logger, error) {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Création d'un multi-writer pour loguer à la fois dans le fichier et sur la console
	mw := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(mw, "", log.LstdFlags)

	return logger, nil
}
