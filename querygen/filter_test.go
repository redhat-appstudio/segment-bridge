package querygen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusConditionFilter(t *testing.T) {
	f := NewStatusConditionFilter("TestType")
	f.opts.reasons = []string{"r1", "r2"}
	f.opts.statuses = []string{"s1", "s2"}
	f.opts.message = "prefix % postfix"

	assert.Equal(t,
		[]string{
			`eval status_condition_index=mvfind('responseObject.status.conditions{}.type', "TestType")`,
			`where isnotnull(status_condition_index) ` +
				`AND mvindex('responseObject.status.conditions{}.reason', status_condition_index) IN ("r1", "r2") ` +
				`AND mvindex('responseObject.status.conditions{}.status', status_condition_index) IN ("s1", "s2") ` +
				`AND like(mvindex('responseObject.status.conditions{}.message', status_condition_index), "prefix % postfix")`,
		},
		f.Commands(),
	)
}

func TestTektonTaskResultFilter(t *testing.T) {
	f := NewTektonTaskResultFilter("result-a")

	assert.Equal(t,
		[]string{
			`eval tekton_task_result_index=mvfind('responseObject.status.taskResults{}.name', "result-a")`,
			`where isnotnull(tekton_task_result_index)`,
		},
		f.Commands(),
	)
}
