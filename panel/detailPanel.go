package panel

import (
	"fmt"
	"net"
	"time"

	"gocap/packet"

	"github.com/google/gopacket/layers"
	"github.com/jroimartin/gocui"
)

type Detail struct {
	*Gui
	Position
	name           string
	timestamp      time.Time
	capturedLength int
	actualLength   int
	ethernetInfo   ETHERNETInfo
	ipInfo         IPInfo
	tcpInfo        TCPInfo
}

type ETHERNETInfo struct {
	etype  layers.EthernetType
	srcMAC net.HardwareAddr
	dstMAC net.HardwareAddr
}

type IPInfo struct {
	version uint8
	ihl     uint8  // Internet Header Length
	tos     uint8  // type of service
	length  uint16 // total length
	id      uint16 // identification
	// Flags      IPv4Flag
	fragOffset uint16 // fragment offset
	ttl        uint8  // time to live
	protocol   string
	checksum   uint16
	srcIP      net.IP
	dstIP      net.IP
}

type TCPInfo struct {
	srcPort, dstPort layers.TCPPort
	seq              uint32 // Sequence number
	ack              uint32 // Acknowledgement number
	dataOffset       uint8
	flag             string // FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
	windowsize       uint16
	checksum         uint16
	urgent           uint16
}

func (d *Detail) SetView(g *gocui.Gui) error {
	if v, err := g.SetView("detail", d.x, d.y, d.w, d.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Packet Detail"
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
	}

	return nil
}

func ShowDetail(g *gocui.Gui, pd *packet.PacketData) {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View("detail")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintf(v, "%+v\n", pd.TS)
		fmt.Fprintf(v, "[SRC] %s:%s\n", pd.TData.SrcIP, pd.TData.SrcPort)
		fmt.Fprintf(v, "[DST] %s:%s\n", pd.TData.DstIP, pd.TData.DstPort)
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
