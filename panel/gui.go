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
