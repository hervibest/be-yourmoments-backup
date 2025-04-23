package route

import (
	http "be-yourmoments/transaction-svc/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

type Route interface {
	SetupRoute()
}

type route struct {
	app                   *fiber.App
	transactionController http.TransactionController
	bankController        http.BankController
	bankWalletController  http.BankWalletController
	reviewController      http.ReviewController
	withdrawalController  http.WithdrawalController
	walletControler       http.WalletController
	transactionWalletCtrl http.TransactionWalletController
	authMiddleware        fiber.Handler
}

func NewRoute(app *fiber.App, transactionController http.TransactionController,
	bankController http.BankController,
	bankWalletController http.BankWalletController,
	reviewController http.ReviewController,
	withdrawalController http.WithdrawalController,
	walletControler http.WalletController,
	transactionWalletCtrl http.TransactionWalletController,
	authMiddleware fiber.Handler) Route {
	return &route{
		app:                   app,
		transactionController: transactionController,
		bankController:        bankController,
		bankWalletController:  bankWalletController,
		reviewController:      reviewController,
		withdrawalController:  withdrawalController,
		walletControler:       walletControler,
		transactionWalletCtrl: transactionWalletCtrl,
		authMiddleware:        authMiddleware,
	}
}

func (r *route) SetupRoute() {
	r.setupTransactionRoute()
	r.setupWebhookRoute()
	r.setupBankRoute()
	r.setupBankWalletRoute()
	r.setupReviewRoute()
	r.setupWithdrawalRoute()
	r.setupWalletRoute()
	r.setupTransactionWalletRoute()
}
