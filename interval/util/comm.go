package util

import (
	"github.com/pkg/errors"
	"muzzammil.xyz/jsonc"
)

func JsonStrRemoveComments(str string) (string, error) {
	jc := jsonc.ToJSON([]byte(str))
	if jsonc.Valid(jc) {
		return string(jc), nil
	}
	return "", errors.New("Invalid JSON")
}
