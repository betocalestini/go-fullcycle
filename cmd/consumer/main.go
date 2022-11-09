package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/betocalestini/go-fullcyle/internal/order/infra/database"
	"github.com/betocalestini/go-fullcyle/internal/order/usecase"
	"github.com/betocalestini/go-fullcyle/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sql.Open("postgres", "host=postgres port=5432 user=userPostgres password=passPostgres dbname=fullcycle sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := database.NewOrderRepository(db)
	uc := usecase.CalculateFinalPriceUseCase{OrderRepository: repository}

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	fmt.Println("rodando 1")

	defer ch.Close()

	out := make(chan amqp.Delivery) //chanel
	go rabbitmq.Consume(ch, out)    //T2

	for msg := range out {
		var inputDTO usecase.OrderInpuDTO
		err := json.Unmarshal(msg.Body, &inputDTO)
		if err != nil {
			fmt.Println(err)
		}
		outputDTO, err := uc.Execute(inputDTO)
		if err != nil {
			fmt.Println(err)
		}
		msg.Ack(false)
		fmt.Println(outputDTO) //T1
		time.Sleep(100 * time.Millisecond)
	}

}
