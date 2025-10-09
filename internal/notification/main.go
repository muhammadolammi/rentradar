package notification

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/muhammadolammi/rentradar/internal/database"
)

func SendNotifications() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading env. err:" + err.Error())
		return
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Println("empty dbURL")
		return
	}
	rabbitmq_url := os.Getenv("RABBITMQ_URL")
	if rabbitmq_url == "" {
		log.Println("there is no rabbitmq_url provided kindly provide a rabbitmq_url")
		return
	}
	smtp_server := os.Getenv("SMTP_SERVER")
	if smtp_server == "" {
		log.Println("there is no smtp_server provided kindly provide a smtp_server")
		return
	}
	smtp_username := os.Getenv("SMTP_USERNAME")
	if smtp_username == "" {
		log.Println("there is no smtp_username provided kindly provide a smtp_username")
		return
	}
	smtp_password := os.Getenv("SMTP_PASSWORD")
	if smtp_password == "" {
		log.Println("there is no smtp_password provided kindly provide a smtp_password")
		return
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)
	smtpModel := SMTPModel{
		Server:   smtp_server,
		Password: smtp_password,
		UserName: smtp_username,
	}
	config := &Config{
		DB:        dbQueries,
		SMTPModel: smtpModel,
	}
	//  This will start 3 pool to consume messages on the rabbitmq queue (rent_radar) and send notification on each message with retries.
	config.startWorkerPool(3, "rent_radar")
}

func PublishNotifications() {}
func main() {
	SendNotifications()
	PublishNotifications()
}
