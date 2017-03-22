package dns

import (
	"github.com/miekg/dns"
	"net"
)

type QueryContext interface{}
type Handler interface {}

type RawHandler func(ctx QueryContext, m []byte) error
type MsgHandler   func(ctx QueryContext, msg *dns.Msg) error

type sessionWriter struct {
	ctx        QueryContext
	handler    Handler
}

// WriteMsg implements the ResponseWriter.WriteMsg method.
func (w *sessionWriter) WriteMsg(m *dns.Msg) (err error) {
	h1, ok1 := w.handler.(RawHandler)
	h2, ok2 := w.handler.(MsgHandler)

	if ok1 {
		var data []byte
		data, err = m.Pack()
		if err != nil {
			return err
		}
		return h1(w.ctx, data)
	}

	if ok2 {
		return h2(w.ctx, m)
	}

	panic("SessionWriter must initial a callback")
	return nil
}

// Write implements the ResponseWriter.Write method.
func (w *sessionWriter) Write(data []byte) (int, error) {
	length := len(data)

	if h, ok := w.handler.(RawHandler); ok  {
		return length, h(w.ctx, data)
	}

	if h, ok := w.handler.(MsgHandler); ok {
		r := new(dns.Msg)
		err := r.Unpack(data)
		if err != nil {
			return 0, err
		}

		return length, h(w.ctx, r)
	}

	panic("SessionWriter must initial a callback")
	return 0, nil

}

// LocalAddr implements the ResponseWriter.LocalAddr method.
func (w *sessionWriter) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr implements the ResponseWriter.RemoteAddr method.
func (w *sessionWriter) RemoteAddr() net.Addr {
	return nil
}

// TsigStatus implements the ResponseWriter.TsigStatus method.
func (w *sessionWriter) TsigStatus() error {
	return nil
}

// TsigTimersOnly implements the ResponseWriter.TsigTimersOnly method.
func (w *sessionWriter) TsigTimersOnly(b bool) {}

// Hijack implements the ResponseWriter.Hijack method.
func (w *sessionWriter) Hijack() {}

// Close implements the ResponseWriter.Close method
func (w *sessionWriter) Close() error {
	return nil
}
