package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"magego/course-33/cmdb/base/controllers/base"
	"magego/course-33/cmdb/base/errors"
	"magego/course-33/cmdb/config"
	"magego/course-33/cmdb/forms"
	"magego/course-33/cmdb/services"
	"net/http"
)

// AuthController 负责认证的认证控制器
type AuthController struct {
	base.BaseController
}

// Login 认证登录
func (a *AuthController) Login() {

	SessionKey := beego.AppConfig.DefaultString("auth::SessionKey", "user")
	LoginController := beego.AppConfig.DefaultString("auth::HomeAction", "HomeController.Index")

	sessionUser := a.GetSession(SessionKey)
	if sessionUser != nil {
		a.Redirect(beego.URLFor(LoginController), http.StatusFound)
		return
	}

	form := &forms.LoginForm{}
	errs := errors.NewErrors()
	//Get请求直接加载页面
	//Post请求记性数据验证（成功/失败）
	if a.Ctx.Input.IsPost() {
		config.Cache.Incr("login")
		fmt.Println(config.Cache.Get("login"))
		//获取用户提交数据
		if err := a.ParseForm(form); err == nil {
			user := services.UserService.GetByName(form.Name)
			if user == nil {
				//用户不存在
				errs.AddError("default", "用户名或密码错误")
				beego.Error(fmt.Sprintf("用户认证失败：%s", form.Name))
			} else if user.ValidPassWord(form.Password) {
				beego.Informational(fmt.Sprintf("用户认证成功：%s", form.Name))
				//用户密码正确
				//记录用户状态(session 记录在服务器端)
				a.SetSession(SessionKey, user.ID)
				a.Redirect(beego.URLFor(LoginController), http.StatusFound)
			} else {
				//用户密码不正确
				errs.AddError("default", "用户名或密码错误")
				beego.Error(fmt.Sprintf("用户认证失败：%s", form.Name))
			}
		} else {
			errs.AddError("default", "用户名或密码错误")
		}

	}

	a.Data["form"] = form
	a.Data["errors"] = errs
	a.Data["xsrf_token"] = a.XSRFToken()
	//定义加载界面
	a.TplName = "auth/login.html"
}

func (a *AuthController) Logout() {
	a.DestroySession()
	action := beego.AppConfig.DefaultString("auth::LogoutController", "AuthController.Login")
	a.Redirect(beego.URLFor(action), http.StatusFound)
}
