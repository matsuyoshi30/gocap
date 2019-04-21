package panel

import (
	"fmt"
	"net"
	"strings"
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
		v.FgColor = gocui.AttrBold | gocui.ColorWhite
		v.Title = v.Name()

		fmt.Fprintf(v, "%-8s  %-21s    %-21s %-20s", "Protocol", "Src Address", "Dst Address", "Timestamp")
	}

	v, err := g.SetView(l.Name(), l.x, l.y+1, l.w, l.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Frame = false
		v.FgColor = gocui.ColorMagenta
		v.SelBgColor = gocui.ColorWhite
		v.SelFgColor = gocui.AttrBold | gocui.ColorBlack

		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
	}
	l.SetKeyBindings()
	go l.Monitoring(l.stop, l.Gui.Gui, v)

	return nil
}

func (l *List) Monitoring(stop chan int, g *gocui.Gui, v *gocui.View) {
	ticker := time.NewTicker(10 * time.Millisecond)

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

func (l *List) CloseView() {
	l.stop <- 0
	close(l.stop)
}

func NewList(gui *Gui, name string, x, y, w, h int) *List {
	return &List{
		Gui:      gui,
		name:     name,
		Position: Position{x, y, w, h},
		stop:     make(chan int, 1),
	}
}

func (l *List) Name() string {
	return l.name
}

func (l *List) GetPacketList(v *gocui.View) error {
	pds := packet.GetPacket()

	for _, pd := range pds {
		// tp := pd.EData.EType
		srcip := pd.TData.SrcIP
		srcport := strings.Split(pd.TData.SrcPort.String(), "(")[0]
		src := srcip.String() + ":" + srcport

		dstip := pd.TData.DstIP
		dstport := strings.Split(pd.TData.DstPort.String(), "(")[0]
		dst := dstip.String() + ":" + dstport

		ts := pd.TS.Format("2006/1/2 15:04:05")

		if srcip != nil && dstip != nil {
			l.Update(func(g *gocui.Gui) error {
				v, err := l.View(l.name)
				if err != nil {
					return err
				}
				fmt.Fprintf(v, "%-8s  %-21s %s %-21s %-20s", "TCP", src, ">>", dst, ts)
				fmt.Fprintf(v, "\n")
				return nil
			})
		}
	}

	return nil
}

func (l *List) SetKeyBindings() {
	l.SetKeybindingsToPanel(l.name)
}
