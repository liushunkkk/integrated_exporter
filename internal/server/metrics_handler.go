package server

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
)

var metricsHandler *MetricsHandler

type MetricsHandler struct {
	OldMetricsBuffer []*bytes.Buffer
	MetricsBuffers   []*bytes.Buffer
}

func NewMetricsHandler() *MetricsHandler {
	if metricsHandler == nil {
		metricsHandler = &MetricsHandler{}
	}
	return metricsHandler
}

func (mh *MetricsHandler) AddBuffer(buf *bytes.Buffer) {
	mh.MetricsBuffers = append(mh.MetricsBuffers, buf)
}

func (mh *MetricsHandler) ClearBuffer() {
	mh.OldMetricsBuffer = mh.MetricsBuffers
	mh.MetricsBuffers = []*bytes.Buffer{}
}

func (mh *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var result bytes.Buffer
	for _, buf := range mh.OldMetricsBuffer {
		if _, err := result.Write(buf.Bytes()); err != nil {
			http.Error(w, "Failed to combine buffers", http.StatusInternalServerError)
			return
		}
	}

	// 检查请求头是否接受 Gzip 压缩
	if r.Header.Get("Accept-Encoding") == "gzip" {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		gz := gzip.NewWriter(w)
		defer func() {
			if err := gz.Close(); err != nil {
				log.Println("Error closing gzip writer:", err)
			}
		}()

		if _, err := gz.Write(result.Bytes()); err != nil {
			http.Error(w, "Failed to write compressed response", http.StatusInternalServerError)
			log.Println("Error writing compressed response:", err)
			return
		}

		return
	}

	// 不支持 Gzip 时，直接返回数据
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", result.Len()))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(result.Bytes()); err != nil {
		log.Println("Error writing response:", err)
	}
}
