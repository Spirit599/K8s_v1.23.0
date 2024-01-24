package main

import (
	"fmt"
	"math"
)

func pearsonCorrelation(x, y []float64) float64 {
	// n := len(x)

	// 计算均值
	meanX := mean(x)
	meanY := mean(y)

	// 计算协方差
	covXY := covariance(x, y, meanX, meanY)

	// 计算标准差
	stdDevX := stdDev(x, meanX)
	stdDevY := stdDev(y, meanY)

	// 计算相关系数
	pearson := covXY / (stdDevX * stdDevY)

	return pearson
}

func mean(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func covariance(x, y []float64, meanX, meanY float64) float64 {
	sum := 0.0
	for i := 0; i < len(x); i++ {
		sum += (x[i] - meanX) * (y[i] - meanY)
	}
	return sum / float64(len(x))
}

func stdDev(values []float64, meanValue float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += math.Pow(value-meanValue, 2)
	}
	variance := sum / float64(len(values))
	return math.Sqrt(variance)
}

func main() {
	// 示例数据
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{10, 8, 6, 4, 2}

	// 计算相关系数
	pearson := pearsonCorrelation(x, y)

	// 打印结果
	fmt.Printf("Pearson Correlation Coefficient: %.4f\n", pearson)
}
