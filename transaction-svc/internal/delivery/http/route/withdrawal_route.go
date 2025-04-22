package route

func (r *route) setupWithdrawalRoute() {
	reviewRoute := r.app.Group("/api/withdrawal", r.authMiddleware)
	reviewRoute.Patch("/add", r.withdrawalController.CreateWithdrawal)
	reviewRoute.Get("/", r.withdrawalController.FindAllWithdrawal)
	reviewRoute.Get("/:withdrawalId", r.withdrawalController.FindWithdrawalById)
	// reviewRoute.Delete("/delete", r.withdrawalController.DeleteWithdrawal)
}
