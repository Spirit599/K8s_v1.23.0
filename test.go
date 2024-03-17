package main

import (
	"fmt"
	"math"
)

// Fungsi untuk menghitung rata-rata dari sebuah slice float64
func mean(data []float64) float64 {
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	return sum / float64(len(data))
}

// Fungsi untuk menghitung selisih data dengan nilai rata-rata
func subtractMean(data []float64) []float64 {
	meanValue := mean(data)
	result := make([]float64, len(data))
	for i, val := range data {
		result[i] = val - meanValue
	}
	return result
}

// Fungsi untuk menghitung varians dari sebuah slice float64
func variance(data []float64) float64 {
	meanValue := mean(data)
	sumSquaredDiff := 0.0
	for _, val := range data {
		diff := val - meanValue
		sumSquaredDiff += diff * diff
	}
	return sumSquaredDiff / float64(len(data)-1)
}

// Fungsi untuk menghitung autokorelasi
func autocorrelation(data []float64, lag int) float64 {
	n := len(data)
	meanData := mean(data)
	varianceData := variance(data)

	var numerator, denominator float64
	for i := 0; i < n-lag; i++ {
		numerator += (data[i] - meanData) * (data[i+lag] - meanData)
	}
	denominator = float64(n-lag) * varianceData

	return numerator / denominator
}

// Fungsi untuk menghitung koefisien autoregressive (AR)
func calculateARCoefficients(data []float64, p int) []float64 {
	coefficients := make([]float64, p)

	for i := 1; i <= p; i++ {
		coefficients[i-1] = autocorrelation(data, i)
	}

	return coefficients
}

// Fungsi untuk menghitung residual dari model ARIMA
func calculateResidual(data []float64, arCoefficients []float64) []float64 {
	n := len(data)
	p := len(arCoefficients)
	residual := make([]float64, n-p)

	copy(residual, data[:p])

	for i := p; i < n; i++ {
		prediction := float64(0)
		for j := 0; j < p; j++ {
			prediction += arCoefficients[j] * data[i-j-1]
		}
		residual = append(residual, data[i]-prediction)
	}

	return residual
}

// Fungsi untuk menghitung moving average (MA)
func calculateMACoefficients(residual []float64, q int) []float64 {
	coefficients := make([]float64, q)

	for i := 1; i <= q; i++ {
		coefficients[i-1] = autocorrelation(residual, i)
	}

	return coefficients
}

// Fungsi untuk melakukan prediksi menggunakan model ARIMA
func predictARIMA(data []float64, p, d, q, predictionLength int) []float64 {
	// Lakukan differencing jika d > 0
	for i := 0; i < d; i++ {
		diffData := make([]float64, len(data)-1)
		for j := 1; j < len(data); j++ {
			diffData[j-1] = data[j] - data[j-1]
		}
		data = diffData
	}

	arCoefficients := calculateARCoefficients(data, p)

	// Lakukan residual calculation
	residual := calculateResidual(data, arCoefficients)

	// Lakukan moving average jika q > 0
	if q > 0 {
		maCoefficients := calculateMACoefficients(residual, q)

		for i := 0; i < predictionLength; i++ {
			prediction := float64(0)
			for j := 1; j <= q; j++ {
				if len(residual) >= j {
					prediction += maCoefficients[j-1] * residual[len(residual)-j]
				}
			}
			residual = append(residual, prediction)
		}
	}

	// Lakukan inverse differencing jika d > 0
	if d > 0 {
		for i := 0; i < d; i++ {
			inverseDiffData := make([]float64, len(residual))
			inverseDiffData[0] = data[len(data)-1]
			for j := 1; j < len(residual); j++ {
				inverseDiffData[j] = inverseDiffData[j-1] + residual[j-1]
			}
			residual = inverseDiffData
		}
	}

	// Lakukan inverse AR jika p > 0
	if p > 0 {
		residual = inverseARIMA(residual, p, d, arCoefficients)
	}

	return residual[len(residual)-predictionLength:]
}

// Fungsi untuk melakukan inverse ARIMA
func inverseARIMA(residual []float64, p, d int, arCoefficients []float64) []float64 {
	n := len(residual)
	data := make([]float64, n)

	copy(data, residual)

	// Lakukan inverse differencing jika d > 0
	if d > 0 {
		for i := n - 1; i >= d; i-- {
			data[i] = data[i] + data[i-d]
		}
	}

	// Lakukan inverse AR jika p > 0
	if p > 0 {
		for i := n - 1; i >= p; i-- {
			prediction := float64(0)
			for j := 0; j < p; j++ {
				prediction += arCoefficients[j] * data[i-j-1]
			}
			data[i-p] = data[i-p] + prediction
		}
	}

	return data
}

func calculateMAE(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) {
		panic("Length of actual and predicted slices should be the same.")
	}

	sum := 0.0
	for i := 0; i < len(actual); i++ {
		sum += math.Abs(actual[i] - predicted[i])
	}

	return sum / float64(len(actual))
}

// Fungsi untuk menghitung Mean Squared Error (MSE)
func calculateMSE(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) {
		panic("Length of actual and predicted slices should be the same.")
	}

	sum := 0.0
	for i := 0; i < len(actual); i++ {
		diff := actual[i] - predicted[i]
		sum += diff * diff
	}

	return sum / float64(len(actual))
}

// Fungsi untuk menghitung Mean Absolute Percentage Error (MAPE)
func calculateMAPE(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) {
		panic("Length of actual and predicted slices should be the same.")
	}

	sumPercentageError := 0.0
	for i := 0; i < len(actual); i++ {
		percentageError := math.Abs((actual[i] - predicted[i]) / actual[i])
		sumPercentageError += percentageError
	}

	mape := (sumPercentageError / float64(len(actual))) * 100.0
	return mape
}

func main() {
	// Data contoh (gantilah data ini dengan data Anda) (data product perminggu )

	// 	SELECT DATE_FORMAT(DATE_SUB(a.created_at, INTERVAL WEEKDAY(a.created_at) DAY), '%Y-%m-%d') AS start_of_week,
	// 	DATE_FORMAT(DATE_ADD(a.created_at, INTERVAL 6 - WEEKDAY(a.created_at) DAY), '%Y-%m-%d') AS end_of_week,
	// 	SUM(b.qty) AS total_qty
	// FROM transactions a
	// LEFT JOIN transaction_lines b ON a.id = b.transaction_id
	// WHERE b.name = 'Pupuk Kompos (5KG)'
	// AND a.created_at >= '2022-07-01 01:43:26.419'
	// AND a.created_at < NOW()  -- Batas tanggal akhir diganti dengan NOW() untuk data hingga hari ini
	// GROUP BY start_of_week, end_of_week
	// ORDER BY start_of_week ASC;

	data := []float64{

		122,
		113,
		116,
		135,
		90,
		129,
		79,
		118,
		98,
		111,
		153,
		123,
		140,
		144,
		140,
		145,
		73,
		84,
		131,
		137,
		136,
		101,
		135,
		106,
		102,
		128,
		138,
		104,
		116,
		94,
		118,
		126,
		144,
		138,
		109,
		136,
		134,
		121,
		129,
		147,
		113,
		102,
		104,
		70,
		138,
		78,
		86,
		96,
		136,
		89,
		124,
		131,
		93,
		108,
		151,
		116,
		111,
	}

	// Parameter ARIMA (gantilah parameter ini sesuai kebutuhan)
	p := 1 // Order autoregressive (AR)
	d := 0 // Derajat differencing (I)
	q := 1 // Order moving average (MA)

	// Panjang prediksi (dalam minggu) (tolong di bikin jadi params )
	predictionLength := 12 // 12 minggu

	// Lakukan prediksi
	predictions := predictARIMA(data, p, d, q, predictionLength)

	// Tampilkan hasil data asli
	fmt.Println("Data Asli:")
	for i, val := range data {
		fmt.Printf("Minggu %d: %.2f\n", i+1, val)
	}

	// Tampilkan hasil prediksi
	fmt.Println("Hasil Prediksi (Bentuk Asli):")
	for i, pred := range predictions {
		fmt.Printf("Minggu %d: %.f\n", len(data)+i+1, math.Round(pred+data[len(data)-1]))
	}

	// Hitung MAE
	mae := calculateMAE(data[len(data)-predictionLength:], predictions)
	fmt.Printf("Mean Absolute Error (MAE): %.2f\n", mae)

	// Hitung MSE
	mse := calculateMSE(data[len(data)-predictionLength:], predictions)
	fmt.Printf("Mean Squared Error (MSE): %.2f\n", mse)

	// Hitung MAPE
	mape := calculateMAPE(data[len(data)-predictionLength:], predictions)
	fmt.Printf("Mean Absolute Percentage Error (MAPE): %.2f%%\n", mape)

}
