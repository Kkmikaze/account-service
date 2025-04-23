package repository

import (
	"account-service/internal/domain/account/v1/entity"
	"account-service/pkg/orm"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository interface {
	CountTotalAccounts(ctx context.Context) (int64, error)
	CreateAccount(ctx context.Context, payload *entity.Account) error
	GetAccountByAccountNumber(ctx context.Context, accountNumber string) (*entity.Account, error)
	GetBalanceByAccountID(ctx context.Context, id string) (float64, error)
	SavingBalance(ctx context.Context, payload *entity.SavingHistory) error
	WithdrawBalance(ctx context.Context, payload *entity.WithdrawHistory) error
}

type accountRepository struct {
	Sql *orm.Provider
	Log *logrus.Logger
}

func (r *accountRepository) CountTotalAccounts(ctx context.Context) (int64, error) {
	var count int64
	// Get the count of all rows in the "accounts" table
	if err := r.Sql.WithContext(ctx).Model(&entity.Account{}).Count(&count).Error; err != nil {
		r.Log.Errorf("[account repository][func: CountTotalAccounts] Failed to count total accounts: %s", err.Error())
		return 0, status.Error(codes.Internal, "Internal Server Error.")
	}
	return count, nil
}

func (r *accountRepository) CreateAccount(ctx context.Context, payload *entity.Account) error {
	return r.Sql.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64

		// Check existing phone
		if err := tx.Model(&entity.Account{}).
			Where("phone = ? AND deleted_at IS NULL", payload.Phone).
			Count(&count).Error; err != nil {
			return status.Error(codes.Internal, "Failed to validate phone.")
		}
		if count > 0 {
			return status.Error(codes.InvalidArgument, "Phone number already exists.")
		}

		// Check existing citizen_id
		if err := tx.Model(&entity.Account{}).
			Where("citizen_id = ? AND deleted_at IS NULL", payload.CitizenID).
			Count(&count).Error; err != nil {
			return status.Error(codes.Internal, "Failed to validate citizen ID.")
		}
		if count > 0 {
			return status.Error(codes.InvalidArgument, "Citizen ID already exists.")
		}

		if err := tx.WithContext(ctx).Create(payload).Error; err != nil {
			r.Log.Errorf("[account repository][func: CreateAccount] Failed to create Account: %s", err.Error())
			return status.Error(codes.Internal, "Internal Server Error.")
		}

		r.Log.Infof("[account repository][func: CreateAccount] Successfully created Account with ID: %s", payload.ID)
		return nil
	})
}

func (r *accountRepository) GetAccountByAccountNumber(ctx context.Context, accountNumber string) (*entity.Account, error) {
	var row entity.Account
	if err := r.Sql.WithContext(ctx).Model(&entity.Account{}).Preload(clause.Associations).Where("account_number = ?", accountNumber).Take(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warnf("[account repository][func: GetAccountByAccountNumber] Account Not Found: %s", err.Error())
			return nil, status.Error(codes.NotFound, "Account Not Found")
		}

		r.Log.Errorf("[account repository][func: GetAccountByAccountNumber] Failed to get account by account_number: %s", err.Error())
		return nil, status.Error(codes.Internal, "Internal Server Error.")
	}

	balance, err := r.GetBalanceByAccountID(ctx, row.ID)
	if err != nil {
		r.Log.Errorf("[account repository][func: GetAccountByAccountNumber] Failed to get balance: %s", err.Error())
		return nil, err
	}
	row.Balance = balance

	return &row, nil
}

func (r *accountRepository) GetBalanceByAccountID(ctx context.Context, id string) (float64, error) {
	var totalSaving, totalWithdraw float64

	// Sum savings using model
	if err := r.Sql.WithContext(ctx).
		Model(&entity.SavingHistory{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("account_id = ?", id).
		Scan(&totalSaving).Error; err != nil {
		r.Log.Errorf("[account repository][GetBalanceByAccountID] Failed to sum savings: %s", err.Error())
		return 0, status.Error(codes.Internal, "Internal Server Error.")
	}

	// Sum withdrawals using model
	if err := r.Sql.WithContext(ctx).
		Model(&entity.WithdrawHistory{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("account_id = ?", id).
		Scan(&totalWithdraw).Error; err != nil {
		r.Log.Errorf("[account repository][GetBalanceByAccountID] Failed to sum withdrawals: %s", err.Error())
		return 0, status.Error(codes.Internal, "Internal Server Error.")
	}

	balance := totalSaving - totalWithdraw
	return balance, nil
}

func (r *accountRepository) SavingBalance(ctx context.Context, payload *entity.SavingHistory) error {
	return r.Sql.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(payload).Error; err != nil {
			r.Log.Errorf("[account repository][func: SavingBalance] Failed to create Saving History: %s", err.Error())
			return status.Error(codes.Internal, "Internal Server Error.")
		}

		r.Log.Infof("[account repository][func: SavingBalance] Successfully created Saving History with ID: %s", payload.ID)
		return nil
	})
}

func (r *accountRepository) WithdrawBalance(ctx context.Context, payload *entity.WithdrawHistory) error {
	return r.Sql.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(payload).Error; err != nil {
			r.Log.Errorf("[account repository][func: WithdrawBalance] Failed to create Withdraw History: %s", err.Error())
			return status.Error(codes.Internal, "Internal Server Error.")
		}

		r.Log.Infof("[account repository][func: WithdrawBalance] Successfully created Withdraw History with ID: %s", payload.ID)
		return nil
	})
}

func NewAccountRepository(
	logger *logrus.Logger,
	sql *orm.Provider,
) AccountRepository {
	return &accountRepository{
		Sql: sql,
		Log: logger,
	}
}
