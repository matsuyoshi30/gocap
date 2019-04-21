package panel

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Navigate struct {
	*Gui
	Position
	name string
}

func (n *Navigate) SetView(g *gocui.Gui) error {
	if v, err := g.SetView(n.Name(), n.x, n.y, n.w, n.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Frame = false
		v.FgColor = gocui.ColorYellow
		fmt.Fprintf(v, "%s\n", "Stop: Ctrl + Q")
	}

	return nil
}

func NewNavigate(gui *Gui, name string, x, y, w, h int) *Navigate {
	return &Navigate{
		Gui:      gui,
		name:     name,
		Position: Position{x, y, w, h},
	}
}

func (n *Navigate) Name() string {
	return n.name
}