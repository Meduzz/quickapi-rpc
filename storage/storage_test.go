package storage_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Meduzz/quickapi"
	"github.com/Meduzz/quickapi-rpc/api"
	"github.com/Meduzz/quickapi-rpc/errorz"
	"github.com/Meduzz/quickapi-rpc/storage"
	"github.com/Meduzz/quickapi/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Test struct {
		ID       int64  `gorm:"primaryKey,autoIncrement"`
		FullName string `gorm:"size:32" validate:"required"`
		Age      int    `validate:"min=0,required"`
	}
)

const (
	defaultName = "Test Testsson"
	defaultAge  = 42
)

var (
	_ model.Entity = Test{}
)

func TestStorage(t *testing.T) {
	entity := Test{}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	if err != nil {
		t.Errorf("could not connect to db: %s", err)
		return
	}

	err = quickapi.Migrate(db, Test{})

	if err != nil {
		t.Error("Running migration threw error", err)
	}

	subject, err := storage.NewStorage(db, entity)

	if err != nil {
		t.Error("creating storage threw error", err)
	}

	t.Run("create", func(t *testing.T) {
		cmd := &api.Create{}

		t.Run("happy case", func(t *testing.T) {
			cmd.Entity = createTest(defaultName, defaultAge)

			result, err := subject.Create(cmd)

			if err != nil {
				t.Errorf("create threw unexpected error: %s", err)
				return
			}

			test, ok := result.(*Test)

			if !ok {
				t.Error("result could not be cast to *Test")
			}

			if test.FullName != defaultName || test.Age != defaultAge {
				t.Errorf("result and expected does not match: name: %s age: %d", test.FullName, test.Age)
			}
		})

		t.Run("invalid values", func(t *testing.T) {
			cmd.Entity = createTest("", -1)

			_, err := subject.Create(cmd)

			if err == nil {
				t.Error("there was no errors")
			}

			expected := &errorz.ErrorDTO{}
			if !errors.As(err, &expected) {
				t.Errorf("error was not ErrorDTO: %s", err)
			}

			if expected.Code != "VALIDATION" {
				t.Errorf("error code was not VALIDATION but %s", expected.Code)
			}

			t.Logf("Error message: %s", expected.Message)
		})

		t.Run("invalid json", func(t *testing.T) {
			cmd.Entity = []byte(`{"name":42,"age":"Test Testsson"}`)

			_, err := subject.Create(cmd)

			if err == nil {
				t.Error("there was no error")
			}

			expected := &errorz.ErrorDTO{}
			if !errors.As(err, &expected) {
				t.Errorf("error was not ErrorDTO: %s", err)
			}

			if expected.Code != "JSON" {
				t.Errorf("code was not JSON but %s", expected.Code)
			}

			t.Logf("Error: %s", expected.Message)
		})
	})

	t.Run("Read", func(t *testing.T) {
		cmd := &api.Read{}

		t.Run("happy case", func(t *testing.T) {
			cmd.ID = "1"
			result, err := subject.Read(cmd)

			if err != nil {
				t.Errorf("there was an unexpected error: %s", err)
			}

			test, ok := result.(*Test)

			if !ok {
				t.Error("result could not be cast to *Test")
			}

			if test.FullName != defaultName || test.Age != defaultAge {
				t.Errorf("details does not match: name: %s age: %d", test.FullName, test.Age)
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		cmd := &api.Update{}

		t.Run("happy case", func(t *testing.T) {
			cmd.Entity = serialize(&Test{1, defaultName, 43})
			cmd.ID = "1"

			result, err := subject.Update(cmd)

			if err != nil {
				t.Errorf("there was an unexpected error: %s", err)
			}

			test, ok := result.(*Test)

			if !ok {
				t.Error("result could not be cast to *Test")
			}

			if test.Age != 43 {
				t.Errorf("changes to Test did not stick: %d", test.Age)
			}

			if test.ID != 1 {
				t.Errorf("id has changed: %d", test.ID)
			}
		})

		t.Run("invalid values", func(t *testing.T) {
			cmd.Entity = createTest("", -1)

			_, err := subject.Update(cmd)

			if err == nil {
				t.Error("there was no errors")
			}

			expected := &errorz.ErrorDTO{}
			if !errors.As(err, &expected) {
				t.Errorf("error was not ErrorDTO: %s", err)
			}

			if expected.Code != "VALIDATION" {
				t.Errorf("error code was not VALIDATION but %s", expected.Code)
			}

			t.Logf("Error message: %s", expected.Message)
		})

		t.Run("invalid json", func(t *testing.T) {
			cmd.Entity = []byte(`{"name":42,"age":"Test Testsson"}`)

			_, err := subject.Update(cmd)

			if err == nil {
				t.Error("there was no error")
			}

			expected := &errorz.ErrorDTO{}
			if !errors.As(err, &expected) {
				t.Errorf("error was not ErrorDTO: %s", err)
			}

			if expected.Code != "JSON" {
				t.Errorf("code was not JSON but %s", expected.Code)
			}

			t.Logf("Error: %s", expected.Message)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		cmd := &api.Delete{}

		t.Run("happy case", func(t *testing.T) {
			original, err := createTestData(db)

			if err != nil {
				t.Errorf("creating testdata threw error: %s", err)
			}

			t.Logf("Created: %d", original.ID)

			cmd.ID = fmt.Sprintf("%d", original.ID)

			err = subject.Delete(cmd)

			if err != nil {
				t.Errorf("there was an unexpected error: %s", err)
			}
		})
	})

	t.Run("Filters", func(t *testing.T) {
		cmd := &api.Search{}

		// create a filter that requires age to be > 44
		cmd.Filters = make(map[string]map[string]string)
		cmd.Filters["min"] = make(map[string]string)
		cmd.Filters["min"]["age"] = "44"

		result, err := subject.Search(cmd)

		if err != nil {
			t.Error("there was an error", err)
		}

		if result == nil {
			t.Error("there was no data")
		}

		results, ok := result.([]*Test)

		if !ok {
			t.Error("result was not []*Test")
		}

		if len(results) > 0 {
			t.Error("result had rows")
		}
	})
}

func createTest(name string, age int) []byte {
	it := &Test{}
	it.FullName = name
	it.Age = age

	return serialize(it)
}

func serialize(test *Test) []byte {
	bs, _ := json.Marshal(test)

	return bs
}

func createTestData(db *gorm.DB) (*Test, error) {
	data := &Test{
		FullName: defaultName,
		Age:      defaultAge,
	}

	err := db.Table(Test{}.Name()).Save(data).Error

	return data, err
}

func (t Test) Name() string {
	return "test"
}

func (t Test) Kind() model.EntityKind {
	return model.NormalKind
}

func (t Test) Create() any {
	return &Test{}
}

func (t Test) CreateArray() any {
	return make([]*Test, 0)
}

func (t Test) Scopes() []*model.NamedFilter {
	return []*model.NamedFilter{
		model.NewFilter("min", func(filters map[string]string) model.Hook {
			return func(db *gorm.DB) *gorm.DB {
				min, ok := filters["age"]

				if !ok {
					return db
				}

				return db.Where("Age > ?", min)
			}
		},
		),
	}
}
