package progress

import (
	"fmt"
	"os"
	"strings"
	"sync"

	terminal "github.com/codemodify/systemkit-terminal"
	progress "github.com/codemodify/systemkit-terminal-progress"
)

// TreeRenderer -
type TreeRenderer interface {
	progress.Renderer

	RunByID(id string)
	SuccessByID(id string)
	FailByID(id string)
}

// staticTree -
type staticTree struct {
	tree            *Node
	stopChannel     chan bool
	finishedChannel chan bool
	theTerminal     *terminal.Terminal
	savedCursorX    int
	savedCursorY    int
	mutex           sync.Mutex
}

// NewStaticTree -
func NewStaticTree(tree *Node) TreeRenderer {

	// 1. set defaults
	if tree.Config.Writer == nil {
		tree.Config.Writer = os.Stdout
	}

	if tree.ID == "" {
		tree.ID = "-9999999"
	}

	// 2.
	return &staticTree{
		tree:            tree,
		stopChannel:     make(chan bool),
		finishedChannel: make(chan bool),
		theTerminal:     terminal.NewTerminal(tree.Config.Writer),
		savedCursorX:    0,
		savedCursorY:    0,
		mutex:           sync.Mutex{},
	}
}

// Run -
func (thisRef *staticTree) Run() {
	thisRef.mutex.Lock()
	defer thisRef.mutex.Unlock()

	// 0.
	if thisRef.tree.Config.HideCursor {
		thisRef.theTerminal.CursorHide()
	}

	// 1. save cursor position
	thisRef.savedCursorX, thisRef.savedCursorY = thisRef.theTerminal.CursorPositionQuery()

	// 2. set everything as queued
	isQueued := true
	isRunning := false
	isSuccess := false
	updateTreeStatus(thisRef.tree, "", isQueued, isRunning, isSuccess)

	// 3. set first as running
	isQueued = false
	isRunning = true
	isSuccess = false
	updateTreeStatus(thisRef.tree, thisRef.tree.ID, isQueued, isRunning, isSuccess)

	// 4.
	thisRef.drawTree()
}

// Success -
func (thisRef *staticTree) Success() {
	thisRef.mutex.Lock()
	defer thisRef.mutex.Unlock()

	// set everything as success
	isQueued := false
	isRunning := false
	isSuccess := true
	thisRef.stop("", isQueued, isRunning, isSuccess)
}

// Fail -
func (thisRef *staticTree) Fail() {
	thisRef.mutex.Lock()
	defer thisRef.mutex.Unlock()

	// set everything as failed
	isQueued := false
	isRunning := false
	isSuccess := false
	thisRef.stop("", isQueued, isRunning, isSuccess)
}

// RunByID -
func (thisRef *staticTree) RunByID(id string) {
	thisRef.mutex.Lock()
	defer thisRef.mutex.Unlock()

	isQueued := false
	isRunning := true
	isSuccess := true
	updateTreeStatus(thisRef.tree, id, isQueued, isRunning, isSuccess)
	thisRef.drawTree()
}

// SuccessByID -
func (thisRef *staticTree) SuccessByID(id string) {
	thisRef.mutex.Lock()
	defer thisRef.mutex.Unlock()

	isQueued := false
	isRunning := false
	isSuccess := true
	updateTreeStatus(thisRef.tree, id, isQueued, isRunning, isSuccess)
	thisRef.drawTree()
}

// FailByID -
func (thisRef *staticTree) FailByID(id string) {
	thisRef.mutex.Lock()
	defer thisRef.mutex.Unlock()

	isQueued := false
	isRunning := false
	isSuccess := false
	updateTreeStatus(thisRef.tree, id, isQueued, isRunning, isSuccess)
	thisRef.drawTree()
}

func (thisRef *staticTree) stop(nodeID string, isQueued bool, isRunning bool, isSuccess bool) {
	updateTreeStatus(thisRef.tree, nodeID, isQueued, isRunning, isSuccess)

	thisRef.stopChannel <- true
	close(thisRef.stopChannel)

	<-thisRef.finishedChannel
}

func (thisRef *staticTree) eraseTree(renderedTree string) {
	sb := strings.Builder{}
	for _, line := range strings.Split(renderedTree, "\n") {
		sb.WriteString(strings.Repeat(" ", len(line)) + "\n")
	}

	thisRef.theTerminal.CursorMoveToXY(thisRef.savedCursorX, thisRef.savedCursorY)
	fmt.Fprint(thisRef.tree.Config.Writer, "\r"+sb.String()+"\r")
}

func (thisRef *staticTree) drawTree() {
	renderedTree := renderTree(thisRef.tree, 0)

	thisRef.eraseTree(renderedTree)

	thisRef.theTerminal.CursorMoveToXY(thisRef.savedCursorX, thisRef.savedCursorY)
	fmt.Fprintf(thisRef.tree.Config.Writer, renderedTree)

	if !thisRef.tree.isQueued && !thisRef.tree.isRunning {
		if thisRef.tree.Config.HideCursor {
			thisRef.theTerminal.CursorShow()
		}
	}
}
