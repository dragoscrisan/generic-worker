package main

import (
	"fmt"
	"net/http"

	"github.com/taskcluster/generic-worker/win32"
	"github.com/taskcluster/taskcluster-base-go/scopes"
)

type RunAsAdministratorFeature struct {
}

func (feature *RunAsAdministratorFeature) Name() string {
	return "Run As Administrator"
}

func (feature *RunAsAdministratorFeature) Initialise() error {
	return nil
}

func (feature *RunAsAdministratorFeature) PersistState() error {
	return nil
}

func (feature *RunAsAdministratorFeature) IsEnabled(task *TaskRun) bool {
	return task.Payload.Features.RunAsAdministrator
}

type RunAsAdministratorTask struct {
	task *TaskRun
}

func (l *RunAsAdministratorTask) ReservedArtifacts() []string {
	return []string{}
}

func (feature *RunAsAdministratorFeature) NewTaskFeature(task *TaskRun) TaskFeature {
	return &RunAsAdministratorTask{
		task: task,
	}
}

func (l *RunAsAdministratorTask) RequiredScopes() scopes.Required {
	return scopes.Required{{
		"generic-worker:run-as-administrator:" + config.ProvisionerID + "/" + config.WorkerType,
	}}
}

type RunAsAdministratorHandler struct {
}

func (l *RunAsAdministratorTask) Start() *CommandExecutionError {
	for _, c := range l.task.Commands {
		adminToken, err := win32.GetLinkedToken(c.Cmd.SysProcAttr.Token)
		if err != nil {
			return Failure(fmt.Errorf("Could not get auth token to run command as administrator: %v", err))
		}
		c.Cmd.SysProcAttr.Token = adminToken
	}
	return nil
}

func ProcessRequests(s *http.Server) {
	// s.ListenAndServe()
}

func (l *RunAsAdministratorTask) Stop(err *ExecutionErrors) {
}

func (handler *RunAsAdministratorHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
}
