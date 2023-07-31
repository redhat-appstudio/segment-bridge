package querygen

import "testing"

func TestGenDedupEval(t *testing.T) {
	fields := []string{
		"auditID",
		"impersonatedUser.username",
		"user.username",
		"objectRef.resource",
		"objectRef.namespace",
		"objectRef.apiGroup",
		"objectRef.apiVersion",
		"objectRef.name",
		"verb",
		"requestReceivedTimestamp",
		"responseObject.metadata.labels.appstudio.openshift.io/application",
		"responseObject.metadata.labels.appstudio.openshift.io/component",
	}
	expected := `eval ` +
		`auditID=mvindex('auditID', 0),` +
		`impersonatedUser.username=mvindex('impersonatedUser.username', 0),` +
		`user.username=mvindex('user.username', 0),` +
		`objectRef.resource=mvindex('objectRef.resource', 0),` +
		`objectRef.namespace=mvindex('objectRef.namespace', 0),` +
		`objectRef.apiGroup=mvindex('objectRef.apiGroup', 0),` +
		`objectRef.apiVersion=mvindex('objectRef.apiVersion', 0),` +
		`objectRef.name=mvindex('objectRef.name', 0),` +
		`verb=mvindex('verb', 0),` +
		`requestReceivedTimestamp=mvindex('requestReceivedTimestamp', 0),` +
		"responseObject.metadata.labels.appstudio.openshift.io/application" +
		"=mvindex('responseObject.metadata.labels.appstudio.openshift.io/application', 0)," +
		"responseObject.metadata.labels.appstudio.openshift.io/component" +
		"=mvindex('responseObject.metadata.labels.appstudio.openshift.io/component', 0)"
	out := GenDedupEval(fields)
	if out != expected {
		t.Errorf("GenDedupEval() returns:\n  %q, expected\n  %q", out, expected)
	}
}
