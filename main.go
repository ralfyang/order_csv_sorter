package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func processCSV(inputFilePath string) (string, error) {
	// sorter.sh 스크립트를 외부 프로세스로 실행하여 CSV 파일 처리 결과를 받아옵니다.
	outputFilePath := "sorted_result.csv"
	cmd := exec.Command("./sorter.sh", inputFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error running sorter.sh: %v", err)
	}

	// 결과 파일 읽어오기
	file, err := os.Open(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening sorted_result.csv: %v", err)
	}
	defer file.Close()

	// 처리된 결과를 문자열로 읽어옴
	resultData := ""
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("error reading sorted_result.csv: %v", err)
		}
		for _, field := range record {
			resultData += field + ","
		}
		resultData += "\n"
	}

	return resultData, nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only POST method is allowed.")
		return
	}

	file, _, err := r.FormFile("csvfile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to read the uploaded file: %v", err)
		return
	}
	defer file.Close()

	// 업로드된 파일을 임시 파일로 저장
	tempFilePath := "temp_file.csv"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating temp file: %v", err)
		return
	}
	defer tempFile.Close()

	io.Copy(tempFile, file)

	resultData, err := processCSV(tempFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error processing CSV file: %v", err)
		return
	}

	// 처리된 결과를 바로 클라이언트에게 전달
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=sorted_result.csv")
	fmt.Fprint(w, resultData)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "upload.html")
	})

	port := "8000"
	fmt.Printf("Starting server on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Server error: %v", err)
	}
}

