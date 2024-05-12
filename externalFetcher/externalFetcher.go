package externalFetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/logger"

	k8sBatchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var log = logger.Default()

type ExternalFetchConfig struct {
	Type         string
	Id           string
	URL          string
	S3BucketPath string
}

func SpinK8sJobsToFetchExternalData(competitions db.Competitions) {
	var allConfigs = []ExternalFetchConfig{}
	c4fItems := FetchConfigCube4Fun(competitions)
	allConfigs = append(allConfigs, c4fItems...)
	j, _ := json.Marshal(allConfigs)
	str := string(j)
	stringifyConfig := strings.ReplaceAll(str, `\`, `\\`)
	stringifyConfig = strings.ReplaceAll(stringifyConfig, `"`, `\"`)
	stringifyConfig = strings.ReplaceAll(stringifyConfig, `'`, `\'`)

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
							Env:     envs(stringifyConfig),
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
		log.Error("Failed to create K8s job", err)
		panic(err)
	}
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
