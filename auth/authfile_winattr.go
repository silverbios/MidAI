//go:build windows
// +build windows

package auth

import (
	"golang.org/x/sys/windows"
)

// WinFSATTR sets a file's attribute to hidden on Windows.
// This is necessary because the configuration file is stored in the user's home
// directory, which is typically not hidden.
//
// This function is a stub on non-Windows systems because there is no equivalent
// concept of hidden files.
func WinFSATTR(ConfigFile string) error {
	// Convert the ConfigFile string to a UTF-16 pointer
	ptr, err := windows.UTF16PtrFromString(ConfigFile)
	if err != nil {
		return err
	}

	// Get the current file attributes
	attributes, err := windows.GetFileAttributes(ptr)
	if err != nil {
		return err
	}

	// Check if the file already has the hidden attribute
	// If it does, don't set it again
	if attributes&windows.FILE_ATTRIBUTE_HIDDEN != 0 {
		return nil
	}

	// Set the file attribute to hidden
	// This will make the file invisible in the file explorer
	err = windows.SetFileAttributes(ptr, attributes|windows.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		return err
	}

	return nil
}
