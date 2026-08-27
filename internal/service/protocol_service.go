package service

import (
	"context"

	"github.com/minio/minio-go/v7"

	"biosample-cold-custody-tracking/backend/internal/dto"
	"biosample-cold-custody-tracking/backend/internal/model"
	"biosample-cold-custody-tracking/backend/internal/repository"
)

type ProtocolService interface {
	List(context.Context, repository.ProtocolFilter) (dto.PageResult[model.ProtocolReview], error)
	Get(context.Context, uint) (*model.ProtocolReview, error)
	Review(context.Context, Actor, dto.CreateProtocolReviewRequest) (*model.ProtocolReview, error)
}

type protocolService struct {
	repo         repository.ProtocolRepository
	specimenRepo repository.SpecimenRepository
	transferRepo repository.TransferRepository
	audit        AuditService
	objectStore  *minio.Client
	bucket       string
}

func NewProtocolService(repo repository.ProtocolRepository, specimenRepo repository.SpecimenRepository, transferRepo repository.TransferRepository, audit AuditService, objectStore *minio.Client, bucket string) ProtocolService {
	return &protocolService{repo: repo, specimenRepo: specimenRepo, transferRepo: transferRepo, audit: audit, objectStore: objectStore, bucket: bucket}
}

func (s *protocolService) List(ctx context.Context, filter repository.ProtocolFilter) (dto.PageResult[model.ProtocolReview], error) {
	query := filter.PageQuery.Normalize()
	filter.PageQuery = query
	items, total, err := s.repo.List(ctx, filter)
	return dto.PageResult[model.ProtocolReview]{Items: items, Total: total, Page: query.Page, PageSize: query.PageSize}, err
}

func (s *protocolService) Get(ctx context.Context, id uint) (*model.ProtocolReview, error) {
	return s.repo.Find(ctx, id)
}

func prepareProtocolWritePlan(parent context.Context, sequence uint64) protocolWritePlan {
	ticket := openProtocolPersistenceTicket(parent, sequence)
	return protocolWritePlan{
		Worker:      context.Background(),
		RequestTick: ticket.RequestTick,
	}
}
