/*
QueryGen generates Splunk queries for obtaining RHTAP user journey events from
the RHTAP cluster audit logs stored in Splunk.

Usage:

	querygen [flags]

The flags are:

	    --index INDEX
		    Specify the Splunk index to query.
		-0
			Print in a format suitable for `xargs -0`
*/
package main

import (
	"flag"
	"fmt"

	"github.com/redhat-appstudio/segment-bridge.git/querygen"
	"github.com/redhat-appstudio/segment-bridge.git/queryprint"
)

var index = flag.String(
	"index",
	"federated:rh_rhtap_stage_audit",
	"the Splunk index to query",
)
var machinePrint = flag.Bool(
	"0",
	false,
	"Output queries in machine-readable compact format, "+
		"seperated by NULL (\"\x00\") characters",
)

func main() {
	flag.Parse()
	printFunc := queryprint.PrettyPrintQueries
	if *machinePrint {
		printFunc = queryprint.MachinePrintQueries
	}
	fmt.Println(printFunc([]queryprint.QueryDesc{
		{
			Title: "Application events",
			Query: querygen.GenApplicationQuery(*index),
		},
		{
			Title: "Component events",
			Query: querygen.GenComponentQuery(*index),
		},
		{
			Title: "Build PipelineRun creation events",
			Query: querygen.GenBuildPipelineRunCreatedQuery(*index),
		},
		{
			Title: "Build PipelineRun started events",
			Query: querygen.GenBuildPipelineRunStartedQuery(*index),
		},
		{
			Title: "Clair scan TaskRun completion events",
			Query: querygen.GenClairScanCompletedQuery(*index),
		},
		{
			Title: "Build PipelineRun Completed or Failed event",
			Query: querygen.GenBuildPipelineRunCompletedQuery(*index),
		},
		{
			Title: "Release Succeeded or Failed events",
			Query: querygen.GenReleaseCompletedQuery(*index),
		},
	}))
}
