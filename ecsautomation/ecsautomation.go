package ecsautomation

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"sync"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"reflect"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"fmt"
)

type EcsAutomation struct {
	*ecs.ECS
}

var Once sync.Once
var ecsautomation *EcsAutomation

func GetInstance(awsSession *session.Session)  *EcsAutomation {
	Once.Do(func() {
		ecsautomation = CreateEcsInstance(awsSession)
	})
	return ecsautomation
}

func CreateEcsInstance(awsSession *session.Session) *EcsAutomation {
	return &EcsAutomation{
		ECS: ecs.New(awsSession),
	}
}

func (capecs *EcsAutomation) GetClusterNameList() ([]string,error){
	response, err := capecs.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		return nil, err
	}
	clusterName := []string{}
	for _,clusterArn := range response.ClusterArns {
		clusterName = append(clusterName,strings.Split(*clusterArn,"/")[1])
	}
	return clusterName,nil
}

func (capecs *EcsAutomation) GetServicesByClusterNameParams(clusterName string, excludeServices []string) (map[string][]string, error) {
	response, err := capecs.ListServices(&ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
	})
	if err != nil {
		return nil, err
	}
	serviceName := make(map[string][]string)
	serviceNameList := []string{}
	excludeServiceNameList := make(map[string]string)
	for _, excludeServiceName := range excludeServices {
		excludeServiceNameList[excludeServiceName] = excludeServiceName
	}
	for _, serviceArn := range response.ServiceArns {
		if reflect.ValueOf(serviceArn).String() != "" {
			if _, ok := excludeServiceNameList[strings.TrimSpace(strings.Split(*serviceArn, "/")[1])]; !ok {
				serviceNameList = append(serviceNameList, strings.TrimSpace(strings.Split(*serviceArn, "/")[1]))
			}
		}
	}
	if len(serviceNameList) != 0 {
		serviceName[clusterName] = serviceNameList
	}
	return serviceName, err
}

func (capecs *EcsAutomation) GetServicesByClusterName(clusterName string) (map[string][]string, error) {
	response, err := capecs.ListServices(&ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
	})
	if err != nil {
		return nil, err
	}
	serviceName := make(map[string][]string)
	serviceNameList := []string{}
	for _,serviceArn := range response.ServiceArns {
		if reflect.ValueOf(serviceArn).String() != "" {
			serviceNameList = append(serviceNameList,strings.TrimSpace(strings.Split(*serviceArn,"/")[1]))
		}
	}
	if len(serviceNameList) != 0 {
		serviceName[clusterName] = serviceNameList
	}
	return serviceName, err
}
func (capecs *EcsAutomation) GetClusterNameListByParams(excludeClusterNameList []string) ([]string, error) {
	response, err := capecs.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		return nil, err
	}
	excludeCluster := make(map[string]string)
	clusterName := []string{}
	for _, excludeClusterName := range excludeClusterNameList {
		excludeCluster[excludeClusterName] = excludeClusterName
	}
	for _,clusterArn := range response.ClusterArns {
		if _,ok := excludeCluster[strings.Split(*clusterArn,"/")[1]]; ! ok {
			if strings.Split(*clusterArn,"/")[1] != "" {
				clusterName = append(clusterName,strings.Split(*clusterArn,"/")[1])
			}
		}
	}
	return clusterName,nil
}

func (capecs *EcsAutomation) StopTaskByService(clusterName string, serviceName string, desireCount int64, wg *sync.WaitGroup){
	defer  wg.Done()
	response,err :=  capecs.UpdateService(&ecs.UpdateServiceInput{
		Service: aws.String(serviceName),
		Cluster: aws.String(clusterName),
		DesiredCount: aws.Int64(desireCount),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Println("Error Code ", awsErr.Code(),"Error Message", awsErr.Message())
		}
	} else {
		fmt.Println("Successfully Set DesireCount",*response.Service.DesiredCount,"for Service ",*response.Service.ServiceName)
	}
}