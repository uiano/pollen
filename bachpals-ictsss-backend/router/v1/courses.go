package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/httputils"
)

type CanvasCourse struct {
	ID                                int                               `json:"id"`
	SisCourseID                       interface{}                       `json:"sis_course_id"`
	UUID                              string                            `json:"uuid"`
	IntegrationID                     interface{}                       `json:"integration_id"`
	SisImportID                       int                               `json:"sis_import_id"`
	Name                              string                            `json:"name"`
	CourseCode                        string                            `json:"course_code"`
	WorkflowState                     string                            `json:"workflow_state"`
	AccountID                         int                               `json:"account_id"`
	RootAccountID                     int                               `json:"root_account_id"`
	EnrollmentTermID                  int                               `json:"enrollment_term_id"`
	GradingPeriods                    interface{}                       `json:"grading_periods"`
	GradingStandardID                 int                               `json:"grading_standard_id"`
	GradePassbackSetting              string                            `json:"grade_passback_setting"`
	CreatedAt                         string                            `json:"created_at"`
	StartAt                           string                            `json:"start_at"`
	EndAt                             string                            `json:"end_at"`
	Locale                            string                            `json:"locale"`
	Enrollments                       interface{}                       `json:"enrollments"`
	TotalStudents                     int                               `json:"total_students"`
	Calendar                          interface{}                       `json:"calendar"`
	DefaultView                       string                            `json:"default_view"`
	SyllabusBody                      string                            `json:"syllabus_body"`
	NeedsGradingCount                 int                               `json:"needs_grading_count"`
	Term                              interface{}                       `json:"term"`
	CourseProgress                    interface{}                       `json:"course_progress"`
	ApplyAssignmentGroupWeights       bool                              `json:"apply_assignment_group_weights"`
	Permissions                       Permissions                       `json:"permissions"`
	IsPublic                          bool                              `json:"is_public"`
	IsPublicToAuthUsers               bool                              `json:"is_public_to_auth_users"`
	PublicSyllabus                    bool                              `json:"public_syllabus"`
	PublicSyllabusToAuth              bool                              `json:"public_syllabus_to_auth"`
	PublicDescription                 string                            `json:"public_description"`
	StorageQuotaMb                    int                               `json:"storage_quota_mb"`
	StorageQuotaUsedMb                int                               `json:"storage_quota_used_mb"`
	HideFinalGrades                   bool                              `json:"hide_final_grades"`
	License                           string                            `json:"license"`
	AllowStudentAssignmentEdits       bool                              `json:"allow_student_assignment_edits"`
	AllowWikiComments                 bool                              `json:"allow_wiki_comments"`
	AllowStudentForumAttachments      bool                              `json:"allow_student_forum_attachments"`
	OpenEnrollment                    bool                              `json:"open_enrollment"`
	SelfEnrollment                    bool                              `json:"self_enrollment"`
	RestrictEnrollmentsToCourseDates  bool                              `json:"restrict_enrollments_to_course_dates"`
	CourseFormat                      string                            `json:"course_format"`
	AccessRestrictedByDate            bool                              `json:"access_restricted_by_date"`
	TimeZone                          string                            `json:"time_zone"`
	Blueprint                         bool                              `json:"blueprint"`
	BlueprintRestrictions             BlueprintRestrictions             `json:"blueprint_restrictions"`
	BlueprintRestrictionsByObjectType BlueprintRestrictionsByObjectType `json:"blueprint_restrictions_by_object_type"`
	Template                          bool                              `json:"template"`
}

type Permissions struct {
	CreateDiscussionTopic bool `json:"create_discussion_topic"`
	CreateAnnouncement    bool `json:"create_announcement"`
}

type BlueprintRestrictions struct {
	Content           bool `json:"content"`
	Points            bool `json:"points"`
	DueDates          bool `json:"due_dates"`
	AvailabilityDates bool `json:"availability_dates"`
}

type Assignment struct {
	Content bool `json:"content"`
	Points  bool `json:"points"`
}

type WikiPage struct {
	Content bool `json:"content"`
}

type BlueprintRestrictionsByObjectType struct {
	Assignment Assignment `json:"assignment"`
	WikiPage   WikiPage   `json:"wiki_page"`
}

func requestCanvasApi(method string, path string, body io.Reader, asArray bool) (interface{}, error) {
	req, err := http.NewRequest(method, viper.GetString("IKT_STACK_CANVAS_API_URL")+path, body)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("IKT_STACK_CANVAS_API_KEY")))

	client := &http.Client{}
	resp, responseErr := client.Do(req)
	if responseErr != nil {
		return nil, responseErr
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if asArray {
		if reflect.TypeOf(b).Kind() == reflect.Map {
			var d map[string]interface{}
			err = json.Unmarshal(b, &d)
			if err != nil {
				return nil, err
			}
		}

		var d []interface{}
		err = json.Unmarshal(b, &d)
		if err != nil {
			return nil, err
		}

		return d, nil
	}

	var d interface{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// GetCourses godoc
// @Summary		Fetches courses from Canvas
// @Description	Fetches a course or an array of courses from the Canvas API
// @Tags        courses
// @Accept      json
// @Produce     json
// @Success     200 {object}	interface{}
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /courses/	[get]
func GetCourses(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	// https://community.canvaslms.com/t5/Canvas-Question-Forum/Getting-a-list-of-ALL-courses/m-p/185855/highlight/true#M89957
	response, err := requestCanvasApi("GET", "/courses?per_page=100", nil, true)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading canvas api", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", response)
	return
}

// GetCourseStudents godoc
// @Summary		Fetches students from course
// @Description	Fetches the students associated with a course in Canvas
// @Tags        courses
// @Accept      json
// @Produce     json
// @Param       courseId  path    string  true    "Course ID"
// @Success     200 {object}	interface{}
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /courses/:id/users	[get]
func GetCourseStudents(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	courseId := c.Param("id")

	response, err := requestCanvasApi("GET", fmt.Sprintf("/courses/%s/users?per_page=1000", courseId), nil, true)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading canvas api", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", response)
	return
}

// GetCourseGroups godoc
// @Summary		Fetches groups from course
// @Description	Fetches the groups associated with a course in Canvas
// @Tags        courses
// @Accept      json
// @Produce     json
// @Param       courseId  path    string  true    "Course ID"
// @Success     200 {object}	interface{}
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /courses/:id/groups	[get]
func GetCourseGroups(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	courseId := c.Param("id")

	response, err := requestCanvasApi("GET", fmt.Sprintf("/courses/%s/groups?per_page=1000", courseId), nil, true)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading canvas api", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", response)
	return
}

// GetGroupUsers godoc
// @Summary		Fetches students in group
// @Description	Fetches the students associated with a group in a Canvas course
// @Tags        courses
// @Accept      json
// @Produce     json
// @Param       groupId  path    string  true    "Group ID"
// @Success     200 {object}	interface{}
// @Failure     400 {object}    nil
// @Failure     401 {object}    nil
// @Failure     500 {object}    nil
// @Router      /courses/groups/:id/users	[get]
func GetGroupUsers(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	groupId := c.Param("id")

	response, err := requestCanvasApi("GET", fmt.Sprintf("/groups/%s/users?per_page=1000", groupId), nil, false)

	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading canvas api", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", response)
	return
}
