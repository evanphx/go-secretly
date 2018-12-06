package vault

import (
	"fmt"
	"strings"

	secretly "github.com/evanphx/go-secretly"
	vault "github.com/hashicorp/vault/api"
)

type v struct{}

func (_ *v) Get(path string) (string, error) {
	parts := strings.SplitN(path, "?", 2)

	slash := strings.IndexByte(path, '/')
	if slash == -1 {
		return "", fmt.Errorf("Invalid secret, must specify mount point")
	}

	mp := parts[0][:slash+1]

	cfg := vault.DefaultConfig()
	c, err := vault.NewClient(cfg)
	if err != nil {
		return "", err
	}

	paths, err := c.Sys().ListMounts()
	if err != nil {
		return "", err
	}

	opts, ok := paths[mp]
	if !ok {
		return "", fmt.Errorf("Invalid secret, unknown mount point")
	}

	// Upgrade this internally to kv version 1
	if opts.Type == "generic" {
		opts.Type = "kv"
	}

	var data map[string]interface{}

	var sec *vault.Secret

	switch opts.Type {
	case "kv":
		var path string

		if opts.Options["version"] != "2" {
			path = parts[0]
		} else {
			path = parts[0][:slash] + "/data/" + parts[0][slash+1:]
		}

		sec, err = c.Logical().Read(path)
		if err != nil {
			return "", err
		}

		if sec == nil {
			return "", secretly.ErrNoSuchSecret
		}

		if len(sec.Data) == 0 {
			return "", secretly.ErrNoSuchSecret
		}

		v, ok := sec.Data["data"]
		if !ok {
			return "", secretly.ErrAmbigiousSecret
		}

		data, ok = v.(map[string]interface{})
		if !ok {
			return "", secretly.ErrAmbigiousSecret
		}

		if len(parts) == 2 {
			v, ok := data[parts[1]]
			if ok {
				return fmt.Sprintf("%v", v), nil
			}

			return "", secretly.ErrNoSuchSecret
		}

		if len(data) == 1 {
			for _, v := range data {
				return fmt.Sprintf("%v", v), nil
			}
		}

		return "", secretly.ErrAmbigiousSecret
	default:
		return "", fmt.Errorf("unsupported secret type: %s", opts.Type)
	}
}

func (_ *v) Put(path, value string) error {
	parts := strings.SplitN(path, "?", 2)

	slash := strings.IndexByte(path, '/')
	if slash == -1 {
		return fmt.Errorf("Invalid secret, must specify mount point")
	}

	mp := parts[0][:slash+1]

	cfg := vault.DefaultConfig()
	c, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	paths, err := c.Sys().ListMounts()
	if err != nil {
		return err
	}

	opts, ok := paths[mp]
	if !ok {
		return fmt.Errorf("Invalid secret, unknown mount point")
	}

	// Upgrade this internally to kv version 1
	if opts.Type == "generic" {
		opts.Type = "kv"
	}

	key := "value"

	if len(parts) == 2 {
		key = parts[1]
	}

	switch opts.Type {
	case "kv":
		var path string

		if opts.Options["version"] != "2" {
			path = parts[0]
		} else {
			path = parts[0][:slash] + "/data/" + parts[0][slash+1:]
		}

		_, err = c.Logical().Write(
			path,
			map[string]interface{}{
				"data": map[string]interface{}{
					key: value,
				},
			})
	default:
		return fmt.Errorf("unsupport secret type: %s", opts.Type)
	}

	return err
}

func init() {
	secretly.AddBackend("vault", &v{})
}
