package uidmap

func validateUIDMap(m map[string]interface{}) bool {
	if len(m) == 0 {
		return false
	}

	for user := range m {
		if user == "<no value>" || user == "" {
			return false
		}
	}
	return true
}
