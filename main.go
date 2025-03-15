package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"ore-no-db/data"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// ログファイルのパス
const logFilePath = "./ore-no-log.log"

// ログファイルのオープン
func openLogFile() (*os.File, error) {
	// ファイルをオープン
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// ログ書き込み
func writeLog(url string, statusCode int) {
	file, err := openLogFile()
	if err != nil {
		log.Fatalf("ログファイルオープン失敗: %v", err)
	}
	defer file.Close()

	// ログのフォーマット
	logMessage := fmt.Sprintf("%s - [%s] \"%s\" %d\n", time.Now().Format(time.RFC3339), url, "GET", statusCode)
	fmt.Println(logMessage)
	file.WriteString(logMessage)
}

func main() {
	// ルーターの作成
	r := mux.NewRouter()

	// エンドポイントの設定
	r.HandleFunc("/api/records/seed/{count}", SeedRecords).Methods("POST")
	r.HandleFunc("/api/records", InsertRecord).Methods("POST")
	r.HandleFunc("/api/records/{id}", SelectRecord).Methods("GET")
	r.HandleFunc("/api/records/bulk", BulkInsertRecords).Methods("POST")
	r.HandleFunc("/api/records/{id}", DeleteRecord).Methods("DELETE")
	r.HandleFunc("/api/records/{id}", UpdateRecord).Methods("PUT")

	// CORS設定
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// SO_REUSEPORT を有効にしたリスナーを作成
	listener, err := reusePortListener("tcp", ":8080")
	if err != nil {
		log.Fatalf("ポートのリスニングに失敗しました: %v", err)
	}

	// HTTP サーバーを起動
	log.Println("Server is running on port 8080...")
	log.Fatal(http.Serve(listener, cors(r)))
}

// SO_REUSEPORT を設定したリスナーを作成
func reusePortListener(network, address string) (net.Listener, error) {
	// ソケットを作成
	addr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	// `SO_REUSEPORT` を設定
	syscallFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return nil, err
	}

	// SO_REUSEPORT オプションをセット
	if err := syscall.SetsockoptInt(syscallFD, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1); err != nil {
		return nil, err
	}

	// アドレスをバインド
	sockAddr := &syscall.SockaddrInet4{Port: addr.Port}
	copy(sockAddr.Addr[:], addr.IP.To4())

	if err := syscall.Bind(syscallFD, sockAddr); err != nil {
		return nil, err
	}

	// リスニング開始
	if err := syscall.Listen(syscallFD, syscall.SOMAXCONN); err != nil {
		return nil, err
	}

	// ファイルディスクリプタを Go の `net.Listener` に変換
	file := os.NewFile(uintptr(syscallFD), "")
	listener, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// SeedRecords はデータベースに初期データを挿入する
func SeedRecords(w http.ResponseWriter, r *http.Request) {
	// パラメータから count を取得
	vars := mux.Vars(r)
	countStr := vars["count"]
	count, err := strconv.Atoi(countStr) // 文字列を整数に変換
	if err != nil {
		http.Error(w, "無効なデータ数", http.StatusBadRequest)
		return
	}

	// データベースパスとテーブル名
	dbPath := "db"
	tableName := "records"
	// 初期データの数を決めて挿入
	data.Seed(dbPath, tableName, count)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "200", "message": fmt.Sprintf("%d 件の初期データが正常に挿入されました", count)})
	writeLog(r.URL.Path, http.StatusOK)
}

// InsertRecord は単一のレコードを挿入する
func InsertRecord(w http.ResponseWriter, r *http.Request) {
	var record data.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "無効なデータ", http.StatusBadRequest)
		return
	}

	dbPath := "db"
	tableName := "records"
	err := data.Insert(dbPath, tableName, record)
	if err != nil {
		http.Error(w, "データ挿入に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)
	writeLog(r.URL.Path, http.StatusCreated)
}

// SelectRecord はIDでレコードを取得する
func SelectRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	dbPath := "db"
	tableName := "records"
	writeLog(r.URL.Path, http.StatusOK)

	record, err := data.RecordSelect(dbPath, tableName, id)
	if err != nil {
		http.Error(w, "レコードが見つかりません", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

// BulkInsertRecords は複数のレコードを挿入する
func BulkInsertRecords(w http.ResponseWriter, r *http.Request) {
	var records []data.Record
	if err := json.NewDecoder(r.Body).Decode(&records); err != nil {
		http.Error(w, "無効なデータ", http.StatusBadRequest)
		return
	}

	dbPath := "db"
	tableName := "records"
	err := data.BulkInsert(dbPath, tableName, records)
	if err != nil {
		http.Error(w, "複数レコードの挿入に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "複数レコードが正常に挿入されました"})
	writeLog(r.URL.Path, http.StatusCreated)
}

// DeleteRecord はIDでレコードを削除する
func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	dbPath := "db"
	tableName := "records"

	err := data.RecordDelete(dbPath, tableName, id)
	if err != nil {
		http.Error(w, "レコード削除に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // レコード削除成功
}

// UpdateRecord はIDでレコードを更新する
func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var record data.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "無効なデータ", http.StatusBadRequest)
		return
	}

	dbPath := "db"
	tableName := "records"
	err := data.RecordUpdate(dbPath, tableName, id, record)
	if err != nil {
		http.Error(w, "レコード更新に失敗しました", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
	writeLog(r.URL.Path, http.StatusOK)
}
