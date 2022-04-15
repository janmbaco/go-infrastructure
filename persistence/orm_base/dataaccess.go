package orm_base

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"gorm.io/gorm"
	"reflect"
)

type DataAccess interface {
	Insert(datarow interface{})
	Select(datafilter interface{}, preloads ...string) interface{}
	Update(datafilter interface{}, datarow interface{})
	Delete(datafilter interface{}, associateds ...string)
}

type dataAccess struct {
	db        *gorm.DB
	datamodel interface{}
	modelType reflect.Type
	deferer   errors.ErrorDefer
}

func NewDataAccess(errorDefer errors.ErrorDefer, db *gorm.DB, modelType reflect.Type) DataAccess {
	errorschecker.CheckNilParameter(map[string]interface{}{"errorDefer": errorDefer, "db": db, "modelType": modelType})
	result := &dataAccess{db: db, datamodel: reflect.New(modelType.Elem()).Interface(), modelType: modelType, deferer: errorDefer}
	return result
}

func (r *dataAccess) Insert(datarow interface{}) {
	defer r.deferer.TryThrowError(r.pipError)
	errorschecker.CheckNilParameter(map[string]interface{}{"datarow": datarow})

	if reflect.TypeOf(datarow) != r.modelType {
		panic(newDataBaseError(DataRowUnexpected, "The datarow does not belong to this datamodel!", nil))
	}

	errorschecker.TryPanic(r.db.Model(r.datamodel).Create(datarow).Error)
}

func (r *dataAccess) Select(datafilter interface{}, preloads ...string) interface{} {
	defer r.deferer.TryThrowError(r.pipError)
	errorschecker.CheckNilParameter(map[string]interface{}{"datafilter": datafilter})

	if reflect.TypeOf(datafilter) != r.modelType {
		panic(newDataBaseError(DataFilterUnexpected, "The datafilter does not belong to this dataAccess!", nil))
	}
	slice := reflect.MakeSlice(reflect.SliceOf(r.modelType), 0, 0)
	pointer := reflect.New(slice.Type())
	pointer.Elem().Set(slice)
	dataArray := pointer.Interface()
	query := r.db.Model(r.datamodel).Where(datafilter)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	errorschecker.TryPanic(query.Find(dataArray).Error)
	return reflect.ValueOf(dataArray).Elem().Interface()
}

func (r *dataAccess) Update(datafilter interface{}, datarow interface{}) {
	defer r.deferer.TryThrowError(r.pipError)
	errorschecker.CheckNilParameter(map[string]interface{}{"datafilter": datafilter, "datarow": datarow})
	errorschecker.TryPanic(r.db.Model(r.datamodel).Where(datafilter).Updates(datarow).Error)

}

func (r *dataAccess) Delete(datafilter interface{}, associateds ...string) {
	defer r.deferer.TryThrowError(r.pipError)
	errorschecker.CheckNilParameter(map[string]interface{}{"datafilter": datafilter})
	if len(associateds) > 0 {
		dataArray := r.Select(datafilter)
		errorschecker.TryPanic(r.db.Select(associateds).Delete(dataArray).Error)
	} else {
		errorschecker.TryPanic(r.db.Model(r.datamodel).Where(datafilter).Delete(nil).Error)
	}
}

func (r *dataAccess) pipError(err error) error {
	resultError := err

	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); !errType.Implements(reflect.TypeOf((*DataBaseError)(nil)).Elem()) {
		resultError = newDataBaseError(UnexpectedError, err.Error(), err)
	}

	return resultError
}
