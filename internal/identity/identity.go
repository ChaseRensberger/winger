package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"winger/internal/config"
)

const (
	PrivateKeyFile = "identity.key"
	PublicKeyFile  = "identity.pub"
)

func Generate() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate keypair: %w", err)
	}
	return pub, priv, nil
}

func SaveKeys(pub ed25519.PublicKey, priv ed25519.PrivateKey) error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: priv.Seed(),
	})
	if err := os.WriteFile(filepath.Join(dir, PrivateKeyFile), privPEM, 0600); err != nil {
		return fmt.Errorf("could not write private key: %w", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "ED25519 PUBLIC KEY",
		Bytes: pub,
	})
	if err := os.WriteFile(filepath.Join(dir, PublicKeyFile), pubPEM, 0644); err != nil {
		return fmt.Errorf("could not write public key: %w", err)
	}

	return nil
}

func LoadPrivateKey() (ed25519.PrivateKey, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(dir, PrivateKeyFile))
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "ED25519 PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key file")
	}

	return ed25519.NewKeyFromSeed(block.Bytes), nil
}

func LoadPublicKey() (ed25519.PublicKey, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(dir, PublicKeyFile))
	if err != nil {
		return nil, fmt.Errorf("could not read public key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "ED25519 PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key file")
	}

	return ed25519.PublicKey(block.Bytes), nil
}

func Sign(priv ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(priv, message)
}

func Verify(pub ed25519.PublicKey, message, sig []byte) bool {
	return ed25519.Verify(pub, message, sig)
}
