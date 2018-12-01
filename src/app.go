// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"syscall"
	"time"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Foo struct {
	Index		int
	Name		string
	Phone		float64
	Create		time.Time
	Update		time.Time
	Addr			string
	Message	string
	checked bool
}

type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo
}

func NewFooModel() *FooModel {
	m := new(FooModel)
	m.ResetRows()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *FooModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Bar

	case 2:
		return item.Baz

	case 3:
		return item.Quux
	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *FooModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FooModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

// Called by the TableView to sort the model.
func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.Bar < b.Bar)

		case 2:
			return c(a.Baz < b.Baz)

		case 3:
			return c(a.Quux.Before(b.Quux))
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *FooModel) ResetRows() {
	// Create some random data.
	m.items = make([]*Foo, rand.Intn(50000))

	now := time.Now()

	for i := range m.items {
		m.items[i] = &Foo{
			Index: i,
			Bar:   strings.Repeat("*", rand.Intn(5)+1),
			Baz:   rand.Float64() * 1000,
			Quux:  time.Unix(rand.Int63n(now.Unix()), 0),
		}
	}

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

var (
	user32           = syscall.NewLazyDLL("User32.dll")
	getSystemMetrics = user32.NewProc("GetSystemMetrics")
)

func GetSystemMetrics(nIndex int) int {
	index := uintptr(nIndex)
	ret, _, _ := getSystemMetrics.Call(index)
	return int(ret)
}

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)


func main() {
	rand.Seed(time.Now().UnixNano())

	boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)

	labelFont:=Font{Family:"Microsoft YaHei UI",PointSize:9}

 	goodIcon, _ := walk.Resources.Icon("img/check.ico")
	badIcon, _ := walk.Resources.Icon("img/stop.ico")

	barBitmap, err := walk.NewBitmap(walk.Size{100, 1})
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	canvas, err := walk.NewCanvasFromImage(barBitmap)
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	canvas.GradientFillRectangle(walk.RGB(255, 0, 0), walk.RGB(0, 255, 0), walk.Horizontal, walk.Rectangle{0, 0, 100, 1})

	canvas.Dispose()

	model := NewFooModel()

	var tv *walk.TableView

	var queryPB,caculPB *walk.PushButton

	MainWindow{
		Size:   Size{GetSystemMetrics(SM_CXSCREEN)*80/100, GetSystemMetrics(SM_CYSCREEN)*80/100},
		Layout: VBox{},
		Background: SystemColorBrush{walk.SysColorWindow},
		Icon: goodIcon,
		Title: "就诊记录",
		//OnSizeChanged: func() {
		//	fmt.Println("OnSizeChanged")
		//},
		OnMouseUp: func(x, y int, button walk.MouseButton) {
			fmt.Println("OnMouseUp")
		},
		OnMouseDown: func(x, y int, button walk.MouseButton) {
			fmt.Println("OnMouseDown")
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns:10},
				Children:[]Widget{
					Label{
						Text: "姓名:",
						Font: labelFont,
						MaxSize: Size{Width:30},
						MinSize: Size{Width:30},
					},
					LineEdit{
						Text: Bind("Name"),
						MaxSize: Size{Width:200,Height:20},
						MinSize: Size{Width:100,Height:20},
					},
					Label{
						Text: "号码:",
						Font: labelFont,
						MaxSize: Size{Width:30},
						MinSize: Size{Width:30},
					},
					LineEdit{
						Text: Bind("Phone"),
						MaxSize: Size{Width:200,Height:20},
						MinSize: Size{Width:100,Height:20},
					},

					Label{
						Text: "日期:",
						Font:labelFont,
						MaxSize: Size{Width:30},
						MinSize: Size{Width:30},
					},
					DateEdit{
						Date: Bind("Start"),
						MaxSize: Size{Width:200,Height:20},
						MinSize: Size{Width:100,Height:20},
					},
					Label{
						Text: "-",
						MaxSize: Size{Width:5},
						MinSize: Size{Width:5},
					},
					DateEdit{
						Date: Bind("End"),
						MaxSize: Size{Width:200,Height:20},
						MinSize: Size{Width:100,Height:20},
					},
					PushButton{
						AssignTo: &queryPB,
						Text:     "查询",
						Font:labelFont,
						MaxSize: Size{Width:60,Height:20},
						MinSize: Size{Width:60,Height:20},
						OnClicked: func() {
						},
					},
					PushButton{
						AssignTo: &caculPB,
						Text:     "计算",
						Font:labelFont,
						MaxSize: Size{Width:60,Height:20},
						MinSize: Size{Width:60,Height:20},
						OnClicked: func() {
						},
					},
				},
			},

			TableView{
				AssignTo:              &tv,
				AlternatingRowBGColor: walk.RGB(239, 239, 239),
				CheckBoxes:            true,
				ColumnsOrderable:      true,
				MultiSelection:        true,
				Columns: []TableViewColumn{
					{Title: "#"},
					{Title: "姓名"},
					{Title: "电话", Alignment: AlignFar},
					{Format: "2006-01-02 15:04:05", Title: "登记时间"},
					{Format: "2006-01-02 15:04:05", Title: "最新时间"},
					{Title: "就诊费用"},
					{Title: "诊断信息",},
					{Title: "住址", Width: 150},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := model.items[style.Row()]

					if item.checked {
						if style.Row()%2 == 0 {
							style.BackgroundColor = walk.RGB(159, 215, 255)
						} else {
							style.BackgroundColor = walk.RGB(143, 199, 239)
						}
					}

					switch style.Col() {
					case 1:
						if canvas := style.Canvas(); canvas != nil {
							bounds := style.Bounds()
							bounds.X += 2
							bounds.Y += 2
							bounds.Width = int((float64(bounds.Width) - 4) / 5 * float64(len(item.Bar)))
							bounds.Height -= 4
							canvas.DrawBitmapPartWithOpacity(barBitmap, bounds, walk.Rectangle{0, 0, 100 / 5 * len(item.Bar), 1}, 127)

							bounds.X += 4
							bounds.Y += 2
							canvas.DrawText(item.Bar, tv.Font(), 0, bounds, walk.TextLeft)
						}

					case 2:
						if item.Baz >= 900.0 {
							style.TextColor = walk.RGB(0, 191, 0)
							style.Image = goodIcon
						} else if item.Baz < 100.0 {
							style.TextColor = walk.RGB(255, 0, 0)
							style.Image = badIcon
						}

					case 3:
						if item.Quux.After(time.Now().Add(-365 * 24 * time.Hour)) {
							style.Font = boldFont
						}
					}
				},
				Model: model,
				OnSelectedIndexesChanged: func() {
					fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
				},
			},
		},
	}.Run()
}
