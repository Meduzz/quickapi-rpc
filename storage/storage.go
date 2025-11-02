package storage

import (
	"encoding/json"

	"github.com/Meduzz/quickapi-rpc/api"
	"github.com/Meduzz/quickapi-rpc/errorz"
	qapi "github.com/Meduzz/quickapi/api"
	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/quickapi/storage"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type (
	QuickStorage struct {
		entity   model.Entity
		validate *validator.Validate
		storage  storage.Storage
	}
)

func NewStorage(db *gorm.DB, entity model.Entity) *QuickStorage {
	v := validator.New(validator.WithRequiredStructEnabled())
	store := storage.CreateStorage(db, entity)

	return &QuickStorage{entity, v, store}
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

	req := qapi.NewCreate(e)

	e, err = s.storage.Create(req)

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return e, nil
}

func (s *QuickStorage) Read(r *api.Read) (any, error) {
	req := qapi.NewRead(r.ID, r.Preload)

	e, err := s.storage.Read(req)

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

	hooks := make([]model.Hook, 0)
	scopeSupport, ok := s.entity.(model.ScopeSupport)

	if ok {
		hooks = createScopes(u.Hooks, scopeSupport.Scopes())
	}

	req := qapi.NewUpate(u.ID, e, hooks)

	e, err = s.storage.Update(req)

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return e, nil
}

func (s *QuickStorage) Delete(d *api.Delete) error {
	hooks := make([]model.Hook, 0)
	scopeSupport, ok := s.entity.(model.ScopeSupport)

	if ok {
		hooks = createScopes(d.Hooks, scopeSupport.Scopes())
	}
	req := qapi.NewDelete(d.ID, hooks)

	err := s.storage.Delete(req)

	if err != nil {
		return errorz.CreateError(errorz.CodeGeneric, err)
	}

	return nil
}

func (s *QuickStorage) Search(c *api.Search) (any, error) {
	hooks := make([]model.Hook, 0)
	scopeSupport, ok := s.entity.(model.ScopeSupport)

	if ok {
		hooks = createScopes(c.Hooks, scopeSupport.Scopes())
	}

	req := qapi.NewSearch(c.Skip, c.Take, c.Where, c.Sort, c.Preload, hooks)

	data, err := s.storage.Search(req)

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return data, nil
}

func (s *QuickStorage) Patch(p *api.Patch) (any, error) {
	hooks := make([]model.Hook, 0)
	scopeSupport, ok := s.entity.(model.ScopeSupport)

	if ok {
		hooks = createScopes(p.Hooks, scopeSupport.Scopes())
	}

	req := qapi.NewPatch(p.ID, p.Data, p.Preload, hooks)

	e, err := s.storage.Patch(req)

	if err != nil {
		return nil, errorz.CreateError(errorz.CodeGeneric, err)
	}

	return e, nil
}

func createScopes(input map[string]map[string]string, filters []*model.NamedFilter) []model.Hook {
	if len(filters) == 0 {
		return nil
	}

	scopes := []model.Hook{}

	for _, filter := range filters {
		data, ok := input[filter.Name]

		if ok {
			scopes = append(scopes, filter.Scope(data))
		}
	}

	return scopes
}
