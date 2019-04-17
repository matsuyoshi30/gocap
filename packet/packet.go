package packet

import (
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device       string = "en0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 300 * time.Second
	handle       *pcap.Handle
)

type PacketData struct {
	EData EthernetData
	IData IpData
	TData TcpData
}

type EthernetData struct {
	EType  layers.EthernetType
	SrcMAC net.HardwareAddr
	DstMAC net.HardwareAddr
}

type IpData struct {
	Version uint8
	Ihl     uint8  // Internet Header Length
	Tos     uint8  // type of service
	Length  uint16 // total length
	Id      uint16 // identification
	// Flags      IPv4Flag
	FragOffset uint16 // fragment offset
	Ttl        uint8  // time to live
	Protocol   layers.IPProtocol
	Checksum   uint16
	SrcIP      net.IP
	DstIP      net.IP
}

type TcpData struct {
	SrcIP      net.IP
	DstIP      net.IP
	SrcPort    layers.TCPPort
	DstPort    layers.TCPPort
	Seq        uint32
	Ack        uint32
	DataOffset uint8
	Window     uint16
	Checksum   uint16
	Urgent     uint16
}

func GetPacket() []*PacketData {
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// var filter string = "tcp"
	// err = handle.SetBPFFilter(filter)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	PacketDatas := make([]*PacketData, 0)
	i := 0
	for packet := range packetSource.Packets() {
		pd := getPacketInfo(packet)
		PacketDatas = append(PacketDatas, pd)

		i++
		if i == 10 {
			break
		}
	}

	return PacketDatas
}

func getPacketInfo(packet gopacket.Packet) *PacketData {
	pd := &PacketData{}

	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetData, _ := ethernetLayer.(*layers.Ethernet)
		pd.EData.EType = ethernetData.EthernetType
		pd.EData.SrcMAC = ethernetData.SrcMAC
		pd.EData.DstMAC = ethernetData.DstMAC
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ipData, _ := ipLayer.(*layers.IPv4)
		pd.IData.Version = ipData.Version
		pd.IData.Ihl = ipData.IHL
		pd.IData.Tos = ipData.TOS
		pd.IData.Length = ipData.Length
		pd.IData.Id = ipData.Id
		// pd.IData.Flags = ipData.Flags
		pd.IData.FragOffset = ipData.FragOffset
		pd.IData.Ttl = ipData.TTL
		pd.IData.Checksum = ipData.Checksum
		pd.IData.Protocol = ipData.Protocol
		pd.IData.SrcIP = ipData.SrcIP
		pd.IData.DstIP = ipData.DstIP
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		ipData, _ := ipLayer.(*layers.IPv4)
		tcpData, _ := tcpLayer.(*layers.TCP)
		pd.TData.SrcIP = ipData.SrcIP
		pd.TData.DstIP = ipData.DstIP
		pd.TData.SrcPort = tcpData.SrcPort
		pd.TData.DstPort = tcpData.DstPort
		pd.TData.Seq = tcpData.Seq
		pd.TData.Ack = tcpData.Ack
		pd.TData.DataOffset = tcpData.DataOffset
		pd.TData.Window = tcpData.Window
		pd.TData.Checksum = tcpData.Checksum
		pd.TData.Urgent = tcpData.Urgent
	}

	return pd
}
