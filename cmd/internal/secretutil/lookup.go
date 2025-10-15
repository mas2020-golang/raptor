package secretutil

import (
	"log/slog"
	"strings"

	"github.com/mas2020-golang/cryptex/packages/utils"
)

// LookupResult describes the outcome of a secret lookup.
type LookupResult struct {
	Secret *utils.Secret
	Value  string
	Item   string
}

// Lookup resolves the given secret reference (e.g. "foo" or "foo.item") within the
// provided box name. It returns both the secret pointer and the resolved value
// (password or item content) if present.
func Lookup(boxName, name string) (*LookupResult, *utils.Box, error) {
	_, _, box, err := utils.OpenBox(boxName, "")
	if err != nil {
		return nil, nil, err
	}

	secret, value, item := findSecretValue(name, box)
	return &LookupResult{
		Secret: secret,
		Value:  value,
		Item:   item,
	}, box, nil
}

func findSecretValue(name string, box *utils.Box) (secret *utils.Secret, value, item string) {
	var (
		secretName string
	)

	elems := strings.Split(name, ".")
	secretName = elems[0]
	if len(elems) > 1 {
		item = strings.Join(elems[1:], "")
	}

	slog.Debug("secretutil.findSecretValue()", "secretName", secretName, "secretItem", item)

	if box == nil || box.Secrets == nil {
		return nil, "", item
	}

	for _, s := range box.Secrets {
		if secretName != s.Name {
			continue
		}
		// Keep pointer to the matching secret regardless of item status.
		secret = s

		if len(item) > 0 {
			if len(s.Others) == 0 {
				continue
			}
			if v, ok := s.Others[item]; ok {
				value = v
			}
			continue
		}

		value = s.Pwd
	}

	return secret, value, item
}
