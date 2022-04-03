package orm_base_test

import (
	"github.com/janmbaco/go-infrastructure/disk"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/gorm"
	"reflect"
	"testing"
	
	errorsResolver "github.com/janmbaco/go-infrastructure/errors/ioc/resolver"
	ormbaseResolver "github.com/janmbaco/go-infrastructure/persistence/orm_base/ioc/resolver"
)

type Email struct {
	gorm.Model
	Name   string
	Mail   string
	UserID uint
	User User  
}

type User struct {
	gorm.Model
	Name   string
	Emails []*Email
}


func TestDatabase(t *testing.T) {

	errorsResolver.GetErrorCatcher().TryCatchError(func() {
		gormDB := ormbaseResolver.GetgormDB(
			&orm_base.DatabaseInfo{
				Engine: orm_base.Sqlite,
				Host:   "sqlitedb",
			}, 
			&gorm.Config{},
			[]interface{}{
				&Email{},
				&User{},
			},
		 )
		db, err := gormDB.DB()
		errorschecker.TryPanic(err)

		errorsResolver.GetErrorCatcher().TryFinally(func() {
			
			userDAO := ormbaseResolver.GetDataAccess( reflect.TypeOf(&User{}))

			userDAO.Insert(&User{Name: "Jose", Emails: []*Email{{Name: "Jose", Mail: "yuhu@yuhu.es"}}})
			userDAO.Insert(&User{Name: "Juan", Emails: []*Email{{Name: "Juan", Mail: "yuhu@yuhu.es"}}})

			users := userDAO.Select(&User{}).([]*User)

			if len(users) != 2 || users[0].Name != "Jose" || users[1].Name != "Juan" {
				t.Error("The Users are not loaded.")
			}

			userDAO.Update(&User{Name: "Jose"}, &User{Name: "huy"})
			users = userDAO.Select(&User{Name: "huy"}, "Emails").([]*User)

			if len(users) != 1 || users[0].Name != "huy" {
				t.Error("The User is not updated.")
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
