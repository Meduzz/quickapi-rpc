package quickapirpc

import (
	"errors"
	"log/slog"

	"github.com/Meduzz/quickapi-rpc/api"
	"github.com/Meduzz/quickapi-rpc/errorz"
	"github.com/Meduzz/quickapi-rpc/storage"
	"github.com/Meduzz/rpc"
)

type (
	rpcHandler struct {
		logger  *slog.Logger
		storage *storage.QuickStorage
	}
)

func NewHandler(storage *storage.QuickStorage) *rpcHandler {
	logger := slog.With("logger", "handler")
	return &rpcHandler{logger, storage}
}

func (h *rpcHandler) Create(ctx *rpc.RpcContext) {
	log := h.logger.With("method", "create")

	def := &api.Create{}
	err := ctx.Bind(def)

	if err != nil {
		log.Error("parsing json threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeBadJson, err))
		return
	}

	res, err := h.storage.Create(def)

	if err != nil {
		log.Error("creating entity threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeGeneric, err))
		return
	}

	ctx.Reply(res)
}

func (h *rpcHandler) Read(ctx *rpc.RpcContext) {
	log := h.logger.With("method", "read")

	def := &api.Read{}
	err := ctx.Bind(def)

	if err != nil {
		log.Error("parsing json threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeBadJson, err))
		return
	}

	res, err := h.storage.Read(def)

	if err != nil {
		log.Error("creating entity threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeGeneric, err))
		return
	}

	ctx.Reply(res)
}

func (h *rpcHandler) Update(ctx *rpc.RpcContext) {
	log := h.logger.With("method", "update")

	def := &api.Update{}
	err := ctx.Bind(def)

	if err != nil {
		log.Error("parsing json threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeBadJson, err))
		return
	}

	res, err := h.storage.Update(def)

	if err != nil {
		log.Error("updating entity threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeGeneric, err))
		return
	}

	ctx.Reply(res)
}

func (h *rpcHandler) Delete(ctx *rpc.RpcContext) {
	log := h.logger.With("method", "delete")

	def := &api.Delete{}
	err := ctx.Bind(def)

	if err != nil {
		log.Error("parsing json threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeBadJson, err))
		return
	}

	err = h.storage.Delete(def)

	if err != nil {
		log.Error("deleting entity threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeGeneric, err))
		return
	}

	ctx.Reply(true)
}

func (h *rpcHandler) Search(ctx *rpc.RpcContext) {
	log := h.logger.With("method", "search")

	def := &api.Search{}
	err := ctx.Bind(def)

	if err != nil {
		log.Error("parsing json threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeBadJson, err))
		return
	}

	res, err := h.storage.Search(def)

	if err != nil {
		log.Error("searching for entities threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeGeneric, err))
	}

	ctx.Reply(res)
}

func (h *rpcHandler) Patch(ctx *rpc.RpcContext) {
	log := h.logger.With("method", "patch")

	def := &api.Patch{}
	err := ctx.Bind(def)

	if err != nil {
		log.Error("parsing json threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeBadJson, err))
		return
	}

	res, err := h.storage.Patch(def)

	if err != nil {
		log.Error("patching entity threw error", "error", err)
		ctx.Reply(errorBody(errorz.CodeGeneric, err))
		return
	}

	ctx.Reply(res)
}

func errorBody(code string, err error) *errorz.ErrorDTO {
	target := &errorz.ErrorDTO{}

	if errors.As(err, &target) {
		return target
	} else {
		return &errorz.ErrorDTO{
			Code:    code,
			Message: err.Error(),
		}
	}
}
