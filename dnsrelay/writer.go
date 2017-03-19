package dnsrelay

import (
	"github.com/miekg/dns"
	"net"
)

type QueryContext interface{}
type Handler interface {}

type WriteHandler func(ctx QueryContext, m []byte) error
type MsgHandler   func(ctx QueryContext, msg *dns.Msg) error

type sessionWriter struct {
	ctx        QueryContext
	handler    Handler
}

// WriteMsg implements the ResponseWriter.WriteMsg method.
func (w *sessionWriter) WriteMsg(m *dns.Msg) (err error) {
	if h, ok := w.handler.(WriteHandler); ok {
		var data []byte
		data, err = m.Pack()
		if err != nil {
			return err
		}
		return h(w.ctx, data)

	} else if h, ok := w.handler.(MsgHandler); ok{
		return h(w.ctx, m)
	}

	panic("SessionWriter must initial a callback")
	return nil
}

// Write implements the ResponseWriter.Write method.
func (w *sessionWriter) Write(data []byte) (int, error) {
	length := len(data)

	if h, ok := w.handler.(WriteHandler); ok  {
		return length, h(w.ctx, data)

	} else if h, ok := w.handler.(MsgHandler); ok {
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
