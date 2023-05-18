package gelf

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestBufferingMax tests the absolute maximum buffer size will result in a valid write
func TestBufferingMax(t *testing.T) {
	a := assert.New(t)
	const chunkDataSize = MaxChunkSize - 12
	w := &writer{
		conn:             nopConn{Writer: io.Discard},
		chunkSize:        MaxChunkSize,
		chunkDataSize:    chunkDataSize,
		compressionType:  CompressionNone,
		compressionLevel: 9,
	}
	const max = chunkDataSize*MaxChunkCount - 1
	b := bytes.Repeat([]byte{'a'}, max)
	n, err := w.Write(b)
	a.NoError(err)
	a.Equal(max, n)
	const expErr = "need 129 chunks but should be less than or equal to 128"
	b = append(b, 'a')
	n, err = w.Write(b)
	a.Equal(n, 0)
	a.EqualError(err, expErr, "adding one more byte should result in too many chunks")
}

// nopConn implements net.Conn with stubs for tests that use Conn.
// all unused methods panic for obvious failure of future tests.
type nopConn struct{ io.Writer }

func (nopConn) Close() (err error)                 { return }
func (nopConn) LocalAddr() (a net.Addr)            { panic("unimplemented") }
func (nopConn) Read(b []byte) (n int, err error)   { panic("unimplemented") }
func (nopConn) RemoteAddr() net.Addr               { panic("unimplemented") }
func (nopConn) SetDeadline(t time.Time) error      { panic("unimplemented") }
func (nopConn) SetReadDeadline(t time.Time) error  { panic("unimplemented") }
func (nopConn) SetWriteDeadline(t time.Time) error { panic("unimplemented") }

var _ net.Conn = nopConn{}
