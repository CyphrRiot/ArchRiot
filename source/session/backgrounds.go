package session

import "archriot-installer/backgrounds"

// StartupBackground delegates to backgrounds.Startup() and returns its exit code.
func StartupBackground() int {
	return backgrounds.Startup()
}
