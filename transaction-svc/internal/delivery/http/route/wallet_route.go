package route

func (r *route) setupWalletRoute() {
	transactionRoute := r.app.Group("/api/wallet", r.authMiddleware, r.creatorMiddleware)
	transactionRoute.Get("/", r.walletControler.GetWallet)
}
