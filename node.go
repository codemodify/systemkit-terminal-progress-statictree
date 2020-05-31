package progress

import (
	"fmt"
	"strings"

	progress "github.com/codemodify/systemkit-terminal-progress"
)

// Node -
type Node struct {
	ID       string
	Config   *progress.Config
	Children []*Node

	isQueued  bool
	isRunning bool
	isSuccess bool
}

func updateTreeStatus(node *Node, nodeID string, isQueued bool, isRunning bool, isSuccess bool) {
	if node.ID == nodeID || nodeID == "" {
		node.isQueued = isQueued
		node.isRunning = isRunning
		node.isSuccess = isSuccess
	}

	for _, child := range node.Children {
		updateTreeStatus(child, nodeID, isQueued, isRunning, isSuccess)
	}

	if len(node.Children) > 0 {
		allQueued := true
		atLeastOneIsRunning := false
		allSucceeded := true
		for _, child := range node.Children {
			allQueued = allQueued && child.isQueued

			if child.isRunning {
				atLeastOneIsRunning = true
			}

			allSucceeded = allSucceeded && child.isSuccess
		}

		if allQueued {
			node.isQueued = true
		} else if atLeastOneIsRunning {
			node.isQueued = false
			node.isRunning = true
		} else {
			node.isQueued = false
			node.isRunning = false
			node.isSuccess = allSucceeded
		}
	}
}

func renderTree(node *Node, level int) string {
	thisLevelSpacesOnTheLeft := strings.Repeat(" ", 2*level)

	stateGlyph := ""
	stateMessage := ""

	if node.isQueued {
		stateGlyph = node.Config.QueuedGlyph
		stateMessage = node.Config.QueuedMessage
	} else if node.isRunning {
		stateGlyph = node.Config.ProgressGlyphs[0]
		stateMessage = node.Config.ProgressMessage
	} else if node.isSuccess {
		stateGlyph = node.Config.SuccessGlyph
		stateMessage = node.Config.SuccessMessage
	} else {
		stateGlyph = node.Config.FailGlyph
		stateMessage = node.Config.FailMessage
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(
		"%s%s%s%s%s\n",
		thisLevelSpacesOnTheLeft,
		node.Config.Prefix,
		stateGlyph,
		node.Config.Suffix,
		stateMessage,
	))

	for _, child := range node.Children {
		renderedNode := renderTree(child, level+1)
		sb.WriteString(renderedNode)
	}

	return sb.String()
}
