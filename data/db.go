// data/db.go
package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// test Seed data
func Seed(dbPath string, tableName string, count int) error {
	// ダミーデータを生成して保存
	for i := 1; i <= count; i++ {
		record := Record{
			ID:    fmt.Sprintf("%d", i),
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
		}

		// DBの保存場所
		dirPath := filepath.Join(dbPath, record.ID, tableName)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		dataFile := filepath.Join(dirPath, "data.json")
		// JSONファイルにレコードを書き込む
		data, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			return err
		}

		// ファイルにデータを保存
		err = os.WriteFile(dataFile, data, 0644)
		if err != nil {
			return err
		}

		fmt.Println("保存したレコード:", record)
	}

	return nil
}

// Insert は新しいレコードを保存する
func Insert(dbPath string, tableName string, record Record) error {
	// 自動的に次のIDを取得
	id, err := Increment(tableName) // incrementパッケージのIncrementを使う
	if err != nil {
		return err
	}

	// IDをレコードに設定
	record.ID = id

	dirPath := filepath.Join(dbPath, record.ID, tableName)
	dataFile := filepath.Join(dirPath, "data.json")

	fmt.Println("データを保存します:", record)

	// 既にデータが存在している場合はエラーを返す
	if _, err := os.Stat(dataFile); err == nil {
		return errors.New("ID already exists")
	}

	// ディレクトリを作成
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	// JSONデータを保存
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(dataFile, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("データを保存しました:", record)
	return nil
}

// BulkInsert は複数のレコードを保存する
func BulkInsert(dbPath string, tableName string, records []Record) error {
	for _, record := range records {
		// 各レコードを挿入
		err := Insert(dbPath, tableName, record)
		if err != nil {
			// エラーが発生した場合はログを出力して続行
			log.Printf("レコード挿入に失敗しました: %v", err)
		} else {
			// 成功した場合
			fmt.Println("レコードが保存されました:", record)
		}
	}
	return nil
}

// RecordSelect は該当idのレコードを取得する
func RecordSelect(dbPath string, tableName, id string) (Record, error) {
	dirPath := filepath.Join(dbPath, id, tableName)
	dataFile := filepath.Join(dirPath, "data.json")

	var record Record

	// ファイルの存在確認
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return record, errors.New("record not found")
	}

	// データの読み込み
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return record, err
	}

	// JSONデータを構造体に変換
	err = json.Unmarshal(data, &record)
	if err != nil {
		return record, err
	}

	return record, nil
}

// RecordDelete は該当idのレコードを削除する
func RecordDelete(dbPath string, tableName, id string) error {
	dirPath := filepath.Join(dbPath, id, tableName)
	dataFile := filepath.Join(dirPath, "data.json")

	// ファイルの存在確認
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return errors.New("record not found")
	}

	// データの削除
	err := os.Remove(dataFile)
	if err != nil {
		return err
	}

	// ディレクトリの削除
	err = os.RemoveAll(dirPath)
	if err != nil {
		return err
	}

	fmt.Println("レコードを削除しました:", id)
	return nil
}

// RecordUpdate は該当idのレコードを更新する
func RecordUpdate(dbPath string, tableName string, id string, record Record) error {
	dirPath := filepath.Join(dbPath, id, tableName)
	dataFile := filepath.Join(dirPath, "data.json")

	// ファイルの存在確認
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return errors.New("record not found")
	}

	// データの更新
	record.ID = id
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(dataFile, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("レコードを更新しました:", record)
	return nil
}
