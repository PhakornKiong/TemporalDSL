package dsl

import (
	"context"
	"io/ioutil"
	"log"

	dsl "encore.app/dsl/workflow"
	"encore.dev/rlog"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"gopkg.in/yaml.v3"
)

//encore:service
type Service struct {
}

func initService() (*Service, error) {
	workflowData, err := ParseYamlFile("dsl/workflow.yaml")

	if err != nil {
		panic(err)
	}

	customWF := dsl.CustomWorkflow{WfData: workflowData}
	errWF := customWF.Init()

	if errWF != nil {
		panic(errWF)
	}

	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
		Logger:   rlog.With(),
	})

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}

	w := worker.New(c, "dsl", worker.Options{Identity: "temporal-worker"})

	registerOptions := workflow.RegisterOptions{
		Name: customWF.Name,
	}

	w.RegisterWorkflowWithOptions(customWF.WF, registerOptions)

	w.RegisterActivity(&dsl.SampleActivities{})

	err = w.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "dsl_" + uuid.New().String(),
		TaskQueue: "dsl",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, customWF.Name)

	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	rlog.Info("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	return &Service{}, nil
}

func ParseYamlFile(fp string) (dsl.WorkflowData, error) {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return dsl.WorkflowData{}, err
	}

	var workflow dsl.WorkflowData
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return dsl.WorkflowData{}, err
	}

	return workflow, nil
}
