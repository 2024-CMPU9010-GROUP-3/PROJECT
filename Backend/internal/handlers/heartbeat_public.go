//go:build public

package handlers

func init() {
	IsAlive = checkPublic
}

func checkPublic() bool {
	return true
}
