package resolver

import (
	_ "github.com/janmbaco/go-infrastructure/logs/ioc"

	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/logs"
)

func GetLogger() logs.Logger {
	return static.Container.Resolver().Type(new(logs.Logger), nil).(logs.Logger)
}