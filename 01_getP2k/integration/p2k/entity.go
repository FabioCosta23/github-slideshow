package p2k

type Receipt struct {
	//ID                            int           `json:"id"`
	BusinessUnitID                int16         `json:"businessUnit"`
	DistributionCenterOrigin      string        `json:"distributionCenterOrigin"`
	DistributionCenterDestination string        `json:"distributionCenterDestination"`
	CNPJOrigin                    string        `json:"documentOrigin"`
	CNPJDestination               string        `json:"documentDestination"`
	Series                        string        `json:"series"`
	Number                        int32         `json:"number"`
	Type                          string        `json:"type"`
	IssueDate                     string        `json:"issueDate"`
	IssuerDate                    string        `json:"issuerDate"`
	Status                        string        `json:"status"`
	Amount                        float32       `json:"amount"`
	IssuerKey                     string        `json:"issuerKey"`
	SenderID                      string        `json:"invoiceSenderId"`
	MovementType                  string        `json:"-"`
	Items                         []ReceiptItem `json:"items"`
}

type ReceiptItem struct {
	//ReceiptID     int32
	ItemID   int16   `json:"id"`
	Sku      string  `json:"sku"`
	UnitCost float32 `json:"unitCost"`
	Quantity int32   `json:"quantity"`
}

type ReceiptGet struct {
	UniqueKey                     string
	BusinessUnitID                int16
	DistributionCenterOrigin      string
	DistributionCenterDestination string
	CNPJOrigin                    string
	CNPJDestination               string
	Series                        string
	Number                        int32
	Type                          string
	IssueDate                     string
	IssuerDate                    string
	Status                        string
	Amount                        float32
	IssuerKey                     string
	SenderID                      string
	MovementType                  string
	ItemID                        int16
	Sku                           string
	UnitCost                      float32
	Quantity                      int32
}
