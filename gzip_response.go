// Copyright (c) Jeevanandam M (https://github.com/jeevatkm)
// go-aah/ahttp source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package ahttp

import (
	"bufio"
	"compress/gzip"
	"io"
	"net"
	"net/http"

	"aahframework.org/essentials.v0"
)

// GzipResponse extends `ahttp.Response` and provides gzip for response
// bytes before writing them to the underlying response.
type GzipResponse struct {
	r  *Response
	gw *gzip.Writer
}

// interface compliance
var (
	_ http.CloseNotifier = &GzipResponse{}
	_ http.Flusher       = &GzipResponse{}
	_ http.Hijacker      = &GzipResponse{}
	_ io.Closer          = &GzipResponse{}
	_ ResponseWriter     = &GzipResponse{}
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Global methods
//___________________________________

// WrapGzipResponseWriter wraps `http.ResponseWriter`, returns aah framework response
// writer that allows to advantage of response process.
func WrapGzipResponseWriter(w http.ResponseWriter, level int) ResponseWriter {
	rw := WrapResponseWriter(w)

	// Since Gzip level is validated in the framework while loading,
	// so expected to have valid level which between 1 and 9.
	gzw, _ := gzip.NewWriterLevel(rw, level)
	return &GzipResponse{gw: gzw, r: rw.(*Response)}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Response methods
//___________________________________

// Status method returns HTTP response status code. If status is not yet written
// it reurns 0.
func (g *GzipResponse) Status() int {
	return g.r.Status()
}

// WriteHeader method writes given status code into Response.
func (g *GzipResponse) WriteHeader(code int) {
	g.r.WriteHeader(code)
}

// Header method returns response header map.
func (g *GzipResponse) Header() http.Header {
	return g.r.Header()
}

// Write method writes bytes into Response.
func (g *GzipResponse) Write(b []byte) (int, error) {
	g.r.setContentTypeIfNotSet(b)
	g.r.WriteHeader(http.StatusOK)

	size, err := g.gw.Write(b)
	g.r.bytesWritten += size
	return size, err
}

// BytesWritten method returns no. of bytes already written into HTTP response.
func (g *GzipResponse) BytesWritten() int {
	return g.r.BytesWritten()
}

// Close method closes the writer if possible.
func (g *GzipResponse) Close() error {
	ess.CloseQuietly(g.gw)
	g.gw = nil
	return g.r.Close()
}

// Unwrap method returns the underlying `http.ResponseWriter`
func (g *GzipResponse) Unwrap() http.ResponseWriter {
	return g.r.Unwrap()
}

// CloseNotify method calls underlying CloseNotify method if it's compatible
func (g *GzipResponse) CloseNotify() <-chan bool {
	return g.r.CloseNotify()
}

// Flush method calls underlying Flush method if it's compatible
func (g *GzipResponse) Flush() {
	if g.gw != nil {
		_ = g.gw.Flush()
	}

	g.r.Flush()
}

// Hijack method calls underlying Hijack method if it's compatible otherwise
// returns an error. It becomes the caller's responsibility to manage
// and close the connection.
func (g *GzipResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return g.r.Hijack()
}
