package route

func (r *route) setupWithdrawalRoute() {
	reviewRoute := r.app.Group("/api/withdrawal", r.authMiddleware, r.creatorMiddleware, r.walletMiddleware)
	reviewRoute.Post("/create", r.withdrawalController.CreateWithdrawal)
	reviewRoute.Get("/", r.withdrawalController.FindAllWithdrawal)
	reviewRoute.Get("/:withdrawalId", r.withdrawalController.FindWithdrawalById)
	// reviewRoute.Delete("/delete", r.withdrawalController.DeleteWithdrawal)
}
