package v1

import (
	"account-service/internal/domain/account/v1/handler"
	"account-service/internal/domain/account/v1/repository"
	"account-service/internal/domain/account/v1/usecase"
	"account-service/pkg/orm"
	accountv1 "account-service/stubs/account/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterServiceServer(
	rpcServer *grpc.Server,
	logger *logrus.Logger,
	sql *orm.Provider,
) {
	accountRepository := repository.NewAccountRepository(logger, sql)

	accountUsecase := usecase.NewAccountUseCase(
		logger,
		accountRepository,
	)

	accountManagementHandler := handler.NewAccountHandler(
		logger,
		accountUsecase,
	)

	accountv1.RegisterAccountServiceServer(rpcServer, accountManagementHandler)

	grpc_health_v1.RegisterHealthServer(rpcServer, health.NewServer())
}
