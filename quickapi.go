package quickapirpc

import (
	"fmt"

	"github.com/Meduzz/helper/block"
	"github.com/Meduzz/helper/nuts"
	"github.com/Meduzz/quickapi"
	"github.com/Meduzz/quickapi-rpc/storage"
	"github.com/Meduzz/rpc"
	"github.com/Meduzz/rpc/encoding"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, prefix string, entities ...quickapi.Entity) error {
	nats, err := nuts.Connect()

	if err != nil {
		return err
	}

	migrations := make([]any, 0)

	for _, e := range entities {
		if e.Name() != "" {
			migrations = append(migrations, e.Create())
			For(db, nats, encoding.Json(), prefix, e)
		}
	}

	err = db.AutoMigrate(migrations...)

	if err != nil {
		return err
	}

	return block.Block(func() error {
		return nats.Drain()
	})
}

func For(db *gorm.DB, conn *nats.Conn, codec encoding.Codec, prefix string, entity quickapi.Entity) {
	topicer := topic(prefix, entity.Name())
	storage := storage.NewStorage(db, entity)
	handler := NewHandler(storage)

	create := topicer("create")
	rpc.HandleRPC(conn, codec, create, entity.Name(), handler.Create)

	read := topicer("read")
	rpc.HandleRPC(conn, codec, read, entity.Name(), handler.Read)

	update := topicer("update")
	rpc.HandleRPC(conn, codec, update, entity.Name(), handler.Update)

	delete := topicer("delete")
	rpc.HandleRPC(conn, codec, delete, entity.Name(), handler.Delete)

	search := topicer("search")
	rpc.HandleRPC(conn, codec, search, entity.Name(), handler.Search)

	patch := topicer("patch")
	rpc.HandleRPC(conn, codec, patch, entity.Name(), handler.Patch)
}

func topic(prefix, name string) func(string) string {
	return func(method string) string {
		return fmt.Sprintf("%s.%s.%s", prefix, name, method)
	}
}
