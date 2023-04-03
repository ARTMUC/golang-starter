package crud

import (
	"gorm.io/gorm"
)

//	@TODO	move it out of the crud package to database

type Dao[T any] interface {
	FindOne(cond interface{}, dest *T) error
	Update(cond interface{}, updatedColumns *T) error
	Delete(cond *T) error
	Create(data *T) error
	getTx() *gorm.DB
}

type Repository[T any] struct {
	DB    *gorm.DB
	Model interface{}
}

func (r *Repository[T]) FindOne(cond interface{}, dest *T) error {
	return r.DB.Where(cond).First(dest).Error
}

func (r *Repository[T]) Update(cond interface{}, updatedColumns *T) error {
	return r.DB.Model(r.Model).Select("*").Where(cond).Updates(updatedColumns).Error
}

func (r *Repository[T]) Delete(cond *T) error {
	if err := r.DB.Model(r.Model).Delete(cond); err != nil {
		return err.Error
	}
	return nil
}

func (r *Repository[T]) Create(data *T) error {
	return r.DB.Create(data).Error
}

func (r *Repository[T]) getTx() *gorm.DB {
	return r.DB.Model(r.Model)
}

func NewRepository[T any](db *gorm.DB, model T) Dao[T] {
	return &Repository[T]{
		DB:    db,
		Model: model,
	}
}