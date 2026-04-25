//go:build !windows

package anonymizer

func isReparsePoint(path string) (bool, error) {
	return false, nil
}
