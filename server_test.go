package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func TestServerPing(t *testing.T) {
	res, err := http.Get("http://" + getEnvValue("HOST") + ":" + getEnvValue("PORT") + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	b := string(body)
	if !strings.Contains(b, "JSON") {
		t.Fatal()
	}
}

func TestHandleCollect(t *testing.T) {
	tests := []struct {
		name           string
		queryString    string
		expectedFields map[string]string
	}{
		{
			name:        "Basic types (int, float, bool, string)",
			queryString: "struct=order,id=555,price=99.99,active=true,title=Book",
			expectedFields: map[string]string{
				"Id":     "int",
				"Price":  "float64",
				"Active": "bool",
				"Title":  "string",
			},
		},
		{
			name:        "Complex types (JSON array and JSON object)",
			queryString: `struct=user,tags=["admin","vip"],meta={"ip":"127.0.0.1"}`,
			expectedFields: map[string]string{
				"Tags": "[]interface {}",
				"Meta": "map[string]interface {}",
			},
		},
	}

	host, port := getEnvValue("HOST"), getEnvValue("PORT")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedURL, err := url.Parse("http://" + host + ":" + port + "/collect")
			if err != nil {
				t.Fatalf("failed to parse base url: %v", err)
			}

			query := parsedURL.Query()
			query.Set("struct", strings.TrimPrefix(tt.queryString, "struct="))
			parsedURL.RawQuery = query.Encode()
			res, err := http.Get(parsedURL.String())
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status 200, got %d", res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read body: %v", err)
			}

			var actualJSON map[string]any
			if err := json.Unmarshal(body, &actualJSON); err != nil {
				t.Fatalf("failed to unmarshal response JSON: %v. Body: %s", err, string(body))
			}

			for key := range tt.expectedFields {
				jsonKey := strings.ToLower(key)
				if _, ok := actualJSON[jsonKey]; !ok {
					t.Errorf("expected key %s was not found in response JSON", jsonKey)
				}
			}

			if tt.name == "Basic types (int, float, bool, string)" {
				if actualJSON["id"].(float64) != 555 {
					t.Errorf("expected id to be 555, got %v", actualJSON["id"])
				}

				if actualJSON["active"].(bool) != true {
					t.Errorf("expected active to be true")
				}
			}
		})
	}
}

func TestLoadPing(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://" + getEnvValue("HOST") + ":" + getEnvValue("PORT") + "/ping",
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}
