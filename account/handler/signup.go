package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/weslleyrsr/auth-engine/account/model"
	"github.com/weslleyrsr/auth-engine/account/model/apperrors"
	"log"
	"net/http"
)

// signupReq is not exported, hence the lowercase name, it is used for validation and json marshalling
type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Signup handler
func (h *Handler) Signup(c *gin.Context) {
	// Define a variable to which we'll bind incoming json body, {email, password}
	var req signupReq

	// Bind incoming JSON to struct and check for validation errors
	if ok := BindData(c, &req); !ok {
		return
	}

	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.UserService.Signup(c, user)

	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// create token pair as strings
	tokens, err := h.TokenService.NewPairFromUser(c, user, "")

	if err != nil {
		log.Printf("Failed to create token pair: %v\n", err.Error())

		// may eventually implement rollback logic here
		// meaning, if we fail to create tokens after creating a user we make sure to clear/delete the created user in the database

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
