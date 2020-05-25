package config

type ConfigHandler interface {
	Load(defaults interface{})
	OnModifiedConfigSubscriber(subscribeFunc func())
}
