package cmd

import (
	"fmt"
	"os"
	_ "reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/jedib0t/go-pretty/table"
)

type Service struct {
	Name, Cluster     string
	Tasks             []*ecs.Task
	TasksTargetHealth []*elbv2.TargetHealthDescription
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func listECSServices(sess *session.Session, cluster string) error {
	svc := ecs.New(sess)
	params := &ecs.ListServicesInput{
		Cluster:    aws.String(cluster),
		MaxResults: aws.Int64(100),
	}

	pageNum := 0
	service := make([]*string, 0)
	err := svc.ListServicesPages(params,
		func(page *ecs.ListServicesOutput, lastPage bool) bool {
			pageNum++
			service = append(service, page.ServiceArns...)
			return pageNum <= 512
		})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	svc2 := make([]string, 0)
	for _, s := range service {
		svc2 = append(svc2, strings.Split(*s, "/")[1])
	}
	for _, i := range svc2 {
		fmt.Println(i)
	}
	return nil
}

func listTasks(sess *session.Session, cluster, svn string) []*string {
	svc := ecs.New(sess)
	params := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		MaxResults:  aws.Int64(100),
		ServiceName: aws.String(svn),
	}
	pageNum := 0
	tasks := make([]*string, 0)
	err := svc.ListTasksPages(params,
		func(page *ecs.ListTasksOutput, lastPage bool) bool {
			pageNum++
			tasks = append(tasks, page.TaskArns...)
			return pageNum <= 3
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return tasks
}

func describeTask(sess *session.Session, cluster string, tasks []*string) []*ecs.Task {
	svc := ecs.New(sess)
	params := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   tasks,
	}
	resp, err := svc.DescribeTasks(params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return resp.Tasks
}

func describService(sess *session.Session, cluster, svn string) *ecs.Service {

	svc := ecs.New(sess)
	params := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []*string{aws.String(svn)},
	}
	descSvc, err := svc.DescribeServices(params)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return descSvc.Services[0]
}

func listHealth(sess *session.Session, tga string) []*elbv2.TargetHealthDescription {

	elbSvc := elbv2.New(sess)

	param := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(tga),
	}
	resultAlb, err := elbSvc.DescribeTargetHealth(param)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return resultAlb.TargetHealthDescriptions
}

func display(svc Service) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	ec2Table := table.NewWriter()
	ecsTable := table.NewWriter()

	fmt.Println(svc.Cluster)
	fmt.Println(svc.Name)
	if len(svc.TasksTargetHealth) != 0 {
		t.AppendHeader(table.Row{"ECS", "EC2 - ALB"})
		ec2Table.AppendHeader(table.Row{
			"service status",
			"instance id",
			"port",
		})
		ecsTable.AppendHeader(table.Row{
			"Task ID",
			"task definition",
			"port",
			"Status",
			"Image Tag",
		})
		for _, i := range svc.TasksTargetHealth {
			ec2Table.AppendRow([]interface{}{
				*i.TargetHealth.State,
				*i.Target.Id,
				*i.Target.Port,
			})
		}
		for _, i := range svc.Tasks {
			var hostPort string
			if len(i.Containers[0].NetworkBindings) == 0 {
				hostPort = "...?"
			} else {
				hostPort = strconv.FormatInt(*i.Containers[0].NetworkBindings[0].HostPort, 10)
			}
			ecsTable.AppendRow([]interface{}{
				strings.Split(*i.TaskArn, "/")[1],
				strings.Split(*i.TaskDefinitionArn, "/")[1],
				hostPort,
				*i.LastStatus,
				strings.Split(*i.Containers[0].Image, ":")[1],
			})
		}
		t.AppendRow([]interface{}{ecsTable.Render(), ec2Table.Render()})
	} else {
		t.AppendHeader(table.Row{"ECS"})
		ecsTable.AppendHeader(table.Row{
			"Task ID",
			"task definition",
			"Status",
			"Image Tag",
		})
		for _, i := range svc.Tasks {
			ecsTable.AppendRow([]interface{}{
				strings.Split(*i.TaskArn, "/")[1],
				strings.Split(*i.TaskDefinitionArn, "/")[1],
				*i.LastStatus,
				strings.Split(*i.Containers[0].Image, ":")[1],
			})
		}
		t.AppendRow([]interface{}{ecsTable.Render()})
	}
	t.Render()

}

func getSVCInfo(sess *session.Session, cluster, svn string) error {

	var svcInfo Service
	svcInfo.Name = svn
	svcInfo.Cluster = cluster

	test := listTasks(sess, svcInfo.Cluster, svcInfo.Name)
	svcInfo.Tasks = describeTask(sess, svcInfo.Cluster, test)
	service := describService(sess, svcInfo.Cluster, svcInfo.Name)
	if len(service.LoadBalancers) == 1 {
		targetGroupArn := *service.LoadBalancers[0].TargetGroupArn
		svcInfo.TasksTargetHealth = listHealth(sess, targetGroupArn)
	}
	display(svcInfo)
	return nil
}
