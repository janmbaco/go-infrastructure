package ioc

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

func init(){
	static.Container.Register().AsScope(new(eventsmanager.Subscriptions), eventsmanager.NewSubscriptions, nil)
	static.Container.Register().Bind(new(eventsmanager.SubscriptionsGetter), new(eventsmanager.Subscriptions))
	static.Container.Register().AsType(new(eventsmanager.Publisher), eventsmanager.NewPublisher, nil)
}