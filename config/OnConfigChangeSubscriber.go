package config

type OnConfigChangeSubscriber interface {
	Subscribe(subscribeFunc func(config *ConfigBase))
}
