package configz

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault-client-go"
)

type VaultConfig struct {
	Token   string
	Prefix  string
	Address string
	Timeout time.Duration

	client *vault.Client
}

type IVaultConfig interface {
	LoadConfig(ctx context.Context, path string, out any) error
}

func NewVault(vc *VaultConfig) IVaultConfig {
	v, err := vault.New(
		vault.WithAddress(vc.Address),
		vault.WithRequestTimeout(vc.Timeout),
	)
	if err != nil {
		panic(err)
	}

	if err := v.SetToken(vc.Token); err != nil {
		panic(err)
	}

	vc.client = v

	return vc
}

func (v *VaultConfig) LoadConfig(ctx context.Context, path string, out any) error {
	secret, err := v.client.Read(ctx, path)
	if err != nil {
		return err
	}

	secretData, ok := secret.Data["data"].(map[string]any)
	if !ok {
		return fmt.Errorf("no secret for '%s'", path)
	}

	if err := decode(v.Prefix, secretData, out); err != nil {
		return err
	}

	return nil
}
