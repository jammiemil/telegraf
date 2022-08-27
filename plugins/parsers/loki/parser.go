package prometheusremotewrite

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/parsers"

	"github.com/grafana/loki/pkg/logproto"
)

type Parser struct {
	DefaultTags map[string]string
}

func (p *Parser) Parse(buf []byte) ([]telegraf.Metric, error) {
	//Decide if this is Logproto or raw json, use the application header to identify?
	var err error
	var metrics []telegraf.Metric
	var req logproto.PushRequest

	if err := req.Unmarshal(buf); err != nil {
		return nil, fmt.Errorf("unable to unmarshal request body: %s", err)
	}
	rex := regexp.MustCompile(`(.*?)="(.*?)"`)
	now := time.Now()
	//something very similar is done here so i may want to compare https://github.com/grafana/loki/blob/main/clients/pkg/promtail/targets/lokipush/pushtarget.go#L113
	//The equivelent of a new input plugin can be seen here https://github.com/grafana/loki/blob/9e84648f3e176d82780db98e59a34d0e40560d1d/pkg/loghttp/push/push.go#L54
	for _, ts := range req.Streams {
		tags := map[string]string{}
		for key, value := range p.DefaultTags {
			tags[key] = value
		}
		// Labels = `{job="foobar", cluster="foo-central1", namespace="bar", container_name="buzz"}`
		labels := strings.Split(strings.Trim(ts.Labels, "{}"), ", ")
		for _, label := range labels {
			labelSplit := rex.FindAllStringSubmatch(label, -1)
			for _, kv := range labelSplit {
				tags[kv[1]] = kv[2]
			}
		}

		for _, s := range ts.Entries {
			fields := make(map[string]interface{})
			//Loki doesnt have metricNames, so we will just call this message
			//Should we try to parse the metric contents? For example if it looks like logfmt then level=info becomes fields["level"]="info"?
			fields["message"] = s.Line
			// converting to telegraf metric
			if len(fields) > 0 {
				t := now
				if s.Timestamp.UnixNano() > 0 {
					t = s.Timestamp
				}
				m := metric.New("loki", tags, fields, t)
				metrics = append(metrics, m)
			}
		}
	}
	return metrics, err
}

func (p *Parser) ParseLine(line string) (telegraf.Metric, error) {
	metrics, err := p.Parse([]byte(line))
	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return nil, fmt.Errorf("No metrics in line")
	}

	if len(metrics) > 1 {
		return nil, fmt.Errorf("More than one metric in line")
	}

	return metrics[0], nil
}

func (p *Parser) SetDefaultTags(tags map[string]string) {
	p.DefaultTags = tags
}

func (p *Parser) InitFromConfig(config *parsers.Config) error {
	return nil
}

func init() {
	parsers.Add("loki",
		func(defaultMetricName string) telegraf.Parser {
			return &Parser{}
		})
}
