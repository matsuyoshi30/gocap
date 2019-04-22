package panel

import (
	"encoding/hex"
	"fmt"

	"gocap/packet"

	"github.com/jroimartin/gocui"
)

type Data struct {
	*Gui
	Position
	name string
}

func (d *Data) SetView(g *gocui.Gui) error {
	if v, err := g.SetView(DataPanel, d.x, d.y, d.w, d.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Packet Deta"
		v.Wrap = true
		v.Frame = true
		v.FgColor = gocui.ColorBlue

		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
	}
	d.SetKeyBindings()

	return nil
}

func ShowData(g *gocui.Gui, pd *packet.PacketData) {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View(DataPanel)
		if err != nil {
			return err
		}
		v.Clear()

		fmt.Fprintf(v, "%s", hex.Dump(pd.Data))

		return nil
	})
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

func (d *Data) SetKeyBindings() {
	d.SetKeybindingsToPanel(d.name)
}
