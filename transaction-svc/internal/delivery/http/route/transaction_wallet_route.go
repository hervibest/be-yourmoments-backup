package route

func (r *route) setupTransactionWalletRoute() {
	transactionRoute := r.app.Group("/api/wallet/transaction", r.authMiddleware, r.creatorMiddleware, r.walletMiddleware)
	transactionRoute.Get("/", r.transactionWalletCtrl.GetAll)
}
