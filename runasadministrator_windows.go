package main

import (
	"fmt"
	"log"
	"net/http"
	"syscall"

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
	if config.RunTasksAsCurrentUser {
		// already running as LocalSystem with UAC elevation
		return nil
	}
	for _, c := range l.task.Commands {
		log.Printf("User token is %v", c.Cmd.SysProcAttr.Token)
		adminToken, err := win32.GetLinkedToken(syscall.Handle(c.Cmd.SysProcAttr.Token))
		if err != nil {
			return MalformedPayloadError(fmt.Errorf(`Could not obtain UAC elevated auth token; you probably need to add group "Administrators" to task.payload.osGroups: %v`, err))
		}
		c.Cmd.SysProcAttr.Token = syscall.Token(adminToken)
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
