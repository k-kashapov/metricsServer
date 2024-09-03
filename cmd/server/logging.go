package main

import (
	"log"
    "time"
	"net/http"
	"go.uber.org/zap"
)

type (
    respData struct {
        status int
        size int
    }

    loggingRespWriter struct {
        http.ResponseWriter
        responseData *respData
    }
)

func (r *loggingRespWriter) Write(b []byte) (int, error) {
    size, err := r.ResponseWriter.Write(b)
    r.responseData.size += size
    return size, err
}

func (r *loggingRespWriter) WriteHeader(statusCode int) {
    r.ResponseWriter.WriteHeader(statusCode)
    r.responseData.status = statusCode
}


func logHandler(h http.HandlerFunc) http.HandlerFunc {
    logFn := func (w http.ResponseWriter, r *http.Request) {
        logger, err := zap.NewDevelopment()
        if err != nil {
            log.Fatal("Could not create zap logger")
        }

        defer logger.Sync()

        sugar := *logger.Sugar()

        rData := &respData {
            status: 0,
            size: 0,
        }

        lWriter := loggingRespWriter {
            w,
            rData,
        }

        start := time.Now()

        h.ServeHTTP(&lWriter, r)

        duration := time.Since(start)

        sugar.Infoln("REQUEST",
            "Method:", r.Method,
            "URI:", r.RequestURI,
            "Duration:", duration)

        sugar.Infoln("RESPONSE", 
            "DataLen:", rData.size,
            "Status:", rData.status)
    }

    return http.HandlerFunc(logFn)
}
