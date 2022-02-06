package selectable

import (
	"bufio"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

type Cell struct {
	Content  string
	Disabled bool
	Color    *color.Color
}

type Table struct {
	cols          []string
	rows          [][]Cell
	Multiple      bool
	Width         int
	Height        int
	activeRow     int
	activeCol     int
	selected      [][]int
	HoverColor    *color.Color
	SelectedColor *color.Color
	NormalColor   *color.Color
	HeaderColor   *color.Color
}

func (t *Table) DefineCol(colLabel string) {
	t.cols = append(t.cols, colLabel)

}

func (t *Table) AddRow(cellContents []Cell) {
	t.rows = append(t.rows, cellContents)
}

func (t *Table) isSelected(rowIndex, colIndex int) bool {
	for _, item := range t.selected {
		if item[0] == rowIndex && item[1] == colIndex {
			return true
		}
	}
	return false
}

func (t *Table) Run() [][]int {
	go print(t)
	readKey(t)

	return t.selected
}

var reader = bufio.NewReader(os.Stdin)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
func readKey(t *Table) {
	keyboard.Open()
	defer func() {
		_ = keyboard.Close()
	}()
	for {
		_, key, _ := keyboard.GetKey()
		if key == keyboard.KeyArrowDown {
			t.activeRow = min(len(t.rows)-1, t.activeRow+1)
			go print(t)
		} else if key == keyboard.KeyArrowUp {
			t.activeRow = max(0, t.activeRow-1)
			go print(t)
		} else if key == keyboard.KeyArrowLeft {
			t.activeCol = max(0, t.activeCol-1)
			go print(t)
		} else if key == keyboard.KeyArrowRight {
			t.activeCol = min(len(t.cols)-1, t.activeCol+1)
			go print(t)
		} else if key == keyboard.KeySpace {
			col := t.rows[t.activeRow][t.activeCol]
			if col.Disabled {
				continue
			}
			if !t.Multiple {
				t.selected = [][]int{[]int{t.activeRow, t.activeCol}}
			} else if !t.isSelected(t.activeRow, t.activeCol) {
				t.selected = append(t.selected, []int{t.activeRow, t.activeCol})
			} else {
				s := sliceIndex(len(t.selected), func(i int) bool { return t.selected[i][0] == t.activeRow && t.selected[i][1] == t.activeCol })
				t.selected = append(t.selected[:s], t.selected[s+1:]...)
			}
			go print(t)
		} else if key == keyboard.KeyEnter {
			break
		} else if key == keyboard.KeyEsc {
			t.selected = [][]int{}
			break
		}
	}

}

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
