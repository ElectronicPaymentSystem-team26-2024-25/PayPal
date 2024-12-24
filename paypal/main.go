package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"paypay.xws.com/paypal/handler"
	"paypay.xws.com/paypal/model"
	"paypay.xws.com/paypal/repo"
	"paypay.xws.com/paypal/service"
)

func initDB() *gorm.DB {
	dsn := "user=postgres password=super dbname=PspDB host=localhost port=5432 sslmode=disable options='-c search_path=paypal'"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}
	err = database.AutoMigrate(&model.Client{}, &model.Order{})
	if err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}
	return database
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	initEnv()

	serverTLSCert, err := tls.LoadX509KeyPair(os.Getenv("CERT_PATH"), os.Getenv("CERT_KEY_PATH"))
	if err != nil {
		log.Fatalf("Error loading certificate and key file: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
	}

	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}

	startServer(database, tlsConfig)
}

func prepareClient(db *gorm.DB) *handler.ClientHandler {
	clientRepo := &repo.ClientRepo{DbConnection: db}
	clientService := &service.ClientService{Repo: clientRepo}
	clientHandler := &handler.ClientHandler{Service: clientService}
	return clientHandler
}

func preparePayment(db *gorm.DB) *handler.PaymentHandler {
	clientRepo := &repo.ClientRepo{DbConnection: db}
	orderRepo := &repo.OrderRepo{DbConnection: db}
	paymentService := &service.PaymentService{
		ClientRepo: clientRepo,
		OrderRepo:  orderRepo}
	paymentHandler := &handler.PaymentHandler{Service: paymentService}
	return paymentHandler
}

func startServer(db *gorm.DB, tlsConfig *tls.Config) {
	clientHandler := prepareClient(db)
	paymentHandler := preparePayment(db)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(os.Getenv("BASE_API_ENDPOINT")+"/client", clientHandler.CreateClient).Methods("POST")
	router.HandleFunc(os.Getenv("BASE_API_ENDPOINT")+"/payment", paymentHandler.ProcessPayment).Methods("POST")
	router.HandleFunc(os.Getenv("BASE_API_ENDPOINT")+"/payment/{orderId}", paymentHandler.CapturePayment).Methods("PUT")

	server := &http.Server{
		Addr:      os.Getenv("SERVICE_PORT"),
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	println("HTTPS server starting on", os.Getenv("SERVICE_PORT"))
	log.Fatal(server.ListenAndServeTLS("", ""))
}
