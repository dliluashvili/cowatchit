package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coder/websocket"
)

// Http Response
type Response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Http SendJsonData
func SendJson(w http.ResponseWriter, response *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}

type WSResponse struct {
	Type  string  `json:"type"`
	Data  any     `json:"data,omitempty"`
	Error *string `json:"error,omitempty"`
}

func SendWSResponse(ctx context.Context, conn *websocket.Conn, response *WSResponse) error {
	data, err := json.Marshal(response)

	if err != nil {
		fmt.Println("err", err)
		return err
	}

	return conn.Write(ctx, websocket.MessageText, data)
}

func SendWSError(ctx context.Context, conn *websocket.Conn, errorMsg string) error {
	response := &WSResponse{
		Type:  "ERROR",
		Error: &errorMsg,
	}

	return SendWSResponse(ctx, conn, response)
}
