package podautoscaler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os/exec"
	"strconv"
)

func arimaPredate(targetUtilization int64, currentReplicas int32) (replicaCount int32) {

	slice := make([]float64, 10)

	slice = append(slice, getRealRate())

	slice = slice[1:]

	dataBytes, _ := json.Marshal(slice)
	fmt.Println(string(dataBytes))

	cmd := exec.Command("python3", "ARIMA.py", string(dataBytes))

	// 创建一个用于捕获标准输出的缓冲区
	var out bytes.Buffer
	cmd.Stdout = &out

	// 运行命令
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	// 打印命令的输出
	str := out.String()
	str = str[:len(str)-1]
	preRate, _ := strconv.ParseFloat(str, 64)
	fmt.Printf("Command output:%s\n", str)
	fmt.Printf("Command output:%f\n", preRate)

	tU := float64(targetUtilization)
	cR := float64(currentReplicas)

	newReplicas := int32(math.Ceil(preRate * cR / tU))

	return newReplicas
}

func getRealRate() (rate float64) {

	usedRate := make([]float64, 0)
	usedRate = append(usedRate, getCpuRate())
	usedRate = append(usedRate, getMemRate())
	usedRate = append(usedRate, getDiskRate())
	usedRate = append(usedRate, getNetRate())

	sum := 0.0
	for _, value := range usedRate {
		sum += value
	}

	weight := make([]float64, len(usedRate))
	for idx, value := range usedRate {
		weight[idx] = value / sum
	}

	finalRate := 0.0
	for idx, value := range usedRate {
		finalRate += value * weight[idx]
	}

	return finalRate
}

func getCpuRate() (cpuRate float64) {

	var curCpuRate = map[string]float64{}

	query := "sum(irate(container_cpu_usage_seconds_total{container=\"nginx\"}[1m])*100)by(pod)"
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
	// fmt.Println("Prometheus JSON Data:", string(result))

	// 解析 JSON 数据
	var promResp PrometheusResponse
	err = json.Unmarshal(result, &promResp)
	if err != nil {
		fmt.Println("Error unmarshalling Prometheus data:", err)
		return
	}

	sumRate := 0.0
	// 使用 for 循环打印解析后的数据
	for _, r := range promResp.Data.Result {
		podName := r.Metric["pod"]
		// 假设 Value 的第二个元素是我们需要的值
		if len(r.Value) > 1 {
			value, ok := r.Value[1].(string)
			if ok {
				num, _ := strconv.ParseFloat(value, 64)
				// fmt.Printf("%s: %s %f %f %f\n", podName, value, num/mb, num, mb)
				curCpuRate[podName] = num
				sumRate += num
			}
		}
	}

	return sumRate / float64(len(curCpuRate))
}

func getMemRate() (cpuRate float64) {

	var curMemRate = map[string]float64{}

	query := "sum(container_memory_usage_bytes{container=\"nginx\"})by(pod)/sum(container_spec_memory_limit_bytes{container=\"nginx\"})by(pod)*100"
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
	// fmt.Println("Prometheus JSON Data:", string(result))

	// 解析 JSON 数据
	var promResp PrometheusResponse
	err = json.Unmarshal(result, &promResp)
	if err != nil {
		fmt.Println("Error unmarshalling Prometheus data:", err)
		return
	}

	sumRate := 0.0
	// 使用 for 循环打印解析后的数据
	for _, r := range promResp.Data.Result {
		podName := r.Metric["pod"]
		// 假设 Value 的第二个元素是我们需要的值
		if len(r.Value) > 1 {
			value, ok := r.Value[1].(string)
			if ok {
				num, _ := strconv.ParseFloat(value, 64)
				// fmt.Printf("%s: %s %f %f %f\n", podName, value, num/mb, num, mb)
				curMemRate[podName] = num
				sumRate += num
			}
		}
	}

	return sumRate / float64(len(curMemRate))
}

func getDiskRate() (cpuRate float64) {

	var curDiskRate = map[string]float64{}

	query := "sum(rate(container_fs_writes_bytes_total{container=\"nginx\"}[1m]))by(pod)"
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
	// fmt.Println("Prometheus JSON Data:", string(result))

	// 解析 JSON 数据
	var promResp PrometheusResponse
	err = json.Unmarshal(result, &promResp)
	if err != nil {
		fmt.Println("Error unmarshalling Prometheus data:", err)
		return
	}

	sumRate := 0.0
	// 使用 for 循环打印解析后的数据
	for _, r := range promResp.Data.Result {
		podName := r.Metric["pod"]
		// 假设 Value 的第二个元素是我们需要的值
		if len(r.Value) > 1 {
			value, ok := r.Value[1].(string)
			if ok {
				num, _ := strconv.ParseFloat(value, 64)
				// fmt.Printf("%s: %s %f %f %f\n", podName, value, num/mb, num, mb)
				curDiskRate[podName] = num
				sumRate += num
			}
		}
	}

	return sumRate / float64(len(curDiskRate))
}

func getNetRate() (cpuRate float64) {

	var curNetRate = map[string]float64{}

	query := "sum(rate(container_network_receive_bytes_total{container=\"nginx\"}[1m]))by(pod)"
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
	// fmt.Println("Prometheus JSON Data:", string(result))

	// 解析 JSON 数据
	var promResp PrometheusResponse
	err = json.Unmarshal(result, &promResp)
	if err != nil {
		fmt.Println("Error unmarshalling Prometheus data:", err)
		return
	}

	sumRate := 0.0
	// 使用 for 循环打印解析后的数据
	for _, r := range promResp.Data.Result {
		podName := r.Metric["pod"]
		// 假设 Value 的第二个元素是我们需要的值
		if len(r.Value) > 1 {
			value, ok := r.Value[1].(string)
			if ok {
				num, _ := strconv.ParseFloat(value, 64)
				// fmt.Printf("%s: %s %f %f %f\n", podName, value, num/mb, num, mb)
				curNetRate[podName] = num
				sumRate += num
			}
		}
	}

	return sumRate / float64(len(curNetRate))
}
