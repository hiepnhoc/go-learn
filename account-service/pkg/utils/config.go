package utils

// GetConfigPath for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
