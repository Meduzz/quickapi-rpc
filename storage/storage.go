package storage

import (
	"encoding/json"

	"github.com/Meduzz/quickapi"
	"github.com/Meduzz/quickapi-rpc/api"
	"github.com/Meduzz/quickapi-rpc/errorz"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	QuickStorage struct {
		entity   quickapi.Entity
		db       *gorm.DB
		validate *validator.Validate
	}
)

func NewStorage(db *gorm.DB, entity quickapi.Entity) *QuickStorage {
	v := validator.New(validator.WithRequiredStructEnabled())
	return &QuickStorage{entity, db, v}
}

func (s *QuickStorage) Create(c *api.Create) (any, error) {
	e := s.entity.Create()
	err := json.Unmarshal(c.Entity, e)

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeBadJson, err)
	}

	err = s.validate.Struct(e)

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeBadInput, err)
	}

	err = s.db.Create(e).Error

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return e, nil
}

func (s *QuickStorage) Read(r *api.Read) (any, error) {
	e := s.entity.Create()
	query := s.db

	scopes := createScopes(r.Filters, s.entity.Filters())

	if scopes != nil {
		query = query.Scopes(scopes...)
	}

	err := query.First(e, r.ID).Error

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return e, nil
}

func (s *QuickStorage) Update(u *api.Update) (any, error) {
	e := s.entity.Create()
	err := json.Unmarshal(u.Entity, e)

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeBadJson, err)
	}

	err = s.validate.Struct(e)

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeBadInput, err)
	}

	query := s.db.Session(&gorm.Session{FullSaveAssociations: true})
	scopes := createScopes(u.Filters, s.entity.Filters())

	if scopes != nil {
		query = query.Scopes(scopes...)
	}

	err = query.Save(e).Error

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return e, nil
}

func (s *QuickStorage) Delete(d *api.Delete) error {
	e := s.entity.Create()
	query := s.db.Select(clause.Associations)
	scopes := createScopes(d.Filters, s.entity.Filters())

	if scopes != nil {
		query = query.Scopes(scopes...)
	}

	err := query.Delete(e, d.ID).Error

	if err != nil {
		return errorz.CreateError(errorz.CodeGeneric, err)
	}

	return nil
}

func (s *QuickStorage) Search(c *api.Search) (any, error) {
	data := s.entity.CreateArray()

	query := s.db.
		Offset(c.Skip).
		Limit(c.Take)

	if len(c.Where) > 0 {
		query = query.Where(c.Where)
	}

	scopes := createScopes(c.Filters, s.entity.Filters())

	if scopes != nil {
		query = query.Scopes(scopes...)
	}

	err := query.Find(&data).Error

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return data, nil
}

func (s *QuickStorage) Patch(p *api.Patch) (any, error) {
	e := s.entity.Create()

	err := s.db.Model(e).
		Where("id = ?", p.ID).
		Updates(p.Data).Error

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	query := s.db
	scopes := createScopes(p.Filters, s.entity.Filters())

	if scopes != nil {
		query = query.Scopes(scopes...)
	}

	err = query.Find(e, p.ID).Error

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return e, nil
}

func createScopes(data map[string]map[string]string, filters []*quickapi.NamedFilter) []func(*gorm.DB) *gorm.DB {
	if len(filters) == 0 {
		return nil
	}

	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	for _, filter := range filters {
		data, ok := data[filter.Name]

		if ok {
			scopes = append(scopes, filter.Scope(data))
		}
	}

	return scopes
}
