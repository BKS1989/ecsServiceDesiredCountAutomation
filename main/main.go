package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"fmt"
	"github.com/BKS1989/ecsServiceDesiredCountAutomation/ecsautomation"
	"flag"
	"strings"
	"reflect"
	"os"
	"log"
	"sync"
)

var clusterNameList []string
var excludeClusterNameList []string
var excludeServiceNameList []string
var serviceList map[string][]string
var wg sync.WaitGroup
func main() {
	clusterExcludeList := flag.String("clusterExcludeList", "", "cluster exclude listing")
	serviceExcludeList := flag.String("serviceExcludeList","","use coma to exclude multiple service name")
	desireTaskDefinition := flag.Int64("desireTaskDefinition",0,"desire task definition for ecs service")
	region := flag.String("region","","aws region")
	flag.Parse()

	if *region == "" {
		log.Fatal("please provide region name")
	}
	fmt.Print(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Fatal("please provide aws Credi")
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(*region),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			fmt.Println(awserr.Code())
		}
	}
	ecsClient := ecsautomation.GetInstance(sess)
	if reflect.ValueOf(clusterExcludeList).String() != "" {
		excludeClusterNameList = strings.Split(*clusterExcludeList, ",")
	}
	if len(excludeClusterNameList) > 0 {
		clusterNameList, err = ecsClient.GetClusterNameListByParams(excludeClusterNameList)
	} else {
		clusterNameList, err = ecsClient.GetClusterNameList()
	}
	fmt.Println(clusterNameList)
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			fmt.Println("ecs get cluster err", awserr.Code())
		}
	}

	if reflect.ValueOf(serviceExcludeList).String() != "" {
		excludeServiceNameList = strings.Split(*serviceExcludeList,",")
		fmt.Println(excludeServiceNameList)
	}
	for _, clusterName := range clusterNameList {
		if len(excludeServiceNameList) > 0 {
			serviceList, err = ecsClient.GetServicesByClusterNameParams(clusterName,excludeServiceNameList)
		} else {
			serviceList, err = ecsClient.GetServicesByClusterName(clusterName)
		}

		if err != nil {
			if awserr, ok := err.(awserr.Error); ok {
				fmt.Println("ecs get Service err", awserr.Message())
			}
		}
		if services, ok := serviceList[clusterName]; ok {
			for _, serviceName := range services {
				wg.Add(1)
				go ecsClient.StopTaskByService(clusterName,serviceName,*desireTaskDefinition,&wg)
			}
		}
	}
	wg.Wait()
}
