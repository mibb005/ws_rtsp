package rtsp

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"github.com/deepch/vdk/format/rtsp"
	"golang.org/x/net/websocket"
)

// "net/http"
// //草拟吗 自己创建的目录 哈哈哈哈哈    还好我比较聪明  要不然 就完蛋了  麻痹
// "golang.org/x/net/websocket"
// "log"

type Command struct {
	Match string `json:"Match"`
	Path  string `json:"Path"`
}

type Error struct {
	Code  uint32
	Msg   string
	Where string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code = %d ; msg = %s ; where = %s", e.Code, e.Msg, e.Where)
}

func NewError(code int, msg string) *Error {
	// 获取代码位置, 代码就不贴了, 不是重点.
	pc, file, line, ok := runtime.Caller(1)
	pcName := runtime.FuncForPC(pc).Name() //获取函数名
	where := fmt.Sprintf("%v   %s   %d   %t   %s", pc, file, line, ok, pcName)
	return &Error{Code: uint32(code), Msg: msg, Where: where}
}

func StartHttp() {

	http.Handle("/shiming", websocket.Handler(Echo))
	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatal(err)
	}

}

func Echo(w *websocket.Conn) {
	// var start bool
	fmt.Println("start--- ")
	for {
		var c Command
		// if err := websocket.Message.Receive(w, &command); err != nil {
		// 	fmt.Println("不能够接受消息 error==", err)
		// 	break
		// }
		if err := websocket.JSON.Receive(w, &c); err != nil {
			fmt.Println("不能够接受消息 error==", err)
			break
		}
		fmt.Printf("recv:%#v\n", c)

		switch c.Match {
		case "OPEN":
			go func() {
				session, err := rtsp.Dial(c.Path)
				if err != nil {
					fmt.Println("err:%s\n", err)
					return
				}
				codec, err := session.Streams()
				if err != nil {
					e := NewError(500, err.Error())
					fmt.Println("err:%s\n", e.Error())
					websocket.JSON.Send(w, e)
					return
				}

				muxer := mp4f.NewMuxer(nil)
				muxer.WriteHeader(codec)
				meta, init := muxer.GetInit(codec)
				err = websocket.Message.Send(w, append([]byte{9}, meta...))
				if err != nil {
					return
				}
				err = websocket.Message.Send(w, init)
				if err != nil {
					return
				}

				var start bool

				for {
					pkt, err := session.ReadPacket()
					if err != nil {
						e := NewError(500, err.Error())
						fmt.Println("err:%s\n", e.Error())
						websocket.JSON.Send(w, e)
						break
					}

					if pkt.IsKeyFrame {
						start = true
					}
					if !start {
						continue
					}

					ready, buf, _ := muxer.WritePacket(pkt, false)
					if ready {
						w.SetWriteDeadline(time.Now().Add(10 * time.Second))
						err := websocket.Message.Send(w, buf)
						if err != nil {
							return
						}
					}
				}
				session.Close()
			}()
		case "CLOSE":
		default:
		}
	}
	fmt.Println("end--- ")
}
