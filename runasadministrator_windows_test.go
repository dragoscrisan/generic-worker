package main

import (
	"path/filepath"
	"testing"
)

func TestRunAsAdministratorEnabled(t *testing.T) {
	defer setup(t)()
	payload := GenericWorkerPayload{
		Command: []string{
			`whoami /groups | find "12288" > nul`,
		},
		MaxRunTime: 10,
		Features: FeatureFlags{
			RunAsAdministrator: true,
		},
		OSGroups: []string{
			"Administrators",
		},
	}
	td := testTask(t)
	td.Scopes = []string{
		"generic-worker:run-as-administrator:" + td.ProvisionerID + "/" + td.WorkerType,
		"generic-worker:os-group:" + td.ProvisionerID + "/" + td.WorkerType + "/Administrators",
	}

	_ = submitAndAssert(t, td, payload, "completed", "completed")
}

func TestRunAsAdministratorDisabled(t *testing.T) {
	defer setup(t)()
	payload := GenericWorkerPayload{
		Command: []string{
			`whoami /groups | find "12288" > nul`,
		},
		MaxRunTime: 10,
	}
	td := testTask(t)

	_ = submitAndAssert(t, td, payload, "failed", "failed")
}

func TestRunAsAdministratorEnabledMissingScopes(t *testing.T) {
	defer setup(t)()
	payload := GenericWorkerPayload{
		Command: []string{
			`whoami /groups | find "12288" > nul`,
		},
		MaxRunTime: 10,
		Features: FeatureFlags{
			RunAsAdministrator: true,
		},
		OSGroups: []string{
			"Administrators",
		},
	}
	td := testTask(t)
	td.Scopes = []string{
		"generic-worker:os-group:" + td.ProvisionerID + "/" + td.WorkerType + "/Administrators",
	}

	_ = submitAndAssert(t, td, payload, "exception", "malformed-payload")
}

func TestRunAsAdministratorMissingOSGroup(t *testing.T) {
	defer setup(t)()
	payload := GenericWorkerPayload{
		Command: []string{
			`whoami /groups | find "12288" > nul`,
		},
		MaxRunTime: 10,
		OSGroups:   []string{}, // Administrators not included!
		Features: FeatureFlags{
			RunAsAdministrator: true,
		},
	}
	td := testTask(t)
	td.Scopes = []string{
		"generic-worker:run-as-administrator:" + td.ProvisionerID + "/" + td.WorkerType,
		"generic-worker:os-group:" + td.ProvisionerID + "/" + td.WorkerType + "/Administrators",
	}

	_ = submitAndAssert(t, td, payload, "exception", "malformed-payload")
}

func TestChainOfTrustWithRunAsAdministrator(t *testing.T) {
	defer setup(t)()
	payload := GenericWorkerPayload{
		Command: []string{
			`type "` + filepath.Join(cwd, config.SigningKeyLocation) + `"`,
		},
		MaxRunTime: 5,
		OSGroups:   []string{"Administrators"},
		Features: FeatureFlags{
			ChainOfTrust:       true,
			RunAsAdministrator: true,
		},
	}
	td := testTask(t)
	td.Scopes = []string{
		"generic-worker:run-as-administrator:" + td.ProvisionerID + "/" + td.WorkerType,
		"generic-worker:os-group:" + td.ProvisionerID + "/" + td.WorkerType + "/Administrators",
	}

	if config.RunTasksAsCurrentUser {
		// When running as current user, chain of trust key is not private so
		// generic-worker should detect that it isn't secured from task user
		// and cause malformed-payload exception.
		expectChainOfTrustKeyNotSecureMessage(t, td, payload)
		return

	}

	_ = submitAndAssert(t, td, payload, "exception", "malformed-payload")
}

func TestChainOfTrustWithoutRunAsAdministrator(t *testing.T) {
	defer setup(t)()
	payload := GenericWorkerPayload{
		Command: []string{
			`type "` + filepath.Join(cwd, config.SigningKeyLocation) + `"`,
		},
		MaxRunTime: 5,
		OSGroups:   []string{"Administrators"},
		Features: FeatureFlags{
			ChainOfTrust:       true,
			RunAsAdministrator: false, // FALSE !!!!
		},
	}
	td := testTask(t)
	td.Scopes = []string{
		"generic-worker:run-as-administrator:" + td.ProvisionerID + "/" + td.WorkerType,
		"generic-worker:os-group:" + td.ProvisionerID + "/" + td.WorkerType + "/Administrators",
	}

	if config.RunTasksAsCurrentUser {
		// When running as current user, chain of trust key is not private so
		// generic-worker should detect that it isn't secured from task user
		// and cause malformed-payload exception.
		expectChainOfTrustKeyNotSecureMessage(t, td, payload)
		return

	}

	_ = submitAndAssert(t, td, payload, "failed", "failed")
}
