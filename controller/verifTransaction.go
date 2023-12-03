package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fenriz07/Golang-Transbank-WebPay-Rest/pkg/transaction"
)

// LogEntry representa una entrada de registro
type LogEntry struct {
	Message     string `json:"message"`
	NumberOrder string `json:"number_order,omitempty"`
	IdSession   string `json:"id_session,omitempty"`
	Status     string `json:"status"`
	Amount     int    `json:"amount"`
	BuyOrder   string `json:"buy_order"`
	SessionID  string `json:"session_id"`
	AccountingDate     string    `json:"accounting_date"`
	TransactionDate    time.Time `json:"transaction_date"`
	PaymentTypeCode    string    `json:"payment_type_code"`
	CardDetail struct {
		CardNumber string `json:"card_number"`
	} `json:"card_detail"`
	AuthorizationCode  string    `json:"authorization_code"`

}

func VerifTransaction(w http.ResponseWriter, r *http.Request) {
	logEntries := []LogEntry{}

	var token string = ""
	var numberOrder string = ""
	var idSession string = ""

	canceledToken := r.FormValue("TBK_TOKEN")
	
	if len(canceledToken) != 0 {
		token = canceledToken
		numberOrder = r.FormValue("TBK_ORDEN_COMPRA")
		idSession = r.FormValue("TBK_ID_SESION")

		logEntries = append(logEntries, LogEntry{
			Message: "estado de transaccion",
			Status: "ANULADO",
			BuyOrder: numberOrder,
			SessionID: idSession,
		})
	} else {
		token = r.FormValue("token_ws")
	}

	estado, err := transaction.Commit(token)
	if err != nil {
		logEntries = append(logEntries, LogEntry{Message: fmt.Sprintf("GetStatus Error: %s", err.Error())})
	}
	
	if estado.Status == "AUTHORIZED" {
		logEntries = append(logEntries, LogEntry{
			Message: "estado de transaccion",
			Status: "RECHAZADO",
			Amount: estado.Amount,
			BuyOrder: estado.BuyOrder,
			SessionID: estado.SessionID,
			AccountingDate: estado.AccountingDate,
			TransactionDate: estado.TransactionDate,
			PaymentTypeCode: estado.PaymentTypeCode,
			CardDetail: estado.CardDetail,
			AuthorizationCode: estado.AuthorizationCode,
		})
	}
	
	if estado.Status == "FAILED" {
		logEntries = append(logEntries, LogEntry{
			Message: "estado de transaccion",
			Status: "RECHAZADO",
			Amount: estado.Amount,
			BuyOrder: estado.BuyOrder,
			SessionID: estado.SessionID,
			AccountingDate: estado.AccountingDate,
			TransactionDate: estado.TransactionDate,
			PaymentTypeCode: estado.PaymentTypeCode,
			CardDetail: estado.CardDetail,
			AuthorizationCode: estado.AuthorizationCode,
		})
	}

	// Convertir las entradas de registro a JSON
	logJSON, err := json.MarshalIndent(logEntries, "", "  ")
	if err != nil {
		log.Println("Error al convertir el registro a JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Imprimir el registro en formato JSON en lugar de la consola
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(logJSON)
}
