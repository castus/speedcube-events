package externalFetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/printer"
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
	fmt.Println(stringifyConfig)
	printer.PrettyPrint(allConfigs)
	return

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
							Image:   "c4stus/webscraper:latest",
							Command: strings.Split("python /app/main.py", " "),
							Env: []v1.EnvVar{
								{
									Name:  "SCRAP_CONFIG",
									Value: stringifyConfig,
								},
							},
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
