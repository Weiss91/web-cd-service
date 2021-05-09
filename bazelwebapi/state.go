package main

type STATE int

const (
	WAITING STATE = iota
	RUNNING
	DONE
	/* later maybe
	PREPARING
	BUILDING
	TESTING
	PUSHING
	*/
)

func (i STATE) ToString() string {
	switch i {
	case WAITING:
		return "waiting"
	case RUNNING:
		return "running"
	case DONE:
		return "done"
	default:
		return "not implemented"
	}
}
