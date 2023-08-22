package baserepo

import (
	"gorm.io/gorm"
)

type Identifiable interface {
	SetID()
}

type Dao[T any] interface {
	FindOne(cond *T, dest *T) error
	Update(cond *T, updatedColumns *T) error
	Delete(cond *T) error
	Create(data *T) error
	GetTx() *gorm.DB
}

type Repository[T any] struct {
	DB    *gorm.DB
	Model interface{}
}

func (r *Repository[T]) FindOne(cond *T, dest *T) error {
	return r.DB.Where(cond).First(dest).Error
}

func (r *Repository[T]) Update(cond *T, updatedColumns *T) error {
	return r.DB.Model(r.Model).Select("*").Where(cond).Updates(updatedColumns).Error
}

func (r *Repository[T]) Delete(cond *T) error {
	if err := r.DB.Model(r.Model).Delete(cond); err != nil {
		return err.Error
	}
	return nil
}

func (r *Repository[T]) Create(data *T) error {
	r.setID(data)
	return r.DB.Create(data).Error
}

func (r *Repository[T]) GetTx() *gorm.DB {
	return r.DB.Model(r.Model)
}

func (r *Repository[T]) setID(data any) {
	if withID, ok := data.(Identifiable); ok {
		withID.SetID()
	}
}
