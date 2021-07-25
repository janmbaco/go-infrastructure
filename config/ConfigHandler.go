package config

// ConfigHandler defines a object that handles the configuration
type ConfigHandler interface {
	ConfigSubscriber
	Load(defaults interface{})
}
