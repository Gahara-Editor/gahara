package text

const (
	NODE_TEXT = "NODE_TEXT"
)

type TextNode struct {
	// NodeType: NODE_TEXT
	NodeType string `json:"type"`
	// Start: the start of the interval
	Start float64 `json:"start"`
	// End: the end of the interval
	End float64 `json:"end"`
	// The text to be attached
	Text string `json:"text"`
	// Name: the name given by the user to the text node
	Name string `json:"name"`
	// ID: the ID of the text node
	ID string `json:"id"`
}
