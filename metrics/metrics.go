// TODO replace all usage of this by https://godoc.org/github.com/prometheus/common/expfmt#TextParser.TextToMetricFamilies
package metrics

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func Fetch() (map[string]string, error) {
	result := make(map[string]string)
	ts := httptest.NewServer(prometheus.Handler())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		return result, err
	}
	reader := bufio.NewReader(res.Body)
	err = nil
	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		line = strings.Trim(line, "\n")
		parts := strings.Split(line, " ")
		if (parts[0] != "#") && (len(parts) > 1) {
			result[parts[0]] = parts[1]
		}
	}
	if err == io.EOF {
		err = nil
	}
	return result, err
}

func Read(metrics map[string]string, name string, labels prometheus.Labels) (float64, error) {
	// http_request_duration_microseconds{handler="alerts",quantile="0.5"}
	fullName := name
	if len(labels) > 0 {
		pairs := []string{}
		// TODO SORT GOOD
		for k, v := range labels {
			pairs = append(pairs, k+"=\""+v+"\"")
		}
		fullName += "{" + strings.Join(pairs, ",") + "}"
	}

	stringValue := metrics[fullName]
	return strconv.ParseFloat(stringValue, 64)
}
