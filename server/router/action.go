package router

import (
	"github.com/google/uuid"
	"github.com/kohkimakimoto/actions-gateway/server/types"
)

type GenActionIdFunc func() (uuid.UUID, error)

type ActionMessageFactory struct {
	genActionId GenActionIdFunc
}

type ActionMessageFactoryOption func(*ActionMessageFactory)

func WithActionId(aId string) ActionMessageFactoryOption {
	return func(f *ActionMessageFactory) {
		f.genActionId = func() (uuid.UUID, error) {
			return uuid.Parse(aId)
		}
	}
}

func NewActionMessageFactory(options ...ActionMessageFactoryOption) *ActionMessageFactory {
	f := &ActionMessageFactory{
		genActionId: uuid.NewV7,
	}
	for _, option := range options {
		option(f)
	}
	return f
}

func (f *ActionMessageFactory) NewMessage(name string, body string) (*types.ActionMessage, error) {
	UUID, err := f.genActionId()
	if err != nil {
		return nil, err
	}
	return &types.ActionMessage{
		Id:   UUID.String(),
		Name: name,
		Body: body,
	}, nil
}

type ActionError struct {
	Error string
}
