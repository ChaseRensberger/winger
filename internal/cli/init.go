package cli

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"winger/internal/config"
	"winger/internal/identity"
	"winger/internal/plan"
)

func Init() error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("winger is already initialized at %s", dir)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("choose a handle: ")
	handle, _ := reader.ReadString('\n')
	handle = strings.TrimSpace(strings.ToLower(handle))
	if handle == "" {
		return fmt.Errorf("handle cannot be empty")
	}

	fmt.Printf("relay url [%s]: ", config.DefaultRelayURL)
	relayURL, _ := reader.ReadString('\n')
	relayURL = strings.TrimSpace(relayURL)
	if relayURL == "" {
		relayURL = config.DefaultRelayURL
	}

	fmt.Print("generating ed25519 keypair... ")
	pub, priv, err := identity.Generate()
	if err != nil {
		return err
	}
	if err := identity.SaveKeys(pub, priv); err != nil {
		return err
	}
	fmt.Println("done")

	cfg := &config.Config{
		RelayURL: relayURL,
		Handle:   handle,
	}
	if err := config.Save(cfg); err != nil {
		return err
	}

	if err := plan.CreateDefault(handle); err != nil {
		return err
	}

	fmt.Print("registering handle with relay... ")
	if err := registerWithRelay(cfg, pub); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}
	fmt.Println("done")

	planPath, _ := config.PlanPath()
	fmt.Println()
	fmt.Println("winger initialized successfully")
	fmt.Printf("  config:  %s\n", dir)
	fmt.Printf("  plan:    %s\n", planPath)
	fmt.Printf("  handle:  %s\n", handle)
	fmt.Println()
	fmt.Println("edit your ~/.plan and run 'winger sync' to publish it")

	return nil
}

func registerWithRelay(cfg *config.Config, pub []byte) error {
	body := map[string]string{
		"handle": cfg.Handle,
		"pubkey": base64.StdEncoding.EncodeToString(pub),
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post(cfg.RelayURL+"/register", "application/json", strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("could not reach relay: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("handle or public key already registered")
	}
	if resp.StatusCode != http.StatusCreated {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("relay returned %d: %s", resp.StatusCode, errResp["error"])
	}

	return nil
}
