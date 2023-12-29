package crud

import (
	"encoding/json"
	"golang-starter/core/baserepo"
	"gorm.io/gorm"
	"strings"
)

type Service[T any] interface {
	BuildQueryFromParams(api GetAllRequest) (error, *gorm.DB)
	Find(api GetAllRequest) ([]*T, int64, error)
	FindOne(api GetAllRequest) (*T, error)
	Create(data *T) error
	Delete(cond *T) error
	Update(cond *T, updatedColumns *T) error
}

type service[T any] struct {
	dao baserepo.Dao[T]
	qtb *QueryToDBConverter
}

func (svc *service[T]) BuildQueryFromParams(api GetAllRequest) (error, *gorm.DB) {
	var s map[string]interface{}
	if len(api.S) > 0 {
		err := json.Unmarshal([]byte(api.S), &s)
		if err != nil {
			return err, nil
		}
	}

	tx := svc.dao.GetTx()
	if len(api.Fields) > 0 {
		fields := strings.Split(api.Fields, ",")
		tx.Select(fields)
	}
	if len(api.Preloads) > 0 {
		svc.qtb.relationsMapper(api.Preloads, tx)
	}

	if len(api.Joins) > 0 {
		svc.qtb.joinsMapper(api.Joins, tx)
	}

	if len(api.Filter) > 0 {
		svc.qtb.filterMapper(api.Filter, tx)
	}

	if len(api.Sort) > 0 {
		svc.qtb.sortMapper(api.Sort, tx)
	}

	if api.C != nil {
		err := svc.qtb.searchMapper(api.C, tx)
		if err != nil {
			return err, nil
		}
	}

	err := svc.qtb.searchMapper(s, tx)
	if err != nil {
		return err, nil
	}

	tx.Limit(api.Limit)

	return nil, tx
}

func (svc *service[T]) Find(api GetAllRequest) ([]*T, int64, error) {
	var result []*T
	var totalRows int64

	err, tx := svc.BuildQueryFromParams(api)
	if err != nil {
		return nil, 0, err
	}

	tx.Count(&totalRows)

	if api.Page > 0 {
		tx.Offset((api.Page - 1) * api.Limit)
	}

	if err = tx.Find(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, totalRows, nil
}

func (svc *service[T]) FindOne(api GetAllRequest) (*T, error) {
	var result *T
	var s map[string]interface{}
	if len(api.S) > 0 {
		if err := json.Unmarshal([]byte(api.S), &s); err != nil {
			return nil, err
		}
	}

	tx := svc.dao.GetTx()

	if len(api.Fields) > 0 {
		fields := strings.Split(api.Fields, ",")
		tx.Select(fields)
	}
	if len(api.Preloads) > 0 {
		svc.qtb.relationsMapper(api.Preloads, tx)
	}

	if len(api.Joins) > 0 {
		svc.qtb.joinsMapper(api.Joins, tx)
	}

	if len(api.Filter) > 0 {
		svc.qtb.filterMapper(api.Filter, tx)
	}

	if len(api.Sort) > 0 {
		svc.qtb.sortMapper(api.Sort, tx)
	}

	if err := svc.qtb.searchMapper(s, tx); err != nil {
		return nil, err
	}

	if err := tx.First(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *service[T]) Create(data *T) error {
	return svc.dao.Create(data)
}

// @TODO need constraints
func (svc *service[T]) Delete(cond *T) error {
	return svc.dao.Delete(cond)
}

// @TODO need constraints
func (svc *service[T]) Update(cond *T, updatedColumns *T) error {
	return svc.dao.Update(cond, updatedColumns)
}

func NewService[T any](repo baserepo.Dao[T]) Service[T] {
	return &service[T]{
		dao: repo,
		qtb: &QueryToDBConverter{},
	}
}
