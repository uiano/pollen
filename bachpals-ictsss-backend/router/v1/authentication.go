package v1

import (
	"context"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rithujohn191/go-oidc"
	"github.com/spf13/viper"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/httputils"
	"golang.org/x/oauth2"
)

// HandleSignIn godoc
// @Summary     Handles Login
// @Description Handler for authenticating the login of a user
// @Tags        oauth2
// @Accept      json
// @Produce     json
// @Success     200 {object}    nil
// @Failure     400 {object}    nil
// @Failure		406	{object}	nil
// @Failure     500 {object}    nil
// @Router      /oauth2/provider    [get]
func HandleSignIn(c *gin.Context) {
	oauth2Config := oauth2.Config{
		ClientID:     viper.GetString("IKT_STACK_CLIENT_ID"),
		ClientSecret: viper.GetString("IKT_STACK_CLIENT_SECRET"),
		RedirectURL:  viper.GetString("IKT_STACK_REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   viper.GetString("IKT_STACK_AUTH_URL"),
			TokenURL:  viper.GetString("IKT_STACK_TOKEN_URL"),
			AuthStyle: 1,
		},
		Scopes: []string{"userinfo-name", "userid-feide", "email", "openid"},
	}

	// Generate random state key
	buf := make([]byte, 16)
	rand.Read(buf)
	state := b64.URLEncoding.EncodeToString(buf)

	c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state))
	return
}

// HandleLogout godoc
// @Summary     Handles Logout
// @Description Handler for a users logout
// @Tags        oauth2
// @Accept      json
// @Produce     json
// @Success     200 {object}    nil
// @Failure     400 {object}    nil
// @Failure		406	{object}	nil
// @Failure     500 {object}    nil
// @Router      /oauth2/logout    [get]
func HandleLogout(c *gin.Context) {
	c.Redirect(http.StatusFound, viper.GetString("IKT_STACK_LOGOUT_URL"))
	return
}

// HandleUserdata godoc
// @Summary     Generates a console
// @Description Generates a console for a VM
// @Tags        oauth2
// @Accept      json
// @Produce     json
// @Success     200 {object}    nil
// @Failure     400 {object}    nil
// @Failure		406	{object}	nil
// @Failure     500 {object}    nil
// @Router      /oauth2/userdata    [post]
func HandleUserdata(c *gin.Context) {
	authorization := c.GetHeader("Authorization")

	if len(authorization) <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", viper.GetString("IKT_STACK_AUTH_MIDDLEWARE_URL"), nil)

	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reaching endpoint", err)
		return
	}

	req.Header.Add("Authorization", authorization)
	resp, err := client.Do(req)

	// Any status code that is not 200 OK should return false
	if resp.StatusCode == http.StatusUnauthorized {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Missing authorization header", err)
		return
	}

	if resp.StatusCode != 200 {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong!", err)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading data", err)
		return
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &jsonData)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while encoding data", err)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", jsonData)
	return
}

// HandleCallback godoc
// @Summary     Handles Redirect
// @Description Handler for the redirecting of a user after a login
// @Tags        oauth2
// @Accept      json
// @Produce     json
// @Success     200 {object}    nil
// @Failure     400 {object}    nil
// @Failure		406	{object}	nil
// @Failure     500 {object}    nil
// @Router      /oauth2/redirect    [get]
func HandleCallback(c *gin.Context) {
	oauth2Config := oauth2.Config{
		ClientID:     viper.GetString("IKT_STACK_CLIENT_ID"),
		ClientSecret: viper.GetString("IKT_STACK_CLIENT_SECRET"),
		RedirectURL:  viper.GetString("IKT_STACK_REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   viper.GetString("IKT_STACK_AUTH_URL"),
			TokenURL:  viper.GetString("IKT_STACK_TOKEN_URL"),
			AuthStyle: 1,
		},
		Scopes: []string{"userinfo-name", "userid-feide", "email", "openid"},
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, viper.GetString("IKT_STACK_AUTHORITY"))
	if err != nil {
		// handle error
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, c.Request.URL.Query().Get("code"))
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to exchange token", err.Error())
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: oauth2Config.ClientID,
	}

	verifier := provider.Verifier(oidcConfig)

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "No id_token field in oauth2 token.", err.Error())
		return
	}

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Failed to verify ID Token:", err.Error())
		return
	}

	respIdToken := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&respIdToken.IDTokenClaims); err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	//httputils.ResponseJson(c, http.StatusOK, "", respIdToken)
	encodedToken := b64.StdEncoding.EncodeToString([]byte(respIdToken.OAuth2Token.AccessToken))
	formattedUrl := fmt.Sprintf("%s?t=%s", viper.GetString("IKT_STACK_FRONTEND_URL"), url.QueryEscape(encodedToken))
	c.Redirect(http.StatusPermanentRedirect, formattedUrl)
	return
}
