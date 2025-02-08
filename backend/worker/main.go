package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Name string `json:"name"`
}

// Worker function untuk memproses batch
func worker(id int, jobs <-chan []JsonResponse, results chan<- []JsonResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	for batch := range jobs {
		fmt.Printf("👷 Worker %d memproses %d data...\n", id, len(batch))
		results <- batch // Kirim hasil proses ke channel hasil
	}
}

func main() {
	router := gin.Default()

	router.POST("/procces", func(ctx *gin.Context) {
		// Ambil filename dari request JSON
		var request struct {
			Filename string `json:"filename"`
		}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(400, gin.H{"message": "Invalid request", "error": err.Error()})
			return
		}

		// Path file dari folder uploads
		filePath := "uploads/" + request.Filename

		// Buka file CSV
		csvFile, err := os.Open(filePath)
		if err != nil {
			ctx.JSON(500, gin.H{"message": "Failed to open file", "error": err.Error()})
			return
		}
		defer csvFile.Close()

		// Baca isi CSV
		reader := csv.NewReader(csvFile)
		records, err := reader.ReadAll()
		if err != nil {
			ctx.JSON(500, gin.H{"message": "Failed to read CSV", "error": err.Error()})
			return
		}

		// Konfigurasi worker pool
		numWorkers := 5                          // Jumlah worker
		batchSize := 5_000                       // Ukuran batch
		jobs := make(chan []JsonResponse, 10)    // Channel untuk batch data
		results := make(chan []JsonResponse, 10) // Channel untuk hasil
		var wg sync.WaitGroup

		// Jalankan worker pool
		for i := 1; i <= numWorkers; i++ {
			wg.Add(1)
			go worker(i, jobs, results, &wg)
		}

		// Kirim data ke channel dalam batch
		go func() {
			var batch []JsonResponse
			for _, record := range records {
				batch = append(batch, JsonResponse{Name: record[0]})

				// Jika batch penuh, kirim ke channel jobs
				if len(batch) == batchSize {
					jobs <- batch
					batch = nil // Reset batch
				}
			}

			// Kirim sisa batch jika masih ada
			if len(batch) > 0 {
				jobs <- batch
			}

			close(jobs) // Tutup channel setelah semua data dikirim
		}()

		// Tunggu worker selesai
		go func() {
			wg.Wait()
			close(results) // Tutup channel hasil setelah semua worker selesai
		}()

		// Gabungkan semua hasil batch
		var allData []JsonResponse
		for batch := range results {
			allData = append(allData, batch...)
		}

		// Simpan ke file JSON
		jsonData, _ := json.MarshalIndent(allData, "", "  ")
		_ = os.WriteFile("output/all_data.json", jsonData, 0644)

		fmt.Println("🎉 Semua data selesai diproses!")

		ctx.JSON(200, gin.H{
			"message": "OK",
			"file":    request.Filename,
		})
	})

	router.Run(":8080")
}
