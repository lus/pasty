package maps

func ExceedsDimensions(src map[string]any, width, depth int) bool {
	if width < 0 || depth < 1 || len(src) > width {
		return true
	}

	for _, value := range src {
		childMap, ok := value.(map[string]any)
		if !ok {
			continue
		}

		if depth == 1 {
			return true
		}

		if ExceedsDimensions(childMap, width, depth-1) {
			return true
		}
	}

	return false
}
