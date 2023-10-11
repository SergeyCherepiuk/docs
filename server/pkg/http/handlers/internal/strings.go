package internal

import "strings"

func ToSentence(s string) string {
	if len(s) <= 0 {
		return s
	}

	firstLetter := strings.ToUpper(string(s[0]))
	return strings.Join([]string{firstLetter, s[1:]}, "")
}
