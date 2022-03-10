package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	selectable "github.com/mkg0/bouldering/internal/golang-selectable-table"

	"github.com/fatih/color"
	"github.com/golang-module/carbon/v2"
	tsize "github.com/kopoli/go-terminal-size"
)

func getCarbonFromUnix(unixMilliDate int64) carbon.Carbon {
	return carbon.Time2Carbon(time.UnixMilli(unixMilliDate))
}

func printTableOnResize(t *selectable.Table) {
	listener, _ := tsize.NewSizeListener()
	for {
		s := <-listener.Change
		t.Width = s.Width
		t.Height = s.Height
		t.RePrint()
	}
}

func askSlot(slots []Slot, start, end carbon.Carbon, isAutoMode bool, enableFuture bool) []Slot {
	s, _ := tsize.GetSize()
	t := selectable.Table{
		Width:         s.Width,
		Height:        s.Height,
		HoverColor:    color.New(color.FgHiWhite).Add(color.BgHiBlue),
		SelectedColor: color.New(color.FgMagenta).Add(color.BgYellow),
		NormalColor:   color.New(color.FgHiWhite),
		HeaderColor:   color.New(color.FgHiWhite).Add(color.Bold).Add(color.Underline),
		Multiple:      isAutoMode,
	}
	go printTableOnResize(&t)

	var table [][]Slot = [][]Slot{}

	colCount := int(start.DiffInDays(end)) + 1
	// define day cols
	for i := 0; i < colCount; i++ {
		t.DefineCol(strings.ToUpper(start.AddDays(i).Format("l-M j", carbon.Berlin)))
	}

	// add slots to rows
	for _, slot := range slots {
		colIndex := int(start.DiffInDays(getCarbonFromUnix(slot.DateList[0].Start).StartOfDay()))
		for len(table) < (colIndex + 1) {
			table = append(table, []Slot{})
		}
		table[colIndex] = append(table[colIndex], slot)
	}

	// sort slots
	for i, col := range table {
		sort.SliceStable(table[i], func(i, j int) bool {
			return col[i].DateList[0].Start < col[j].DateList[0].Start
		})

	}

	// add leading margin rows
	// longestRowCount := len(getLongestArr(table))
	// for i, col := range table {
	// 	if len(col) < longestRowCount {
	// 		table[i] = append(make([]Slot, longestRowCount-len(col)), col...)
	// 	}
	// }

	//add rows to table
	for colIndex := 0; colIndex < len(getLongestArr(table)); colIndex++ {
		var cols []selectable.Cell
		for rowIndex := 0; rowIndex < colCount; rowIndex++ {
			isOutRange := rowIndex > len(table)-1 || colIndex > len(table[rowIndex])-1
			if isOutRange || table[rowIndex][colIndex].DateList == nil {
				cols = append(cols, selectable.Cell{Content: "-", Disabled: true})
				continue
			}
			slot := table[rowIndex][colIndex]
			date := slot.DateList[0]
			if date.Start < time.Now().UnixMilli() {
				cols = append(cols, selectable.Cell{Content: "too late", Disabled: true})
				continue
			}
			free := slot.MaxCourseParticipantCount - slot.CurrentCourseParticipantCount
			time := fmt.Sprintf("%s (%v)", getCarbonFromUnix(date.Start).Format("H:i", carbon.Berlin), free)
			disabled := false
			if slot.State != "BOOKABLE" {
				disabled = true
			}

			if isAutoMode {
				cols = append(cols, selectable.Cell{Content: time, Disabled: !disabled})
			} else {
				if enableFuture && slot.State == "NOT_YET_BOOKABLE" {
					disabled = false
				}
				cols = append(cols, selectable.Cell{Content: time, Disabled: disabled})
			}
		}
		t.AddRow(cols)
	}

	res := t.Run()
	var slotsToBook []Slot
	for _, v := range res {
		slotsToBook = append(slotsToBook, table[v[1]][v[0]])
	}
	return slotsToBook
}
