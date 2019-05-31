package middleware

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"

	"github.com/atlanhq/sc-go-message/boot"
)

var (
	ORGID = "orgId"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserID     string          `bson:"userId" json:"userId"`
	OrgID      string          `bson:"orgId" json:"orgId"`
	Name       string          `bson:"name" json:"name"`
	AuthID     string          `bson:"authId" json:"authId"`
	OrgName    string          `bson:"orgName" json:"orgName"`
	OrgSlug    string          `bson:"orgSlug" json:"orgSlug"`
	Email      string          `bson:"email" json:"email"`
	GroupArray []bson.ObjectId `bson:"groupArray" json:"groupArray"`
	OrgAlias   string          `bson:"orgAlias" json:"orgAlias"`
	UserAlias  string          `bson:"userAlias" json:"userAlias"`
}

type SCContext struct {
	echo.Context
	Claims *JWTClaims
	Config *boot.Config
}

func (c *SCContext) GetAuthId() string {
	return c.Claims.AuthID
}

func (c *SCContext) GetUserId() string {
	return c.Claims.UserID
}

func (c *SCContext) GetOrgId() string {
	return c.Claims.OrgID
}

func (c *SCContext) GetOrgName() string {
	return c.Claims.OrgName
}

func (c *SCContext) GetOrgSlug() string {
	return c.Claims.OrgSlug
}

func (c *SCContext) GetUserName() string {
	return c.Claims.Name
}

func (c *SCContext) GetUserEmail() string {
	return c.Claims.Email
}

func (c *SCContext) GetGroupArray() []bson.ObjectId {
	return c.Claims.GroupArray
}

func (c *SCContext) GetRequestMethod() string {
	return c.Request().Method
}

func (c *SCContext) isAuth() bool {
	return c.GetUserId() != ""
}

func (c *SCContext) GetSessionToken() string {
	return c.Get("sessionToken").(string)
}

func (c *SCContext) GetContextVariable(variable string) string {

	switch variable {
	case ORGID:
		return c.Claims.OrgID
	default:
		return ""
	}
}
