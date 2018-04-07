package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

type promResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

type Result struct {
	Metric string
	Values map[string]string
}

func Query(host string, start time.Time, end time.Time, query string, step time.Duration) ([]Result, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	if step == 0 {
		step = steps(end.Sub(start))
	}

	u.Path = "/api/v1/query_range"
	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.Unix()))
	q.Set("end", fmt.Sprintf("%d", end.Unix()))
	q.Set("step", fmt.Sprintf("%d", step / time.Second))
	u.RawQuery = q.Encode()

	response, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("didn't return 200 OK but %s: %s", response.Status, u.String())
	}

	var resp promResponse
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, err
	}

	if resp.Data.ResultType != "matrix" {
		return nil, fmt.Errorf("result type isn't of type matrix: %s", resp.Data.ResultType)
	}

	if len(resp.Data.Result) == 0 {
		return nil, fmt.Errorf(color.YellowString("no timeseries found"))
	}

	var results []Result
	for _, res := range resp.Data.Result {
		r := Result{}
		r.Metric = metricName(res.Metric)

		values := make(map[string]string)
		for _, vals := range res.Values {
			timestamp := vals[0].(float64)
			value := vals[1].(string)
			values[fmt.Sprintf("%.f", timestamp)] = value
		}
		r.Values = values

		results = append(results, r)
	}

	return results, nil
}

func steps(dur time.Duration) time.Duration {
	if dur < 15*time.Minute {
		return 1 * time.Second
	}
	if dur < 30*time.Minute {
		return 3 * time.Second
	}
	return time.Duration(dur / 254)
}

func metricName(metric map[string]string) string {
	if len(metric) == 0 {
		return "{}"
	}

	out := ""
	var inner []string
	for key, value := range metric {
		if key == "__name__" {
			out = value
			continue
		}
		inner = append(inner, fmt.Sprintf(`%s="%s"`, key, value))
	}

	if len(inner) == 0 {
		return out
	}

	sort.Slice(inner, func(i, j int) bool {
		return inner[i] < inner[j]
	})

	return out + "{" + strings.Join(inner, ",") + "}"
}
