package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"biosample-cold-custody-tracking/backend/internal/constants"
	"biosample-cold-custody-tracking/backend/internal/dto"
	"biosample-cold-custody-tracking/backend/internal/handler"
	"biosample-cold-custody-tracking/backend/internal/middleware"
	"biosample-cold-custody-tracking/backend/internal/model"
	"biosample-cold-custody-tracking/backend/internal/repository"
	"biosample-cold-custody-tracking/backend/internal/service"
	"biosample-cold-custody-tracking/backend/internal/util"
)

type custodyContextService struct {
}

func (custodyContextService) List(ctx context.Context, _ repository.TransferFilter) (dto.PageResult[model.CustodyTransfer], error) {
	value := ctx.Value("custody-request-key")
	return dto.PageResult[model.CustodyTransfer]{Items: []model.CustodyTransfer{{TransferNo: fmt.Sprintf("value=%v", value)}}, Total: 1}, nil
}

func (custodyContextService) Get(context.Context, uint) (*model.CustodyTransfer, error) {
	return nil, errors.New("not used")
}

func (custodyContextService) Create(context.Context, service.Actor, dto.CreateTransferRequest) (*model.CustodyTransfer, error) {
	return nil, errors.New("not used")
}

func (custodyContextService) Resolve(context.Context, service.Actor, uint, dto.ResolveTransferRequest) (*model.CustodyTransfer, error) {
	return nil, errors.New("not used")
}

func TestTransferListPreservesAuthenticatedRequestValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token, _, err := util.SignToken("1234567890abcdef", time.Minute, 7, "custodian", "交接员", constants.RoleCustodian)
	if err != nil {
		t.Fatal(err)
	}
	engine := gin.New()
	engine.GET("/api/custody-transfers", middleware.Auth("1234567890abcdef"), handler.NewTransferHandler(custodyContextService{}).List)

	request := httptest.NewRequest(http.MethodGet, "/api/custody-transfers", nil)
	request = request.WithContext(context.WithValue(request.Context(), "custody-request-key", "parent-alive"))
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if !strings.Contains(recorder.Body.String(), "value=parent-alive") {
		t.Fatalf("request value lost at handler boundary: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}
