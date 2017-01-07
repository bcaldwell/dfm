package dfm

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func uniqueSliceTransform(a []string) (output []string) {

	for _, s := range a {
		output = appendIfUnique(output, s)
	}

	return output
}

func appendIfUnique(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
