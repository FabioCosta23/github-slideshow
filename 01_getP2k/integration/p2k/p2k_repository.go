package p2k

import (
	"database/sql"
	"fmt"
	"io"
	"strconv"
	"time"

	_ "github.com/sijms/go-ora"
	"github.com/vjeantet/jodaTime"
)

// const P2KConnectString = "oracle://app_stock_worker:rkrCuEw3j@172.16.154.105/p2k"

// Get recepits on P2K front-end database
// DistributionCenterDestination, CNPJOrigin and CNPJDestination not available in data origin.
func GetStoreReceipt(distributionCenters map[string]string) ([]Receipt, error) {

	var receiptList []Receipt

	transactionDate, err := strconv.Atoi(jodaTime.Format("YYYYMMdd", time.Now()))
	if err != nil {
		return receiptList, err
	}

	transactionDate = 20210531

	P2KConnectString := fmt.Sprintf("%s://%s:%s@%s/%s", GetEnv("DB_P2K_DRIVER"), GetEnv("DB_P2K_USER"), GetEnv("DB_P2K_PASSWORD"), GetEnv("DB_P2K_HOST"), GetEnv("DB_P2K_NAME"))

	fmt.Println(ErrPrefix, "Conn: ", P2KConnectString)

	sqlQuery := fmt.Sprintf(`select 1 as businessUnit
	                              , trim(to_char(a.codigo_loja, '0000')) as distr_origin
                                  , (case when a.serie_nfe = 0 then 'ECF' else to_char(a.serie_nfe) end) as receipt_series
                                  , (case a.tipo_venda when 9 then a.numero_nfe else a.numero_cupom end) as receipt_number
                                  , (case when a.tipo_autorizacao_fiscal = 'NFC-e' then 65 else 99 end) as type
                                  , to_char(a.data_impr_fechamento_cupom, 'YYYY-MM-DD HH24:MM:SS') as issue_date
                                  , to_char(a.data_impr_fechamento_cupom, 'YYYY-MM-DD HH24:MM:SS') as issuer_date
                                  , b.status_item
                                  , a.valor_total_venda as amount
                                  , a.chave_acesso_nfe as issuerKey
                                  , 'S' as mov_type
                                  , b.num_seq_produto as item_id
                                  , b.codigo_produto as sku
                                  , b.valor_unitario_produto as unit_cost
                                  , b.qtd_vendida as quantity
                             from dbcsi_p2k_cent_prod.p2k_cab_transacao a, dbcsi_p2k_cent_prod.p2k_item_transacao b
                             where a.codigo_loja = b.codigo_loja
                               and a.numero_componente = b.numero_componente
                               and a.data_transacao = b.data_transacao
                               and a.nsu_transacao = b.nsu_transacao
                               and a.tipo_venda in (1,9)
                               and b.status_item = 'V'
			            	   and a.data_transacao = %d
                             order by (case a.tipo_venda when 9 then a.numero_nfe else a.numero_cupom end), b.num_seq_produto `, transactionDate) // order by required

	// Formato com parametros
	//  and a.data_transacao = $1
	//  and a.data_transacao = %d and a.data_impr_fechamento_cupom >= sysdate - (interval '$2' minute)

	fmt.Println(ErrPrefix, "Params: ", transactionDate, GetEnv("RECEIPT_GET_INTERVAL"))

	conn, err := sql.Open("oracle", P2KConnectString)
	if err != nil {
		fmt.Println(ErrPrefix, "Can't open connection:")
		return receiptList, err
	}

	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Println(ErrPrefix, "Can't ping connection:")
		return receiptList, err
	}

	fmt.Println("P2K Successfully connected.")
	stmt, err := conn.Prepare(sqlQuery)
	if err != nil {
		fmt.Println(ErrPrefix, "Can't prepare query:")
		return receiptList, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(ErrPrefix, "Can't create query:")
		return receiptList, err
	}

	defer rows.Close()

	fmt.Println("Querying P2K database...")
	fmt.Println(sqlQuery)

	defer rows.Close()

	var receipt Receipt
	var item ReceiptItem
	var lastUniqueKey string

	for rows.Next() {
		var receiptGet ReceiptGet
		item = ReceiptItem{}

		err := rows.Scan(&receiptGet.BusinessUnitID,
			&receiptGet.DistributionCenterOrigin,
			&receiptGet.Series,
			&receiptGet.Number,
			&receiptGet.Type,
			&receiptGet.IssueDate,
			&receiptGet.IssuerDate,
			&receiptGet.Status,
			&receiptGet.Amount,
			&receiptGet.IssuerKey,
			&receiptGet.MovementType,
			&receiptGet.ItemID,
			&receiptGet.Sku,
			&receiptGet.UnitCost,
			&receiptGet.Quantity)

		if err != nil {
			fmt.Println(ErrPrefix, "ERROR: Nota: ", receiptGet.DistributionCenterOrigin, "/", receiptGet.Series, "/", receiptGet.Number, " | ", err.Error())
			break
		}

		receiptGet.UniqueKey = fmt.Sprintf("%s_%d_%s_%s", receiptGet.DistributionCenterOrigin, receiptGet.Number, receiptGet.Series, receiptGet.Type)
		// fmt.Println(ErrPrefix, "KEY: ", receiptGet.UniqueKey)

		if receiptGet.UniqueKey != lastUniqueKey {
			lastUniqueKey = fmt.Sprintf("%s_%d_%s_%s", receiptGet.DistributionCenterOrigin, receiptGet.Number, receiptGet.Series, receiptGet.Type)
			if receipt.Number > 0 {
				receiptList = append(receiptList, receipt)
				receipt.Items = nil
			}
		}

		receipt.BusinessUnitID = receiptGet.BusinessUnitID
		receipt.DistributionCenterOrigin = receiptGet.DistributionCenterOrigin
		receipt.DistributionCenterDestination = receiptGet.DistributionCenterDestination
		receipt.CNPJOrigin = distributionCenters[receiptGet.DistributionCenterOrigin]
		receipt.Series = receiptGet.Series
		receipt.Number = receiptGet.Number
		receipt.Type = receiptGet.Type
		receipt.IssueDate = receiptGet.IssueDate
		receipt.IssuerDate = receiptGet.IssuerDate
		receipt.Status = receiptGet.Status
		receipt.Amount = receiptGet.Amount
		receipt.IssuerKey = receiptGet.IssuerKey
		receipt.SenderID = GetEnv("SENDER_ID")
		receipt.MovementType = receiptGet.MovementType

		item.ItemID = receiptGet.ItemID
		item.Sku = receiptGet.Sku
		item.UnitCost = receiptGet.UnitCost
		item.Quantity = receiptGet.Quantity

		receipt.Items = append(receipt.Items, item)
	}
	if rows.Err() != nil && rows.Err() != io.EOF {
		if err != nil {
			fmt.Println(ErrPrefix, "Can't fetch row:")
			return receiptList, err
		}
	} else {
		receiptList = append(receiptList, receipt)
	}

	return receiptList, nil
}
