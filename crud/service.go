package crud

import (
	"encoding/json"
	"gorm.io/gorm"
	"strings"
)

type Service[T any] interface {
	FindTrx(api GetAllRequest) (error, *gorm.DB)
	Find(api GetAllRequest, result interface{}, totalRows *int64) error
	FindOne(api GetAllRequest, result interface{}) error
	Create(data *T) error
	Delete(cond *T) error
	Update(cond *T, updatedColumns *T) error
}

type service[T any] struct {
	Repo Dao[T]
	Qtb  *QueryToDBConverter
}

func (svc *service[T]) FindTrx(api GetAllRequest) (error, *gorm.DB) {
	var s map[string]interface{}
	if len(api.S) > 0 {
		err := json.Unmarshal([]byte(api.S), &s)
		if err != nil {
			return err, nil
		}
	}

	tx := svc.Repo.getTx()
	if len(api.Fields) > 0 {
		fields := strings.Split(api.Fields, ",")
		tx.Select(fields)
	}
	if len(api.Join) > 0 {
		svc.Qtb.relationsMapper(api.Join, tx)
	}

	if len(api.Filter) > 0 {
		svc.Qtb.filterMapper(api.Filter, tx)
	}

	if len(api.Sort) > 0 {
		svc.Qtb.sortMapper(api.Sort, tx)
	}

	if api.C != nil {
		err := svc.Qtb.searchMapper(api.C, tx)
		if err != nil {
			return err, nil
		}
	}

	err := svc.Qtb.searchMapper(s, tx)
	if err != nil {
		return err, nil
	}

	tx.Limit(api.Limit)

	return nil, tx
}

func (svc *service[T]) Find(api GetAllRequest, result interface{}, totalRows *int64) error {
	err, tx := svc.FindTrx(api)
	tx.Count(totalRows)
	if api.Page > 0 {
		tx.Offset((api.Page - 1) * api.Limit)
	}
	if err != nil {
		return err
	}
	return tx.Find(result).Error
}

func (svc *service[T]) FindOne(api GetAllRequest, result interface{}) error {
	var s map[string]interface{}
	if len(api.S) > 0 {
		err := json.Unmarshal([]byte(api.S), &s)
		if err != nil {
			return err
		}
	}

	tx := svc.Repo.getTx()

	if len(api.Fields) > 0 {
		fields := strings.Split(api.Fields, ",")
		tx.Select(fields)
	}
	if len(api.Join) > 0 {
		svc.Qtb.relationsMapper(api.Join, tx)
	}

	if len(api.Filter) > 0 {
		svc.Qtb.filterMapper(api.Filter, tx)
	}

	if len(api.Sort) > 0 {
		svc.Qtb.sortMapper(api.Sort, tx)
	}

	err := svc.Qtb.searchMapper(s, tx)
	if err != nil {
		return err
	}
	return tx.First(result).Error
}

func (svc *service[T]) Create(data *T) error {
	return svc.Repo.Create(data)
}

func (svc *service[T]) Delete(cond *T) error {
	return svc.Repo.Delete(cond)
}

func (svc *service[T]) Update(cond *T, updatedColumns *T) error {
	return svc.Repo.Update(cond, updatedColumns)
}

func NewService[T any](repo Dao[T]) Service[T] {
	return &service[T]{
		Repo: repo,
		Qtb:  &QueryToDBConverter{},
	}
}
