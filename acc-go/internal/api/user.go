package api

import (
	"acc/internal/consts"
	"acc/internal/model"
	"acc/internal/pkg/r"
	"github.com/gin-gonic/gin"
)

type user struct {
	Nickname  string `form:"nickname" binding:"required"`
	Username  string `form:"username" binding:"required"`
	Password  string `form:"password" binding:"required"`
	Agreement bool   `form:"agreement" binding:"required"`
}

// SignUp 注册
func SignUp(c *gin.Context) {
	var p user
	if err := r.BindAndValid(c, &p); err != nil {
		r.RenderFail(c, err)
		return
	}

	if !p.Agreement {
		r.RenderFail(c, consts.ErrUserDisagreement)
		return
	}

	ok, err := userService.Exist(p.Username)
	if err != nil {
		r.RenderFail(c, err)
		return
	}
	if ok {
		r.RenderFail(c, consts.ErrUserDuplicateUsername)
		return
	}

	user := model.User{
		Nickname:  p.Nickname,
		Username:  p.Username,
		Password:  p.Password,
		State:     consts.UserStateNormal,
		Agreement: 1,
	}

	err = userService.SignUp(&user)
	r.Render(c, nil, err)
}
