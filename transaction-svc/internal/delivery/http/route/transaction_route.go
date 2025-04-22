package route

func (r *route) setupTransactionRoute() {
	transactionRoute := r.app.Group("/api/transaction", r.authMiddleware)
	transactionRoute.Post("/create", r.transactionController.CreateTransaction)
}

func (r *route) setupWebhookRoute() {
	webhookRoute := r.app.Group("/api/webhook")
	webhookRoute.Post("/notify", r.transactionController.Notify)
}
