package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/math2001/money/db"
)

type ErrInvalidPayment error

type Payment map[string]interface{}

func (api *API) AddPayment(u *db.User, serializedpayment []byte) error {

	var payment Payment
	if err := json.Unmarshal(serializedpayment, &payment); err != nil {
		return ErrInvalidPayment(fmt.Errorf("unmarshaling json payment: %s", err))
	}

	if err := isValidPayment(payment); err != nil {
		return ErrInvalidPayment(err)
	}

	content, err := u.Load("/payments")
	var payments []Payment
	var patherr *os.PathError

	if errors.As(err, &patherr) && os.IsNotExist(patherr) {
		content = []byte("[]")
	} else if err != nil {
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

func (api *API) ListPayments(u *db.User) ([]Payment, error) {
	content, err := u.Load("/payments")
	if err != nil {
		return nil, fmt.Errorf("loading payments: %s", err)
	}

	var ps []Payment
	if err := json.Unmarshal(content, &ps); err != nil {
		return nil, fmt.Errorf("parsing payments: %s", err)
	}

	return ps, nil
}

// isValidPayment makes sure the required keys are present, and of the right
// type
func isValidPayment(p Payment) error {
	// this should combine errors (ie find as many errors as possible)
	if _, ok := p["name"]; !ok {
		return errors.New("need 'name' field")
	}
	if _, ok := p["amount"]; !ok {
		return errors.New("need 'amount' field")
	}
	if _, ok := p["date"]; !ok {
		return errors.New("need 'date' field")
	}
	if _, ok := p["amount"].(float64); !ok {
		return errors.New("'amount' should be a float")
	}

	return nil
}

// func (api *API) GetBalance(u *db.User)
