package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luvnyen/student-research-paper-storage/pkg/dto"
	"github.com/luvnyen/student-research-paper-storage/pkg/models"
	"github.com/luvnyen/student-research-paper-storage/pkg/utils"
	"github.com/luvnyen/student-research-paper-storage/service"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Failed to process request", errDTO.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	_, err := ctx.Request.Cookie("token")
	if err == nil {
		response := utils.BuildErrorResponse("Failed to process request", "You are already logged in", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(models.Student); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken

		cookie := &http.Cookie{
			Name:       "token",
			Value:      generatedToken,
			Expires:    time.Now().Add(time.Hour * 24),
			Path:       "/",
			HttpOnly:   true,
			Secure:     false,
			SameSite:   http.SameSiteStrictMode,
			MaxAge:     86400,
			RawExpires: "86400",
		}
		http.SetCookie(ctx.Writer, cookie)

		response := utils.BuildResponse(true, "Login successful", v)
		ctx.JSON(200, response)
		return
	}
	response := utils.BuildErrorResponse("Login failed", "Invalid email or password", utils.EmptyObj{})
	ctx.AbortWithStatusJSON(401, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Failed to process request", errDTO.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := utils.BuildErrorResponse("Failed to process request", "Email already exists", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
	} else {
		createdStudent := c.authService.CreateStudent(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdStudent.ID, 10))
		createdStudent.Token = token
		response := utils.BuildResponse(true, "Register successful", createdStudent)
		ctx.JSON(200, response)
	}
}

func (c *authController) Logout(ctx *gin.Context) {
	_, err := ctx.Request.Cookie("token")
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", "You are not logged in", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(ctx.Writer, cookie)
	response := utils.BuildResponse(true, "Logout successful", utils.EmptyObj{})
	ctx.JSON(200, response)
}
