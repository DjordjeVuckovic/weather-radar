package util

func IsJSONFormat(output string) bool {
	return len(output) > 0 && output[0] == '{'
}
