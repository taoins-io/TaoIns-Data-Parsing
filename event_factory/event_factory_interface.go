package event_factory

import (
	"tao/vo"
)

type TypeEvent interface {
	PullEvent(eventNode vo.EventNode)
}
