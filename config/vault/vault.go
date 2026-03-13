package vault

import (
	"context"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

const (
	roleIDKey   = "role_id"
	secretIDKey = "secret_id"
)

func GetConfig(
	ctx context.Context,
	options ...Option,
) (vaultConfig map[string]interface{}, err error) {
	var (
		secretID   = os.Getenv("VAULT_SECRET_ID")
		vaultAddr  = os.Getenv("VAULT_ADDR")
		secretPath = os.Getenv("VAULT_PATH")
	)

	vaultConfig = make(map[string]interface{})

	opts := newDefaultOptions()
	for _, o := range options {
		o.apply(opts)
	}

	if err = opts.Validate(
		secretPath,
		vaultAddr,
		secretID,
	); err != nil {
		return nil, fmt.Errorf("failed to validate vault init opts: %s", err.Error())
	}

	data := map[string]interface{}{
		roleIDKey:   opts.roleID,
		secretIDKey: secretID,
	}

	client, err := vaultapi.NewClient(nil)
	if err != nil {
		return vaultConfig, err
	}

	err = client.SetAddress(vaultAddr)
	if err != nil {
		return vaultConfig, err
	}

	resp, err := client.Logical().Write(opts.loginPath, data)
	if err != nil {
		return vaultConfig, err
	}

	token := resp.Auth.ClientToken
	client.SetToken(token)

	secret, err := client.KVv2(opts.mountPath).Get(ctx, opts.secretPath)
	if err != nil {
		return vaultConfig, err
	}

	vaultConfig, ok := secret.Data[opts.cfgKey].(map[string]interface{})
	if !ok {
		return vaultConfig, fmt.Errorf("invalid config format")
	}

	return vaultConfig, nil
}
