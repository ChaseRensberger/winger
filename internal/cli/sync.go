package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"winger/internal/config"
	"winger/internal/identity"
	"winger/internal/plan"
	"winger/internal/relay"
)

func Sync() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("not initialized: run 'winger init' first")
	}

	priv, err := identity.LoadPrivateKey()
	if err != nil {
		return err
	}

	pub, err := identity.LoadPublicKey()
	if err != nil {
		return err
	}

	content, err := plan.Read()
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	message := fmt.Sprintf("%s:%s:%d", cfg.Handle, content, timestamp)
	sig := identity.Sign(priv, []byte(message))

	req := relay.SyncRequest{
		Handle:    cfg.Handle,
		Content:   content,
		Signature: base64.StdEncoding.EncodeToString(sig),
		Timestamp: timestamp,
		Pubkey:    base64.StdEncoding.EncodeToString(pub),
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(cfg.RelayURL+"/sync", "application/json", strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("could not reach relay: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("sync failed: %s", errResp["error"])
	}

	fmt.Printf("synced %s to %s\n", cfg.Handle, cfg.RelayURL)
	return nil
}
