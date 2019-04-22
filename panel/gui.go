package panel

import (
	"log"

	"github.com/jroimartin/gocui"
)

const (
	ListHeaderPanel = "Packet List"
	ListPanel       = "Packet List scroll"
	DetailPanel     = "Detail scroll"
	DataPanel       = "Data scroll"
	NavigatePanel   = "Navigate"
)

type Gui struct {
	*gocui.Gui
	Panels     map[string]Panel
	PanelNames []string
}

type Panel interface {
	SetView(*gocui.Gui) error
	Name() string
}

type Position struct {
	x, y int
	w, h int
}

func New(mode gocui.OutputMode) *Gui {
	g, err := gocui.NewGui(mode)
	if err != nil {
		log.Panicln(err)
	}
	// defer g.Close()

	g.Cursor = true
	gui := &Gui{
		Gui:        g,
		Panels:     make(map[string]Panel),
		PanelNames: []string{},
	}

	gui.init()

	return gui
}

func (gui *Gui) init() {
	maxX, maxY := gui.Size()

	gui.storePanel(NewList(gui, ListPanel, 0, 0, maxX/2-1, maxY-1-3))
	gui.storePanel(NewDetail(gui, DetailPanel, maxX/2, 0, maxX-1, maxY/2-1))
	gui.storePanel(NewData(gui, DataPanel, maxX/2, maxY/2, maxX-1, maxY-1-3))
	gui.storePanel(NewNavigate(gui, NavigatePanel, 0, maxY-1-3, maxX-1, maxY-1))

	for _, panel := range gui.Panels {
		if err := panel.SetView(gui.Gui); err != nil {
			log.Panicln(err)
		}
	}

	gui.SwitchPanel(ListPanel)
}

func (gui *Gui) addPanelName(panel Panel) {
	name := panel.Name()
	gui.PanelNames = append(gui.PanelNames, name)
}

func (gui *Gui) storePanel(panel Panel) {
	gui.Panels[panel.Name()] = panel

	storeTarget := map[string]bool{
		ListPanel:     true,
		DetailPanel:   true,
		DataPanel:     true,
		NavigatePanel: true,
	}

	if storeTarget[panel.Name()] {
		gui.addPanelName(panel)
	}
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}

	return nil
}

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if _, err := v.Line(cy + 1); err != nil {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Gui) SwitchPanel(next string) *gocui.View {
	v := g.CurrentView()
	if v != nil {
		v.Highlight = false
	}

	v, err := SetCurrentPanel(g.Gui, next)
	if err != nil {
		panic(err)
	}

	return v
}

func SetCurrentPanel(g *gocui.Gui, name string) (*gocui.View, error) {
	v, err := g.SetCurrentView(name)

	if err != nil {
		return nil, err
	}

	v.Highlight = true

	return g.SetViewOnTop(name)
}

func SwitchToListPanel(g *gocui.Gui, v *gocui.View) error {
	_, err := SetCurrentPanel(g, ListPanel)

	return err
}

func SwitchToDetailPanel(g *gocui.Gui, v *gocui.View) error {
	_, err := SetCurrentPanel(g, DetailPanel)

	return err
}

func SwitchToDataPanel(g *gocui.Gui, v *gocui.View) error {
	_, err := SetCurrentPanel(g, DataPanel)

	return err
}

func (g *Gui) SetKeybindingsToPanel(panel string) {
	if err := g.SetKeybinding(panel, gocui.KeyCtrlN, gocui.ModNone, CursorDown); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding(panel, gocui.KeyArrowDown, gocui.ModNone, CursorDown); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding(panel, gocui.KeyCtrlP, gocui.ModNone, CursorUp); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding(panel, gocui.KeyArrowUp, gocui.ModNone, CursorUp); err != nil {
		panic(err)
	}

	if err := g.SetKeybinding(panel, 'l', gocui.ModNone, SwitchToListPanel); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding(panel, 'd', gocui.ModNone, SwitchToDetailPanel); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding(panel, 't', gocui.ModNone, SwitchToDataPanel); err != nil {
		panic(err)
	}
}
