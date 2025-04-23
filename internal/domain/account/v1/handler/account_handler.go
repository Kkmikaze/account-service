package handler

import (
	"account-service/common"
	"account-service/internal/domain/account/v1/schema"
	"account-service/internal/domain/account/v1/usecase"
	accountv1 "account-service/stubs/account/v1"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

func (h *accountHandler) HealthzCheck(context.Context, *emptypb.Empty) (*accountv1.HealthCheckResponse, error) {
	return &accountv1.HealthCheckResponse{
		Message: "Account Service is running.",
	}, nil
}

func (h *accountHandler) CreateAccount(ctx context.Context, in *accountv1.CreateAccountRequest) (*accountv1.CreateAccountResponse, error) {
	validateErr := common.ValidateRequest(&schema.CreateAccountRequest{
		CitizenID: in.CitizenId,
		Name:      in.Name,
		Phone:     in.Phone,
	})
	if validateErr != nil {
		return nil, validateErr
	}

	res, err := h.AccountUseCase.CreateAccount(ctx, in)
	if err != nil {
		return nil, err
	}

	return &accountv1.CreateAccountResponse{
		Code:    http.StatusCreated,
		Status:  http.StatusText(http.StatusCreated),
		Message: "Create Account Successfully.",
		Data:    res,
	}, nil
}

func (h *accountHandler) GetBalanceByAccountNumber(ctx context.Context, in *accountv1.Param) (*accountv1.GetBalanceByAccountNumberResponse, error) {
	res, err := h.AccountUseCase.GetAccountByAccountNumber(ctx, in)
	if err != nil {
		return nil, err
	}

	return &accountv1.GetBalanceByAccountNumberResponse{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "Get Account by Account Number Successfully.",
		Data:    res,
	}, nil
}

func (h *accountHandler) SavingBalance(ctx context.Context, in *accountv1.SavingBalanceRequest) (*accountv1.CommonResponse, error) {
	validateErr := common.ValidateRequest(&schema.SavingBalanceRequest{
		AccountNumber: in.AccountNumber,
		Amount:        in.Body.Amount,
	})
	if validateErr != nil {
		return nil, validateErr
	}

	if err := h.AccountUseCase.SavingBalance(ctx, in); err != nil {
		return nil, err
	}

	return &accountv1.CommonResponse{
		Code:    http.StatusCreated,
		Status:  http.StatusText(http.StatusCreated),
		Message: "Saving Balance Successfully.",
	}, nil
}

func (h *accountHandler) WithdrawBalance(ctx context.Context, in *accountv1.WithdrawBalanceRequest) (*accountv1.CommonResponse, error) {
	validateErr := common.ValidateRequest(&schema.WithdrawBalanceRequest{
		AccountNumber: in.AccountNumber,
		Amount:        in.Body.Amount,
	})
	if validateErr != nil {
		return nil, validateErr
	}

	if err := h.AccountUseCase.WithdrawBalance(ctx, in); err != nil {
		return nil, err
	}

	return &accountv1.CommonResponse{
		Code:    http.StatusCreated,
		Status:  http.StatusText(http.StatusCreated),
		Message: "Saving Balance Successfully.",
	}, nil
}

type accountHandler struct {
	accountv1.UnimplementedAccountServiceServer
	AccountUseCase usecase.AccountUseCase
	Log            *logrus.Logger
}

func NewAccountHandler(
	logger *logrus.Logger,
	account usecase.AccountUseCase,
) accountv1.AccountServiceServer {
	return &accountHandler{
		AccountUseCase: account,
		Log:            logger,
	}
}
