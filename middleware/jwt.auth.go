package middleware

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/luvnyen/student-research-paper-storage/pkg/utils"
	"github.com/luvnyen/student-research-paper-storage/service"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			response := utils.BuildErrorResponse("Failed to process request", "You are not logged in", utils.EmptyObj{})
			c.AbortWithStatusJSON(401, response)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtService.GetSecret()), nil
		})
		if err != nil {
			response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
			c.AbortWithStatusJSON(401, response)
			return
		}

		validatedToken, err := jwtService.ValidateToken(token.Raw)
		if err != nil {
			response := utils.BuildErrorResponse("Invalid token", "Invalid token", utils.EmptyObj{})
			c.JSON(401, response)
			return
		}

		if validatedToken.Valid {
			claims := validatedToken.Claims.(jwt.MapClaims)
			log.Println("User ID: ", claims["user_id"])
		} else {
			log.Println(err)
			response := utils.BuildErrorResponse("Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(401, response)
		}
	}
}
