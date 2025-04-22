package route

func (r *route) setupBankRoute() {
	bankRoute := r.app.Group("/api/bank", r.authMiddleware)
	bankRoute.Post("/create", r.bankController.CreateBank)
	bankRoute.Delete("/delete/:bankId", r.bankController.DeleteBank)
	bankRoute.Get("/", r.bankController.FindAllBank)
	bankRoute.Get("/:bankId", r.bankController.FindBankById)
}
