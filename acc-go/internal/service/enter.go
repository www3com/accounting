package service

import "acc/internal/model"

var (
	accountDao     = new(model.AccountDao)
	ledgerDao      = new(model.LedgerDao)
	memberDao      = new(model.MemberDao)
	projectDao     = new(model.ProjectDao)
	supplierDao    = new(model.SupplierDao)
	userDao        = new(model.UserDao)
	tplLedgerDao   = new(model.TplLedgerDao)
	transactionDao = new(model.TransactionDao)

	ledgerService = new(LedgerService)
)
