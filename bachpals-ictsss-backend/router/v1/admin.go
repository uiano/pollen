package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database/repositories"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/httputils"
)

type RequestBodyAdminCreate struct {
	Name string `json:"name"`
	RequestBodyUserId
}

type RequestBodyAdminUpdate struct {
	Name      string `json:"name"`
	UserId    string `json:"user_id"`
	UpdatedId string `json:"updated_id"`
}

func IsAdmin(c *gin.Context) bool {
	id := c.MustGet("user_id")

	// Check if user performing the request is an admin
	admin := repositories.ReadAdminById(id.(string))
	if admin == nil {
		return false
	}
	return true
}

// GetAdministrator godoc
// @Summary		Fetches an admin user
// @Description	Fetches an admin user from the DB by id
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param		adminId	path	string	true	"Admin ID"
// @Success     200 {object}	database.Admin
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /admin/:id	[get]
func GetAdministrator(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	adminId := c.Param("id")

	user := strings.Split(adminId, "@")
	formattedUser := user[0] + "@uia.no"

	if len(adminId) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Invalid id", nil)
		return
	}

	adminUser := repositories.ReadAdminById(formattedUser)

	if adminUser == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong reading database!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", adminUser)
	return
}

// ListAdministrators godoc
// @Summary		Fetches all admin users
// @Description	Fetches all admin users from DB
// @Tags        admin
// @Accept      json
// @Produce     json
// @Success     200 {array}	[]database.Admin
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /admin/	[get]
func ListAdministrators(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	admins := repositories.ReadAdmins()

	if admins == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong reading database!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", admins)
	return
}

// AddAdministrator godoc
// @Summary		Adds a new admin
// @Description	Adds a new admin user to DB
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param		requestBody	body	RequestBodyAdminCreate	true	"Request Body"
// @Success     200 {array}	[]database.Admin
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router		/admin/	[post]
func AddAdministrator(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	var requestBody RequestBodyAdminCreate

	err := c.BindJSON(&requestBody)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Invalid json body", nil)
		return
	}

	insertAdmin := repositories.InsertAdmin(requestBody.UserId, requestBody.Name)

	if insertAdmin == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong reading database!", nil)
		return
	}

	admins := repositories.ReadAdmins()

	if admins == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong while reading data!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", admins)
	return
}

// DelAdministrator godoc
// @Summary		Deletes an admin
// @Description	Deletes an admin from DB
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param		requestBody	body	RequestBodyUserId	true	"Request Body"
// @Success     200 {array}	[]database.Admin
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router		/admin/	[delete]
func DelAdministrator(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	var requestBody RequestBodyUserId

	err := c.BindJSON(&requestBody)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Invalid json body", nil)
		return
	}

	admins := repositories.ReadAdmins()

	if admins == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong reading database!", nil)
		return
	}

	if len(admins) == 1 {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Cannot delete last system administrator!", nil)
		return
	}

	count := repositories.DeleteAdminById(requestBody.UserId)

	if count == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong reading database!", nil)
		return
	}

	admins = repositories.ReadAdmins()

	if admins == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong reading database!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "Administrator deleted successfully!", admins)
	return
}

// UpdateAdministrator godoc
// @Summary		Updates an admin
// @Description	Updates an admin user in DB
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param		requestBody	body	RequestBodyAdminUpdate	true	"Request Body"
// @Success     200 {array}	[]database.Admin
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /admin/	[put]
func UpdateAdministrator(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	var requestBody RequestBodyAdminUpdate
	err := c.BindJSON(&requestBody)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Invalid json body", nil)
		return
	}

	isUpdated := repositories.UpdateAdminById(requestBody.UserId, requestBody.UpdatedId, requestBody.Name)
	if !isUpdated {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong while updating!", nil)
		return
	}

	admins := repositories.ReadAdmins()

	if admins == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Something went wrong while reading data!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", admins)
	return
}
