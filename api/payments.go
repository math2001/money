package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
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
	var ps []Payment
	content, err := u.Load("/payments")
	var patherr *os.PathError
	if errors.As(err, &patherr) && os.IsNotExist(patherr) {
		return ps, nil // no payments
	} else if err != nil {
		return nil, fmt.Errorf("loading payments: %s", err)
	}

	if err := json.Unmarshal(content, &ps); err != nil {
		return nil, fmt.Errorf("parsing payments: %s", err)
	}

	return ps, nil
}

// Scan requires user just to make sure that only members use this expensive
// feature
func (api *API) Scan(user *db.User, header *multipart.FileHeader, img image.Image) (*Payment, error) {
	log.Printf("start scan job for %s: %q %d", user.Email, header.Filename, header.Size)
	defer log.Printf("done scan job for %s: %q %d", user.Email, header.Filename, header.Size)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		writer.Close()
		return nil, fmt.Errorf("create form file part: %s", err)
	}
	png.Encode(part, img)
	writer.Close()

	u := &url.URL{
		Scheme: "http", // FIXME: use https
		Host:   api.ocrserver,
		Path:   "/file",
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, fmt.Errorf("creating request: %s", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing scan request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("reading request body: %s", err)
		} else {
			log.Printf("body: %q", body)
		}
		return nil, fmt.Errorf("scan request, invalid response code %d instead of 200", resp.StatusCode)
	}

	// TODO: Contribute to ocrserver, because this isn't right (with version
	// 0.2.0, I actually get text/plain)
	if false && resp.Header.Get("Content-Type") != "application/json" {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("reading request body: %s", err)
		} else {
			log.Printf("body: %q", body)
		}
		return nil, fmt.Errorf("scan request, invalid response content type %q instead of \"application/json\"", resp.Header.Get("Content-Type"))
	}

	var payment *Payment
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&payment); err != nil {
		return nil, fmt.Errorf("scan request, decode response: %s", err)
	}

	return payment, nil
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
