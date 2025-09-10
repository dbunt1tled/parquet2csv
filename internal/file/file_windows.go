//go:build windows

package file

// На Windows нет `uid/gid`, проверка владельца не имеет смысла.
// Будем считать, что если директория доступна — значит writable.
func isWritableByOwner(path string) (bool, error) {
	return true, nil
}