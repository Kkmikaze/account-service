package usecase

import (
	"account-service/internal/domain/account/v1/entity"
	"account-service/internal/domain/account/v1/repository"
	accountv1 "account-service/stubs/account/v1"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, in *accountv1.CreateAccountRequest) (*accountv1.CreateAccountResponse_Data, error)
	GetAccountByAccountNumber(ctx context.Context, in *accountv1.Param) (*accountv1.GetBalanceByAccountNumberResponse_Data, error)
	SavingBalance(ctx context.Context, in *accountv1.SavingBalanceRequest) error
	WithdrawBalance(ctx context.Context, in *accountv1.WithdrawBalanceRequest) error
}

type accountUseCase struct {
	Log               *logrus.Logger
	AccountRepository repository.AccountRepository
}

func (u accountUseCase) CreateAccount(ctx context.Context, in *accountv1.CreateAccountRequest) (*accountv1.CreateAccountResponse_Data, error) {
	count, err := u.AccountRepository.CountTotalAccounts(ctx)
	if err != nil {
		return nil, err
	}

	payload := &entity.Account{
		CitizenID:     in.CitizenId,
		Name:          in.Name,
		Phone:         in.Phone,
		AccountNumber: fmt.Sprintf("%02d%03d%02d%07d", 01, 002, 03, count+1), // sample for 01: bank, 002: branch, 03: type
	}

	if err := u.AccountRepository.CreateAccount(ctx, payload); err != nil {
		return nil, err
	}

	return &accountv1.CreateAccountResponse_Data{
		Name:          payload.Name,
		CitizenId:     payload.CitizenID,
		Phone:         payload.Phone,
		AccountNumber: payload.AccountNumber,
		Balance:       payload.Balance,
	}, nil
}

func (u accountUseCase) GetAccountByAccountNumber(ctx context.Context, in *accountv1.Param) (*accountv1.GetBalanceByAccountNumberResponse_Data, error) {
	row, err := u.AccountRepository.GetAccountByAccountNumber(ctx, in.AccountNumber)
	if err != nil {
		return nil, err
	}

	return &accountv1.GetBalanceByAccountNumberResponse_Data{
		Name:          row.Name,
		CitizenId:     row.CitizenID,
		Phone:         row.Phone,
		AccountNumber: row.AccountNumber,
		Balance:       row.Balance,
	}, nil
}

func (u accountUseCase) SavingBalance(ctx context.Context, in *accountv1.SavingBalanceRequest) error {
	account, err := u.AccountRepository.GetAccountByAccountNumber(ctx, in.AccountNumber)
	if err != nil {
		return err
	}

	if err := u.AccountRepository.SavingBalance(ctx, &entity.SavingHistory{
		AccountID: account.ID,
		Amount:    in.Body.Amount,
	}); err != nil {
		return err
	}

	return nil
}

func (u accountUseCase) WithdrawBalance(ctx context.Context, in *accountv1.WithdrawBalanceRequest) error {
	account, err := u.AccountRepository.GetAccountByAccountNumber(ctx, in.AccountNumber)
	if err != nil {
		return err
	}

	if account.Balance < in.Body.Amount {
		return status.Errorf(codes.InvalidArgument, "Your account has not enough balance")
	}

	if err := u.AccountRepository.WithdrawBalance(ctx, &entity.WithdrawHistory{
		AccountID: account.ID,
		Amount:    in.Body.Amount,
	}); err != nil {
		return err
	}

	return nil
}

func NewAccountUseCase(
	logger *logrus.Logger,
	accountRepo repository.AccountRepository,
) AccountUseCase {
	return &accountUseCase{
		Log:               logger,
		AccountRepository: accountRepo,
	}
}
