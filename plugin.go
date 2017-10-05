package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/devicefarm"
)

// Plugin defines the Device farm plugin parameters.
type Plugin struct {
	Key          string
	Secret       string
	Region       string
	YamlVerified bool

	TestProject string
	RunName     string
}

// Exec runs the plugin
func (p *Plugin) Exec() error {
	fmt.Println("Begin ")
	fmt.Println("Plugin.Key ", p.Key)
	fmt.Println("Plugin.Secret ", p.Secret)
	fmt.Println("Plugin.TestProject ", p.TestProject)
	fmt.Println("Plugin.RunName ", p.RunName)

	// create the configuration
	conf := &aws.Config{
		Region: aws.String(p.Region),
	}

	// Use key and secret if provided otherwise fall back to ec2 instance profile
	if p.Key != "" && p.Secret != "" {
		conf.Credentials = credentials.NewStaticCredentials(p.Key, p.Secret, "")
	}
	//create Device Farm service
	svc := devicefarm.New(session.New(), conf)

	//Get AWS device farm Test project
	project, err := getTestProject(p.TestProject, svc)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		//Get the run to see the status
		run, err := getRun(p.RunName, project, svc)
		if err != nil {
			return err
		}
		fmt.Println("run", run)

		if *run.Status == "COMPLETED" && (*run.Result == "ERRORED" || *run.Result == "FAILED") {
			return fmt.Errorf("The test run has failed")
		}
		if *run.Status == "COMPLETED" && *run.Result == "PASSED" {
			break
		}
	}

	return nil
}

func getRun(runName string, project *devicefarm.Project, svc *devicefarm.DeviceFarm) (*devicefarm.Run, error) {

	var result *devicefarm.ListRunsOutput
	var listRunInput devicefarm.ListRunsInput
	var err error

	for {
		listRunInput.Arn = aws.String(*project.Arn)

		if result != nil && result.NextToken != nil {
			listRunInput.NextToken = aws.String(*result.NextToken)
		}
		result, err = svc.ListRuns(&listRunInput)

		if err != nil {
			return nil, err
		}

		for _, run := range result.Runs {
			if *run.Name == runName {
				return run, nil
			}
		}

		if result.NextToken == nil {
			return nil, fmt.Errorf("There was no Run with the name %s", runName)
		}
	}
}

func getTestProject(testProjectName string, svc *devicefarm.DeviceFarm) (*devicefarm.Project, error) {
	result, err := svc.ListProjects(nil)
	if err != nil {
		return nil, err
	}

	for _, project := range result.Projects {
		if *project.Name == testProjectName {
			return project, nil
		}
	}

	return nil, fmt.Errorf("There was no project with the name %s", testProjectName)
}
