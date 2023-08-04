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

// GenPipelineRunQuery returns a Splunk query for generating Segment events
// representing creation of AppStudio build PipelineRuns.
func GenPipelineRunQuery(index string) string {
	q, _ := UserJourneyQueryGen(
		index,
		`verb=create `+
			`"responseStatus.code" IN (200, 201) `+
			`"objectRef.apiGroup"="tekton.dev" `+
			`"objectRef.resource"="pipelineruns" `+
			`"responseObject.metadata.labels.pipelines.appstudio.openshift.io/type"=build`+
			`"responseObject.metadata.resourceVersion"="*"`,
		[]string{"namespace", "application", "component"},
	)
	return q
}
