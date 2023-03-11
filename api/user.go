package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/crackz/simple-bank/db/sqlc"
	"github.com/crackz/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
)

type createUserDto struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func newUserResponse(user *db.User) *userResponse {
	return &userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) registerUser(ctx *gin.Context) {
	var createDto createUserDto

	if err := ctx.ShouldBindJSON(&createDto); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(createDto.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       createDto.Username,
		FullName:       createDto.FullName,
		Email:          createDto.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("username or email is duplicated")))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, newUserResponse(&user))
}

type loginUserDto struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"accessToken"`
	*userResponse
}

func (server *Server) loginUser(ctx *gin.Context) {
	var loginUserDto loginUserDto

	if err := ctx.ShouldBindJSON(&loginUserDto); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, loginUserDto.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidUsernameOrPassword))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.CheckPassword(loginUserDto.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrInvalidUsernameOrPassword))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.JwtDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, &loginUserResponse{
		accessToken,
		newUserResponse(&user),
	})

}
