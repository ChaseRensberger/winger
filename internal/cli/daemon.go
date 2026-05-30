package cli

import (
	"fmt"
	"log"
	"time"

	"winger/internal/config"

	"github.com/fsnotify/fsnotify"
)

func Daemon() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("not initialized: run 'winger init' first")
	}

	planPath, err := config.PlanPath()
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("could not create file watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(planPath); err != nil {
		return fmt.Errorf("could not watch %s: %w", planPath, err)
	}

	fmt.Printf("watching %s for changes (relay: %s, handle: %s)\n", planPath, cfg.RelayURL, cfg.Handle)

	var debounce *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Write) {
				if debounce != nil {
					debounce.Stop()
				}
				debounce = time.AfterFunc(500*time.Millisecond, func() {
					fmt.Printf("change detected, syncing... ")
					if err := Sync(); err != nil {
						log.Printf("sync error: %v", err)
					}
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("watcher error: %v", err)
		}
	}
}
