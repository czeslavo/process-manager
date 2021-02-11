package activities

import (
	"time"

	"github.com/czeslavo/process-manager/2_temporal/worker"
	"go.temporal.io/sdk/workflow"
)

var DefaultActivityOptions = workflow.ActivityOptions{
	TaskQueue:              worker.TaskQueue,
	ScheduleToCloseTimeout: time.Second * 60,
	ScheduleToStartTimeout: time.Second * 60,
	StartToCloseTimeout:    time.Second * 60,
	HeartbeatTimeout:       time.Second * 10,
	WaitForCancellation:    false,
}
