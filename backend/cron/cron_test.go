package cron

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestCronJobExecution(t *testing.T) {
	c := cron.New()

	executed := make(chan bool, 1)


	_, err := c.AddFunc("@every 1s", func() {
		executed <- true
	})

	if err != nil {
		t.Fatalf("L'Erreur lors de l'ajout de la fonction: %v", err)
	}
	c.Start()
	defer c.Stop()

	select {
	case <-executed:
	
	case <-time.After(2 * time.Second):
		t.Fatal("La tâche cron n'a pas été executé")	
	}
}