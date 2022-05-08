package v1

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/remoteconsoles"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/startstop"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/spf13/viper"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/database/repositories"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/gopher"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/httputils"
	"gitlab.internal.uia.no/dat304-g-22v/ictsss/bachpals-ictsss-backend/utils"
)

type RequestBodyUserId struct {
	UserId string `json:"user_id"`
}

type RequestBodyVmOrder struct {
	ServerName  string   `json:"server_name"`
	ServerImage string   `json:"server_image"`
	Users       []string `json:"users"`
	GroupName   string   `json:"group_name"`
}

type RequestBodyVmOrderAll struct {
	ServerName     string   `json:"server_name"`
	ServerImage    string   `json:"server_image"`
	Users          []string `json:"users"`
	GroupName      string   `json:"group_name"`
	Everyone       string   `json:"everyone"`
	IncludeTa      string   `json:"include_ta"`
	IncludeTeacher string   `json:"include_teacher"`
	CourseCode     string   `json:"course_code"`
}

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// GetAllVMs godoc
// @Summary     Retrieves list of all VMs from DB
// @Description Gets all VMs
// @Tags        vms
// @Accept      json
// @Produce     json
// @Success     200 {object}    []database.VirtualMachine
// @Router      /vms/all   [get]
func GetAllVms(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	var virtualMachines []database.VirtualMachine
	virtualMachines = repositories.GetVMS()

	if virtualMachines == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Could not load virtual machines!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", virtualMachines)
	return
}

// GetVMs godoc
// @Summary     Retrieves list of users VMs from DB
// @Description Gets all VMs assosciated with a user ID
// @Tags        vms
// @Accept      json
// @Produce     json
// @Success     200 {object}    []database.VirtualMachine
// @Router      /vms/   [get]
func GetVMs(c *gin.Context) {
	var virtualMachines interface{}

	virtualMachines = repositories.GetVMByUserId(c.MustGet("user_id"))

	if virtualMachines == nil {
		httputils.AbortWithStatusJSON(c, http.StatusNotAcceptable, "Could not load virtual machines!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", virtualMachines)
	return
}

// StatusVM godoc
// @Summary     Updates status of VM
// @Description Retrieves the status of the VM then updates it in the DB
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       id  path    string  true    "Server ID"
// @Success     200 {object}    database.VirtualMachine
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/:id/status [get]
func StatusVM(c *gin.Context) {
	id := c.Param("id")

	if len(id) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Missing id!", nil)
		return
	}

	if !IsAdmin(c) {
		if !repositories.CheckIfOwnsVm(id, c.MustGet("user_id").(string)) {
			httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "You don't own this virtual machine!", nil)
			return
		}
	}

	client := gopher.GetClient()

	res := servers.Get(client, id)
	r, err := res.Extract()

	if err != nil {
		var l map[string]interface{}
		m, _ := json.Marshal(res.Err)
		_ = json.Unmarshal(m, &l)
		v := reflect.ValueOf(l["Actual"])

		// Remove vm if it doesn't exists
		if v.Float() == http.StatusNotFound {
			_, err = repositories.DeleteVMById(id)
			if err != nil {
				httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error deleting virtual machine from database!", nil)
				return
			}

			var virtualMachines interface{}

			if IsAdmin(c) {
				virtualMachines = repositories.GetVMS()
			} else {
				virtualMachines = repositories.GetVMByUserId(c.MustGet("user_id"))
			}

			httputils.AbortWithStatusJSON(c, http.StatusOK, "", virtualMachines)
			return
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading virtual machine status!", nil)
		return
	}

	isUpdated := repositories.UpdateVMStatusById(id, r.Status)
	if !isUpdated {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error updating virtual machine status!", nil)
		return
	}

	vm, err := repositories.GetVMById(id)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading database!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "!", vm)
	return
}

// StartVM godoc
// @Summary     Starts a VM
// @Description If the VM is SHUTOFF, tries to get it ACTIVE
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       id  path    string  true    "Server ID"
// @Success     200 {object}    database.VirtualMachine
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/:id/start  [post]
func StartVM(c *gin.Context) {
	id := c.Param("id")

	if len(id) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Missing id!", nil)
		return
	}

	if !IsAdmin(c) {
		if !repositories.CheckIfOwnsVm(id, c.MustGet("user_id").(string)) {
			httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "You don't own this virtual machine!", nil)
			return
		}
	}

	client := gopher.GetClient()

	response := startstop.Start(client, id)

	err := response.Err

	if err != nil {
		log.Println("Result: ", response.ExtractErr())
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while starting virtual machine!", nil)
		return
	}

	// Wait until server status changes to "active"
	if err := servers.WaitForStatus(client, id, database.VirtualMachineStatusActive, database.ServerStatusPollingTime); err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while waiting for virtual machine to become active!", nil)
		return
	}

	r, err := servers.Get(client, id).Extract()
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading virtual machine status!", nil)
		return
	}

	isUpdated := repositories.UpdateVMStatusById(id, r.Status)
	if !isUpdated {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error updating virtual machine status!", nil)
		return
	}

	vm, err := repositories.GetVMById(id)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading database!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", vm)
	return
}

// StopVM godoc
// @Summary     Stops a VM
// @Description If the VM is ACTIVE, tries to get it SHUTOFF
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       id  path    string  true    "Server ID"
// @Success     200 {object}    database.VirtualMachine
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/:id/stop   [post]
func StopVM(c *gin.Context) {
	id := c.Param("id")

	if len(id) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Missing id!", nil)
		return
	}

	if !IsAdmin(c) {
		if !repositories.CheckIfOwnsVm(id, c.MustGet("user_id").(string)) {
			httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "You don't own this virtual machine!", nil)
			return
		}
	}

	client := gopher.GetClient()

	response := startstop.Stop(client, id)

	err := response.Err

	if err != nil {
		log.Println("Result: ", response.ExtractErr())
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while stopping virtual machine!", nil)
		return
	}

	// Wait until server status changes to "stopped"
	if err := servers.WaitForStatus(client, id, database.VirtualMachineStatusInactive, database.ServerStatusPollingTime); err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while waiting for virtual machine to shutdown!", nil)
		return
	}

	r, err := servers.Get(client, id).Extract()

	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading virtual machine status!", nil)
		return
	}

	isUpdated := repositories.UpdateVMStatusById(id, r.Status)
	if !isUpdated {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error updating virtual machine status!", nil)
		return
	}

	vm, err := repositories.GetVMById(id)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading database!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", vm)
	return
}

// RebootVM godoc
// @Summary     Reboots a VM
// @Description Turns a VM off and on again.
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       id  path    string  true    "Server ID"
// @Success     200 {object}    database.VirtualMachine
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/:id/reboot [post]
func RebootVM(c *gin.Context) {
	id := c.Param("id")

	if len(id) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Missing id!", nil)
		return
	}

	if !IsAdmin(c) {
		if !repositories.CheckIfOwnsVm(id, c.MustGet("user_id").(string)) {
			httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "You don't own this virtual machine!", nil)
			return
		}
	}

	client := gopher.GetClient()
	response := servers.Reboot(client, id, servers.RebootOpts{Type: "soft"})

	err := response.Err

	if err != nil {
		log.Println("Result: ", response.ExtractErr())
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while rebooting virtual machine!", nil)
		return
	}

	r, err := servers.Get(client, id).Extract()

	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading virtual machine status!", nil)
		return
	}

	isUpdated := repositories.UpdateVMStatusById(id, r.Status)
	if !isUpdated {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error updating virtual machine status!", nil)
		return
	}

	// Wait until server status changes to "active"
	if err := servers.WaitForStatus(client, id, database.VirtualMachineStatusActive, database.ServerStatusPollingTime); err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while waiting for virtual machine to become active!", nil)
		return
	}

	r, err = servers.Get(client, id).Extract()

	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading virtual machine status!", nil)
		return
	}

	isUpdated = repositories.UpdateVMStatusById(id, r.Status)
	if !isUpdated {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error updating virtual machine status!", nil)
		return
	}

	vm, err := repositories.GetVMById(id)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading database!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", vm)
	return
}

// RespawnVM godoc
// @Summary     Delete and recreate a VM
// @Description Deletes a VM and recreates it with the same parameters.
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       requestStruct   body    RequestBodyVmOrder   true   "Request Body"
// @Success     200 {object}    database.VirtualMachine
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/:id/respawn   [post]
func RespawnVM(c *gin.Context) {
	id := c.Param("id")

	if len(id) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Missing id!", nil)
		return
	}

	if !IsAdmin(c) {
		if !repositories.CheckIfOwnsVm(id, c.MustGet("user_id").(string)) {
			httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "You don't own this virtual machine!", nil)
			return
		}
	}

	vm, err := repositories.GetVMById(id)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading database!", nil)
		return
	}

	// Extract Metadata from VM
	//createOpts := vm.Metadata
	// There is more data to get than just Metadata, remember to get all of it.
	var requestStruct RequestBodyVmOrder
	requestStruct.ServerName = vm.ServerName
	requestStruct.ServerImage = vm.ServerImage

	client := gopher.GetClient()

	// Copied from DeleteVM
	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: vm.ServerIp,
	}

	err = floatingips.DisassociateInstance(client, vm.ServerId, disassociateOpts).ExtractErr()
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error unassigning floating ip from virtual machine!", nil)
		return
	}

	// Use RequestBodyVmOrder to get what is needed from existing VM,
	//  then do as OrderVM.
	if len(requestStruct.ServerImage) == 0 {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Missing server image id!", nil)
		return
	}

	imageInfo := repositories.GetImageByImageId(requestStruct.ServerImage)

	if imageInfo == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading images!", nil)
		return
	}

	vmUsers := repositories.GetVmGroupMembers(id)
	var tmp []string
	for _, vmUsers := range vmUsers {
		tmp = append(tmp, vmUsers+"@uia.no")
	}

	users := strings.Join(tmp, ",")

	var userData []byte
	if len(imageInfo.ImageConfig) > 0 {
		userData = utils.GenerateUserData(imageInfo.ImageConfig, users)
	}

	serverName := ""
	if len(requestStruct.ServerName) > 0 {
		serverName = strings.ToUpper(requestStruct.ServerName)
	} else {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Missing server name!", nil)
		return
	}

	blockDevices := []bootfromvolume.BlockDevice{
		bootfromvolume.BlockDevice{
			DeleteOnTermination: true,
			DestinationType:     bootfromvolume.DestinationVolume,
			SourceType:          bootfromvolume.SourceImage,
			UUID:                requestStruct.ServerImage,
			VolumeSize:          viper.GetInt("IKT_STACK_VM_VOLUME_SIZE"),
		},
	}

	serverCreateOpts := servers.CreateOpts{
		Name:      serverName,
		FlavorRef: viper.GetString("IKT_STACK_VM_FLAVOR_ID"),
		UserData:  userData,
		Networks: []servers.Network{
			{
				UUID: viper.GetString("IKT_STACK_VM_NETWORK_ID"),
			},
		},
		Metadata: map[string]string{
			"VM_IMAGE_ID":            requestStruct.ServerImage,
			"VM_FLAVOR_ID":           viper.GetString("IKT_STACK_VM_FLAVOR_ID"),
			"VM_KEY_NAME":            serverName,
			"VM_VOLUME_SIZE":         viper.GetString("IKT_STACK_VM_VOLUME_SIZE"),
			"VM_NETWORK_ID":          viper.GetString("IKT_STACK_VM_NETWORK_ID"),
			"VM_FLOATING_NETWORK_ID": viper.GetString("IKT_STACK_VM_FLOATING_NETWORK_ID"),
		},
	}

	serverCreateOptsExt := keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           viper.GetString("IKT_STACK_VM_KEY_NAME"),
	}

	createOpts := bootfromvolume.CreateOptsExt{
		CreateOptsBuilder: serverCreateOptsExt,
		BlockDevice:       blockDevices,
	}

	server, err := bootfromvolume.Create(client, createOpts).Extract()

	// Unable to create a virtual machine
	if err != nil {
		fmt.Println(err)
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to create a virtual machine!", nil)
		return
	}

	// Wait until server status changes to "active"
	if err := servers.WaitForStatus(client, server.ID, database.VirtualMachineStatusActive, database.ServerStatusPollingTime); err != nil {
		fmt.Println(err)

		// If we can't get status in time, remove hanging virtual machine.
		result := servers.Delete(client, server.ID)

		// Do we want to print a response or fail silently?
		// This should probably be logged into some kind of logging system.
		if result.ExtractErr() != nil {
			fmt.Println("This error occurred while deleting vm,"+
				"after database failed insertion of a new vm.", result.ExtractErr())
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while waiting for virtual machine to become active!", nil)
		return
	}

	createFloatingIpOpts := floatingips.CreateOpts{
		Pool: viper.GetString("IKT_STACK_VM_FLOATING_NETWORK_ID"),
	}

	fip, err := floatingips.Create(client, createFloatingIpOpts).Extract()
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while creating floating ip for virtual machine!", nil)
		return
	}

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: fip.IP,
		FixedIP:    fip.FixedIP,
	}

	err = floatingips.AssociateInstance(client, server.ID, associateOpts).ExtractErr()
	if err != nil {

		err = floatingips.Delete(client, fip.ID).ExtractErr()
		if err != nil {
			fmt.Println("This error occurred deleting floating ip," +
				"after failed attempt to assign it to a vm.")
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while assigning floating ip to virtual machine!", nil)
		return
	}

	err = secgroups.AddServer(client, server.ID, viper.GetString("IKT_STACK_VM_SECURITY_GROUP_ID")).ExtractErr()
	if err != nil {
		fmt.Println(err)

		// If we can't add security group to a virtual machine.
		result := servers.Delete(client, server.ID)

		// Do we want to print a response or fail silently?
		// This should probably be logged into some kind of logging system.
		if result.ExtractErr() != nil {
			fmt.Println("This error occurred while adding a security to a vm.", result.ExtractErr())
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while adding a network security group to a virtual machine!", nil)
		return
	}

	insertResponse := repositories.InsertMultipleVms(server, fip.IP, serverName, users, imageInfo.ImageId)
	if insertResponse == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to save virtual machine!", nil)
		return
	}

	result := servers.Delete(client, id)

	if result.ExtractErr() != nil {

		forceDeleteResult := servers.ForceDelete(client, id)
		forceDeleteResultError := forceDeleteResult.ExtractErr()

		if forceDeleteResultError != nil {
			httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to delete virtual machine!", nil)
			return
		}
	}

	_, err = repositories.DeleteVMById(id)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error deleting virtual machine from database!", nil)
		return
	}

	usersVms := repositories.GetVMByUserId(c.MustGet("user_id"))
	if usersVms == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to read virtual machines!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "Virtual machine respawned successfully!", usersVms)
	return
}

func OrderVM(c *gin.Context) {
	// Read request body
	var requestStruct RequestBodyVmOrder
	err := c.BindJSON(&requestStruct)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Something is wrong with request body!", nil)
		return
	}

	client := gopher.GetClient()

	if len(requestStruct.ServerImage) == 0 {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Missing server image id!", nil)
		return
	}

	var groupMembers []string
	users := c.MustGet("user_id").(string)
	if len(requestStruct.Users) > 0 {
		for _, v := range requestStruct.Users {
			groupMembers = append(groupMembers, v)
		}

		users += "," + strings.Join(groupMembers, ",")
	}

	imageInfo := repositories.GetImageByImageId(requestStruct.ServerImage)

	if imageInfo == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading images!", nil)
		return
	}

	var userData []byte
	if len(imageInfo.ImageConfig) > 0 {
		userData = utils.GenerateUserData(imageInfo.ImageConfig, users)
	}

	serverName := ""
	if len(requestStruct.ServerName) > 0 {
		serverName = strings.ToUpper(requestStruct.ServerName)
	} else {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Missing server name!", nil)
		return
	}

	if len(requestStruct.GroupName) > 0 {
		serverName = serverName + "-" + strings.ToUpper(requestStruct.GroupName)
	} else {
		cleanUserName := ""

		if strings.Contains(c.MustGet("user_id").(string), "@uia.no") {
			cleanUserName = strings.Replace(c.MustGet("user_id").(string), "@uia.no", "", -1)
		} else if strings.Contains(c.MustGet("user_id").(string), "@student.uia.no") {
			cleanUserName = strings.Replace(c.MustGet("user_id").(string), "@student.uia.no", "", -1)
		}

		serverName = serverName + "-" + cleanUserName
	}

	blockDevices := []bootfromvolume.BlockDevice{
		bootfromvolume.BlockDevice{
			DeleteOnTermination: true,
			DestinationType:     bootfromvolume.DestinationVolume,
			SourceType:          bootfromvolume.SourceImage,
			UUID:                requestStruct.ServerImage,
			VolumeSize:          viper.GetInt("IKT_STACK_VM_VOLUME_SIZE"),
		},
	}

	serverCreateOpts := servers.CreateOpts{
		Name:      serverName,
		FlavorRef: viper.GetString("IKT_STACK_VM_FLAVOR_ID"),
		UserData:  userData,
		Networks: []servers.Network{
			{
				UUID: viper.GetString("IKT_STACK_VM_NETWORK_ID"),
			},
		},
		Metadata: map[string]string{
			"VM_IMAGE_ID":            requestStruct.ServerImage,
			"VM_FLAVOR_ID":           viper.GetString("IKT_STACK_VM_FLAVOR_ID"),
			"VM_KEY_NAME":            serverName,
			"VM_VOLUME_SIZE":         viper.GetString("IKT_STACK_VM_VOLUME_SIZE"),
			"VM_NETWORK_ID":          viper.GetString("IKT_STACK_VM_NETWORK_ID"),
			"VM_FLOATING_NETWORK_ID": viper.GetString("IKT_STACK_VM_FLOATING_NETWORK_ID"),
		},
	}

	serverCreateOptsExt := keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           viper.GetString("IKT_STACK_VM_KEY_NAME"),
	}

	createOpts := bootfromvolume.CreateOptsExt{
		CreateOptsBuilder: serverCreateOptsExt,
		BlockDevice:       blockDevices,
	}

	server, err := bootfromvolume.Create(client, createOpts).Extract()

	// Unable to create a virtual machine
	if err != nil {
		fmt.Println(err)
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to create a virtual machine!", nil)
		return
	}

	// Wait until server status changes to "active"
	if err := servers.WaitForStatus(client, server.ID, database.VirtualMachineStatusActive, database.ServerStatusPollingTime); err != nil {
		fmt.Println(err)

		// If we can't get status in time, remove hanging virtual machine.
		result := servers.Delete(client, server.ID)

		// Do we want to print a response or fail silently?
		// This should probably be logged into some kind of logging system.
		if result.ExtractErr() != nil {
			fmt.Println("This error occurred while deleting vm,"+
				"after database failed insertion of a new vm.", result.ExtractErr())
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while waiting for virtual machine to become active!", nil)
		return
	}

	createFloatingIpOpts := floatingips.CreateOpts{
		Pool: viper.GetString("IKT_STACK_VM_FLOATING_NETWORK_ID"),
	}

	fip, err := floatingips.Create(client, createFloatingIpOpts).Extract()
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while creating floating ip for virtual machine!", nil)
		return
	}

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: fip.IP,
		FixedIP:    fip.FixedIP,
	}

	err = floatingips.AssociateInstance(client, server.ID, associateOpts).ExtractErr()
	if err != nil {

		err = floatingips.Delete(client, fip.ID).ExtractErr()
		if err != nil {
			fmt.Println("This error occurred deleting floating ip," +
				"after failed attempt to assign it to a vm.")
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while assigning floating ip to virtual machine!", nil)
		return
	}

	err = secgroups.AddServer(client, server.ID, viper.GetString("IKT_STACK_VM_SECURITY_GROUP_ID")).ExtractErr()
	if err != nil {
		fmt.Println(err)

		// If we can't add security group to a virtual machine.
		result := servers.Delete(client, server.ID)

		// Do we want to print a response or fail silently?
		// This should probably be logged into some kind of logging system.
		if result.ExtractErr() != nil {
			fmt.Println("This error occurred while adding a security to a vm.", result.ExtractErr())
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while adding a network security group to a virtual machine!", nil)
		return
	}

	insertResponse := repositories.InsertMultipleVms(server, fip.IP, serverName, users, imageInfo.ImageId)
	if insertResponse == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to save virtual machine!", nil)
		return
	}

	usersVms := repositories.GetVMByUserId(c.MustGet("user_id"))
	if usersVms == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to read virtual machines!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "Virtual machine created successfully!", usersVms)
	return
}

// DeleteVM godoc
// @Summary     Deletes a VM
// @Description Deletes the VM both in OpenStack and in DB
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       id  path    string  true    "Server ID"
// @Success     200 {object}    database.VirtualMachine
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/:id    [delete]
func DeleteVM(c *gin.Context) {
	id := c.Param("id")

	if len(id) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Missing id!", nil)
		return
	}

	if !IsAdmin(c) {
		if !repositories.CheckIfOwnsVm(id, c.MustGet("user_id").(string)) {
			httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "You don't own this virtual machine!", nil)
			return
		}
	}

	vm, err := repositories.GetVMById(id)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading database!", nil)
		return
	}

	client := gopher.GetClient()

	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: vm.ServerIp,
	}

	err = floatingips.DisassociateInstance(client, vm.ServerId, disassociateOpts).ExtractErr()
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error unassigning floating ip from virtual machine!", nil)
		return
	}

	result := servers.Delete(client, id)

	if result.ExtractErr() != nil {

		forceDeleteResult := servers.ForceDelete(client, id)
		forceDeleteResultError := forceDeleteResult.ExtractErr()

		if forceDeleteResultError != nil {
			httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to delete virtual machine!", nil)
			return
		}
	}

	_, err = repositories.DeleteVMById(id)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error deleting virtual machine from database!", nil)
		return
	}

	var virtualMachines interface{}

	if IsAdmin(c) {
		virtualMachines = repositories.GetVMS()
	} else {
		virtualMachines = repositories.GetVMByUserId(c.MustGet("user_id"))
	}

	httputils.ResponseJson(c, http.StatusOK, "Virtual machine deleted successfully!", virtualMachines)
	return
}

// GenerateConsoleUrl godoc
// @Summary     Generates a console
// @Description Generates a console for a VM
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       id  path    string  true    "Server ID"
// @Success     200 {object}    interface{}
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/:id/console    [get]
func GenerateConsoleUrl(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Missing id!", nil)
		return
	}

	if !IsAdmin(c) {
		if !repositories.CheckIfOwnsVm(id, c.MustGet("user_id").(string)) {
			httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "You don't own this virtual machine!", nil)
			return
		}
	}

	client := gopher.GetClient()
	client.Microversion = "2.6"

	createOpts := remoteconsoles.CreateOpts{
		Protocol: remoteconsoles.ConsoleProtocolVNC,
		Type:     remoteconsoles.ConsoleTypeNoVNC,
	}

	remoteConsole, err := remoteconsoles.Create(client, id, createOpts).Extract()
	if err != nil {
		fmt.Println("This error occurred while adding a security to a vm.", err)
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while generating access link!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", remoteConsole)
	return
}

// OrderVMFromCanvasAllStudents godoc
// @Summary     Creates a new VM
// @Description Handles a request to create a new VM, for all students in a canvas course
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       requestStruct   body    RequestBodyVmOrderAll   true    "Request Body"
// @Success     200 {object}    []database.VirtualMachine
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/canvas/all   [post]
func OrderVMFromCanvasAllStudents(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	// Read request body
	var requestStruct RequestBodyVmOrderAll
	err := c.BindJSON(&requestStruct)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Something is wrong with request body!", nil)
		return
	}

	courseId := requestStruct.CourseCode

	response, canvasErr := requestCanvasApi("GET", fmt.Sprintf("/courses/%s/users?per_page=1000&enrollment_type=student", courseId), nil, true)
	if canvasErr != nil {
		fmt.Println(canvasErr)
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading canvas api", nil)
		return
	}

	data := response.([]interface{})

	if requestStruct.IncludeTeacher == "true" {
		response, canvasErr := requestCanvasApi("GET", fmt.Sprintf("/courses/%s/users?per_page=1000&enrollment_type=teacher", courseId), nil, true)
		if canvasErr != nil {
			fmt.Println(canvasErr)
			httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading canvas api", nil)
			return
		}

		d := response.([]interface{})
		data = append(data, d...)
	}

	if requestStruct.IncludeTa == "true" {
		response, canvasErr := requestCanvasApi("GET", fmt.Sprintf("/courses/%s/users?per_page=1000&enrollment_type=ta", courseId), nil, true)
		if canvasErr != nil {
			fmt.Println(canvasErr)
			httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error reading canvas api", nil)
			return
		}

		d := response.([]interface{})
		data = append(data, d...)
	}

	if len(data) == 0 {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error no users in response", nil)
		return
	}

	client := gopher.GetClient()

	if len(requestStruct.ServerImage) == 0 {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Missing server image id!", nil)
		return
	}

	imageInfo := repositories.GetImageByImageId(requestStruct.ServerImage)

	for _, v := range data {
		val := v.(map[string]interface{})

		go func(val map[string]interface{}) {

			// If a user is invited but has not accepted invitation to the course, their login_id is nil by default
			if val["login_id"] == nil {
				return
			}

			if imageInfo == nil {
				return
			}

			cleanUserName := ""
			if strings.Contains(val["login_id"].(string), "@uia.no") {
				cleanUserName = strings.Replace(val["login_id"].(string), "@uia.no", "", -1)
			} else if strings.Contains(val["login_id"].(string), "@student.uia.no") {
				cleanUserName = strings.Replace(val["login_id"].(string), "@student.uia.no", "", -1)
			}

			var groupMembers []string
			groupMembers = append(groupMembers, val["login_id"].(string))
			users := strings.Join(groupMembers, ",")

			var userData []byte
			if len(imageInfo.ImageConfig) > 0 {
				userData = utils.GenerateUserData(imageInfo.ImageConfig, users)
			}

			serverName := ""
			if len(requestStruct.ServerName) > 0 {
				serverName = strings.ToUpper(requestStruct.ServerName)
			} else {
				return
			}

			serverName = serverName + "-" + cleanUserName

			blockDevices := []bootfromvolume.BlockDevice{
				bootfromvolume.BlockDevice{
					DeleteOnTermination: true,
					DestinationType:     bootfromvolume.DestinationVolume,
					SourceType:          bootfromvolume.SourceImage,
					UUID:                requestStruct.ServerImage,
					VolumeSize:          viper.GetInt("IKT_STACK_VM_VOLUME_SIZE"),
				},
			}

			serverCreateOpts := servers.CreateOpts{
				Name:      serverName,
				FlavorRef: viper.GetString("IKT_STACK_VM_FLAVOR_ID"),
				UserData:  userData,
				Networks: []servers.Network{
					{
						UUID: viper.GetString("IKT_STACK_VM_NETWORK_ID"),
					},
				},
				Metadata: map[string]string{
					"VM_IMAGE_ID":            requestStruct.ServerImage,
					"VM_FLAVOR_ID":           viper.GetString("IKT_STACK_VM_FLAVOR_ID"),
					"VM_KEY_NAME":            serverName,
					"VM_VOLUME_SIZE":         viper.GetString("IKT_STACK_VM_VOLUME_SIZE"),
					"VM_NETWORK_ID":          viper.GetString("IKT_STACK_VM_NETWORK_ID"),
					"VM_FLOATING_NETWORK_ID": viper.GetString("IKT_STACK_VM_FLOATING_NETWORK_ID"),
				},
			}

			serverCreateOptsExt := keypairs.CreateOptsExt{
				CreateOptsBuilder: serverCreateOpts,
				KeyName:           viper.GetString("IKT_STACK_VM_KEY_NAME"),
			}

			createOpts := bootfromvolume.CreateOptsExt{
				CreateOptsBuilder: serverCreateOptsExt,
				BlockDevice:       blockDevices,
			}

			server, err := bootfromvolume.Create(client, createOpts).Extract()

			// Unable to create a virtual machine
			if err != nil {
				fmt.Println(err)
				return
			}

			// Wait until server status changes to "active"
			if err := servers.WaitForStatus(client, server.ID, database.VirtualMachineStatusActive, database.ServerStatusPollingTime); err != nil {
				fmt.Println(err)

				// If we can't get status in time, remove hanging virtual machine.
				result := servers.Delete(client, server.ID)

				// Do we want to print a response or fail silently?
				// This should probably be logged into some kind of logging system.
				if result.ExtractErr() != nil {
					fmt.Println("This error occurred while deleting vm,"+
						"after database failed insertion of a new vm.", result.ExtractErr())
				}
				return
			}

			createFloatingIpOpts := floatingips.CreateOpts{
				Pool: viper.GetString("IKT_STACK_VM_FLOATING_NETWORK_ID"),
			}

			fip, err := floatingips.Create(client, createFloatingIpOpts).Extract()
			if err != nil {
				return
			}

			associateOpts := floatingips.AssociateOpts{
				FloatingIP: fip.IP,
				FixedIP:    fip.FixedIP,
			}

			err = floatingips.AssociateInstance(client, server.ID, associateOpts).ExtractErr()
			if err != nil {

				err = floatingips.Delete(client, fip.ID).ExtractErr()
				if err != nil {
					fmt.Println("This error occurred deleting floating ip," +
						"after failed attempt to assign it to a vm.")
				}
				return
			}

			err = secgroups.AddServer(client, server.ID, viper.GetString("IKT_STACK_VM_SECURITY_GROUP_ID")).ExtractErr()
			if err != nil {
				fmt.Println(err)

				// If we can't add security group to a virtual machine.
				result := servers.Delete(client, server.ID)

				// Do we want to print a response or fail silently?
				// This should probably be logged into some kind of logging system.
				if result.ExtractErr() != nil {
					fmt.Println("This error occurred while adding a security to a vm.", result.ExtractErr())
				}
				return
			}

			insertResponse := repositories.InsertMultipleVms(server, fip.IP, serverName, users, imageInfo.ImageId)
			if insertResponse == nil {
				return
			}
			return
		}(val)
	}

	var virtualMachines []database.VirtualMachine
	virtualMachines = repositories.GetVMS()
	if virtualMachines == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to read virtual machines!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "Virtual machine created successfully!", virtualMachines)
	return
}

// OrderVMFromCanvas godoc
// @Summary     Creates a new VM
// @Description Handles a request to create a new VM, for a student in a canvas course.
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       requestStruct   body    RequestBodyVmOrder   true   "Request Body"
// @Success     200 {object}    database.VirtualMachine
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/canvas   [post]
func OrderVMFromCanvas(c *gin.Context) {
	if !IsAdmin(c) {
		httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "Administrator credentials required!", nil)
		return
	}

	// Read request body
	var requestStruct RequestBodyVmOrder
	err := c.BindJSON(&requestStruct)
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Something is wrong with request body!", nil)
		return
	}

	client := gopher.GetClient()

	if len(requestStruct.ServerImage) == 0 {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Missing server image id!", nil)
		return
	}

	var groupMembers []string
	for _, v := range requestStruct.Users {
		groupMembers = append(groupMembers, v)
	}
	users := strings.Join(groupMembers, ",")

	imageInfo := repositories.GetImageByImageId(requestStruct.ServerImage)

	if imageInfo == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while reading images!", nil)
		return
	}

	var userData []byte
	if len(imageInfo.ImageConfig) > 0 {
		userData = utils.GenerateUserData(imageInfo.ImageConfig, users)
	}

	serverName := ""
	if len(requestStruct.ServerName) > 0 {
		serverName = strings.ToUpper(requestStruct.ServerName)
	} else {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Missing server name!", nil)
		return
	}

	if len(requestStruct.GroupName) > 0 {
		serverName = serverName + "-" + strings.ToUpper(requestStruct.GroupName)
	} else {
		cleanUserName := ""

		if strings.Contains(users, "@uia.no") {
			cleanUserName = strings.Replace(users, "@uia.no", "", -1)
		} else if strings.Contains(users, "@student.uia.no") {
			cleanUserName = strings.Replace(users, "@student.uia.no", "", -1)
		}

		serverName = serverName + "-" + cleanUserName
	}

	blockDevices := []bootfromvolume.BlockDevice{
		bootfromvolume.BlockDevice{
			DeleteOnTermination: true,
			DestinationType:     bootfromvolume.DestinationVolume,
			SourceType:          bootfromvolume.SourceImage,
			UUID:                requestStruct.ServerImage,
			VolumeSize:          viper.GetInt("IKT_STACK_VM_VOLUME_SIZE"),
		},
	}

	serverCreateOpts := servers.CreateOpts{
		Name:      serverName,
		FlavorRef: viper.GetString("IKT_STACK_VM_FLAVOR_ID"),
		UserData:  userData,
		Networks: []servers.Network{
			{
				UUID: viper.GetString("IKT_STACK_VM_NETWORK_ID"),
			},
		},
		Metadata: map[string]string{
			"VM_IMAGE_ID":            requestStruct.ServerImage,
			"VM_FLAVOR_ID":           viper.GetString("IKT_STACK_VM_FLAVOR_ID"),
			"VM_KEY_NAME":            serverName,
			"VM_VOLUME_SIZE":         viper.GetString("IKT_STACK_VM_VOLUME_SIZE"),
			"VM_NETWORK_ID":          viper.GetString("IKT_STACK_VM_NETWORK_ID"),
			"VM_FLOATING_NETWORK_ID": viper.GetString("IKT_STACK_VM_FLOATING_NETWORK_ID"),
		},
	}

	serverCreateOptsExt := keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           viper.GetString("IKT_STACK_VM_KEY_NAME"),
	}

	createOpts := bootfromvolume.CreateOptsExt{
		CreateOptsBuilder: serverCreateOptsExt,
		BlockDevice:       blockDevices,
	}

	server, err := bootfromvolume.Create(client, createOpts).Extract()

	// Unable to create a virtual machine
	if err != nil {
		fmt.Println(err)
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to create a virtual machine!", nil)
		return
	}

	// Wait until server status changes to "active"
	if err := servers.WaitForStatus(client, server.ID, database.VirtualMachineStatusActive, database.ServerStatusPollingTime); err != nil {
		fmt.Println(err)

		// If we can't get status in time, remove hanging virtual machine.
		result := servers.Delete(client, server.ID)

		// Do we want to print a response or fail silently?
		// This should probably be logged into some kind of logging system.
		if result.ExtractErr() != nil {
			fmt.Println("This error occurred while deleting vm,"+
				"after database failed insertion of a new vm.", result.ExtractErr())
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while waiting for virtual machine to become active!", nil)
		return
	}

	createFloatingIpOpts := floatingips.CreateOpts{
		Pool: viper.GetString("IKT_STACK_VM_FLOATING_NETWORK_ID"),
	}

	fip, err := floatingips.Create(client, createFloatingIpOpts).Extract()
	if err != nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while creating floating ip for virtual machine!", nil)
		return
	}

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: fip.IP,
		FixedIP:    fip.FixedIP,
	}

	err = floatingips.AssociateInstance(client, server.ID, associateOpts).ExtractErr()
	if err != nil {

		err = floatingips.Delete(client, fip.ID).ExtractErr()
		if err != nil {
			fmt.Println("This error occurred deleting floating ip," +
				"after failed attempt to assign it to a vm.")
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while assigning floating ip to virtual machine!", nil)
		return
	}

	err = secgroups.AddServer(client, server.ID, viper.GetString("IKT_STACK_VM_SECURITY_GROUP_ID")).ExtractErr()
	if err != nil {
		fmt.Println(err)

		// If we can't add security group to a virtual machine.
		result := servers.Delete(client, server.ID)

		// Do we want to print a response or fail silently?
		// This should probably be logged into some kind of logging system.
		if result.ExtractErr() != nil {
			fmt.Println("This error occurred while adding a security to a vm.", result.ExtractErr())
		}

		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while adding a network security group to a virtual machine!", nil)
		return
	}

	insertResponse := repositories.InsertMultipleVms(server, fip.IP, serverName, users, imageInfo.ImageId)
	if insertResponse == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to save virtual machine!", nil)
		return
	}

	var virtualMachines []database.VirtualMachine
	virtualMachines = repositories.GetVMS()
	if virtualMachines == nil {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Unable to read virtual machines!", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "Virtual machine created successfully!", virtualMachines)
	return
}

// GetPassword godoc
// @Summary     Fetches VM password
// @Description Fetches the password for a VM
// @Tags        vms
// @Accept      json
// @Produce     json
// @Param       id  path    string  true    "Server ID"
// @Success     200 {object}    interface{}
// @Failure     400 {object}    nil
// @Failure     406 {object}    nil
// @Failure     500 {object}    nil
// @Router      /vms/:id/password    [get]
func GetPassword(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		httputils.AbortWithStatusJSON(c, http.StatusBadRequest, "Missing id!", nil)
		return
	}

	if !IsAdmin(c) {
		if !repositories.CheckIfOwnsVm(id, c.MustGet("user_id").(string)) {
			httputils.AbortWithStatusJSON(c, http.StatusUnauthorized, "You don't own this virtual machine!", nil)
			return
		}
	}

	client := gopher.GetClient()

	key := utils.ReadPrivateKey()
	if reflect.TypeOf(key).Kind() == reflect.String {
		httputils.ResponseJson(c, http.StatusInternalServerError, key.(string), nil)
		return
	}

	password, err := servers.GetPassword(client, id).ExtractPassword(key.(*rsa.PrivateKey))

	if err != nil {
		fmt.Println("This error occurred while adding security to a vm.", err)
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Error while generating access link!", nil)
		return
	}

	if len(password) == 0 {
		httputils.AbortWithStatusJSON(c, http.StatusInternalServerError, "Password is not set", nil)
		return
	}

	httputils.ResponseJson(c, http.StatusOK, "", password)
	return
}
