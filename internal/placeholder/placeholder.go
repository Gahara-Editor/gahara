package placeholder

import (
	"strings"

	"github.com/google/uuid"
)

const (
	NODE_PLACEHOLDER = "NODE_PLACEHOLDER"
)

type PlaceholderNode struct {
	// NodeType: NODE_PLACEHOLDER
	NodeType string `json:"type"`
	// Start: the start of the interval
	PlaceholderStart float64 `json:"start"`
	// End: the end of the interval
	PlaceholderEnd float64 `json:"end"`
	// RID: empty for placeholders
	PlaceholderRID string `json:"rid"`
	// ID: the ID of the placeholder node
	PlaceholderID string `json:"id"`
	// Name: the name given by the user to the placeholder node
	PlaceholderName string `json:"name"`
}

func CreateNode(rid string, name string, start, end float64) *PlaceholderNode {
	return &PlaceholderNode{
		NodeType:         NODE_PLACEHOLDER,
		PlaceholderRID:   rid,
		PlaceholderID:    strings.Replace(uuid.New().String(), "-", "", -1),
		PlaceholderName:  name,
		PlaceholderStart: start,
		PlaceholderEnd:   end,
	}
}

func (p *PlaceholderNode) Type() string {
	return NODE_PLACEHOLDER
}
func (p *PlaceholderNode) RID() string {
	return p.PlaceholderRID
}
func (p *PlaceholderNode) ID() string {
	return p.PlaceholderID
}

func (p *PlaceholderNode) Name() string {
	return p.PlaceholderName
}

func (p *PlaceholderNode) Start() float64 {
	return p.PlaceholderStart
}

func (p *PlaceholderNode) End() float64 {
	return p.PlaceholderEnd
}

func (p *PlaceholderNode) Rename(name string) {
	p.PlaceholderName = name
}
