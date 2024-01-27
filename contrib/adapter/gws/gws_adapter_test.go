package gws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lxzan/gws"
	"github.com/lxzan/xray"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var _port = int64(10086)

func nextAddr() string {
	port := atomic.AddInt64(&_port, 1)
	return fmt.Sprintf("127.0.0.1:%d", port)
}

type websocketHandler struct {
	onMessage func(socket *gws.Conn, message *gws.Message)
	onPing    func(socket *gws.Conn, payload []byte)
	onPong    func(socket *gws.Conn, payload []byte)
	onClose   func(socket *gws.Conn, err error)
	onOpen    func(socket *gws.Conn)
}

func (c *websocketHandler) OnOpen(socket *gws.Conn) {
	if c.onOpen != nil {
		c.onOpen(socket)
	}
}

func (c *websocketHandler) OnClose(socket *gws.Conn, err error) {
	if c.onClose != nil {
		c.onClose(socket, err)
	}
}

func (c *websocketHandler) OnPing(socket *gws.Conn, payload []byte) {
	if c.onPing != nil {
		c.onPing(socket, payload)
	}
}

func (c *websocketHandler) OnPong(socket *gws.Conn, payload []byte) {
	if c.onPong != nil {
		c.onPong(socket, payload)
	}
}

func (c *websocketHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	if c.onMessage != nil {
		c.onMessage(socket, message)
	}
}

func TestNewAdapter(t *testing.T) {
	var as = assert.New(t)

	t.Run("normal", func(t *testing.T) {
		t.Run("encode", func(t *testing.T) {
			addr := nextAddr()
			go func() {
				gws.NewServer(&gws.BuiltinEventHandler{}, nil).Run(addr)
			}()
			time.Sleep(100 * time.Millisecond)

			var router = xray.New()
			var wg = &sync.WaitGroup{}
			wg.Add(1)
			router.On("greet", func(ctx *xray.Context) {
				wg.Done()
			})

			var adapter = NewAdapter(router).SetHeaderCodec(xray.TextMapHeader)

			socket, _, err := gws.NewClient(&gws.BuiltinEventHandler{}, &gws.ClientOption{
				Addr: "ws://" + addr,
			})
			as.NoError(err)

			adapter.ServeWebSocket(socket, &gws.Message{
				Opcode: gws.OpcodeText,
				Data:   bytes.NewBufferString(`0033{"X-Path":"/greet","X-Action":""}{"hello":"world!"}`),
			})
			wg.Wait()
		})

		t.Run("decode", func(t *testing.T) {
			var wg = &sync.WaitGroup{}
			wg.Add(2)

			addr := nextAddr()
			go func() {
				handler := &websocketHandler{}
				handler.onOpen = func(socket *gws.Conn) {
					socket.WriteString(`0033{"X-Path":"/greet","X-Action":""}{"hello":"world!"}`)
				}
				handler.onMessage = func(socket *gws.Conn, message *gws.Message) {
					var req = make(map[string]any)
					message.Data.Next(4)
					json.NewDecoder(message).Decode(&req)
					as.Equal(req["ack"], "1")
					wg.Done()
				}
				gws.NewServer(handler, nil).Run(addr)
			}()
			time.Sleep(100 * time.Millisecond)

			var router = xray.New()
			var adapter = NewAdapter(router).SetHeaderCodec(xray.TextMapHeader)
			clientHandler := &websocketHandler{}
			clientHandler.onMessage = func(socket *gws.Conn, message *gws.Message) {
				adapter.ServeWebSocket(socket, message)
			}
			socket, _, err := gws.NewClient(clientHandler, &gws.ClientOption{
				Addr: "ws://" + addr,
			})
			as.NoError(err)
			go func() {
				socket.ReadLoop()
			}()

			router.On("greet", func(ctx *xray.Context) {
				var req = make(map[string]any)
				ctx.BindJSON(&req)
				as.Equal(req["hello"], "world!")
				ctx.WriteJSON(int(gws.OpcodeText), xray.Any{"ack": "1"})
				wg.Done()
			})

			wg.Wait()
		})
	})
}

func TestOthers(t *testing.T) {
	var w = newWriterPool().Get()
	w.Reset(&gws.Conn{}, gws.OpcodeText, xray.TextMapHeader)
	assert.Equal(t, xray.ProtocolWebSocket, w.Protocol())
	w.RawResponseWriter()
}
