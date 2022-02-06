package selectable

import (
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
)

var disabledColor = color.New(color.FgHiWhite, color.Faint)

// var disabledSelectedColor = color.New(color.FgWhite, color.Faint, color.Bold)
var disabledSelectedColor = color.New(color.FgWhite, color.Faint, color.Bold, color.BgHiBlack)

func print(t *Table) {
	colWidth := math.Round(float64(t.Width / len(t.cols)))
	var output string
	//header
	for _, label := range t.cols {
		output = output + t.HeaderColor.Sprint(label)
		output = output + t.HeaderColor.Sprint(strings.Repeat(" ", int(colWidth)-len(label)))
	}
	output = output + t.HeaderColor.Sprintln("")

	// content
	for i := 0; i < len(t.rows); i++ {
		for i2 := 0; i2 < len(t.cols); i2++ {
			cell := t.rows[i][i2]
			if cell.Content == "" {
				cell.Content = "-"
			}
			var marker = t.NormalColor
			isSelected := t.isSelected(i, i2)
			isHovered := t.activeRow == i && t.activeCol == i2
			if isHovered {
				marker = t.HoverColor
			}
			if isSelected {
				marker = t.SelectedColor
			}
			if cell.Disabled && isHovered {
				marker = disabledSelectedColor
			} else if cell.Disabled == true {
				marker = disabledColor
			}

			output = output + marker.Sprint(cell.Content)
			output = output + marker.Sprint(strings.Repeat(" ", int(colWidth)-len(cell.Content)))
		}
		output = output + t.NormalColor.Sprintln("")
	}

	// add empty rows to fill screen
	for i := 0; i < t.Height-2-len(t.rows); i++ {
		output = output + t.NormalColor.Sprintln(strings.Repeat(" ", int(colWidth)))
	}

	fmt.Print(output)
}
