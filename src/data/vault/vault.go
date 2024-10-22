// vault.go
package vault

import (
    "encoding/base64"
    "fmt"
    "github.com/hashicorp/vault/api"
)

type VaultClient struct {
    client *api.Client
}

// NewVaultClient creates a new Vault client
func NewVaultClient(address, token string) (*VaultClient, error) {
    config := api.DefaultConfig()
    config.Address = address
    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("error creating Vault client: %v", err)
    }

    client.SetToken(token)
    return &VaultClient{client: client}, nil
}

// Encrypt data using Vault
func (vc *VaultClient) Encrypt(plaintext string, key string) (string, error) {
    encryptPath := fmt.Sprintf("transit/encrypt/%s", key)
    encryptData := map[string]interface{}{
        "plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
    }
    encryptSecret, err := vc.client.Logical().Write(encryptPath, encryptData)
    if err != nil {
        return "", fmt.Errorf("error encrypting data: %v", err)
    }

    ciphertext, ok := encryptSecret.Data["ciphertext"].(string)
    if !ok {
        return "", fmt.Errorf("ciphertext not found in Vault response")
    }
    return ciphertext, nil
}

// Decrypt data using Vault
func (vc *VaultClient) Decrypt(ciphertext string, key string) (string, error) {
    decryptPath := fmt.Sprintf("transit/decrypt/%s", key)
    decryptData := map[string]interface{}{
        "ciphertext": ciphertext,
    }
    decryptSecret, err := vc.client.Logical().Write(decryptPath, decryptData)
    if err != nil {
        return "", fmt.Errorf("error decrypting data: %v", err)
    }

    decodedBase64, ok := decryptSecret.Data["plaintext"].(string)
    if !ok {
        return "", fmt.Errorf("plaintext not found in Vault response")
    }

    decodedBytes, err := base64.StdEncoding.DecodeString(decodedBase64)
    if err != nil {
        return "", fmt.Errorf("error decoding base64: %v", err)
    }

    return string(decodedBytes), nil
}
