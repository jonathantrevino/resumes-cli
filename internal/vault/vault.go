package vault

import (
	"errors"
	"strings"
	"unicode"
)

func Slug(s string) (string, error) {

	// lowercase everything
	s = strings.ToLower(s)

	trimmed := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	joined := strings.Join(trimmed, "-")
	s = joined

	if s == "" {
		return "", errors.New("String is empty")
	}

	return s, nil
}
