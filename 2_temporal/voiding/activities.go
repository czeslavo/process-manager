package voiding

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func MarkDocumentAsVoidedInReports(ctx context.Context, documentUUID string) error {
	fmt.Println("marking as voided in reports", documentUUID)

	form := url.Values{}
	form.Add("id", documentUUID)

	request, err := http.NewRequest("POST", "http://localhost:8080/reports/markDocumentAsVoided", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	if _, err := client.Do(request); err != nil {
		return err
	}

	return nil
}

func VoidDocumentInBilling(ctx context.Context, documentUUID string) error {
	fmt.Println("voiding document in billing", documentUUID)

	form := url.Values{}
	form.Add("id", documentUUID)

	request, err := http.NewRequest("POST", "http://localhost:8080/documents/void", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	if _, err := client.Do(request); err != nil {
		return err
	}

	return nil
}
