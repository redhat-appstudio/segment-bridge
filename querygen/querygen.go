// Package querygen is used to generate Splunk queries for fetching user journey
// events from the RHTAP K8s event log.
package querygen

// GenApplicationQuery returns a Splunk query for generating Segment events
// representing AppStudio Application object events.
func GenApplicationQuery(index string) string {
	q, _ := UserJourneyQueryGen(
		index,
		`verb=create `+
			`"responseStatus.code" IN (200, 201) `+
			`"objectRef.apiGroup"="appstudio.redhat.com" `+
			`"objectRef.resource"="applications" `+
			`("impersonatedUser.username"="*" OR (user.username="*" AND NOT user.username="system:*")) `+
			`(verb!=create OR "responseObject.metadata.resourceVersion"="*")`,
		[]string{"name", "userId"},
	)
	return q
}

// GenBuildPipelineRunCreatedQuery returns a Splunk query for generating Segment events
// representing creation of AppStudio build PipelineRuns.
func GenBuildPipelineRunCreatedQuery(index string) string {
	q, _ := UserJourneyQueryGen(
		index,
		`verb=create `+
			`"responseStatus.code" IN (200, 201) `+
			`"objectRef.apiGroup"="tekton.dev" `+
			`"objectRef.resource"="pipelineruns" `+
			`"responseObject.metadata.labels.pipelines.appstudio.openshift.io/type"=build `+
			`"responseObject.metadata.resourceVersion"="*" `+
			`| eval event="Build PipelineRun created" `,
		[]string{"namespace", "application", "component"},
	)
	return q
}

// GenBuildPipelineRunStartedQuery returns a Splunk query for generating Segment events
// representing the start of AppStudio build PipelineRuns.
func GenBuildPipelineRunStartedQuery(index string) string {
	q, _ := UserJourneyQueryGen(
		index,
		`verb=update `+
			`"responseStatus.code"=200 `+
			`"objectRef.apiGroup"="tekton.dev" `+
			`"objectRef.resource"="pipelineruns" `+
			`"objectRef.subresource"="status" `+
			`"responseObject.metadata.labels.pipelines.appstudio.openshift.io/type"=build `+
			`"responseObject.metadata.resourceVersion"="*" `+
			`"responseObject.status.startTime"="*" `+
			`| eval event="Build PipelineRun started" `+
			`| spath path=responseObject.status.conditions{} output=conditions `+
			`| mvexpand conditions `+
			`| spath input=conditions `+
			`| where type="Succeeded" AND reason="Running" AND like(message, "Tasks Completed: 0 %") `,
		[]string{"namespace", "application", "component"},
	)
	return q
}

// GenClairScanCompletedQuery returns a Splunk query for generating Segment events
// when the clair-scan task completes.
func GenClairScanCompletedQuery(index string) string {
	q, _ := UserJourneyQueryGen(
		index,
		`verb=update `+
			`"responseStatus.code"=200 `+
			`"objectRef.apiGroup"="tekton.dev" `+
			`"objectRef.resource"="taskruns" `+
			`"objectRef.subresource"="status" `+
			`"requestObject.metadata.labels.tekton.dev/pipelineTask"="clair-scan" `+
			`"responseObject.status.completionTime"="*" `+
			`| spath responseObject.status.conditions{} `+
			`| mvexpand responseObject.status.conditions{} `+
			`| search responseObject.status.conditions{}.type="Succeeded" `+
			`responseObject.status.conditions{}.reason="Succeeded" `+
			`responseObject.status.conditions{}.status="True" `+
			`| spath responseObject.status.taskResults{} `+
			`| mvexpand responseObject.status.taskResults{} `+
			`| search responseObject.status.taskResults{}.name="CLAIR_SCAN_RESULT" `+
			`| spath input=responseObject.status.taskResults{}.value path=vulnerabilities.critical output=clair_scan_result.vulnerabilities.critical `+
			`| spath input=responseObject.status.taskResults{}.value path=vulnerabilities.high output=clair_scan_result.vulnerabilities.high `+
			`| spath input=responseObject.status.taskResults{}.value path=vulnerabilities.medium output=clair_scan_result.vulnerabilities.medium `+
			`| spath input=responseObject.status.taskResults{}.value path=vulnerabilities.low output=clair_scan_result.vulnerabilities.low `+
			`| eval event="Clair scan TaskRun completed"`,
		[]string{"namespace", "application", "component", "vulnerabilities_critical", "vulnerabilities_high", "vulnerabilities_medium", "vulnerabilities_low"},
	)
	return q
}
