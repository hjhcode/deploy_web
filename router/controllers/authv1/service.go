package authv1

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hjhcode/deploy-web/managers"
	"github.com/hjhcode/deploy-web/router/controllers/base"
)

func RegisterService(router *gin.RouterGroup) {
	router.POST("service/add", httpHandlerServiceAdd)
	router.POST("service/del", httpHandlerServiceDel)
	router.POST("service/update", httpHandlerServiceUpdate)
	router.POST("service/deploy", httpHandlerServiceDeploy)
	router.GET("service/show", httpHandlerServiceShow)
	router.GET("service/search", httpHandlerServiceSearch)
	router.GET("service/detail", httpHandlerServiceDetail)
}

type ServiceParam struct {
	ServiceId       int64  `form:"service_id" json:"service_id"`
	ServiceName     string `form:"service_name" json:"service_name" binding:"required"`
	ServiceDescribe string `form:"service_describe" json:"service_describe" binding:"required"`
	HostList        string `form:"host_list" json:"host_list" binding:"required"`
	MirrorList      string `form:"mirror_list" json:"mirror_list" binding:"required"`
	DockerConfig    string `form:"docker_config" json:"docker_config" binding:"required"`
	ServiceMember   string `form:"service_member" json:"service_member" binding:"required"`
}

type ServiceIdParam struct {
	ServiceId int64 `json:"service_id" binding:"required"`
}

func httpHandlerServiceAdd(c *gin.Context) {
	var service ServiceParam
	err := c.BindJSON(&service)
	if err != nil {
		panic(err.Error())
	}
	accountId := base.UserId(c)
	if result, mess := managers.AddNewService(accountId, service.ServiceName, service.ServiceDescribe,
		service.HostList, service.MirrorList, service.DockerConfig, service.ServiceMember); !result {
		c.JSON(http.StatusOK, base.Fail(mess))
		return
	}

	c.JSON(http.StatusOK, base.Success())
}

func httpHandlerServiceDel(c *gin.Context) {
	var service ServiceIdParam
	err := c.BindJSON(&service)
	if err != nil {
		panic(err.Error())
	}
	accountId := base.UserId(c)
	if result, mess := managers.DelService(service.ServiceId, accountId); !result {
		c.JSON(http.StatusOK, base.Fail(mess))
		return
	}

	c.JSON(http.StatusOK, base.Success())
}

func httpHandlerServiceUpdate(c *gin.Context) {
	var service ServiceParam
	err := c.BindJSON(&service)
	if err != nil {
		panic(err.Error())
	}
	accountId := base.UserId(c)
	if result, mess := managers.UpdateService(accountId, service.ServiceId, service.ServiceName, service.ServiceDescribe,
		service.HostList, service.MirrorList, service.DockerConfig, service.ServiceMember); !result {
		c.JSON(http.StatusOK, base.Fail(mess))
		return
	}

	c.JSON(http.StatusOK, base.Success())
}

func httpHandlerServiceDeploy(c *gin.Context) {
	var service ServiceIdParam
	err := c.BindJSON(&service)
	if err != nil {
		panic(err.Error())
	}
	accountId := base.UserId(c)
	//var accountId int64 = 1
	if result, mess := managers.DeployService(service.ServiceId, accountId); !result {
		c.JSON(http.StatusOK, base.Fail(mess))
		return
	}

	c.JSON(http.StatusOK, base.Success())
}

func httpHandlerServiceShow(c *gin.Context) {
	nums := c.Query("")
	//nums := c.Query("size")
	size, _ := strconv.Atoi(nums)
	requestPage := c.Query("")
	//requestPage := c.Query("requestpage")
	start, _ := strconv.Atoi(requestPage)
	serviceList, num := managers.GetAllService(size, start)
	response := map[string]interface{}{
		"request_page": requestPage,
		"total_page":   num,
		"data":         serviceList,
	}

	c.JSON(http.StatusOK, base.Success(response))
}

func httpHandlerServiceSearch(c *gin.Context) {
	nums := c.Query("")
	size, _ := strconv.Atoi(nums)
	requestPage := c.Query("")
	start, _ := strconv.Atoi(requestPage)
	serviceName := c.Query("")
	serviceList, num := managers.GetServiceByParam(serviceName, size, start)
	response := map[string]interface{}{
		"request_page": requestPage,
		"total_page":   num,
		"data":         serviceList,
	}

	c.JSON(http.StatusOK, base.Success(response))
}

func httpHandlerServiceDetail(c *gin.Context) {
	id := c.Query("")
	//id := c.Query("id")
	serviceId, _ := strconv.ParseInt(id, 10, 64)
	service := managers.GetOneService(serviceId)
	c.JSON(http.StatusOK, base.Success(service))
}
