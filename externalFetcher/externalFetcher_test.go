package externalFetcher

import (
	"strings"
	"testing"
)

func TestK8SJobsConfig_PPO(t *testing.T) {
	mock := []K8SConfigPPODTO{
		{
			Type: "PPO",
			Id:   "ppo-id",
			URL:  "http://google.com",
		},
	}
	sut := stringifyConfig([]K8SConfigCube4FunDTO{}, mock)
	out := `
	[
	{
		"Type": "PPO",
		"Id": "ppo-id",
		"URL": "http://google.com",
		"S3BucketPath": "PPO/ppo-id/info.html",
		"PageType": "info"
	},
	{
		"Type": "PPO",
		"Id": "ppo-id",
		"URL": "http://google.com",
		"S3BucketPath": "PPO/ppo-id/competitors.html",
		"PageType": "competitors"
	}
	]
	`
	out = oneLine(out)
	if out != sut {
		t.Error("Expect PPO config, got", out)
	}
}

func TestK8SJobsConfig_Cube4Fun(t *testing.T) {
	mock := []K8SConfigCube4FunDTO{
		{
			Type: "Cube4Fun",
			Id:   "c4f-id",
			URL:  "http://google.com",
		},
	}
	sut := stringifyConfig(mock, []K8SConfigPPODTO{})
	out := `
	[
	{
		"Type": "Cube4Fun",
		"Id": "c4f-id",
		"URL": "http://google.com/info",
		"S3BucketPath": "Cube4Fun/c4f-id/info.html",
		"PageType": "info"
	}
	]
	`
	out = oneLine(out)
	if out != sut {
		t.Error("Expect Cube4Fun config, got", out)
	}
}

func oneLine(multiLine string) string {
	multiLine = strings.ReplaceAll(multiLine, "\n", "")
	multiLine = strings.ReplaceAll(multiLine, "\t", "")
	multiLine = strings.ReplaceAll(multiLine, " ", "")

	return multiLine
}
