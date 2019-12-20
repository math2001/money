package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/math2001/money/db"
)

type Payment map[string]interface{}

func (api *API) AddPayment(u *db.User, p map[string]interface{}) error {

	payment := Payment(p)

	content, err := u.Load("/payments")
	var payments []Payment
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading existing payments: %s", err)
	}

	err = json.Unmarshal(content, &payments)
	if err != nil {
		return fmt.Errorf("parsing existing payments: %s", err)
	}

	payments = append(payments, payment)

	content, err = json.Marshal(payments)
	if err != nil {
		return fmt.Errorf("json encoding payments: %s", err)
	}

	err = u.Save("/payments", content)
	if err != nil {
		return fmt.Errorf("saving payments to db: %s", err)
	}

	return nil
}
