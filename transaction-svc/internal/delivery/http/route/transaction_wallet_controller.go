package route

func (r *route) setupTransactionWalletRoute() {
	transactionRoute := r.app.Group("/api/wallet/transaction", r.authMiddleware)
	transactionRoute.Get("/", r.transactionWalletCtrl.GetAll)
}
