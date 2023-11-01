// Package queryprint contains utilities for printing one or more Splunk queries
package queryprint

import "strings"

// QueryDesc includes a printable description of a Splunk query: A descriptive
// title for it and the query string itself.
type QueryDesc struct {
	Title string
	Query string
}

// PrettyPrintQueries prints the given set of queries in a human-readable format
func PrettyPrintQueries(queries []QueryDesc) string {
	var builder strings.Builder
	for i, query := range queries {
		if i > 0 {
			builder.WriteString("\n")
			builder.WriteString("\n")
		}
		builder.WriteString(query.Title)
		builder.WriteString("\n")
		builder.WriteString(strings.Repeat("-", len(query.Title)))
		builder.WriteString("\n")
		builder.WriteString(prettyPrintQuery(query.Query, "    "))
	}
	return builder.String()
}

func prettyPrintQuery(query string, indent string) string {
	var builder strings.Builder
	if len(query) < 60 {
		builder.WriteString(indent)
		builder.WriteString(query)
	} else {
		for i, line := range strings.Split(query, "|") {
			if i > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(indent)
			if i > 0 {
				builder.WriteString("|")
			}
			builder.WriteString(line)
		}
	}
	return builder.String()
}

// MachinePrintQueries prints the given set of queries in a compact,
// machine-readable format where queries are separated by NULL ("\x00")
// characters
func MachinePrintQueries(queries []QueryDesc) string {
	var builder strings.Builder
	for i, query := range queries {
		if i > 0 {
			builder.WriteRune('\x00')
		}
		builder.WriteString(query.Query)
	}
	return builder.String()
}
