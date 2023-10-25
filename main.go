// Copyright (c) 2015-2023 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/minio/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logFile    string
	maxSize    int
	maxBackups int
	maxAge     int
	compress   bool
	address    string
	authToken  = env.Get("WEBHOOK_AUTH_TOKEN", "")
)

func main() {
	flag.StringVar(&logFile, "log-file", "", "path to the file where webhook will log incoming events")
	flag.IntVar(&maxSize, "maxSize", 1024*5, "MaxSize is the maximum size in megabytes of the log file before it gets rotated. It defaults to 100 megabytes.")
	flag.IntVar(&maxBackups, "maxBackups", 5, "MaxBackups is the maximum number of old log files to retain. The default is to retain all old log files (though MaxAge may still cause them to get deleted.)")
	flag.IntVar(&maxAge, "maxAge", 30, "MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename. Note that a day is defined as 24 hours and may not exactly correspond to calendar days due to daylight savings, leap seconds, etc. The default is not to remove old log files based on age.")
	flag.BoolVar(&compress, "compress", false, "Compress determines if the rotated log files should be compressed using gzip. The default is not to perform compression.")
	flag.StringVar(&address, "address", ":8080", "bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname")

	flag.Parse()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""  // 关闭时间戳
	encoderCfg.LevelKey = "" // 关闭日志级别
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg), // 选择控制台编码器
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}),
		zap.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()

	if logFile == "" {
		log.Fatalln("--log-file must be specified")
	}
	log.Println("listen addr", address)
	err := http.ListenAndServe(address, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authToken != "" {
			if authToken != r.Header.Get("Authorization") {
				http.Error(w, "authorization header missing", http.StatusBadRequest)
				return
			}
		}
		switch r.Method {
		case http.MethodPost:
			// 读取请求body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()
			logger.Info(string(body))
		default:
		}
	}))
	if err != nil {
		log.Fatal(err)
	}
}
