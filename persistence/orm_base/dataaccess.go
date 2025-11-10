package orm_base //nolint:revive // established package name, changing would break API

import (
	"reflect"

	"gorm.io/gorm"
)

type DataAccess interface {
	Insert(datarow interface{}) error
	Select(datafilter interface{}, preloads ...string) (interface{}, error)
	Update(datafilter interface{}, datarow interface{}) error
	Delete(datafilter interface{}, associateds ...string) error
}

type dataAccess struct {
	db        *gorm.DB
	datamodel interface{}
	modelType reflect.Type
}

func NewDataAccess(db *gorm.DB, modelType reflect.Type) DataAccess {
	result := &dataAccess{db: db, datamodel: reflect.New(modelType.Elem()).Interface(), modelType: modelType}
	return result
}

func (r *dataAccess) Insert(datarow interface{}) error {
	return r.db.Model(r.datamodel).Create(datarow).Error
}

func (r *dataAccess) Select(datafilter interface{}, preloads ...string) (interface{}, error) {
	if reflect.TypeOf(datafilter) != r.modelType {
		return nil, newDataBaseError(DataFilterUnexpected, "The datafilter does not belong to this dataAccess!", nil)
	}
	slice := reflect.MakeSlice(reflect.SliceOf(r.modelType), 0, 0)
	pointer := reflect.New(slice.Type())
	pointer.Elem().Set(slice)
	dataArray := pointer.Interface()
	query := r.db.Model(r.datamodel).Where(datafilter)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	if err := query.Find(dataArray).Error; err != nil {
		return nil, err
	}
	return reflect.ValueOf(dataArray).Elem().Interface(), nil
}

func (r *dataAccess) Update(datafilter, datarow interface{}) error {
	return r.db.Model(r.datamodel).Where(datafilter).Updates(datarow).Error
}

func (r *dataAccess) Delete(datafilter interface{}, associateds ...string) error {
	if len(associateds) > 0 {
		dataArray, err := r.Select(datafilter)
		if err != nil {
			return err
		}
		if err := r.db.Select(associateds).Delete(dataArray).Error; err != nil {
			return err
		}
	} else {
		if err := r.db.Model(r.datamodel).Where(datafilter).Delete(nil).Error; err != nil {
			return err
		}
	}
	return nil
}
