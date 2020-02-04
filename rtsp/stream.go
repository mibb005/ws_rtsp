package rtsp

import (
	"fmt"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtsp"
)

func StartStream(path string) {
	if ch == nil {
		ch = make(chan Packet)
	}
	session, err := rtsp.Dial(path)
	if err != nil {
		fmt.Println("err:%s\n", err)
	}
	// codec, err := session.Streams()
	// if err != nil {
	// 	fmt.Println("err:%s\n", err)
	// }
	var keyFrame av.Packet
	var start bool

	for {
		fmt.Println("3\n")
		pkt, err := session.ReadPacket()
		if err != nil {
			fmt.Println("err:%s\n", err)
			break
		}
		var tmp Packet
		tmp.d = pkt
		if pkt.IsKeyFrame {
			keyFrame = pkt
			start = true
		}
		if !start {
			continue
		}
		tmp.f = keyFrame
		ch <- tmp
	}
	fmt.Println("5\n")
	session.Close()
}
