package route

func (r *route) setupBankWalletRoute() {
	bankWalletRoute := r.app.Group("/api/wallet/bank", r.authMiddleware)
	bankWalletRoute.Post("/add", r.bankWalletController.CreateBankWallet)
	bankWalletRoute.Get("/", r.bankWalletController.FindAllBankWallet)
	bankWalletRoute.Delete("/delete/:bankWalletId", r.bankWalletController.DeleteBankWallet)
}
