package utils

import "encoding/json"

// verifica se é um json válido
func IsJson(s string) error {
	var js struct{}

	if err := json.Unmarshal([]byte(s), &js); err != nil {
		return err
	}

	return nil
}
