package business

import (
	"fmt"
	"strconv"
	"time"

	"github.com/grupo-sbf/go_Sandbox/01_getP2k/integration/p2k"
	"github.com/vjeantet/jodaTime"
)

// Obter notas fiscais de fornecedor no portal e enviar para o stock-invoice-consumer
// Roda pelo schedule do kubernetes
func ProducerP2k() {
	fmt.Println("GET Docs", jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()))

	ReceiptList, err := p2k.GetStoreReceipt()

	if err != nil {
		fmt.Println("(02)Erro ao pesquisar notas no Portal", err)
	}

	for i, docs := range ReceiptList {
		fmt.Printf("NF[%d]: %v / Qtd Itens: %d\r\n", i, docs.Number, len(docs.Items))
		fmt.Println(docs)
	}
	fmt.Println("")

	fmt.Println("Finalizada pesquisa no P2K..: ", jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()))

	fmt.Println("Notas carregadas......: ", strconv.Itoa(len(ReceiptList)))
	fmt.Println("Processo finalizado...: ", jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()))
}
