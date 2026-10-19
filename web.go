package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ===============共通の処理================
// ステータスコードとデータを受け取りJSONを返す
func writeJSON(w http.ResponseWriter, status int, data any) {
	// Content-Type ヘッダーをapplication/jsonに設定
	w.Header().Set("Content-type", "application/json")
	//CORS
	w.Header().Set("Access-Control-Allow-Origin", "*") // 不要？
	// status
	w.WriteHeader(status)

	// マップをJSONにエンコードしてレスポンスとして送信
	json.NewEncoder(w).Encode(data)

}

// エラーメッセージを返す
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// ========/api/hello================
func helloHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")

	// レスポンス用のマップを作成
	response := map[string]string{
		"message": "Hello " + name,
	}
	writeJSON(w, http.StatusOK, response) // ←共通の関数呼び出しへ変更

	// Content-Type ヘッダーをapplication/jsonに設定 ←共通の処理へ
	// w.Header().Set("Content-type", "application/json")

	// マップをJSONにエンコードしてレスポンスとして送信h←共通の処理へ
	// json.NewEncoder(w).Encode(response)
}

// ========/api/categories =============
// カテゴリの配列を含むJSONを返す
// {
// 	"categories": ["本", "映画", "音楽"]
// }

// type Category struct {

// }
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories := []string{"本", "映画", "音楽"}
	response := map[string][]string{
		"categories": categories,
	}
	writeJSON(w, http.StatusOK, response)
}

// ===============/api/caluculator=============
// レスポンスの構造体
type CalculatorResponse struct {
	Operator string  `json:"operator"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Result   float64 `json:"result"`
}

// Goは大文字だと別のパッケージからアクセスできるため、
// キーは大文字でないといけないがJSONは小文字がいい

// 表記ゆれを揃える エラー防止
func normalizeOperator(o string) (string, bool) {
	switch o {
	case "+", " ", "add":
		return "+", true
	case "-", "sub":
		return "-", true
	case "*", "mul":
		return "*", true
	case "/", "÷", "div":
		return "/", true
	default:
		return "", false
	}
}

// ＋はスペースの意味になるので%2B
func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// 	writeJSON(w, http.StatusOK, map[string]string{
	// 	"o": query.Get("o"),
	// 	"x": query.Get("x"),
	// 	"y": query.Get("y"),
	// })

	op, ok := normalizeOperator(query.Get("o"))
	if !ok {
		writeError(w, http.StatusBadRequest, "パラメータには+,-,*,/のいずれかを使用してください")
	}

	// 数値に変換するにはstrconv
	x, err := strconv.ParseFloat(query.Get("x"), 64)

	if err != nil {
		http.Error(w, "xは数値にしてください", http.StatusBadRequest)
		return
	}
	y, err := strconv.ParseFloat(query.Get("y"), 64)
	if err != nil {
		http.Error(w, "yは数値にしてください", http.StatusBadRequest)
		return
	}

	// 計算する
	var result float64
	switch op {
	case "+":
		result = x + y
	case "-":
		result = x - y
	case "*":
		result = x * y
	case "/":
		if y == 0 {
			writeError(w, http.StatusBadRequest, "0では割れません")
			return
		}
		result = x / y

	}

	writeJSON(w, http.StatusOK, CalculatorResponse{
		Operator: op,
		X:        x,
		Y:        y,
		Result:   result,
	})
}
func main() {
	fmt.Println("Starting the server")

	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/categories", categoriesHandler)
	http.HandleFunc("/api/calculator", calculatorHandler)

	http.ListenAndServe(":8000", nil)
}
