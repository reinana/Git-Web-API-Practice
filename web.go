package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
func writeError (w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}


// ========/api/hello================
func helloHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")

	// レスポンス用のマップを作成
	response := map[string]string {
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
	response := map[string][]string {
		"categories": categories,
	}
	writeJSON(w, http.StatusOK, response)
}
func main() {
	fmt.Println("Starting the server")

	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/categories", categoriesHandler)

	http.ListenAndServe(":8000", nil)
}