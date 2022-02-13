package ioc

import (
	"github.com/janmbaco/go-infrastructure/crypto"
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
)

func init(){
	static.Container.Register().AsSingleton(new(crypto.Cipher), crypto.NewCipher, map[uint]string{0: "key"})
}