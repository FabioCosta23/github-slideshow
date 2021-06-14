package p2k

type ReceiptGet struct {
	UniqueKey                     string
	BusinessUnitID                int16
	DistributionCenterOrigin      string
	DistributionCenterDestination string
	Series                        string
	Number                        int32
	Type                          string
	IssueDate                     string
	IssuerDate                    string
	Status                        string
	Amount                        float32
	IssuerKey                     string
	MovementType                  string
	ItemID                        int16
	Sku                           string
	UnitCost                      float32
	Quantity                      int32
}
