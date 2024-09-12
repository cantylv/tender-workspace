// Copyright © ivanlobanov. All rights reserved.
package middlewares

import (
	"context"
	"net/http"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	"tender-workspace/internal/utils/recorder"
	"time"

	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type AccessLogStart struct {
	UserAgent      string
	RealIp         string
	ContentLength  int64
	URI            string
	Method         string
	StartTimeHuman string
	RequestId      string
	Logger         *zap.Logger
}

type AccessLogEnd struct {
	LatencyMs      int64
	ResponseSize   string // in bytes
	ResponseStatus int
	EndTimeHuman   string
	RequestId      string
	Logger         *zap.Logger
}

// Access
// Middleware that logs the start and end of request handling.
func Access(h http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.NewV4().String()
		ctx := context.WithValue(r.Context(), mc.ContextKey(mc.RequestID), requestId)
		r = r.WithContext(ctx)

		rec := recorder.NewResponseWriter(w)

		timeNow := time.Now()
		// nginx will proxy headers, like "User-Agent", "X-Real-IP", "Content-Length"
		startLog := AccessLogStart{
			UserAgent:      r.UserAgent(),
			RealIp:         r.RemoteAddr,
			ContentLength:  r.ContentLength,
			URI:            r.RequestURI,
			Method:         r.Method,
			StartTimeHuman: f.FormatTime(timeNow),
			RequestId:      requestId,
			Logger:         logger,
		}
		LogInitRequest(startLog)

		h.ServeHTTP(rec, r)

		timeEnd := time.Now()
		endLog := AccessLogEnd{
			LatencyMs:      timeEnd.Sub(timeNow).Milliseconds(),
			ResponseSize:   w.Header().Get("Content-Length"),
			ResponseStatus: rec.StatusCode,
			EndTimeHuman:   f.FormatTime(timeEnd),
			RequestId:      requestId,
			Logger:         logger,
		}
		LogEndRequest(endLog)
	})
}

// LogInitRequest
// Logs user-agent, real-ip and etc..
func LogInitRequest(data AccessLogStart) {
	data.Logger.Info("init request",
		zap.String("user-agent", data.UserAgent),
		zap.String("real-ip", data.RealIp),
		zap.Int64("content-length", data.ContentLength),
		zap.String("uri", data.URI),
		zap.String("method", data.Method),
		zap.String("start-time-human", data.StartTimeHuman),
		zap.String(mc.RequestID, data.RequestId),
	)
}

// LogEndRequest
// Logs latency in ms, response size and etc..
func LogEndRequest(data AccessLogEnd) {
	data.Logger.Info("end of request",
		zap.Int64("latensy-ms", data.LatencyMs),
		zap.String("response-size", data.ResponseSize),
		zap.Int("response-status", data.ResponseStatus),
		zap.String("end-time-human", data.EndTimeHuman),
		zap.String(mc.RequestID, data.RequestId),
	)
}
