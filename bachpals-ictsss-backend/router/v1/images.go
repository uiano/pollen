package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/spf13/viper"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database/repositories"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/gopher"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/httputils"
)

type ImageIdStruct struct {
	Id string `json:"id"`
}

type GetServerImagesMapping struct {
	ImageId string `bson:"id"`
	Name    string `bson:"name"`
}

type AddImageStruct struct {
	Published             string `json:"published"`
	ImageId               string `json:"image_id"`
	ImageName             string `json:"image_name"`
	ImageDescription      string `json:"image_description"`
	ImageDisplayName      string `json:"image_display_name"`
	ImageConfig           string `json:"image_config"`
	ImageReadRootPassword bool   `json:"image_read_root_password"`
}

type UpdateImageStruct struct {
	Id                    string `json:"id"`
	Published             string `json:"published"`
	ImageId               string `json:"image_id"`
	ImageName             string `json:"image_name"`
	ImageDescription      string `json:"image_description"`
	ImageDisplayName      string `json:"image_display_name"`
	ImageConfig           string `json:"image_config"`
	ImageReadRootPassword bool   `json:"image_read_root_password"`
}

type PublishedImagesStruct struct {
	ImageId          string `bson:"image_id"`
	ImageDisplayName string `bson:"image_display_name"`
}

// GetServerImages godoc
// @Summary     Fetches images
// @Description Fetches images from OpenStack
// @Tags        image
// @Accept      json
// @Produce     json
// @Success     200 {object}    []database.Images
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /image/server   [get]
func GetServerImages(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	client := gopher.GetClient()

	allPages, err := images.List(client, nil).AllPages()
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Could not read images!", nil)
		return
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		fmt.Println("Could not load images from instance!", err)
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Could not extract information from server!", nil)
		return
	}

	var imagesList []GetServerImagesMapping
	for _, v := range allImages {
		imagesList = append(imagesList, GetServerImagesMapping{
			Name:    v.Name,
			ImageId: v.ID,
		})
	}

	httputils.ResponseJson(c, http.StatusOK, "", imagesList)
	return
}

// GetImage godoc
// @Summary Retrieves an image
// @Description Retrieves a specific image
// @Tags        image
// @Accept      json
// @Produce     json
// @Param       id  path    string  true    "Image ID"
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Success     501 {object}    nil
// @Router      /image/:id  [get]
func GetImage(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	id := c.Param("id")

	if len(id) == 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Invalid id", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusNotImplemented, "", nil)
	return
}

// GetImages godoc
// @Summary     Retrieve list of images
// @Description Retrieves a list of images that can be used by admins
// @Tags        image
// @Accept      json
// @Produce     json
// @Success     200 {object}    []database.Images
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /image/ [get]
func GetImages(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	imagesList := repositories.GetImages()
	if imagesList == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading images!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", imagesList)
	return
}

// AddImage godoc
// @Summary     Add a new image
// @Description Adds a new image for use in the service
// @Tags        image
// @Accept      json
// @Produce     json
// @Param       image   body    AddImageStruct  true    "Request Body"
// @Success     200 {object}    []database.Images
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /image/ [post]
func AddImage(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	var image AddImageStruct
	err := c.BindJSON(&image)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Bad request, invalid json data!", nil)
		return
	}

	data := make(map[string]interface{})
	data["ImageId"] = image.ImageId
	data["ImageName"] = image.ImageName
	data["ImageDescription"] = image.ImageDescription
	data["ImageDisplayName"] = image.ImageDisplayName
	data["Published"] = image.Published
	data["ImageConfig"] = image.ImageConfig
	data["ImageReadRootPassword"] = image.ImageReadRootPassword

	inserted := repositories.InsertImage(data)
	if inserted == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while inserting image!", nil)
		return
	}

	imagesList := repositories.GetImages()
	if imagesList == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading images!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", imagesList)
	return
}

// DeleteImage godoc
// @Summary     Deletes an image
// @Description Deletes an image from the service
// @Tags        image
// @Accept      json
// @Produce     json
// @Param       body    body    ImageIdStruct   true    "Request Body"
// @Success     200 {object}    []database.Images
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /image/ [delete]
func DeleteImage(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	var body ImageIdStruct
	err := c.BindJSON(&body)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Bad request, invalid json data!", nil)
		return
	}

	deleted := repositories.DeleteImageById(body.Id)
	if deleted == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while deleting images!", nil)
		return
	}

	imagesList := repositories.GetImages()
	if imagesList == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading images!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", imagesList)
	return
}

// UpdateImage godoc
// @Summary     Updates an image
// @Description Handles the updates on an image
// @Tags        image
// @Accept      json
// @Produce     json
// @Param       image   body    UpdateImageStruct   true    "Request Body"
// @Success     200 {object}    database.Images
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /image/ [put]
func UpdateImage(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	var image UpdateImageStruct
	err := c.BindJSON(&image)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Bad request, invalid json data!", nil)
		return
	}

	data := make(map[string]interface{})
	data["Id"] = image.Id
	data["ImageId"] = image.ImageId
	data["ImageName"] = image.ImageName
	data["ImageDescription"] = image.ImageDescription
	data["ImageDisplayName"] = image.ImageDisplayName
	data["Published"] = image.Published
	data["ImageConfig"] = image.ImageConfig
	data["ImageReadRootPassword"] = image.ImageReadRootPassword

	updated := repositories.UpdateImageById(data)
	if updated == nil {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Could not update image!", nil)
		return
	}

	imageUpdated := repositories.GetImageById(data["Id"])
	if imageUpdated == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading images!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", imageUpdated)
	return
}

// GetImagesConfig godoc
// @Summary     Fetches config files
// @Description Fetches config files from hard drive
// @Tags        image
// @Accept      json
// @Produce     json
// @Success     200 {object}    []string
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /image/config   [get]
func GetImagesConfig(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	dir, err := ioutil.ReadDir(viper.GetString("IKT_STACK_TEMPLATES_USERDATA_DIR"))
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading dir!", nil)
		return
	}

	var files []string
	for _, v := range dir {
		files = append(files, v.Name())
	}

	httputils.ResponseJson(c, http.StatusOK, "", files)
	return
}

// GetPublishedImages godoc
// @Summary     Fetches published images
// @Description Fetches images marked as published from DB
// @Tags        image
// @Accept      json
// @Produce     json
// @Success     200 {object}    []PublishedImagesStruct
// @Failure     500 {object}    nil
// @Router      /image/published    [get]
func GetPublishedImages(c *gin.Context) {
	imagesList := repositories.GetPublishedImages()
	if imagesList == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading images!", nil)
		return
	}

	var publicImages []PublishedImagesStruct
	for _, v := range imagesList {
		publicImages = append(publicImages, PublishedImagesStruct{
			ImageId:          v.ImageId,
			ImageDisplayName: v.ImageDisplayName,
		})
	}

	httputils.ResponseJson(c, http.StatusOK, "", publicImages)
	return
}
