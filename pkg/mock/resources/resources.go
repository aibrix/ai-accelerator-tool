package resources

import "embed"

//go:embed injectiond.so
var Resources embed.FS

// GetInjectionLibrary returns the embedded injection library as bytes
func GetInjectionLibrary() ([]byte, error) {
	return Resources.ReadFile("injectiond.so")
}
