package p2k

import (
	"fmt"
)

// Get store document from ID
func GetDistributionCenters() (distributionCenters map[string]string, err error) {

	sqlQueryParams := `select document, id from distribution_center`

	rows, err := dbInstance.Query(sqlQueryParams)
	if err != nil {
		err = fmt.Errorf("failed to get distribution centers, err: %s", err.Error())
		return distributionCenters, err
	}
	defer rows.Close()

	distributionCenters = map[string]string{}

	for rows.Next() {
		var document, id string
		if err = rows.Scan(&document, &id); err != nil {
			err = fmt.Errorf("failed to fetch scan distribution centers row, err: %s", err.Error())
			return distributionCenters, err
		}
		distributionCenters[id] = document
	}

	return distributionCenters, nil
}
