package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

var spec *openapi3.T

func main() {
	// OpenAPI ファイルのパスを指定
	specPath := "path/to/openapi/spec.yaml"

	// OpenAPI ファイルの読み込み
	specBytes, err := ioutil.ReadFile(specPath)
	if err != nil {
		log.Fatalf("failed to read spec: %v", err)
	}

	// OpenAPI スキーマのパース
	spec, err = openapi3.NewLoader().LoadFromData(specBytes)
	if err != nil {
		log.Fatalf("failed to parse spec: %v", err)
	}

	// ルーターの設定
	r, err := router.NewRouter(context.Background(), spec)
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	// ルートのハンドラーの設定
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// OpenAPI 仕様に従ってリクエストをルーティングし、ハンドラーを設定
		route, pathParams, err := r.FindRoute(r.Method, r.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// リクエストのバリデーション
		requestValidationInput := &router.RequestValidationInput{
			Request:    r,
			Route:      route,
			PathParams: pathParams,
		}
		if err := router.ValidateRequest(context.Background(), requestValidationInput); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// ハンドラーの実行
		route.Handler.ServeHTTP(w, r)
	})

	// サーバーの起動
	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ハンドラー関数
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// レスポンスの作成
	response := map[string]string{"message": "Hello, world!"}

	// JSON 形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
