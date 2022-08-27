package prometheusremotewrite

import (
	"fmt"
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
	var err error
	var metrics []telegraf.Metric
	var req logproto.PushRequest

	if err := req.Unmarshal(buf); err != nil {
		return nil, fmt.Errorf("unable to unmarshal request body: %s", err)
	}

	now := time.Now()

	for _, ts := range req.Streams {
		tags := map[string]string{}
		for key, value := range p.DefaultTags {
			tags[key] = value
		}

		//What format is Labels stored in, Appears to be a string so does the string need to be unpacked somehow?
		// Labels = `{job="foobar", cluster="foo-central1", namespace="bar", container_name="buzz"}`
		//This is ugly but i dont see an easy to use function for this in logproto
		labels := strings.Split(ts.Labels[1:len(ts.Labels)-1], ",")
		for _, label := range labels {
			labelSplit := strings.Split(label, "=") //This will break if any label name or value contains "="
			tags[labelSplit[0]] = labelSplit[1]
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
				m := metric.New("Loki", tags, fields, t)
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
