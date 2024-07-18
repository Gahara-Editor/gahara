package timeline

const (
	NODE_VIDEO = iota + 1
	NODE_AUDIO
	NODE_TEXT
	NODE_PLACEHOLDER
)

type TimelineNode interface {
	Type() string
	RID() string
	ID() string
	Name() string
	Start() float64
	End() float64
	Rename(string)
}

type Timeline struct {
	// Nodes: all the nodes of the timeline
	Nodes [][]TimelineNode `json:"timeline"`
}
