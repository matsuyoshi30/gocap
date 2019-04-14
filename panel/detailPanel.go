package panel

import (
	"fmt"
	"net"
	"time"

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
	arpInfo        ARPInfo
	ipInfo         IPInfo
	tcpInfo        TCPInfo
	udpInfo        UDPInfo
}

type ETHERNETInfo struct {
	etype  string
	srcMAC string
	dstMAC string
}

type ARPInfo struct {
	addrType        layers.LinkType
	protocol        layers.EthernetType
	hwAddressSize   uint8
	protAddressSize uint8
	operation       uint16
	srcHwAddress    string
	srcProtAddress  string
	dstHwAddress    string
	dstProtAddress  string
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

type UDPInfo struct {
	srcPort, dstPort layers.UDPPort
	length           uint16
	checksum         uint16
}

func (d *Detail) SetView(g *gocui.Gui) error {
	if v, err := g.SetView("detail", d.x, d.y, d.w, d.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Packet Detail"
		fmt.Fprintf(v, "Print Packet Detail Information")
	}

	return nil
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
