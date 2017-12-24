# ecsServiceDesiredCountAutomation

set desire count for ecs service on aws                                                     
go run main.go -region AWS_REGION -desireTaskDefinition DESIRECOUNT

exclude ecs cluster to set desire count for service                                              
go run main.go -clusterExcludeList="EcsCluster1,EcsCluster2,EcsCluster3" -region AWS_REGION -desireTaskDefinition DESIRECOUNT

exclude ecs service to set desire count for service

go run main.go -serviceExcludeList="<EcsService1>,<EcsService2>,<EcsService3>" -region AWS_REGION -desireTaskDefinition DESIRECOUNT

