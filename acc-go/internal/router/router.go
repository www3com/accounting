package routers

import (
	"acc/internal/api"
	"acc/internal/middleware"
	"acc/internal/pkg/r"
	"github.com/gin-gonic/gin"
	"github.com/upbos/go-saber/e"
	"net/http"
)

func InitRouter() *gin.Engine {
	gin.SetMode("release")
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.NoRoute(func(c *gin.Context) {
		r.RenderFail(c, e.ErrNotFound)
	})

	// 静态资源
	engine.StaticFS("/web", http.Dir("./web"))

	// 注册和登录
	engine.GET("api/sign-in", api.SignIn)  // 登录
	engine.POST("api/sign-up", api.SignUp) // 注册

	// 权限校验
	root := engine.Group("/api")
	root.Use(middleware.Auth())

	// 账本
	root.GET("/ledger", api.ListLedger) // 查询账本

	// 账户
	root.GET("/accounts", api.ListAccounts)                      // 查询账户列表
	root.GET("/account/overview", api.Overview)                  // 查询账号概要信息
	root.GET("/account/incomes", api.ListIncomes)                //  查询收入科目
	root.GET("/account/expenses", api.ListExpenses)              // 查询支出科目
	root.GET("/account/income-expenses", api.ListIncomeExpenses) // 查询收入支出科目
	root.POST("/account", api.SaveAccount)                       // 创建账户或者更新账户
	root.DELETE("/account", api.DeleteAccount)                   // 删除账户
	root.PUT("/account/name", api.UpdateName)                    // 调整账户名称
	root.PUT("/account/remark", api.UpdateRemark)                // 调整账户描述
	root.PUT("/account/balance", api.UpdateBalance)              // 调整账户余额

	// 交易
	root.GET("/transactions", api.ListTransaction)           // 查询交易
	root.GET("/transaction/total", api.ListTotalTransaction) // 查询总额、结余
	root.POST("/transaction", api.SaveTransaction)           // 保存交易

	// 成员
	root.GET("/members", api.ListMembers) // 查询成员列表

	// 项目
	root.GET("/projects", api.ListProjects) // 查询项目列表

	// 供应商
	root.GET("/suppliers", api.ListSupplier) // 查询供应商列表

	return engine
}
