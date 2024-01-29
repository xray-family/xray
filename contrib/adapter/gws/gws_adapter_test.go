package gws

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lxzan/gws"
	"github.com/lxzan/xray"
	"github.com/lxzan/xray/codec"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
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

			var adapter = NewAdapter(router)

			socket, _, err := gws.NewClient(&gws.BuiltinEventHandler{}, &gws.ClientOption{
				Addr: "ws://" + addr,
			})
			as.NoError(err)

			adapter.ServeWebSocket(socket, &gws.Message{
				Opcode: gws.OpcodeText,
				Data:   bytes.NewBufferString(`0021[["X-Path","/greet"]]{"hello":"world!"}`),
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
					socket.WriteString(`0021[["X-Path","/greet"]]{"hello":"world!"}`)
				}
				handler.onMessage = func(socket *gws.Conn, message *gws.Message) {
					var req = make(map[string]any)
					var header = &xray.SliceHeader{}
					unmarshalHeader(codec.StdJsonCodec, message.Data, header)
					json.NewDecoder(message).Decode(&req)
					as.Equal(req["ack"], "1")
					wg.Done()
				}
				gws.NewServer(handler, nil).Run(addr)
			}()
			time.Sleep(100 * time.Millisecond)

			var router = xray.New()
			var adapter = NewAdapter(router).SetHeaderTpl(&xray.SliceHeader{})
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
				ctx.Writer.Header().Set("x-real-ip", "127.0.0.1")
				ctx.WriteJSON(int(gws.OpcodeText), xray.Any{"ack": "1"})
				wg.Done()
			})

			wg.Wait()
		})

		t.Run("decode error", func(t *testing.T) {
			addr := nextAddr()
			go func() {
				handler := &websocketHandler{}
				handler.onOpen = func(socket *gws.Conn) {
					socket.WriteString(`00xx`)
				}
				gws.NewServer(handler, nil).Run(addr)
			}()
			time.Sleep(100 * time.Millisecond)

			wg := &sync.WaitGroup{}
			wg.Add(1)
			adapter := NewAdapter(xray.New())
			clientHandler := &websocketHandler{}
			clientHandler.onMessage = func(socket *gws.Conn, message *gws.Message) {
				assert.Error(t, adapter.ServeWebSocket(socket, message))
				wg.Done()
			}
			socket, _, err := gws.NewClient(clientHandler, &gws.ClientOption{
				Addr: "ws://" + addr,
			})
			as.NoError(err)
			go func() {
				socket.ReadLoop()
			}()
			wg.Wait()
		})
	})

	t.Run("encode error", func(t *testing.T) {
		addr := nextAddr()
		go func() {
			handler := &websocketHandler{}
			gws.NewServer(handler, nil).Run(addr)
		}()
		time.Sleep(100 * time.Millisecond)

		clientHandler := &websocketHandler{}
		socket, _, err := gws.NewClient(clientHandler, &gws.ClientOption{
			Addr: "ws://" + addr,
		})
		as.NoError(err)

		writer := newWriterPool().Get()
		writer.codec = codec.StdJsonCodec
		writer.conn = socket
		writer.header = &mockHeader{xray.SliceHeader{}}
		writer.buf = bytes.NewBufferString("")
		writer.Header().Set(xray.XPath, "/")
		writer.Code(int(gws.OpcodeText))
		writer.Write([]byte("hello"))
		as.NotNil(writer.RawResponseWriter())
		as.Error(writer.Flush())
	})
}

type mockHeader struct{ xray.SliceHeader }

func (c *mockHeader) MarshalJSON() ([]byte, error) {
	return nil, errors.New("1")
}

func TestAdapter_Marshal(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var w = bytes.NewBufferString("")
		var header = &xray.SliceHeader{}
		header.Set("x-path", "/greet")
		header.Set("x-method", "post")
		var err = marshalHeader(codec.StdJsonCodec, w, header)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(w.String()), `0042[["x-path","/greet"],["x-method","post"]]`)
	})

	t.Run("", func(t *testing.T) {
		var w = bytes.NewBufferString("")
		var header = &xray.SliceHeader{}
		var err = marshalHeader(codec.StdJsonCodec, w, header)
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(w.String()), `0002[]`)
	})

	t.Run("", func(t *testing.T) {
		var w = bytes.NewBufferString("")
		var header = &mockHeader{xray.SliceHeader{}}
		header.Set("x-path", "/")
		var err = marshalHeader(codec.StdJsonCodec, w, header)
		assert.Error(t, err)
	})
}

func TestAdapter_Unmarshal(t *testing.T) {
	t.Run("", func(t *testing.T) {
		reader := bytes.NewBufferString(`0041[["x-path","/greet"],["x-method","post"]]`)
		header := &xray.SliceHeader{}
		err := unmarshalHeader(codec.StdJsonCodec, reader, header)
		assert.NoError(t, err)
	})

	t.Run("", func(t *testing.T) {
		reader := bytes.NewBufferString(`0048[["x-path","/greet"],["x-method","post"]]`)
		header := &xray.SliceHeader{}
		err := unmarshalHeader(codec.StdJsonCodec, reader, header)
		assert.True(t, errors.Is(err, io.ErrShortBuffer))
	})

	t.Run("", func(t *testing.T) {
		reader := bytes.NewBufferString(`00`)
		header := &xray.SliceHeader{}
		err := unmarshalHeader(codec.StdJsonCodec, reader, header)
		assert.True(t, errors.Is(err, io.ErrShortBuffer))
	})

	t.Run("", func(t *testing.T) {
		reader := bytes.NewBufferString(`00xx`)
		header := &xray.SliceHeader{}
		err := unmarshalHeader(codec.StdJsonCodec, reader, header)
		assert.Error(t, err)
	})
}
