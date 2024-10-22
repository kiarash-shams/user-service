package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"user-service/api/helper"
	"user-service/config"
	"user-service/constant"
	"user-service/pkg/service_errors"
	"user-service/services"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	var TokenService = services.NewTokenService(cfg)

	return func(c *gin.Context) {
		var err error
		claimMap := map[string]interface{}{}
		auth := c.GetHeader(constant.AuthorizationHeaderKey)

		token := strings.Split(auth, " ")
		fmt.Println(token)
		if auth == "" || len(token) < 2 {
			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
		} else {
			claimMap, err = TokenService.GetClaims(token[1])
			fmt.Println("claimMap: ", claimMap)
			fmt.Println("claimMap user_id: ", claimMap["UserId"])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
				default:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
				}
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.GenerateBaseResponseWithError(
				nil, false, helper.AuthError, err,
			))
			return
		}



		c.Set(constant.UserIdKey, claimMap[constant.UserIdKey])
		c.Set(constant.UID, claimMap[constant.UID])
		c.Set(constant.Level, claimMap[constant.Level])
		c.Set(constant.Otp, claimMap[constant.Otp])
		c.Set(constant.EmailKey, claimMap[constant.EmailKey])
		c.Set(constant.State, claimMap[constant.State])
		c.Set(constant.RolesKey, claimMap[constant.RolesKey])
		c.Set(constant.ExpireTimeKey, claimMap[constant.ExpireTimeKey])
		c.Set(constant.IssuedAtKey, claimMap[constant.IssuedAtKey])
		

		c.Next()
	}
}

func Authorization(validRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Keys) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
			return
		}
		rolesVal := c.Keys[constant.RolesKey]

		if rolesVal == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
			return
		}
		roles := rolesVal.([]interface{})
		val := map[string]int{}
		for _, item := range roles {
			val[item.(string)] = 0
		}

		for _, item := range validRoles {
			if _, ok := val[item]; ok {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, helper.GenerateBaseResponse(nil, false, helper.ForbiddenError))
	}
}
