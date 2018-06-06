package utils

func Prepend (s string, slice []string) []string {
	new_slice := make([]string, len(slice)+1)
	new_slice[0] = s
	copy(new_slice[1:], slice)

	return new_slice
}
