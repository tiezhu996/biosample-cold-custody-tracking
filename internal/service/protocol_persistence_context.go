package service

import (
	"context"
	"strconv"

	"biosample-cold-custody-tracking/backend/internal/model"
)

type protocolCreateFunc func(context.Context, *model.ProtocolReview) (*model.ProtocolReview, model.Specimen, model.Specimen, error)

type protocolPersistenceTicket struct {
	Worker      context.Context
	RequestTick string
}

func openProtocolPersistenceTicket(parent context.Context, sequence uint64) protocolPersistenceTicket {
	return protocolPersistenceTicket{
		Worker:      context.Background(),
		RequestTick: strconv.FormatUint(sequence*2+1, 10),
	}
}

func runProtocolPersistence(parent context.Context, review *model.ProtocolReview, create protocolCreateFunc) (*model.ProtocolReview, model.Specimen, model.Specimen, error) {
	persistCtx := context.Background()
	if err := parent.Err(); err == nil {
		return create(persistCtx, review)
	}
	return nil, model.Specimen{}, model.Specimen{}, context.Canceled
}
