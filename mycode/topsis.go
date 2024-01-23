package main

import (
	"fmt"
	"math"
)

// 归一化决策矩阵
func normalize(decisionMatrix [][]float64) [][]float64 {
	numRows := len(decisionMatrix)
	numCols := len(decisionMatrix[0])
	normalizedMatrix := make([][]float64, numRows)

	for i := range normalizedMatrix {
		normalizedMatrix[i] = make([]float64, numCols)
	}

	for j := 0; j < numCols; j++ {
		var sum float64 = 0
		for i := 0; i < numRows; i++ {
			sum += decisionMatrix[i][j] * decisionMatrix[i][j]
		}
		sum = math.Sqrt(sum)
		for i := 0; i < numRows; i++ {
			normalizedMatrix[i][j] = decisionMatrix[i][j] / sum
		}
	}

	return normalizedMatrix
}

// 计算加权归一化决策矩阵
func weightedNormalizedMatrix(normalizedMatrix [][]float64, weights []float64) [][]float64 {
	weightedMatrix := make([][]float64, len(normalizedMatrix))

	for i, row := range normalizedMatrix {
		weightedMatrix[i] = make([]float64, len(row))
		for j, value := range row {
			weightedMatrix[i][j] = value * weights[j]
		}
	}

	return weightedMatrix
}

// 计算理想解和负理想解
func idealSolutions(weightedMatrix [][]float64) ([]float64, []float64) {
	numCols := len(weightedMatrix[0])
	ideal := make([]float64, numCols)
	negativeIdeal := make([]float64, numCols)

	for j := 0; j < numCols; j++ {
		ideal[j] = math.SmallestNonzeroFloat64
		negativeIdeal[j] = math.MaxFloat64
		for i := 0; i < len(weightedMatrix); i++ {
			if weightedMatrix[i][j] > ideal[j] {
				ideal[j] = weightedMatrix[i][j]
			}
			if weightedMatrix[i][j] < negativeIdeal[j] {
				negativeIdeal[j] = weightedMatrix[i][j]
			}
		}
	}

	return ideal, negativeIdeal
}

// 计算与理想解和负理想解的距离
func distanceAndRank(weightedMatrix [][]float64, ideal []float64, negativeIdeal []float64) []float64 {
	scores := make([]float64, len(weightedMatrix))
	for i, row := range weightedMatrix {
		var idealDist, negIdealDist float64 = 0, 0
		for j := range row {
			idealDist += math.Pow(row[j]-ideal[j], 2)
			negIdealDist += math.Pow(row[j]-negativeIdeal[j], 2)
		}
		scores[i] = math.Sqrt(negIdealDist) / (math.Sqrt(idealDist) + math.Sqrt(negIdealDist))
	}
	return scores
}

func main() {
	// 示例决策矩阵
	decisionMatrix := [][]float64{
		{89, 1},
		{60, 3},
		{74, 2},
		{99, 0},
	}
	// 示例权重
	weights := []float64{1, 1, 1, 1}

	normalizedMatrix := normalize(decisionMatrix)
	fmt.Println("Scores:", normalizedMatrix)
	weightedMatrix := weightedNormalizedMatrix(normalizedMatrix, weights)
	ideal, negativeIdeal := idealSolutions(weightedMatrix)
	scores := distanceAndRank(weightedMatrix, ideal, negativeIdeal)

	fmt.Println("Scores:", scores)
}
