package panel

import (
	"fmt"
	"strings"
	"time"

	"gocap/packet"

	"github.com/jroimartin/gocui"
)

type Detail struct {
	*Gui
	Position
	name           string
	timestamp      time.Time
	capturedLength int
	actualLength   int
}

func (d *Detail) SetView(g *gocui.Gui) error {
	if v, err := g.SetView(DetailPanel, d.x, d.y, d.w, d.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Packet Detail"
		v.Wrap = true
		v.Frame = true
		v.FgColor = gocui.ColorCyan

		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
	}
	d.SetKeyBindings()

	return nil
}

func ShowDetail(g *gocui.Gui, pd *packet.PacketData) {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View(DetailPanel)
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintf(v, "[TIMESTAMP] %s\n", pd.TS.Format("2006/1/2(Mon) 15:04:05"))
		fmt.Fprintf(v, "\n")

		fmt.Fprintf(v, "[[ETHERNET LAYER]]\n")
		fmt.Fprintf(v, " [Ethernet Type] %s\n", pd.EData.EType)
		fmt.Fprintf(v, " [Source]        %s\n", pd.EData.SrcMAC)
		fmt.Fprintf(v, " [Destination]   %s\n", pd.EData.DstMAC)
		fmt.Fprintf(v, "\n")

		fmt.Fprintf(v, "[[IP LAYER]]\n")
		fmt.Fprintf(v, " [Version]                %d\n", pd.IData.Version)
		fmt.Fprintf(v, " [Internet Header Length] %d\n", pd.IData.Ihl)
		fmt.Fprintf(v, " [Type Of Service]        %d\n", pd.IData.Tos)
		fmt.Fprintf(v, " [Length]                 %d\n", pd.IData.Length)
		fmt.Fprintf(v, " [Identification]         %d\n", pd.IData.Id)
		fmt.Fprintf(v, " [Fragment Offset]        %d\n", pd.IData.FragOffset)
		fmt.Fprintf(v, " [Time to Live]           %d\n", pd.IData.Ttl)
		fmt.Fprintf(v, " [Checksum]               %d\n", pd.IData.Checksum)
		fmt.Fprintf(v, " [Protocol]               %s\n", pd.IData.Protocol)
		fmt.Fprintf(v, " [Source]                 %s\n", pd.IData.SrcIP)
		fmt.Fprintf(v, " [Destination]            %s\n", pd.IData.DstIP)
		fmt.Fprintf(v, "\n")

		fmt.Fprintf(v, "[[TCP LAYER]]\n")
		fmt.Fprintf(v, " [Source]                 %s:%s\n", pd.TData.SrcIP,
			strings.Split(pd.TData.SrcPort.String(), "(")[0])
		fmt.Fprintf(v, " [Destination]            %s:%s\n", pd.TData.DstIP,
			strings.Split(pd.TData.DstPort.String(), "(")[0])
		fmt.Fprintf(v, " [Sequence Number]        %d\n", pd.TData.Seq)
		fmt.Fprintf(v, " [Acknowledgement Number] %d\n", pd.TData.Ack)
		fmt.Fprintf(v, " [Data Offset]            %d\n", pd.TData.DataOffset)
		fmt.Fprintf(v, " [Window size]            %d\n", pd.TData.Window)
		fmt.Fprintf(v, " [Checksum]               %d\n", pd.TData.Checksum)
		fmt.Fprintf(v, " [Urgent pointer]         %d\n", pd.TData.Urgent)

		return nil
	})
}

func NewDetail(gui *Gui, name string, x, y, w, h int) *Detail {
	return &Detail{
		Gui:      gui,
		name:     name,
		Position: Position{x, y, w, h},
	}
}

func (d *Detail) Name() string {
	return d.name
}

func (d *Detail) SetKeyBindings() {
	d.SetKeybindingsToPanel(d.name)
}
