package vault

import (
	"errors"
)

type vaultOptions struct {
	mountPath  string
	secretPath string
	cfgKey     string
	loginPath  string
	roleID     string
}

func (opts *vaultOptions) Validate(
	secretPath,
	vaultAddr,
	secretID string,
) error {
	if opts.roleID == "" {
		return errors.New("role id were not set, use WithRoleID option to set it")
	}

	if opts.secretPath == "" {
		if secretPath == "" {
			return errors.New("secret path were not set, use WithSecretPath option or set it to the environment variables")
		}

		opts.secretPath = secretPath
	}

	if secretID == "" {
		return errors.New("secret id were not set, set it to the environment variables, check k8s secret")
	}
	if vaultAddr == "" {
		return errors.New("vault address were not set, set it to the environment variables")
	}

	return nil
}

type Option interface {
	apply(*vaultOptions)
}

type mountPathOption struct {
	path string
}

func (o *mountPathOption) apply(opts *vaultOptions) {
	opts.mountPath = o.path
}

func WithMountPath(path string) Option {
	return &mountPathOption{path: path}
}

type cfgPathOption struct {
	secretPath string
}

func (o *cfgPathOption) apply(opts *vaultOptions) {
	opts.secretPath = o.secretPath
}

func WithSecretPath(cfgPath string) Option {
	return &cfgPathOption{secretPath: cfgPath}
}

type loginPathOption struct {
	loginPath string
}

func (o *loginPathOption) apply(opts *vaultOptions) {
	opts.loginPath = o.loginPath
}

func WithLoginPath(loginPath string) Option {
	return &loginPathOption{loginPath: loginPath}
}

type roleIdOption struct {
	roleID string
}

func (o roleIdOption) apply(opts *vaultOptions) {
	opts.roleID = o.roleID
}

func WithVaultConfigKey(cfgKey string) Option {
	return &cfgKeyOption{cfgKey: cfgKey}
}

type cfgKeyOption struct {
	cfgKey string
}

func (o cfgKeyOption) apply(opts *vaultOptions) {
	opts.cfgKey = o.cfgKey
}

func WithRoleId(roleID string) Option {
	return &roleIdOption{roleID: roleID}
}

const (
	defaultMountPoint = "ap"
	defaultConfigPath = "cfg"
	defaultLoginPath  = "auth/approle/login"
)

func newDefaultOptions() *vaultOptions {
	return &vaultOptions{
		mountPath: defaultMountPoint,
		cfgKey:    defaultConfigPath,
		loginPath: defaultLoginPath,
	}
}
