package utils

import(
	"safebase/config"
)

// Ce fichier fournit des fonctions utilitaires pour enregistrer des messages d'infos et d'erreur
func LogInfo(msg string) {
	config.AppLogger.Println("INFO:", msg)
}

func LogError(msg string, err error) {
	if err != nil {
		config.AppLogger.Println("ERROR:", msg, "|", err.Error())
	} else {
		config.AppLogger.Println("ERROR:", msg)
	}
}