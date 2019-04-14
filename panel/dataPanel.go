package panel

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Data struct {
	*Gui
	Position
	name string
}

func (d *Data) SetView(g *gocui.Gui) error {
	if v, err := g.SetView("data", d.x, d.y, d.w, d.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Packet Deta"
		fmt.Fprintf(v, "Print Packet Data")
	}

	return nil
}

func NewData(gui *Gui, name string, x, y, w, h int) *Data {
	return &Data{
		Gui:      gui,
		name:     name,
		Position: Position{x, y, w, h},
	}
}

func (d *Data) Name() string {
	return d.name
}
