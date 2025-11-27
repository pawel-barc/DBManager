package config

import (
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Ce fichier initialise un système de logs avec rotation automatique créant un logger global qui enregistre les événements dans un fichier app.log
var AppLogger *log.Logger

func InitLogger() {
	os.MkdirAll("logs", 0755)

	logFile := &lumberjack.Logger{
		Filename: "logs/app.log",
		MaxSize: 10,
		MaxBackups: 5,
		MaxAge: 30,
		Compress: true,
	}

	AppLogger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}