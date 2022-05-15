package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/models"
	"github.com/hexcraft-biz/controller"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	*controller.Prototype
}

func NewUsers(cfg config.ConfigInterface) *Users {
	return &Users{
		Prototype: controller.New("users", cfg.GetDB()),
	}
}

type TargetUser struct {
	ID string `uri:"id" binding:"required"`
}

func (ctrl *Users) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var targetUser TargetUser
		if err := c.ShouldBindUri(&targetUser); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByID(targetUser.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				if absRes, absErr := entityRes.GetAbsUser(); absErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
					return
				} else {
					c.AbortWithStatusJSON(http.StatusOK, absRes)
					return
				}
			}
		}
	}
}

type updatePwdParams struct {
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Users) UpdatePwd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var targetUser TargetUser
		if err := c.ShouldBindUri(&targetUser); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var params updatePwdParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		usersEngine := models.NewUsersTableEngine(ctrl.DB)

		if entityRes, err := usersEngine.GetByID(targetUser.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				// TODO next version about password log
				saltedPwd := append([]byte(params.Password), entityRes.Salt...)
				compareErr := bcrypt.CompareHashAndPassword(entityRes.Password, saltedPwd)
				if compareErr == nil {
					c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
					return
				}

				if _, err := usersEngine.ResetPwd(entityRes.ID, params.Password, entityRes.Salt); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				} else {
					c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": http.StatusText(http.StatusNoContent)})
					return
				}
			}
		}
	}
}

type updateStatusParams struct {
	Status string `json:"status" binding:"required,oneof='enabled' 'disabled' 'suspended'"`
}

func (ctrl *Users) UpdateStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var targetUser TargetUser
		if err := c.ShouldBindUri(&targetUser); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var params updateStatusParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		usersEngine := models.NewUsersTableEngine(ctrl.DB)

		if entityRes, err := usersEngine.GetByID(targetUser.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				if _, err := usersEngine.UpdateStatus(entityRes.ID, params.Status); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				} else {
					if entityRes, err := usersEngine.GetByID(targetUser.ID); err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
						return
					} else {
						if entityRes == nil {
							c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
							return
						} else {
							if absRes, absErr := entityRes.GetAbsUser(); absErr != nil {
								c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
								return
							} else {
								c.AbortWithStatusJSON(http.StatusOK, absRes)
								return
							}
						}
					}
				}
			}
		}
	}
}
