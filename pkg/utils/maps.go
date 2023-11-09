package utils

func CompareMaps(a map[string]string, b map[string]string) bool {
	/*if len(a) != len(b) {
		return false
	}*/

	for k := range a {
		if v, ok := b[k]; !ok {
			return false
		} else {
			if v != a[k] {
				return false
			}
		}
	}

	return true
}
