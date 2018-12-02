// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	. "medic/src/data"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn	int
	sortOrder	walk.SortOrder
	search		*Search
	items		[]*Foo
	sItems		[]*Foo
	sum 		float64
	sSum		float64
	lSum		float64
	SumLabel	*walk.Label
	SSumLabel	*walk.Label
	LSumLabel	*walk.Label
}

func NewFooModel() *FooModel {
	m := new(FooModel)
	m.sortColumn=3
	m.sortOrder = 0
	m.sItems = Read()
	m.items = m.sItems[0:]
	copy(m.items,m.sItems)
	m.search = new(Search)
	m.SumLabel = new(walk.Label)
	m.SSumLabel = new(walk.Label)
	m.LSumLabel = new(walk.Label)
	for _,item := range m.items{
		m.sum=m.sum+item.PaidFee
		if item.Create.Month() == time.Now().Month(){
			m.sSum = m.sSum+item.PaidFee
		}
		if item.RealFee>item.PaidFee{
			m.lSum = m.lSum+item.RealFee-item.PaidFee
		}
	}
	m.search.Start, _ = time.Parse("2006-01-02 15:04:05","2018-12-01 00:00:00")
	m.search.End= time.Now().Add(time.Hour*24)
	m.ResetRows()
	return m
}

func (m *FooModel) GetSearch() *Search {
	return m.search
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *FooModel) RowCount() int {
	return len(m.sItems)
}

func (m *FooModel) Head(foo *Foo) {
	m.items=append(append([]*Foo{},foo), m.items...)
	m.sItems=append(append([]*Foo{},foo), m.sItems...)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *FooModel) Value(row, col int) interface{} {
	item := m.sItems[row]
	switch col {
	case 0:
		return ""
	case 1:
		return item.Name
	case 2:
		return item.Phone
	case 3:
		return item.Sex
	case 4:
		return item.Age
	case 5:
		return item.AllFee
	case 6:
		return item.RealFee
	case 7:
		return item.PaidFee
	case 8:
		return item.Create
	case 9:
		return item.Update
	case 10:
		return item.Diagnosed
	case 11:
		return item.Program
	case 12:
		return item.Address
	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *FooModel) Checked(row int) bool {
	return m.sItems[row].Checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FooModel) SetChecked(row int, checked bool) error {
	m.sItems[row].Checked = checked
	return nil
}

// Called by the TableView to sort the model.
func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.sItems, func(i, j int) bool {
		a, b := m.sItems[i], m.sItems[j]

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
			return c(a.Name < b.Name)
		case 2:
			return c(a.Phone < b.Phone)
		case 3:
			return c(a.Sex < b.Sex)
		case 4:
			return c(a.Age < b.Age)
		case 5:
			return c(a.AllFee < b.AllFee)
		case 6:
			return c(a.RealFee < b.RealFee)
		case 7:
			return c(a.PaidFee < b.PaidFee)
		case 8:
			return c(a.Create.Before(b.Create))
		case 9:
			return c(a.Update.Before(b.Update))
		case 10:
			return c(a.Diagnosed < b.Diagnosed)
		case 11:
			return c(a.Program < b.Program)
		case 12:
			return c(a.Address < b.Address)
		}
		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *FooModel) ResetRows() {
	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
	_ = m.Sort(m.sortColumn, m.sortOrder)
}

func (m *FooModel) Search() {
	sItems := []*Foo{}
	start:=m.search.Start.Format("2006-01-02 15:04:05")
	end:=m.search.End.Format("2006-01-02 15:04:05")
	for _,item:=range m.items{
		create:=item.Create.Format("2006-01-02 15:04:05")
		if strings.Contains(item.Name,m.search.Name) && create>start && create<end && strings.Contains(item.Phone,m.search.Phone) {
			sItems = append(sItems,item)
		}
	}
	m.sItems = sItems
	m.ResetRows()
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

var model = NewFooModel()

type Search struct {
	Name 	string
	Phone 	string
	Start 	time.Time
	End		time.Time
}

func main() {
	fmt.Println("hello world!")
	walk.FocusEffect, _ = walk.NewBorderGlowEffect(walk.RGB(0, 63, 255))
	walk.InteractionEffect, _ = walk.NewDropShadowEffect(walk.RGB(63, 63, 63))
	walk.ValidationErrorEffect, _ = walk.NewBorderGlowEffect(walk.RGB(255, 0, 0))

	with:=GetSystemMetrics(SM_CXSCREEN)
	height:=GetSystemMetrics(SM_CYSCREEN)
	//boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)

	inputWidth:=(with*80/100-30*3-60*2-150)/4

	labelFont:=Font{Family:"Microsoft YaHei UI",PointSize:9}

 	goodIcon, _ := walk.Resources.Icon("img/check.ico")
	//badIcon, _ := walk.Resources.Icon("img/stop.ico")

	barBitmap, err := walk.NewBitmap(walk.Size{Width: 100, Height: 1})
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	canvas, err := walk.NewCanvasFromImage(barBitmap)
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	_ = canvas.GradientFillRectangle(walk.RGB(255, 0, 0), walk.RGB(0, 255, 0), walk.Horizontal, walk.Rectangle{0, 0, 100, 1})

	canvas.Dispose()

	var tv *walk.TableView
	var db *walk.DataBinder
	var queryPB,addPB *walk.PushButton
	var mw *walk.MainWindow
	_, _ = MainWindow{
		AssignTo:   &mw,
		Size:       Size{Width: with * 90 / 100, Height: height - 150},
		Layout:     VBox{},
		Background: SystemColorBrush{Color: walk.SysColorWindow},
		Icon:       goodIcon,
		Title:      "就诊记录",
		Children: []Widget{
			Composite{
				DataBinder: DataBinder{
					AssignTo:       &db,
					Name:           "search",
					DataSource:     model.GetSearch(),
					ErrorPresenter: ToolTipErrorPresenter{},
				},
				Layout:  Grid{Columns: 10},
				MaxSize: Size{Width: with * 80 / 100, Height: 40},
				MinSize: Size{Width: with * 80 / 100, Height: 40},
				Children: []Widget{
					Label{
						Text:    "姓名:",
						Font:    labelFont,
						MaxSize: Size{Width: 30},
						MinSize: Size{Width: 30},
					},
					LineEdit{
						Text:    Bind("Name"),
						MaxSize: Size{Width: inputWidth, Height: 20},
						MinSize: Size{Width: inputWidth, Height: 20},
					},
					Label{
						Text:    "号码:",
						Font:    labelFont,
						MaxSize: Size{Width: 30},
						MinSize: Size{Width: 30},
					},
					LineEdit{
						Text:    Bind("Phone"),
						MaxSize: Size{Width: inputWidth, Height: 20},
						MinSize: Size{Width: inputWidth, Height: 20},
					},

					Label{
						Text:    "日期:",
						Font:    labelFont,
						MaxSize: Size{Width: 30},
						MinSize: Size{Width: 30},
					},
					DateEdit{
						Date:    Bind("Start"),
						MaxSize: Size{Width: inputWidth, Height: 20},
						MinSize: Size{Width: inputWidth, Height: 20},
					},
					Label{
						Text:    "-",
						MaxSize: Size{Width: 5},
						MinSize: Size{Width: 5},
					},
					DateEdit{
						Date:    Bind("End"),
						MaxSize: Size{Width: inputWidth, Height: 20},
						MinSize: Size{Width: inputWidth, Height: 20},
					},
					PushButton{
						AssignTo: &queryPB,
						Text:     "查询",
						Font:     labelFont,
						MaxSize:  Size{Width: 60, Height: 20},
						MinSize:  Size{Width: 60, Height: 20},
						OnClicked: func() {
							if err := db.Submit(); err == nil {
								model.Search()
							}
						},
					},
					PushButton{
						AssignTo: &addPB,
						Text:     "登记",
						Font:     labelFont,
						MaxSize:  Size{Width: 60, Height: 20},
						MinSize:  Size{Width: 60, Height: 20},
						OnClicked: func() {
							if cmd, err := AddDialog(mw, nil); err == nil && cmd == walk.DlgCmdOK {
								model.ResetRows()
							}
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
				//MinSize:Size{Width:with*75/100,Height:height-300},
				//MaxSize:Size{Width:with,Height:height-100},
				Columns: []TableViewColumn{
					{Title: "#", Width: 22},
					{Title: "姓名", Alignment: AlignNear, Width: 50},
					{Title: "电话", Alignment: AlignNear, Width: 100},
					{Title: "性别", Alignment: AlignNear, Width: 40},
					{Title: "年龄", Alignment: AlignFar, Width: 40},
					{Title: "诊费", Alignment: AlignFar, Width: 80},
					{Title: "实收", Alignment: AlignFar, Width: 80},
					{Title: "已付", Alignment: AlignFar, Width: 80},
					{Title: "登记时间", Alignment: AlignCenter, Format: "2006-01-02", Width: 110},
					{Title: "最新时间", Alignment: AlignCenter, Format: "2006-01-02", Width: 110},
					{Title: "病理诊断", Alignment: AlignNear, Width: 130},
					{Title: "治疗方案", Alignment: AlignNear, Width: 130},
					{Title: "住址", Width: 130},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := model.sItems[style.Row()]
					if item.Checked {
						style.BackgroundColor = walk.RGB(143, 199, 239)
					}
					switch style.Col() {
					case 7:
						if item.PaidFee < item.RealFee {
							style.TextColor = walk.RGB(255, 0, 0)
						}
					}
				},
				Model: model,
				OnItemActivated: func() {
					if cmd, err := AddDialog(mw, model.sItems[tv.SelectedIndexes()[0]]); err == nil && cmd == walk.DlgCmdOK {
						model.ResetRows()
					}
				},
			},
			Composite{
				Layout: Grid{Columns: 13},
				Children: []Widget{
					Label{
						Text:    "合计收入:",
						Font:    labelFont,
						MaxSize: Size{Width: 50},
						MinSize: Size{Width: 50},
					},
					Label{
						Text:     strconv.FormatFloat(model.sum, 'f', 1, 64) + " 元",
						AssignTo: &model.SumLabel,
						Font:     labelFont,
						MaxSize:  Size{Width: 100},
						MinSize:  Size{Width: 100},
					},
					Label{
						Text:    "当月总计:",
						Font:    labelFont,
						MaxSize: Size{Width: 50},
						MinSize: Size{Width: 50},
					},
					Label{
						Text:     strconv.FormatFloat(model.sSum, 'f', 1, 64) + " 元",
						AssignTo: &model.LSumLabel,
						Font:     labelFont,
						MaxSize:  Size{Width: 100},
						MinSize:  Size{Width: 100},
					},
					Label{
						Text:    "差额总计:",
						Font:    labelFont,
						MaxSize: Size{Width: 50},
						MinSize: Size{Width: 50},
					},
					Label{
						Text:     strconv.FormatFloat(model.lSum, 'f', 1, 64) + " 元",
						AssignTo: &model.SSumLabel,
						Font:     labelFont,
						MaxSize:  Size{Width: 100},
						MinSize:  Size{Width: 100},
					},
					/*PushButton{
						Text: "首页",
						OnClicked: func() {

						},
					},
					PushButton{
						Text: "上一页",
						OnClicked: func() {
						},
					},
					LineEdit{
						Text:    Bind("CurPage"),
						MaxSize: Size{Width: inputWidth, Height: 20},
						MinSize: Size{Width: inputWidth, Height: 20},
					},
					PushButton{
						Text: "跳转",
						OnClicked: func() {

						},
					},
					PushButton{
						Text: "下一页",
						OnClicked: func() {

						},
					},
					PushButton{
						Text: "末页",
						OnClicked: func() {

						},
					},*/
				},
			},
		},
	}.Run()
}

const (
	SexWoman Sex = "女"
	SexMan Sex = "男"
)

func AddDialog(owner walk.Form,foo *Foo) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton
	addIcon, _ := walk.Resources.Icon("img/plus.png")
	addFlag := foo==nil
	if addFlag{
		foo = new(Foo)
		foo.Sex=SexMan
	}
	return Dialog{
		AssignTo:      &dlg,
		Icon: addIcon,
		Background: SystemColorBrush{Color:walk.SysColorWindow},
		Title:         Bind("添加病例"),
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "foo",
			DataSource:     foo,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{Width:300},
		MaxSize: Size{Width:400},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "姓名:",
					},
					LineEdit{
						Text: Bind("Name"),
					},

					Label{
						Text: "联系电话:",
					},
					LineEdit{
						Text: Bind("Phone"),
					},

					RadioButtonGroupBox{
						ColumnSpan: 2,
						Title:      "性别:",
						Layout:     HBox{},
						DataMember: "Sex",
						Buttons: []RadioButton{
							{Text: "男", Value: SexMan},
							{Text: "女", Value: SexWoman},
						},
					},

					Label{
						Text: "年龄:",
					},
					NumberEdit{
						Value:    Bind("Age"),
						Decimals: 0,
					},

					Label{
						Text: "病人住址:",
					},
					LineEdit{
						Text: Bind("Address"),
					},

					Label{
						Text: "病因诊断:",
					},
					TextEdit{
						ColumnSpan: 2,
						MinSize:    Size{Width: 100, Height: 80},
						Text:       Bind("Diagnosed"),
					},

					Label{
						Text: "治疗方案:",
					},
					TextEdit{
						ColumnSpan: 2,
						MinSize:    Size{Width: 100, Height: 80},
						Text:       Bind("Program"),
					},

					Label{
						Text: "就诊费用:",
					},
					NumberEdit{
						Value:    Bind("AllFee"),
						Suffix:   "$",
						Decimals: 1,
					},

					Label{
						Text: "实收费用:",
					},
					NumberEdit{
						Value:    Bind("RealFee"),
						Suffix:   "$",
						Decimals: 1,
					},

					Label{
						Text: "已付费用:",
					},
					NumberEdit{
						Value:    Bind("PaidFee"),
						Suffix:   "$",
						Decimals: 1,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "保存",
						OnClicked: func() {
							if err := db.Submit(); err == nil {
								foo.Update = time.Now()
								if addFlag {
									foo.Create = foo.Update
									model.Head(foo)
								}
								Write(model.items)
							}
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
