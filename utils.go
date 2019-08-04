package gosql

type stringlist []string

func (list stringlist) contains(s string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == s {
			return true
		}
	}
	return false
}
