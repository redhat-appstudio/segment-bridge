/*
QueryGen generates Splunk queries for obtaining RHTAP user journey events from
the RHTAP cluster audit logs stored in Splunk.

Usage:

    querygen [flags]

The flags are:

    --index INDEX
	    Specify the Splunk index to query.
*/
package main

import (
	"flag"
	"fmt"

	"github.com/redhat-appstudio/segment-bridge.git/querygen"
)

var index = flag.String(
	"index",
	"federated:rh_rhtap_stage_audit",
	"the Splunk index to query",
)

func main() {
	flag.Parse()
	fmt.Println(querygen.GenApplicationQuery(*index))
}
