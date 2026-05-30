package cli

import (
	"fmt"
	"io"
	"net/http"

	"winger/internal/config"
)

func Lookup(handle string) error {
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{RelayURL: config.DefaultRelayURL}
	}

	req, err := http.NewRequest("GET", cfg.RelayURL+"/"+handle, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not reach relay: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no plan found for %s", handle)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("relay returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Print(string(body))
	return nil
}
