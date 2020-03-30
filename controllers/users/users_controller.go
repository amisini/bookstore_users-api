package users

import (
	"fmt"
	"net/http"

	"github.com/amisini/bookstore_users-api/domain/users"
	"github.com/amisini/bookstore_users-api/services"
	"github.com/amisini/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveError := services.CreateUser(user)
	if saveError != nil {
		c.JSON(saveError.Status, saveError)
		return
	}
	fmt.Println(user)
	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement me")
}

func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implement me")
}
