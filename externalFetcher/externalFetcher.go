package externalFetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/castus/speedcube-events/logger"

	k8sBatchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var PageTypes = struct {
	Info        string
	Competitors string
}{
	Info:        "info",
	Competitors: "competitors",
}

var log = logger.Default()

type ExternalFetchConfig struct {
	Type         string
	Id           string
	URL          string
	S3BucketPath string
	PageType     string
}

type K8SConfigCube4FunDTO struct {
	Type string
	Id   string
	URL  string
}

type K8SConfigPPODTO struct {
	Type string
	Id   string
	URL  string
}

func PrintK8sJobsConfig(c4f []K8SConfigCube4FunDTO, ppo []K8SConfigPPODTO) string {
	return stringifyConfig(c4f, ppo)
}

func SpinK8sJobsToFetchExternalData(c4f []K8SConfigCube4FunDTO, ppo []K8SConfigPPODTO) {
	jobConfig := stringifyConfig(c4f, ppo)
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	k8s, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	jobs := k8s.BatchV1().Jobs("speedcube-events")
	var backOffLimit int32 = 0
	var ttlSecondsAfterFinished int32 = 300
	now := time.Now()
	jobSpec := &k8sBatchV1.Job{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      fmt.Sprintf("scrape-web-%s", now.Format("20060102150405")),
			Namespace: "speedcube-events",
		},
		Spec: k8sBatchV1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    "run-webscraper",
							Image:   "c4stus/speedcube-web-scraper:latest",
							Command: strings.Split("python /app/main.py", " "),
							Env:     envs(jobConfig),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err = jobs.Create(context.TODO(), jobSpec, k8sMetaV1.CreateOptions{})
	if err != nil {
		log.Error("Failed to create K8s job", "error", err)
		panic(err)
	}
}

func stringifyConfig(c4f []K8SConfigCube4FunDTO, ppo []K8SConfigPPODTO) string {
	var allConfigs = []ExternalFetchConfig{}
	for _, item := range c4f {
		allConfigs = append(allConfigs, cube4FunConfig(PageTypes.Info, item))
	}
	for _, item := range ppo {
		allConfigs = append(allConfigs, ppoConfig(PageTypes.Info, item))
		allConfigs = append(allConfigs, ppoConfig(PageTypes.Competitors, item))
	}
	j, _ := json.Marshal(allConfigs)
	str := string(j)
	stringifyConfig := strings.ReplaceAll(str, `\`, `\\`)

	return stringifyConfig
}

func envs(config string) []v1.EnvVar {
	var envs = []v1.EnvVar{}
	envs = append(envs, v1.EnvVar{
		Name:  "SCRAP_CONFIG",
		Value: config,
	})

	return append(envs, secrets()...)
}

func secrets() []v1.EnvVar {
	var secrets = []v1.EnvVar{}
	keys := []string{
		"S3_WEB_DATA_BUCKET_NAME",
		"AWS_S3_API_SECRET",
		"AWS_S3_API_KEY",
	}
	name := "speedcube-events-secrets"
	for _, key := range keys {
		secrets = append(secrets, v1.EnvVar{
			Name: key,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: name,
					},
					Key: key,
				},
			},
		})
	}

	return secrets
}

func cube4FunConfig(page string, item K8SConfigCube4FunDTO) ExternalFetchConfig {
	return ExternalFetchConfig{
		Type:         item.Type,
		Id:           item.Id,
		URL:          fmt.Sprintf("%s/%s", item.URL, page),
		S3BucketPath: fmt.Sprintf("%s/%s/%s.html", item.Type, item.Id, page),
		PageType:     page,
	}
}

func ppoConfig(page string, item K8SConfigPPODTO) ExternalFetchConfig {
	return ExternalFetchConfig{
		Type:         item.Type,
		Id:           item.Id,
		URL:          item.URL,
		S3BucketPath: fmt.Sprintf("%s/%s/%s.html", item.Type, item.Id, page),
		PageType:     page,
	}
}
