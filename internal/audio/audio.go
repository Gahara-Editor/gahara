package audio

import (
	"strings"

	"github.com/google/uuid"
)

const (
	NODE_AUDIO = "NODE_AUDIO"
)

type AudioNode struct {
	// NodeType: NODE_AUDIO
	NodeType string `json:"type"`
	// Start: the start of the interval
	AudioStart float64 `json:"start"`
	// End: the end of the interval
	AudioEnd float64 `json:"end"`
	// RID: the root ID of the node, that is, the original audio from which this nodes derives
	AudioRID string `json:"rid"`
	// ID: the ID of the audio node
	AudioID string `json:"id"`
	// Name: the name given by the user to the audio node
	AudioName string `json:"name"`
}

func CreateNode(rid string, name string, start, end float64) *AudioNode {
	return &AudioNode{
		NodeType:   NODE_AUDIO,
		AudioRID:   rid,
		AudioID:    strings.Replace(uuid.New().String(), "-", "", -1),
		AudioName:  name,
		AudioStart: start,
		AudioEnd:   end,
	}
}

func (a *AudioNode) Type() string {
	return NODE_AUDIO
}
func (a *AudioNode) RID() string {
	return a.AudioRID
}
func (a *AudioNode) ID() string {
	return a.AudioID
}

func (a *AudioNode) Name() string {
	return a.AudioName
}

func (a *AudioNode) Start() float64 {
	return a.AudioStart
}

func (a *AudioNode) End() float64 {
	return a.AudioEnd
}

func (a *AudioNode) Rename(name string) {
	a.AudioName = name
}
