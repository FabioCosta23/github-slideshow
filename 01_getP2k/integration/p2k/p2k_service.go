package p2k

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
)

// Obter notas fiscais de fornecedor no portal e enviar para o stock-invoice-consumer
// Roda pelo schedule do kubernetes
func LoadReceiptP2K() {
	fmt.Println("GET Docs", jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()))

	distributionCenters, err := GetDistributionCenters()
	if err != nil {
		fmt.Println(ErrPrefix, err)
	}

	ReceiptList, err := GetStoreReceipt(distributionCenters)
	if err != nil {
		fmt.Println(ErrPrefix, err)
	}

	// for i, docs := range ReceiptList {
	// 	fmt.Printf("NF[%d]: %v / Qtd Itens: %d\r\n", i, docs.Number, len(docs.Items))
	// 	fmt.Println(docs)
	// }

	// fmt.Println("\nFinalizada pesquisa no P2K..: ", jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()))
	// fmt.Println("Notas carregadas......: ", strconv.Itoa(len(ReceiptList)))
	// fmt.Println("Processo finalizado...: ", jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()))

	err = sendReceiptStockConsumer(ReceiptList)
	if err != nil {
		fmt.Println(ErrPrefix, err)
	}
}

func sendReceiptStockConsumer(receiptList []Receipt) error {

	writer := &kafka.Writer{
		Addr:         kafka.TCP(GetEnv("KAFKA_HOST")),
		Topic:        GetEnv("KAFKA_INVOICE_TOPIC_NAME"),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	if GetEnv("ENVIRONMENT") != "local" {
		logrus.Info("using SASL Plain in topic connection")

		t := writer.Transport.(*kafka.Transport)
		t.SASL = plain.Mechanism{
			Username: GetEnv("KAFKA_USERNAME"),
			Password: GetEnv("KAFKA_PASSWORD"),
		}
		t.TLS = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	for _, receipt := range receiptList {
		objJson, err := json.Marshal(receipt)
		if err != nil {
			return err
		}
		recepitKey := []byte(fmt.Sprintf("%s_%s_%s_%d", GetEnv("SENDER_ID"), receipt.DistributionCenterOrigin, strings.TrimSpace(receipt.Series), receipt.Number))

		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   recepitKey,
				Value: []byte(objJson),
			},
		)
		fmt.Println("JSON enviado para o topico: ", recepitKey)
	}
	if err := writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	return nil
}

// const (
// 	DISTRIBUTION_CENTER_ID int = 3333
// 	STOCK_TYPE_ID          int = 1
// )

// var (
// 	gerestRepository GerestRepository
// )

// func init() {
// 	gerestRepository = &GerestRepositoryImpl{}
// }
