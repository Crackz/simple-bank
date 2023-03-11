package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/crackz/simple-bank/token"
	"github.com/gin-gonic/gin"
)

var (
	authorizationHeaderName = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "auth_payload"
)

func authMiddleware(tokenMaster token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderName)

		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header is not provided")))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header format is invalid")))
			return
		}

		gotAuthorizationType := strings.ToLower(fields[0])
		if gotAuthorizationType != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization type is unsupported")))
			return
		}

		token := fields[1]
		payload, err := tokenMaster.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
