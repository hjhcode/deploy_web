package managers

import (
	"encoding/json"
	"time"

	"github.com/hjhcode/deploy-web/common/components"
	"github.com/hjhcode/deploy-web/models"
)

type Data struct {
	Stage []stage
}

type stage struct {
	Status  int
	Machine []machine
}

type machine struct {
	Id     int64
	Status int
}

func GetAllDeploy() ([]map[string]interface{}, int) {
	deployList, err := models.Deploy{}.QueryAllDeploy()
	if err != nil {
		panic(err.Error())
	}

	if deployList == nil {
		return nil, 0
	}

	var deployLists []map[string]interface{}
	for i := 0; i < len(deployList); i++ {
		startTime := time.Unix(deployList[i].DeployStart, 0).Format("2006-01-02 15:04:05")
		endTime := time.Unix(deployList[i].DeployEnd, 0).Format("2006-01-02 15:04:05")
		deploy := make(map[string]interface{})
		deploy["id"] = deployList[i].Id
		deploy["account_id"] = getCreator(deployList[i].AccountId)
		deploy["service_id"] = getServiceName(deployList[i].ServiceId)
		deploy["deploy_start"] = startTime
		deploy["deploy_end"] = endTime
		deploy["deploy_statu"] = deployList[i].DeployStatu
		deployLists = append(deployLists, deploy)
	}

	count, err := models.Deploy{}.CountAllDeployByPage()
	if err != nil {
		panic(err.Error())
	}

	return deployLists, int(count)
}

func GetDeployServiceData(deployId int64) map[string]interface{} {
	deploy, err := models.Deploy{}.GetById(deployId)
	if err != nil {
		panic(err.Error())
	}

	if deploy == nil {
		return nil
	}

	service, err := models.Service{}.GetById(deploy.ServiceId)
	if err != nil {
		panic(err.Error())
	}

	if service == nil {
		return nil
	}

	lists := make(map[string]interface{})
	createtime := time.Unix(service.CreateDate, 0).Format("2006-01-02 15:04:05")
	updatetime := time.Unix(service.UpdateDate, 0).Format("2006-01-02 15:04:05")
	lists["service_name"] = service.ServiceName
	lists["service_describe"] = service.ServiceDescribe
	lists["account_name"] = getCreator(service.AccountId)
	lists["host_list"] = deploy.MirrorList
	lists["mirror_list"] = getMirrorName(deploy.MirrorList)
	lists["create_date"] = createtime
	lists["update_date"] = updatetime
	lists["docker_config"] = deploy.DockerConfig
	lists["deploy_statu"] = deploy.DeployStatu
	lists["deploy_log"] = deploy.DeployLog

	return lists
}

func StartDeployService(accountId int64, deployId int64, groupId int) (bool, string) {

	deploy, err := models.Deploy{}.GetById(deployId)
	if err != nil {
		panic(err.Error())
	}

	result := isServiceMember(deploy.ServiceId, accountId)
	if !result {
		return false, "You have no authority"
	}

	var data Data
	if err := json.Unmarshal([]byte(deploy.MirrorList), &data); err != nil {
		panic(err.Error())
	}

	data.Stage[groupId-1].Status = 1 //该分组处于部署状态

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	deployRecord := &models.Deploy{
		Id:          deployId,
		DeployStatu: 1, //部署中
		DeployStart: time.Now().Unix(),
		HostList:    string(jsonBytes),
	}

	errs := models.Deploy{}.Update(deployRecord)
	if errs != nil {
		panic(err.Error())
	}

	mess := &components.SendMess{OrderType: 1, DataId: deployId}
	components.Send("deploy", mess)

	return true, ""
}

func BackDeployService(accountId int64, deployId int64) (bool, string) {

	deploy, err := models.Deploy{}.GetById(deployId)
	if err != nil {
		panic(err.Error())
	}

	result := isServiceMember(deploy.ServiceId, accountId)
	if !result {
		return false, "You have no authority"
	}

	deployRecord := &models.Deploy{
		Id:          deployId,
		DeployStatu: 2, //回滚
	}

	errs := models.Deploy{}.Update(deployRecord)
	if errs != nil {
		panic(err.Error())
	}

	mess := &components.SendMess{OrderType: 2, DataId: deployId}
	components.Send("deploy", mess)

	return true, ""
}
