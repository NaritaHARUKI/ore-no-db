// data/increment.go
package data

import (
	"os"
	"path/filepath"
	"strconv"
)

// Increment は自動的に次のIDを取得する
func Increment(tableName string) (string, error) {
	// db/counters ディレクトリが存在しない場合は作成
	counterDirPath := filepath.Join("db", "counters")
	err := os.MkdirAll(counterDirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	// counter.txt ファイルが存在するかチェック
	counterFilePath := filepath.Join(counterDirPath, tableName+"-increment.txt")
	if _, err := os.Stat(counterFilePath); os.IsNotExist(err) {
		// ファイルが存在しなければ、カウンターを 1 から開始
		err := os.WriteFile(counterFilePath, []byte("1"), 0644)
		if err != nil {
			return "", err
		}
		return "1", nil
	}

	// ファイルが存在すればカウンターを読み込んでインクリメント
	data, err := os.ReadFile(counterFilePath)
	if err != nil {
		return "", err
	}

	// カウンターの値を整数に変換
	counter, err := strconv.Atoi(string(data))
	if err != nil {
		return "", err
	}

	// インクリメント
	counter++

	// 新しいカウンターの値をファイルに書き込む
	err = os.WriteFile(counterFilePath, []byte(strconv.Itoa(counter)), 0644)
	if err != nil {
		return "", err
	}

	// インクリメント後のIDを返す
	return strconv.Itoa(counter), nil
}
