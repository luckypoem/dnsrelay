package dnsrelay

import (
	"github.com/miekg/dns"
	"net"
)

type WriteHandler func(m []byte) (int, error)
type MsgHandler func(msg *dns.Msg) error

type sessionWriter struct {
	writeFunc  WriteHandler
	msgHandler MsgHandler
}

// WriteMsg implements the ResponseWriter.WriteMsg method.
func (w *sessionWriter) WriteMsg(m *dns.Msg) (err error) {
	if w.writeFunc != nil {
		var data []byte
		data, err = m.Pack()
		if err != nil {
			return err
		}

		_, err = w.Write(data)
		return err
	} else if w.msgHandler != nil {
		return w.msgHandler(m)
	}

	panic("SessionWriter must initial a callback")
	return nil
}

// Write implements the ResponseWriter.Write method.
func (w *sessionWriter) Write(data []byte) (int, error) {
	if w.writeFunc != nil {
		return w.writeFunc(data)
	} else if w.msgHandler != nil {
		r := new(dns.Msg)
		err := r.Unpack(data)
		if err != nil {
			return 0, err
		}

		return len(data), w.msgHandler(r)
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
