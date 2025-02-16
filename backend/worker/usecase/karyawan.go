package usecase

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"worker/database"
	"worker/models"
)

func CreateKaryawan(karyawan []models.Karyawan, batchSize int) error {
	db := database.DB
	tx := db.Begin()

	err := tx.CreateInBatches(&karyawan, batchSize).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func WorkerPool(id int, batch <-chan []models.Karyawan, result chan<- []models.Karyawan, wg *sync.WaitGroup, batchSize int) {
	defer wg.Done()

	for row := range batch {
		fmt.Printf("Worker %d memproses %d data...\n", id, len(row))

		if err := CreateKaryawan(row, batchSize); err != nil {
			fmt.Printf("Worker %d gagal insert data: %v\n", id, err)
		} else {
			result <- row
		}
	}
}

func ProcessCsv(filename string) ([]models.Karyawan, error) {
	path := "uploads/" + filename
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %v", err)
	}
	defer file.Close()

	csvNewReader := csv.NewReader(file)
	read, err := csvNewReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca CSV: %v", err)
	}

	batch := 200
	worker := 4
	var wg sync.WaitGroup
	batchChannel := make(chan []models.Karyawan, 10)
	resultChannel := make(chan []models.Karyawan, 10)

	for i := 0; i < worker; i++ {
		wg.Add(1)
		go WorkerPool(i, batchChannel, resultChannel, &wg, batch)
	}

	go func() {
		var data []models.Karyawan
		for _, v := range read {
			salaryInt, _ := strconv.Atoi(v[3])
			data = append(data, models.Karyawan{
				Name:       v[0],
				Position:   v[1],
				Department: v[2],
				Salary:     salaryInt,
			})

			if len(data) >= batch {
				batchChannel <- data
				data = nil
			}
		}
		if len(data) > 0 {
			batchChannel <- data
		}
		close(batchChannel)
	}()

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	var finalData []models.Karyawan
	for batch := range resultChannel {
		finalData = append(finalData, batch...)
	}

	return finalData, nil
}
