package panel

import (
	"fmt"
	"net"
	"time"

	"gocap/packet"

	"github.com/google/gopacket/layers"
	"github.com/jroimartin/gocui"
)

type List struct {
	*Gui
	Position
	name    string
	Packets []*Packet
	stop    chan int
}

type Packet struct {
	Type       layers.EthernetType
	SrcAddress net.IP
	DstAddress net.IP
}

func (l *List) SetView(g *gocui.Gui) error {
	if v, err := g.SetView(ListHeaderPanel, l.x, l.y, l.w, l.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Frame = true
		v.Title = v.Name()

		fmt.Fprintf(v, "Print Packet List Header")
	}

	v, err := g.SetView(l.Name(), l.x, l.y+1, l.w, l.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Frame = false

		go l.Monitoring(l.stop, l.Gui.Gui, v)
	}

	return nil
}

func (l *List) Monitoring(stop chan int, g *gocui.Gui, v *gocui.View) {
	ticker := time.NewTicker(1 * time.Second)

LOOP:
	for {
		select {
		case <-ticker.C:
			l.GetPacketList(v)
		case <-stop:
			ticker.Stop()
			break LOOP
		}
	}
}

func NewList(gui *Gui, name string, x, y, w, h int) *List {
	return &List{
		Gui:      gui,
		name:     name,
		Position: Position{x, y, w, h},
	}
}

func (l *List) Name() string {
	return l.name
}

func (l *List) GetPacketList(v *gocui.View) error {
	pds := packet.GetPacket()

	for _, pd := range pds {
		tp := pd.EData.EType
		src := pd.EData.SrcMAC
		dst := pd.EData.DstMAC

		l.Update(func(g *gocui.Gui) error {
			v, err := l.View(l.name)
			if err != nil {
				return err
			}
			fmt.Fprintf(v, "%s: %s -> %s", tp, src, dst)
			fmt.Fprintf(v, "\n")
			return nil
		})
	}

	return nil
}
