// Package handler defines the HTTP handlers for the spend categoriser API.
// All handlers follow the GoFr signature: func(*gofr.Context) (any, error).
package handler

import (
	"io"
	"mime/multipart"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"

	"github.com/spend-categoriser/backend/internal/categorizer"
	"github.com/spend-categoriser/backend/internal/extractor"
	"github.com/spend-categoriser/backend/internal/models"
	"github.com/spend-categoriser/backend/internal/parser"
)

// uploadForm binds the multipart "files" field to a file header.
type uploadForm struct {
	File *multipart.FileHeader `file:"files"`
}

// Health returns a simple status check response.
func Health(_ *gofr.Context) (any, error) {
	return response.Raw{Data: map[string]string{"status": "ok"}}, nil
}

// ParsePDF accepts multipart PDF file uploads, extracts transactions from each,
// and returns the categorised results. The "files" form field may contain one or
// more PDF files. Returns a JSON array of ParseResponse objects.
func ParsePDF(ctx *gofr.Context) (any, error) {
	var form uploadForm
	if err := ctx.Bind(&form); err != nil {
		return nil, errBadRequest("failed to parse upload: " + err.Error())
	}

	if form.File == nil {
		return nil, errBadRequest("no files uploaded")
	}

	var results []models.ParseResponse

	for _, fh := range []*multipart.FileHeader{form.File} {
		f, err := fh.Open()
		if err != nil {
			return nil, errBadRequest("failed to read file: " + fh.Filename)
		}

		const maxPDFSize = 50 << 20 // 50 MB
		data, err := io.ReadAll(io.LimitReader(f, maxPDFSize+1))
		f.Close()
		if err != nil {
			return nil, errBadRequest("failed to read file: " + fh.Filename)
		}
		if int64(len(data)) > maxPDFSize {
			return nil, errEntityTooLarge("file too large: " + fh.Filename)
		}

		text, err := parser.ExtractText(data)
		if err != nil {
			return nil, errUnprocessable("failed to parse PDF: " + fh.Filename + ": " + err.Error())
		}

		transactions, period := extractor.ExtractTransactions(text)

		var total float64
		for _, t := range transactions {
			if t.Amount > 0 {
				total += t.Amount
			}
		}

		results = append(results, models.ParseResponse{
			Transactions: transactions,
			BillPeriod:   period,
			CardName:     "",
			TotalAmount:  total,
			FileName:     fh.Filename,
		})
	}

	return response.Raw{Data: results}, nil
}

// GetCategories returns all available spend categories and their subcategories.
func GetCategories(_ *gofr.Context) (any, error) {
	cats := categorizer.GetCategories()
	return response.Raw{Data: cats}, nil
}
