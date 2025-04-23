package route

func (r *route) setupWalletRoute() {
	transactionRoute := r.app.Group("/api/wallet", r.authMiddleware)
	transactionRoute.Get("/", r.walletControler.GetWallet)
}
