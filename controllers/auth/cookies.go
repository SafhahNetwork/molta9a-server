package authController

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetGinCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("molta9a_biscuit", token, 60*60*STAY_LOGGED_HOURS, "/", "localhost", false, true)
}
