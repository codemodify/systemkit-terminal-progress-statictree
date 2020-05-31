package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	progress "github.com/codemodify/systemkit-terminal-progress"
	treeProgress "github.com/codemodify/systemkit-terminal-progress-statictree"
)

func main() {
	tree := &treeProgress.Node{
		ID:     "root",
		Config: progress.NewDefaultConfig("Root"),
		Children: []*treeProgress.Node{
			{
				ID:     "child-1",
				Config: progress.NewDefaultConfig("Child - 1"),
				Children: []*treeProgress.Node{
					{
						ID:     "child-1-1",
						Config: progress.NewDefaultConfig("Child - 1-1"),
					},
					{
						ID:     "child-1-2",
						Config: progress.NewDefaultConfig("Child - 1-2"),
					},
					{
						ID:     "child-1-3",
						Config: progress.NewDefaultConfig("Child - 1-3"),
					},
					{
						ID:     "child-1-4",
						Config: progress.NewDefaultConfig("Child - 1-4"),
					},
				},
			},
			{
				ID:     "child-2",
				Config: progress.NewDefaultConfig("Child - 2"),
				Children: []*treeProgress.Node{
					{
						ID:     "child-2-1",
						Config: progress.NewDefaultConfig("Child - 2-1"),
					},
					{
						ID:     "child-2-2",
						Config: progress.NewDefaultConfig("Child - 2-2"),
					},
					{
						ID:     "child-2-3",
						Config: progress.NewDefaultConfig("Child - 2-3"),
					},
					{
						ID:     "child-2-4",
						Config: progress.NewDefaultConfig("Child - 2-4"),
					},
				},
			},
			{
				ID:     "child-3",
				Config: progress.NewDefaultConfig("Child - 3"),
				Children: []*treeProgress.Node{
					{
						ID:     "child-3-1",
						Config: progress.NewDefaultConfig("Child - 3-1"),
					},
					{
						ID:     "child-3-2",
						Config: progress.NewDefaultConfig("Child - 3-2"),
					},
					{
						ID:     "child-3-3",
						Config: progress.NewDefaultConfig("Child - 3-3"),
					},
					{
						ID:     "child-3-4",
						Config: progress.NewDefaultConfig("Child - 3-4"),
					},
				},
			},
			{
				ID:     "child-4",
				Config: progress.NewDefaultConfig("Child - 4"),
				Children: []*treeProgress.Node{
					{
						ID:     "child-4-1",
						Config: progress.NewDefaultConfig("Child - 4-1"),
					},
					{
						ID:     "child-4-2",
						Config: progress.NewDefaultConfig("Child - 4-2"),
					},
					{
						ID:     "child-4-3",
						Config: progress.NewDefaultConfig("Child - 4-3"),
					},
					{
						ID:     "child-4-4",
						Config: progress.NewDefaultConfig("Child - 4-4"),
					},
				},
			},
		},
	}

	tp := treeProgress.NewStaticTree(tree)
	tp.Run()

	wg := sync.WaitGroup{}
	rand.Seed(100000000)
	for i := 1; i <= 4; i++ {
		for j := 1; j <= 4; j++ {

			wg.Add(1)
			go func(a int, b int) {
				nodeID := fmt.Sprintf("child-%d-%d", a, b)

				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				tp.RunByID(nodeID)

				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				if (rand.Intn(10) % 2) == 0 {
					tp.SuccessByID(nodeID)
				} else {
					tp.FailByID(nodeID)
				}

				wg.Done()
			}(i, j)
		}
	}

	wg.Wait()
}
