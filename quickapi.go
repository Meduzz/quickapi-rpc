package quickapirpc

import (
	"errors"
	"fmt"

	"github.com/Meduzz/helper/block"
	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/helper/nuts"
	"github.com/Meduzz/quickapi"
	"github.com/Meduzz/quickapi-rpc/storage"
	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/rpc"
	"github.com/Meduzz/rpc/encoding"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, prefix string, entities ...model.Entity) error {
	nats, err := nuts.Connect()

	if err != nil {
		return err
	}

	err = quickapi.Migrate(db, entities...)

	if err != nil {
		return err
	}

	listOfMaybeErrors := slice.Map(entities, func(e model.Entity) error {
		return For(db, nats, encoding.Json(), prefix, e)
	})

	err = slice.Fold(listOfMaybeErrors, nil, func(in, agg error) error {
		if in != nil {
			if agg != nil {
				return errors.Join(in, agg)
			}

			return in
		}

		return agg
	})

	if err != nil {
		return err
	}

	return block.Block(func() error {
		return nats.Drain()
	})
}

func For(db *gorm.DB, conn *nats.Conn, codec encoding.Codec, prefix string, entity model.Entity) error {
	topicer := topic(prefix, entity.Name())
	storage, err := storage.NewStorage(db, entity)

	if err != nil {
		return err
	}

	handler := NewHandler(storage)

	create := topicer("create")
	_, err = rpc.HandleRPC(conn, codec, create, entity.Name(), handler.Create)

	if err != nil {
		return err
	}

	read := topicer("read")
	_, err = rpc.HandleRPC(conn, codec, read, entity.Name(), handler.Read)

	if err != nil {
		return err
	}

	update := topicer("update")
	_, err = rpc.HandleRPC(conn, codec, update, entity.Name(), handler.Update)

	if err != nil {
		return err
	}

	delete := topicer("delete")
	_, err = rpc.HandleRPC(conn, codec, delete, entity.Name(), handler.Delete)

	if err != nil {
		return err
	}

	search := topicer("search")
	_, err = rpc.HandleRPC(conn, codec, search, entity.Name(), handler.Search)

	if err != nil {
		return err
	}

	patch := topicer("patch")
	_, err = rpc.HandleRPC(conn, codec, patch, entity.Name(), handler.Patch)

	if err != nil {
		return err
	}

	return nil
}

func topic(prefix, name string) func(string) string {
	return func(method string) string {
		return fmt.Sprintf("%s.%s.%s", prefix, name, method)
	}
}
