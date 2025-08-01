package collector

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ecs_exporter/config"
	"ecs_exporter/token"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type ECSCollector struct {
	apiCfg      *config.APIConfig
	tokenMgr    *token.TokenManager
	descRunning *prometheus.Desc
	descRAM     *prometheus.Desc
}

type ECSRequest struct {
	Dimensions []struct {
		Field string `json:"field"`
		Index int    `json:"index"`
	} `json:"dimensions"`
	Metrics []struct {
		Field   string `json:"field"`
		AggType string `json:"aggType"`
	} `json:"metrics"`
}

func NewECSCollector(cfg *config.APIConfig, tm *token.TokenManager) *ECSCollector {
	return &ECSCollector{
		apiCfg:   cfg,
		tokenMgr: tm,
		descRunning: prometheus.NewDesc(
			"ecs_instance_running_time_seconds",
			"Running time of ECS instance",
			[]string{"name", "vdc_level2", "status", "os_version", "flavor_name", "azone", "cluster", "project"},
			nil,
		),
		descRAM: prometheus.NewDesc(
			"ecs_instance_ram_gb",
			"RAM size of ECS instance (GB)",
			[]string{"name"},
			nil,
		),
	}
}

func (c *ECSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.descRunning
	ch <- c.descRAM
}

func (c *ECSCollector) Collect(ch chan<- prometheus.Metric) {
	pageNo := 1
	pageSize := c.apiCfg.PageSize
	if pageSize == 0 {
		pageSize = 100
	}

	body := map[string]interface{}{
		"dimensions": []map[string]interface{}{
			{
				"field": "uuid",
				"index": 1,
			},
			{
				"field": "name",
				"index": 2,
			},
		},
		"metrics": []map[string]interface{}{
			{
				"field":   "cpusize",
				"aggType": "sum",
			},
			{
				"field":   "ramsize",
				"aggType": "avg",
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		logrus.Errorf("construct request failed: %v", err)
		return
	}

	for {
		url := fmt.Sprintf("%s?pageNo=%d&pageSize=%d", c.apiCfg.Endpoint, pageNo, pageSize)

		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
		if err != nil {
			logrus.Errorf("new request error: %v", err)
			return
		}

		req.Header.Set("X-Auth-Token", c.tokenMgr.GetToken())

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Errorf("api request error: %v", err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			logrus.Errorf("api request returned status: %v", resp.Status)
			resp.Body.Close()
			return
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			logrus.Errorf("read response body error: %v", err)
			return
		}

		var response struct {
			TotalSize int                      `json:"totalSize"`
			PageNo    int                      `json:"pageNo"`
			PageSize  int                      `json:"pageSize"`
			Datas     []map[string]interface{} `json:"datas"`
		}

		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			logrus.Errorf("json unmarshal error: %v", err)
			return
		}
		logrus.Infof("dataset: %++v", response)
		for _, item := range response.Datas {
			// 取出各字段, 你需要根据实际字段名调整
			name, _ := item["object.name"].(string)
			vdcLevel2, _ := item["vdc.vdcLevel2"].(string)
			status, _ := item["otherInfo.status"].(string)
			osVersion, _ := item["otherInfo.osVersion"].(string)
			flavorName, _ := item["otherInfo.flavorName"].(string)
			azone, _ := item["logicLoc.azoneName"].(string)
			cluster, _ := item["logicLoc.clusterName"].(string)
			project, _ := item["tenant.projectName"].(string)

			runningTime, _ := item["runningTime"].(float64)
			ramSize, _ := item["ramSize"].(float64)

			ch <- prometheus.MustNewConstMetric(
				c.descRunning,
				prometheus.GaugeValue,
				runningTime,
				name,
				vdcLevel2,
				status,
				osVersion,
				flavorName,
				azone,
				cluster,
				project,
			)

			ch <- prometheus.MustNewConstMetric(
				c.descRAM,
				prometheus.GaugeValue,
				ramSize,
				name,
			)
		}

		if (pageNo * pageSize) >= response.TotalSize {
			break
		}
		pageNo++
	}
}
