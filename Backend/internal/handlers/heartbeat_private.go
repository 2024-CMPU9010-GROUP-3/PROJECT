//go:build private

package handlers

func init() {
	IsAlive = checkPrivate
}

func checkPrivate() bool {
	return true
}
