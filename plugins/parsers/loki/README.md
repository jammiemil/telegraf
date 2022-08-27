# Prometheus Remote Write Parser Plugin

Converts Loki Push Streams directly into Telegraf metrics. It can
be used with [http_listener_v2](/plugins/inputs/http_listener_v2). There are no
additional configuration options for Loki Push Streams.

## Configuration

```toml
[[inputs.http_listener_v2]]
  ## Address and port to host HTTP listener on
  service_address = ":3100"

  ## Paths to listen to.
  paths = ["/loki/api/v1/push"]

  ## Data format to consume.
  data_format = "loki"
```

## Example Input

```json
protobuf:

logproto.PushRequest{
        Streams: []*logproto.Streams{
            {
                Labels: []*logproto.Labels{
                    {Name: "source", Value: "applicationLog"},
                    {Name: "instance", Value: "localhost"},
                    {Name: "job", Value: "promtail"},
                },
                Entries: []logproto.Entries{
                    {Timestamp: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC).UnixNano(), Entry: "This is a log Line"},
                    {Timestamp: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC).UnixNano(), Entry: "This is another log Line"}
                },
            },
        },
    }

Json body:

{
  "streams": [
    {
      "stream": {
        "source": "applicationLog",
        "instance": "localhost",
        "job": "promtail"
      },
      "values": [
          [ time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC).UnixNano(), "This is a log Line" ],
          [ time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC).UnixNano(), "This is another log Line" ]
      ]
    }
  ]
}

```

## Example Output

```text
loki,source=applicationLog,instance=localhost,job=prometheus message="This is a log line" 1614889298859000000
loki,source=applicationLog,instance=localhost,job=prometheus message="This is another log Line" 1614889298859000000
```
