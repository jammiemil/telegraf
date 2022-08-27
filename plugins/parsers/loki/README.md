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

logproto.PushRequest{
		Streams: []logproto.Stream{
			{
				Labels: `{job="foobar", cluster="foo-central1", namespace="bar", container_name="buzz"}`,
				Entries: []logproto.Entry{
					{Timestamp: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC), Line: `level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg="compact blocks" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources="[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]" duration=2.897213221s`},
					{Timestamp: time.Date(2020, 4, 1, 0, 0, 1, 0, time.UTC), Line: `level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg="compact blocks" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources="[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]" duration=2.897213221s`},
					{Timestamp: time.Date(2020, 4, 1, 0, 0, 2, 0, time.UTC), Line: `level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg="compact blocks" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources="[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]" duration=2.897213221s`},
					{Timestamp: time.Date(2020, 4, 1, 0, 0, 3, 0, time.UTC), Line: `level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg="compact blocks" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources="[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]" duration=2.897213221s`},
				},
			},
		},
	}

```

## Example Output

```text
loki,job=foobar,cluster=foo-central1,namespace=bar,container_name=buzz message="level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg=\"compact blocks\" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources=\"[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]\" duration=2.897213221s" 1585699200000000000
loki,source=applicationLog,instance=localhost,job=prometheus message="level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg=\"compact blocks\" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources=\"[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]\" duration=2.897213221s" 1585699201000000000
loki,job=foobar,cluster=foo-central1,namespace=bar,container_name=buzz message="level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg=\"compact blocks\" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources=\"[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]\" duration=2.897213221s" 1585699202000000000
loki,source=applicationLog,instance=localhost,job=prometheus message="level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg=\"compact blocks\" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources=\"[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]\" duration=2.897213221s" 1585699203000000000
```
