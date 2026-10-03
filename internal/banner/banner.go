// Package banner holds the program's banner, embedded from a file.
package banner

import _ "embed"

//go:embed banner.txt
var text string

// Text returns the banner.
func Text() string {
	return text
}
