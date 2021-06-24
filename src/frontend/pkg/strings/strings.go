package strings

func Substr(s string, from, length int) string {
	runes := []rune(s)
	to := from + length

	if to > len(runes) {
		to = len(runes)
	}

	if from > len(runes) {
		from = len(runes)
	}

	return string(runes[from:to])
}
