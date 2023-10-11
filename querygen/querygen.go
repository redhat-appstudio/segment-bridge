// Package querygen is used to generate Splunk queries for fetching user journey
// events from the RHTAP K8s event log.
package querygen

import "fmt"

// GenApplicationQuery returns a Splunk query for generating Segment events
// representing AppStudio Application object events.
func GenApplicationQuery(index string) string {
	q, _ := NewUserJourneyQuery(index).
		WithPredicate(
			`verb=create `+
				`"responseStatus.code" IN (200, 201) `+
				`"objectRef.apiGroup"="appstudio.redhat.com" `+
				`"objectRef.resource"="applications" `+
				`("impersonatedUser.username"="*" OR (user.username="*" AND NOT user.username="system:*")) `+
				`(verb!=create OR "responseObject.metadata.resourceVersion"="*")`,
		).
		WithFields("name", "userId").
		String()
	return q
}

// GenComponentQuery returns a Splunk query for generating Segment events
// representing AppStudio Component object events.
func GenComponentQuery(index string) string {
	q, _ := NewUserJourneyQuery(index).
		WithPredicate(
			`verb IN (create, update, delete, patch) `+
				`"responseStatus.code" IN (200, 201) `+
				`"objectRef.apiGroup"="appstudio.redhat.com" `+
				`"objectRef.resource"="components" `+
				`("impersonatedUser.username"="*" OR (user.username="*" AND NOT user.username="system:*")) `+
				`(verb!=create OR "responseObject.metadata.resourceVersion"="*")`,
		).
		WithFields("name", "userId", "application", "component", "src_url", "src_revision").
		String()
	return q
}

// GenBuildPipelineRunCreatedQuery returns a Splunk query for generating Segment events
// representing creation of AppStudio build PipelineRuns.
func GenBuildPipelineRunCreatedQuery(index string) string {
	q, _ := NewUserJourneyQuery(index).
		WithPredicate(
			`verb=create `+
				`"responseStatus.code" IN (200, 201) `+
				`"objectRef.apiGroup"="tekton.dev" `+
				`"objectRef.resource"="pipelineruns" `+
				`"responseObject.metadata.labels.pipelines.appstudio.openshift.io/type"=build `+
				`"responseObject.metadata.resourceVersion"="*"`,
		).
		WithEventExpr(`"Build PipelineRun created"`).
		WithFields("namespace", "application", "component", "repo", "commit_sha", "target_branch",
			"git_trigger_event_type", "git_trigger_provider", "pipeline_log_url").
		String()
	return q
}

// GenBuildPipelineRunStartedQuery returns a Splunk query for generating Segment events
// representing the start of AppStudio build PipelineRuns.
func GenBuildPipelineRunStartedQuery(index string) string {
	statusFilter := NewStatusConditionFilter("Succeeded")
	statusFilter.opts.reasons = []string{"Running"}
	statusFilter.opts.message = "Tasks Completed: 0 %"

	q, _ := NewUserJourneyQuery(index).
		WithPredicate(
			`verb=update `+
				`"responseStatus.code"=200 `+
				`"objectRef.apiGroup"="tekton.dev" `+
				`"objectRef.resource"="pipelineruns" `+
				`"objectRef.subresource"="status" `+
				`"responseObject.metadata.labels.pipelines.appstudio.openshift.io/type"=build `+
				`"responseObject.metadata.resourceVersion"="*" `+
				`"responseObject.status.startTime"="*"`,
		).
		WithFilter(statusFilter).
		WithEventExpr(`"Build PipelineRun started"`).
		WithFields("namespace", "application", "component",
			"git_trigger_event_type", "git_trigger_provider", "pipeline_log_url").
		String()
	return q
}

// GenClairScanCompletedQuery returns a Splunk query for generating Segment events
// when the clair-scan task completes.
func GenClairScanCompletedQuery(index string) string {
	statusFilter := NewStatusConditionFilter("Succeeded")
	statusFilter.opts.reasons = []string{"Succeeded"}
	statusFilter.opts.statuses = []string{"True"}

	trFilter := NewTektonTaskResultFilter("CLAIR_SCAN_RESULT")

	q, _ := NewUserJourneyQuery(index).
		WithPredicate(
			`verb=update `+
				`"responseStatus.code"=200 `+
				`"objectRef.apiGroup"="tekton.dev" `+
				`"objectRef.resource"="taskruns" `+
				`"objectRef.subresource"="status" `+
				`"requestObject.metadata.labels.tekton.dev/pipelineTask"="clair-scan" `+
				`"responseObject.status.completionTime"="*"`,
		).
		WithFilter(statusFilter).
		WithFilter(trFilter).
		WithCommands(
			fmt.Sprintf(`eval clair_scan_result=%s`, trFilter.FieldSet()["tekton_task_result"].srcExpr),
			`spath input=clair_scan_result, path=vulnerabilities.critical output=clair_scan_result.vulnerabilities.critical`,
			`spath input=clair_scan_result, path=vulnerabilities.high output=clair_scan_result.vulnerabilities.high`,
			`spath input=clair_scan_result, path=vulnerabilities.medium output=clair_scan_result.vulnerabilities.medium`,
			`spath input=clair_scan_result, path=vulnerabilities.low output=clair_scan_result.vulnerabilities.low`,
		).
		WithEventExpr(`"Clair scan TaskRun completed"`).
		WithFields(
			"namespace", "application", "component",
			"vulnerabilities_critical", "vulnerabilities_high",
			"vulnerabilities_medium", "vulnerabilities_low",
		).
		String()
	return q
}

// GenBuildPipelineRunCompletedQuery returns a Splunk query for generating Segment events
// representing success or failure of AppStudio build PipelineRuns.
func GenBuildPipelineRunCompletedQuery(index string) string {
	statusFilter := NewStatusConditionFilter("Succeeded")
	statusFilter.opts.reasons = []string{"Completed", "Failed"}

	eventExpr := fmt.Sprintf(
		`"Build PipelineRun ".%s`,
		statusFilter.FieldSet()["status_reason"].srcExpr,
	)

	q, _ := NewUserJourneyQuery(index).
		WithPredicate(
			`verb=update `+
				`"responseStatus.code"=200 `+
				`"objectRef.apiGroup"="tekton.dev" `+
				`"objectRef.resource"="pipelineruns" `+
				`"objectRef.subresource"="status" `+
				`"responseObject.metadata.labels.pipelines.appstudio.openshift.io/type"=build `+
				`"responseObject.metadata.resourceVersion"="*" `+
				`"responseObject.status.completionTime"="*"`,
		).
		WithFilter(statusFilter).
		WithEventExpr(eventExpr).
		WithFields(
			"namespace", "application", "component",
			"status_message", "status_reason",
			"repo", "commit_sha", "target_branch",
			"git_trigger_event_type", "git_trigger_provider",
			 "pipeline_log_url").
		String()

	return q
}

// GenReleaseCompletedQuery returns a Splunk query for generating Segment events
// representing the Release resource success/failure state changes.
func GenReleaseCompletedQuery(index string) string {
	statusFilter := NewStatusConditionFilter("Released")
	statusFilter.opts.reasons = []string{"Succeeded", "Failed"}

	eventExpr := fmt.Sprintf(
		`"Release ".%s`,
		statusFilter.FieldSet()["status_reason"].srcExpr,
	)

	q, _ := NewUserJourneyQuery(index).
		WithPredicate(
			`verb=patch `+
				`"responseStatus.code"=200 `+
				`"objectRef.apiGroup"="appstudio.redhat.com" `+
				`"objectRef.resource"="releases" `+
				`"objectRef.subresource"="status" `+
				`"responseObject.metadata.resourceVersion"="*" `+
				`"responseObject.status.completionTime"="*"`,
		).
		WithFilter(statusFilter).
		WithEventExpr(eventExpr).
		WithFields("namespace", "name", "status_reason", "status_message").
		String()
	return q
}

// GenPullRequestCreatedQuery returns a Splunk query for generating Segment events
// whenever a Pull request is created in the users GitHub repository.
func GenPullRequestCreatedQuery(index string) string {
	q, _ := NewUserJourneyQuery(index).
		WithPredicate(
			`verb=update `+
				`"responseStatus.code"=200 `+
				`"objectRef.apiGroup"="appstudio.redhat.com" `+
				`"objectRef.resource"="components" `+
				`"user.username"="system:serviceaccount:build-service:build-service-controller-manager" `+
				`"responseObject.metadata.annotations.build.appstudio.openshift.io/status"="*pac*" `+
				`(NOT "responseObject.metadata.annotations.build.appstudio.openshift.io/request"="*")`,
		).
		WithCommands(
			`spath input="responseObject.metadata.annotations.build.appstudio.openshift.io/status", path=pac.state output=build_status.pac.state`,
			`search "build_status.pac.state"="enabled"`,
			`spath input="responseObject.metadata.annotations.build.appstudio.openshift.io/status", path=pac.merge-url output=build_status.pac.merge-url`,
			`dedup build_status.pac.merge-url sortby +_time`,
		).
		WithEventExpr(`"Pull request created"`).
		WithFields("namespace", "name", "application", "component", "merge_url", "src_url", "src_revision").
		String()
	return q
}
