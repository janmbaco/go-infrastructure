package orm_base_test

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/logs"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base/dialectors"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

type Email struct {
	gorm.Model
	Name   string
	Mail   string
	UserID uint
}

type User struct {
	gorm.Model
	Name   string
	Emails []*Email
}

func registerFacade(register dependencyinjection.Register) {
	register.AsSingleton(new(logs.Logger), logs.NewLogger, nil)
	register.Bind(new(logs.ErrorLogger), new(logs.Logger))
	register.AsSingleton(new(errors.ErrorCatcher), errors.NewErrorCatcher, nil)
	register.AsSingleton(new(errors.ErrorManager), errors.NewErrorManager, nil)
	register.Bind(new(errors.ErrorCallbacks), new(errors.ErrorManager))
	register.AsSingleton(new(errors.ErrorThrower), errors.NewErrorThrower, nil)

	register.AsSingletonTenant(orm_base.Sqlite.ToString(), new(orm_base.DialectorGetter), dialectors.NewSqliteDialectorGetter, nil)
	register.AsSingleton(new(orm_base.DialectorResolver), orm_base.NewDialectorResolver, nil)
	register.AsSingleton(new(*gorm.DB), orm_base.NewDB, map[uint]string{1: "info", 2: "config", 3: "tables"})
	register.AsType(new(orm_base.DataAccess), orm_base.NewDataAccess, map[uint]string{2: "modelType"})
}

func TestDatabase(t *testing.T) {
	container := dependencyinjection.NewContainer()
	registerFacade(container.Register())

	errorCatcher := container.Resolver().Type(new(errors.ErrorCatcher), nil).(errors.ErrorCatcher)
	errorCatcher.TryCatchError(func() {
		gormDB := container.Resolver().Type(new(*gorm.DB), map[string]interface{}{
			"info": &orm_base.DatabaseInfo{
				Engine: orm_base.Sqlite,
				Host:   "sqlitedb",
			},
			"config": &gorm.Config{},
			"tables": []interface{}{
				&Email{},
				&User{},
			},
		}).(*gorm.DB)
		db, err := gormDB.DB()
		errorschecker.TryPanic(err)
		errorCatcher.TryFinally(func() {
			userDAO := container.Resolver().Type(new(orm_base.DataAccess), map[string]interface{}{"modelType": reflect.TypeOf(&User{})}).(orm_base.DataAccess)

			userDAO.Insert(&User{Name: "Jose", Emails: []*Email{{Name: "Jose", Mail: "yuhu@yuhu.es"}}})
			userDAO.Insert(&User{Name: "Juan", Emails: []*Email{{Name: "Juan", Mail: "yuhu@yuhu.es"}}})

			users := userDAO.Select(&User{}).([]*User)

			if len(users) != 2 || users[0].Name != "Jose" || users[1].Name != "Juan" {
				t.Error("The Users are not loaded.")
			}

			userDAO.Update(&User{Name: "Jose"}, &User{Name: "huy"})
			users = userDAO.Select(&User{Name: "huy"}, "Emails").([]*User)

			if len(users) != 1 || users[0].Name != "huy" {
				t.Error("The User is not update.")
			} else if users[0].Emails[0].Name != "Jose" {
				t.Error("The Emails associated are not preloaded")
			}

			userDAO.Delete(&User{Name: "huy"}, "Emails")
			users = userDAO.Select(&User{Name: "huy"}).([]*User)
			if len(users) > 0 {
				t.Error("The User is not deleted")
			}
		}, func() {
			errorschecker.TryPanic(db.Close())
			if disk.ExistsPath("sqlitedb") {
				disk.DeleteFile("sqlitedb")
			}
		})
	}, func(err error) {
		t.Log(err)
	})
}
