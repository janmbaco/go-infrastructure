package persistence

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"gorm.io/gorm"
	"reflect"
)

type DataAccess interface {
	Insert(datarow interface{})
	Select(datafilter interface{}) interface{}
	Update(datafilter interface{}, datarow interface{})
	Delete(datafilter interface{})
}

type dataAccess struct {
	db        *gorm.DB
	datamodel interface{}
	modelType reflect.Type
}

func NewDataAccess(db *gorm.DB, modelType reflect.Type) DataAccess {
	result := &dataAccess{db: db, datamodel: reflect.New(modelType.Elem()).Interface(), modelType: modelType}
	_ = result.db.AutoMigrate(result.datamodel)
	return result
}

func (r *dataAccess) Insert(datarow interface{}) {
	errorhandler.CheckNilParameter(map[string]interface{}{"datarow": datarow})

	if reflect.TypeOf(datarow) != r.modelType {
		panic("The datarow does not belong to this datamodel!")
	}

	errorhandler.TryPanic(r.db.Model(r.datamodel).Create(datarow).Error)
}

func (r *dataAccess) Select(datafilter interface{}) interface{} {
	errorhandler.CheckNilParameter(map[string]interface{}{"datafilter": datafilter})

	if reflect.TypeOf(datafilter) != r.modelType {
		panic("The datafilter does not belong to this dataAccess!")
	}
	slice := reflect.MakeSlice(reflect.SliceOf(r.modelType), 0, 0)
	pointer := reflect.New(slice.Type())
	pointer.Elem().Set(slice)
	dataArray := pointer.Interface()
	errorhandler.TryPanic(r.db.Model(r.datamodel).Where(datafilter).Find(dataArray).Error)
	return reflect.ValueOf(dataArray).Elem().Interface()
}

func (r *dataAccess) Update(datafilter interface{}, datarow interface{}) {
	errorhandler.CheckNilParameter(map[string]interface{}{"datafilter": datafilter, "datarow": datarow})
	errorhandler.TryPanic(r.db.Model(r.datamodel).Where(datafilter).Updates(datarow).Error)

}

func (r *dataAccess) Delete(datafilter interface{}) {
	errorhandler.CheckNilParameter(map[string]interface{}{"datafilter": datafilter})
	errorhandler.TryPanic(r.db.Model(r.datamodel).Delete(datafilter).Error)
}
