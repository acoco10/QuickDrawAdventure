package utilityfunctions

import ()

func StringInSlice(testString string, testList []string) bool {
	for _, b := range testList {
		if b == testString {
			return true
		}
	}
	return false
}

// return slice of keys for maps with string keys and value ints
func MapKeysSI(dict map[string]int) []string {

	i := 0

	keys := make([]string, len(dict))

	for key := range dict {
		keys[i] = key
		i++
	}

	return keys
}

// return slice of keys for maps with string keys and skill values
func MapKeysSSk(dict map[string]battleStats.Skill) []string {
	i := 0

	keys := make([]string, len(dict))

	for key := range dict {
		keys[i] = key
		i++
	}

	return keys
}
