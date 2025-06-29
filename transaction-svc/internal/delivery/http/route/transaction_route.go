package route

func (r *route) setupTransactionRoute() {
	transactionRoute := r.app.Group("/api/transaction", r.authMiddleware)
	transactionRoute.Post("/create", r.creatorMiddleware, r.transactionController.CreateTransaction)
	transactionRoute.Get("/:transactionID", r.transactionController.GetUserTransactionWithDetail)
	transactionRoute.Get("/", r.transactionController.GetAllUserTransaction)

	transactionRouteV2 := r.app.Group("/api/v2/transaction", r.authMiddleware)
	transactionRouteV2.Post("/create", r.creatorMiddleware, r.transactionController.CreateTransactionV2)
}

func (r *route) setupWebhookRoute() {
	webhookRoute := r.app.Group("/api/webhook")
	webhookRoute.Post("/notify", r.transactionController.Notify)
}
