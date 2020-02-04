package rtsp

import "github.com/deepch/vdk/av"

var ch chan Packet

type Packet struct {
	d av.Packet
	f av.Packet
}
