//go:build !windows
// +build !windows

package auth

// WinFSATTR is a no-op on non-Windows systems.
//
// On Windows, this function sets the file attribute of the specified file
// to hidden. This is necessary because the configuration file is stored in
// the user's home directory, which is typically not hidden.
//
// This function is a stub on non-Windows systems because there is no
// equivalent concept of hidden files. The function simply returns nil
// without doing anything.
func WinFSATTR(_ string) error {
	return nil
}
