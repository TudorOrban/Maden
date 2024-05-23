package scheduler

import "maden/pkg/shared"

type Scheduler interface {
	SchedulePod(pod *shared.Pod) error
}