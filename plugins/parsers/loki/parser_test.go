package prometheusremotewrite

import (
	"testing"
	"time"

	"github.com/grafana/loki/pkg/logproto"
	"github.com/stretchr/testify/require"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/testutil"
)

var (
	line        = `level=info ts=2019-12-12T15:00:08.325Z caller=compact.go:441 component=tsdb msg="compact blocks" count=3 mint=1576130400000 maxt=1576152000000 ulid=01DVX9ZHNM71GRCJS7M34Q0EV7 sources="[01DVWNC6NWY1A60AZV3Z6DGS65 01DVWW7XXX75GHA6ZDTD170CSZ 01DVX33N5W86CWJJVRPAVXJRWJ]" duration=2.897213221s`
	streamInput = logproto.Stream{
		Labels: `{job="foobar", cluster="foo-central1", namespace="bar", container_name="buzz"}`,
		Entries: []logproto.Entry{
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC), Line: line},
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 1, 0, time.UTC), Line: line},
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 2, 0, time.UTC), Line: line},
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 3, 0, time.UTC), Line: line},
		},
	}
	StreamAdapterInput = logproto.StreamAdapter{
		Labels: `{job="foobar", cluster="foo-central1", namespace="bar", container_name="buzz"}`,
		Entries: []logproto.EntryAdapter{
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC), Line: line},
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 1, 0, time.UTC), Line: line},
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 2, 0, time.UTC), Line: line},
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 3, 0, time.UTC), Line: line},
		},
	}
)

func TestParseStreamAdapter(t *testing.T) {
	//Logproto
	inoutBytes, err := StreamAdapterInput.Marshal()
	require.NoError(t, err)

	expected := []telegraf.Metric{
		testutil.MustMetric(
			"loki",
			map[string]string{
				"job":       "foobar",
				"cluster":   "foo-central1",
				"namespace": "bar", "container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
		),
		testutil.MustMetric(
			"loki",
			map[string]string{
				"job":       "foobar",
				"cluster":   "foo-central1",
				"namespace": "bar", "container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 1, 0, time.UTC),
		),
		testutil.MustMetric(
			"loki",
			map[string]string{
				"job":       "foobar",
				"cluster":   "foo-central1",
				"namespace": "bar", "container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 2, 0, time.UTC),
		),
		testutil.MustMetric(
			"loki",
			map[string]string{
				"job":       "foobar",
				"cluster":   "foo-central1",
				"namespace": "bar", "container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 3, 0, time.UTC),
		),
	}

	parser := Parser{
		DefaultTags: map[string]string{},
	}

	metrics, err := parser.Parse(inoutBytes)
	require.NoError(t, err)
	require.Len(t, metrics, 4)
	testutil.RequireMetricsEqual(t, expected, metrics, testutil.IgnoreTime(), testutil.SortMetrics())
}

func TestParseStream(t *testing.T) {
	//Logproto
	inoutBytes, err := streamInput.Marshal()
	require.NoError(t, err)

	expected := []telegraf.Metric{
		testutil.MustMetric(
			"loki",
			map[string]string{
				"job":       "foobar",
				"cluster":   "foo-central1",
				"namespace": "bar", "container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
		),
		testutil.MustMetric(
			"loki",
			map[string]string{
				"job":       "foobar",
				"cluster":   "foo-central1",
				"namespace": "bar", "container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 1, 0, time.UTC),
		),
		testutil.MustMetric(
			"loki",
			map[string]string{
				"job":       "foobar",
				"cluster":   "foo-central1",
				"namespace": "bar", "container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 2, 0, time.UTC),
		),
		testutil.MustMetric(
			"loki",
			map[string]string{
				"job":       "foobar",
				"cluster":   "foo-central1",
				"namespace": "bar", "container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 3, 0, time.UTC),
		),
	}

	parser := Parser{
		DefaultTags: map[string]string{},
	}

	metrics, err := parser.Parse(inoutBytes)
	require.NoError(t, err)
	require.Len(t, metrics, 4)
	testutil.RequireMetricsEqual(t, expected, metrics, testutil.IgnoreTime(), testutil.SortMetrics())
}

func TestDefaultTags(t *testing.T) {
	logprotoInput := logproto.StreamAdapter{
		Labels: `{job="foobar", cluster="foo-central1", namespace="bar", container_name="buzz"}`,

		Entries: []logproto.EntryAdapter{
			{Timestamp: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC), Line: line},
		},
	}

	inoutBytes, err := logprotoInput.Marshal()
	require.NoError(t, err)

	expected := []telegraf.Metric{
		testutil.MustMetric(
			"prometheus_remote_write",
			map[string]string{
				"defaultTag":     "defaultTagValue",
				"job":            "foobar",
				"cluster":        "foo-central1",
				"namespace":      "bar",
				"container_name": "buzz",
			},
			map[string]interface{}{
				"message": line,
			},
			time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
		),
	}

	parser := Parser{
		DefaultTags: map[string]string{
			"defaultTag": "defaultTagValue",
		},
	}

	metrics, err := parser.Parse(inoutBytes)
	require.NoError(t, err)
	require.Len(t, metrics, 1)
	testutil.RequireMetricsEqual(t, expected, metrics, testutil.IgnoreTime(), testutil.SortMetrics())
}
