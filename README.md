# Audit/Logger Webhook

`webhook` listens for incoming events from MinIO server and logs these events to a log file.

Usage:
```
webhook --log-file <logfile>
```

Environment only settings:

| ENV                | Description                                                        |
|--------------------|--------------------------------------------------------------------|
| WEBHOOK_AUTH_TOKEN | Authorization token optional to authenticate/trust incoming events |

The webhook service can be setup as a systemd service using the `webhook.service` shipped with
this project

To send Audit logs from MinIO server, please configure MinIO using the command:
```
mc admin config set myminio audit_webhook endpoint=http://webhookendpoint:8080 auth_token=webhooksecret
```

To send server logs from MinIO server, please configure MinIO using the command:
```
mc admin config set myminio logger_webhook endpoint=http://webhookendpoint:8081 auth_token=webhooksecret
```

> NOTE: audit_webhook and logger_webhook should *not* be configured to send events to the same webhook instance.

```
$ ./webhook -h     
Usage of ./webhook:
  -address string
        bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname (default ":8080")
  -compress
        Compress determines if the rotated log files should be compressed using gzip. The default is not to perform compression.
  -log-file string
        path to the file where webhook will log incoming events
  -maxAge int
        MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename. Note that a day is defined as 24 hours and may not exactly correspond to calendar days due to daylight savings, leap seconds, etc. (default 30)
  -maxBackups int
        MaxBackups is the maximum number of old log files to retain. (default 5)
  -maxSize int
        MaxSize is the maximum size in megabytes of the log file before it gets rotated. It defaults to 1024*5 megabytes. (default 5120)
```
