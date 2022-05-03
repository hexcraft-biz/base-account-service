package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/models"
	"github.com/hexcraft-biz/controller"
)

type Users struct {
	*controller.Prototype
}

func NewUsers(cfg *config.Config) *Users {
	return &Users{
		Prototype: controller.New("users", cfg.DB),
	}
}

type TargetUser struct {
	ID string `uri:"id" binding:"required"`
}

func (ctrl *Users) Get() gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO get JWT from header
		var targetUser TargetUser
		if err := c.ShouldBindUri(&targetUser); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest), "results": err.Error()})
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByID(targetUser.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError), "results": err.Error()})
		} else {
			if absRes, absErr := entityRes.GetAbsUser(); absErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError), "results": absErr.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK), "results": absRes})
			}
		}
	}
}

func (ctrl *Users) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}
}

func (ctrl *Users) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}
}
