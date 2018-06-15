package main

import (
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

func (l *RunAsAdministratorTask) Start() *CommandExecutionError {
	return nil
}

func (l *RunAsAdministratorTask) Stop(err *ExecutionErrors) {
}
