package middleware

import (
	token "github.com/wlady3190/ecommerce/tokens"
	"net/http"
	"github.com/gin-gonic/gin"

)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("Authorization")
		if ClientToken == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no header autorization"})
			c.Abort()
			return
		}
		if len(ClientToken) > 7 && ClientToken[:7] == "Bearer " {
            ClientToken = ClientToken[7:]
			return
        }
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next() //procede al siguiente paso, como llamar a la api
	}
	
}


