package service

import (
	"acc/internal/consts"
	acc "acc/internal/consts/account"
	t "acc/internal/consts/transaction"
	"acc/internal/model"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	TradingTime int64           `form:"tradingTime" binding:"required"`
	Type        int             `form:"type" binding:"required"`
	AccountId   int64           `form:"accountId" binding:"required"`
	CpAccountId int64           `form:"cpAccountId" binding:"required"`
	Amount      decimal.Decimal `form:"amount" binding:"required"`
	Remark      string          `form:"remark"`
	ProjectId   int64           `form:"projectId"`
	MemberId    int64           `form:"memberId"`
	SupplierId  int64           `form:"supplierId"`
}

type TransactionBO struct {
	Accounts   []int64 `form:"accounts"`
	CpAccounts []int64 `form:"cpAccounts"`
	ProjectId  int64   `form:"project"`
	MemberId   int64   `form:"memberId"`
	SupplierId int64   `form:"supplierId"`
	StartTime  int64   `form:"startTime"`
	EndTime    int64   `form:"endTime"`
	Remark     string  `form:"remark"`
}

type TransactionVO struct {
	TradingTime int64           `json:"tradingTime"`
	Type        string          `json:"type"`
	Account     string          `json:"account"`
	CpAccount   string          `json:"cpAccount"`
	Project     string          `json:"project"`
	Member      string          `json:"member"`
	Supplier    string          `json:"supplier"`
	Amount      decimal.Decimal `json:"amount"`
	Remark      string          `form:"remark"`
}

type TransactionService struct{}

func (s *TransactionService) InsertBalance(tx *gorm.DB, ledgerId int64, userId int64, accountId int64, amount decimal.Decimal) error {
	current, err := accountDao.Get(ledgerId, accountId)
	if err != nil {
		return errors.WithMessage(err, "get account when adjusting balance")
	}

	now := time.Now().UnixMicro()
	transDO := model.Transaction{
		CreateTime: now,
		LedgerId:   ledgerId,
		Amount:     amount.Abs(),
		UserId:     userId,
	}

	switch current.Type {
	case acc.TypeAsset, acc.TypeReceivables:
		// 如果是正数
		if amount.IsPositive() {
			transDO.AccountId = consts.AccountBalanceIncome
			transDO.CpAccountId = accountId
			transDO.Type = t.TypeIncome
		} else {
			transDO.AccountId = consts.AccountBalanceExpenses
			transDO.CpAccountId = accountId
			transDO.Type = t.TypeExpenses
		}
	case acc.TypeDebt, acc.TypeEquity:
		// 如果是正数
		if amount.IsPositive() {
			transDO.AccountId = consts.AccountBalanceExpenses
			transDO.CpAccountId = accountId
			transDO.Type = t.TypeExpenses
		} else {
			transDO.AccountId = consts.AccountBalanceIncome
			transDO.CpAccountId = accountId
			transDO.Type = t.TypeIncome
		}
	}
	if err := transactionDao.Insert(tx, &transDO); err != nil {
		return err
	}
	return nil
}

func (s *TransactionService) Insert(tx *gorm.DB, ledgerId int64, userId int64, transBO *Transaction) error {
	now := time.Now().UnixMicro()
	transDO := model.Transaction{
		CreateTime:  now,
		TradingTime: transBO.TradingTime,
		LedgerId:    ledgerId,
		Type:        transBO.Type,
		AccountId:   transBO.AccountId,
		CpAccountId: transBO.CpAccountId,
		Amount:      transBO.Amount,
		Remark:      transBO.Remark,
		ProjectId:   transBO.ProjectId,
		MemberId:    transBO.MemberId,
		SupplierId:  transBO.SupplierId,
		UserId:      userId,
	}

	if err := transactionDao.Insert(tx, &transDO); err != nil {
		return err
	}
	return nil
}

func (s *TransactionService) List(ledgerId int64, trans *TransactionBO) ([]*TransactionVO, error) {
	dos, err := transactionDao.List(ledgerId, trans)
	if err != nil {
		return nil, err
	}
	accounts, err := accountDao.List(ledgerId)
	if err != nil {
		return nil, err
	}
	projects, err := projectDao.ListAll(ledgerId)
	if err != nil {
		return nil, err
	}
	members, err := memberDao.ListAll(ledgerId)
	if err != nil {
		return nil, err
	}
	suppliers, err := supplierDao.ListAll(ledgerId)
	if err != nil {
		return nil, err
	}
	var vos []*TransactionVO
	for _, do := range dos {
		tran := TransactionVO{
			TradingTime: do.TradingTime,
			Type:        translateType(do.Type),
			Account:     translateAccount(do.AccountId, accounts),
			CpAccount:   translateAccount(do.CpAccountId, accounts),
			Project:     translateProject(do.ProjectId, projects),
			Member:      translateMember(do.MemberId, members),
			Supplier:    translateSupplier(do.SupplierId, suppliers),
			Amount:      do.Amount,
			Remark:      do.Remark,
		}
		vos = append(vos, &tran)
	}

	return vos, nil
}

func translateType(typeId int) string {
	return t.Mapper[typeId]
}

func translateAccount(accountId int64, accounts []*model.Account) string {
	for i := 0; i < len(accounts); i++ {
		if accountId == accounts[i].ID {
			return accounts[i].Name
		}
	}
	return ""
}

func translateProject(projectId int64, projects []*model.Project) string {
	for i := 0; i < len(projects); i++ {
		if projectId == projects[i].ID {
			return projects[i].Name
		}
	}
	return ""
}

func translateMember(memberId int64, members []*model.Member) string {
	for i := 0; i < len(members); i++ {
		if memberId == members[i].ID {
			return members[i].Name
		}
	}
	return ""
}

func translateSupplier(supplierId int64, suppliers []*model.Supplier) string {
	for i := 0; i < len(suppliers); i++ {
		if supplierId == suppliers[i].ID {
			return suppliers[i].Name
		}
	}
	return ""
}
