package leastrequested

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func getDiskIOData() {
	// query := "(1-sum(rate(node_cpu_seconds_total{mode=\"idle\"}[1m]))by(instance)/sum(rate(node_cpu_seconds_total[1m]))by(instance))*100"
	query := "sum by(instance) (rate(node_disk_written_bytes_total[1m]))"
	// 对查询进行 URL 编码
	encodedQuery := url.QueryEscape(query)

	// 构建 Prometheus 查询 URL
	prometheusURL := fmt.Sprintf("http://192.168.0.60:30003/api/v1/query?query=%s&format=json", encodedQuery)

	// 获取 Prometheus 查询结果
	result, err := getPrometheusData(prometheusURL)
	if err != nil {
		fmt.Println("Error getting Prometheus data:", err)
		return
	}

	// 打印获取的 JSON 数据
	fmt.Println("Prometheus JSON Data:", string(result))

	// 解析 JSON 数据
	var promResp PrometheusResponse
	err = json.Unmarshal(result, &promResp)
	if err != nil {
		fmt.Println("Error unmarshalling Prometheus data:", err)
		return
	}

	// 使用 for 循环打印解析后的数据
	for _, r := range promResp.Data.Result {
		nodeName := r.Metric["instance"]
		// 假设 Value 的第二个元素是我们需要的值
		if len(r.Value) > 1 {
			value, ok := r.Value[1].(string)
			if ok {
				num, _ := strconv.ParseFloat(value, 64)
				// fmt.Printf("%s: %s %f %f %f\n", nodeName, value, num/mb, num, mb)
				curDiskIO[nodeName] = num / mb
			}
		}
	}
}
