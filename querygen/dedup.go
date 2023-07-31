package querygen

import (
	"fmt"
	"strings"
)

// GenDedupEval generates a Splunk "eval" command for converting multi value
// fields into single-value fields. See RHTAPWATCH-293 for why we need this.
func GenDedupEval(fields []string) string {
	var builder strings.Builder
	builder.WriteString("eval ")
	for i, field := range fields {
		if i > 0 {
			builder.WriteString(",")
		}
		fmt.Fprintf(&builder, `%s=mvindex('%s', 0)`, field, field)
	}
	return builder.String()
}
