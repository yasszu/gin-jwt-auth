package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"gin-jwt-auth/conf"
	"gin-jwt-auth/jwt"
	"gin-jwt-auth/model"
	"gin-jwt-auth/repository"
	"gin-jwt-auth/util"
)

type IAccountHandler interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Verify(c *gin.Context)
}

type AccountHandler struct {
	accountRepository repository.IAccountRepository
	conf              *conf.Conf
}

func NewAccountHandler(repository repository.IAccountRepository, conf *conf.Conf) *AccountHandler {
	return &AccountHandler{accountRepository: repository, conf: conf}
}

func (h AccountHandler) RegisterRoot(e *gin.Engine) {
	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)
	e.POST("/logout", h.Logout)
}

func (h AccountHandler) RegisterV1(v1 *gin.RouterGroup) {
	v1.GET("/me", h.Me)
}

// Signup POST /signup
func (h *AccountHandler) Signup(c *gin.Context) {
	var form model.SignUpForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account model.Account
	if err := account.Populate(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.accountRepository.CreateAccount(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.signJWT(&account, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Login POST /login
func (h *AccountHandler) Login(c *gin.Context) {
	var form model.LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountRepository.GetAccountByEmail(form.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := util.ComparePassword(account.PasswordHash, form.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.signJWT(account, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout POST /logout
func (h *AccountHandler) Logout(c *gin.Context) {
	util.DeleteAuthorizationCookie(c)
	c.String(http.StatusOK, "Logout success")
}

// Me  GET /v1/me
func (h *AccountHandler) Me(c *gin.Context) {
	accountID := jwt.BindUser(c).AccountID
	account, err := h.accountRepository.GetAccountById(accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.NewAccountResponse(account))
}

func (h AccountHandler) signJWT(account *model.Account, c *gin.Context) (*string, error) {
	token, err := jwt.Sign(account.Email, account.ID, h.conf.JWT.Secret)
	if err != nil {
		return nil, err
	}
	util.SaveAuthorizationCookie(token, c)
	return &token, nil
}
