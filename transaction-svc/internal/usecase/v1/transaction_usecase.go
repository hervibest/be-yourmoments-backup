package v1

// Transaction Wrapper

// // Penambahan
// err := repository.BeginTransaction(u.db, ctx, u.logs, func(tx *sqlx.Tx) error {
// 	if err := u.transactionRepository.UpdateCallback(ctx, tx, updateTransaction); err != nil {
// 		return helper.WrapInternalServerError(u.logs, "failed to update transaction callback in database", err)
// 	}
// 	return nil
// })
// if err != nil {
// 	return err
// }
