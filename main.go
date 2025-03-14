package main

import (
	"fmt"
	"log"
	"ore-no-db/data"
)

func main() {
	dbPath := "db"
	tableName := "records"
	action := "seed"

	switch action {
	case "seed":
		// データベースに初期データを挿入
		data.Seed(dbPath, tableName, 1000)
		fmt.Println("初期データが正常に挿入されました。")
	case "insert":
		// 単一のレコードを挿入
		record := data.Record{Name: "ore daesu", Email: "ore@mail.com"}
		err := Insert(dbPath, tableName, record)
		if err != nil {
			log.Fatalf("データ挿入に失敗しました: %v", err)
		}

	case "select":
		// レコードを取得
		id := "1" // 取得するIDを指定
		_, err := RecordSelect(dbPath, tableName, id)
		if err != nil {
			log.Fatalf("レコード取得に失敗しました: %v", err)
		}

	case "bulk":
		records := []data.Record{
			{Name: "ore daesu2", Email: "ore2@mail.com"},
			{Name: "john doe", Email: "john@mail.com"},
			{Name: "jane smith", Email: "jane@mail.com"},
		}
		// 複数のレコードを挿入
		err := BulkInsert(dbPath, tableName, records)
		if err != nil {
			log.Fatalf("複数レコード挿入に失敗しました: %v", err)
		}
	case "delete":
		// レコードを削除
		id := "1" // 削除するIDを指定
		err := data.RecordDelete(dbPath, tableName, id)
		if err != nil {
			log.Fatalf("レコード削除に失敗しました: %v", err)
		}
		fmt.Println("レコードが正常に削除されました。")
	case "update":
		// レコードを更新
		id := "1" // 更新するIDを指定
		record := data.Record{Name: "ore daesu", Email: "update@mail.com"}
		err := data.RecordUpdate(dbPath, tableName, id, record)
		if err != nil {
			log.Fatalf("レコード更新に失敗しました: %v", err)
		}
		fmt.Println("レコードが正常に更新されました。", record)
	default:
		fmt.Println("アクションが指定されていません。")
	}
}

// Insert は単一のレコードを挿入する
func Insert(dbPath, tableName string, record data.Record) error {
	err := data.Insert(dbPath, tableName, record)
	if err != nil {
		return err
	}
	fmt.Println("レコードが正常に保存されました。")
	return nil
}

// RecordSelect はIDでレコードを取得する
func RecordSelect(dbPath, tableName, id string) (data.Record, error) {
	record, err := data.RecordSelect(dbPath, tableName, id)
	if err != nil {
		return data.Record{}, err
	}
	fmt.Println("取得したレコード:", record)
	return record, nil
}

// BulkInsert は複数のレコードを挿入する
func BulkInsert(dbPath, tableName string, records []data.Record) error {
	err := data.BulkInsert(dbPath, tableName, records)
	if err != nil {
		return err
	}
	fmt.Println("複数のレコードが正常に保存されました。")
	return nil
}
