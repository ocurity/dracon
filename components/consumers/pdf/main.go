// Package main of the pdf consumer implements a simple consumer for
// applying a go-template to a dracon scan, converting the result to pdf and then
// uploading the result to the S3 bucket passed as an argument
// the consumer expects the environment variables
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY
// to be set along with the "bucket" and "region" arguments to be passed
package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-errors/errors"
	playwright "github.com/playwright-community/playwright-go"

	"github.com/ocurity/dracon/components/consumers"
)

var (
	bucket         string
	region         string
	reportTemplate string
)

func main() {
	flag.StringVar(&bucket, "bucket", "", "s3 bucket name")
	flag.StringVar(&region, "region", "", "s3 bucket region")
	flag.StringVar(&reportTemplate, "template", "", "report html template location")

	if err := consumers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	var responses any
	if consumers.Raw {
		r, err := consumers.LoadToolResponse()
		if err != nil {
			log.Fatal("could not load raw results, file malformed: ", err)
		}
		responses = r
	} else {
		r, err := consumers.LoadEnrichedToolResponse()
		if err != nil {
			log.Fatal("could not load enriched results, file malformed: ", err)
		}
		responses = r
	}
	result, pdfBytes, err := buildPdf(responses)
	if err != nil {
		log.Fatal(err)
	}

	if err = sendToS3(result, bucket, region, pdfBytes); err != nil {
		log.Fatal(err)
	}
}

func sendToS3(filename, bucket, region string, pdfBytes []byte) error {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return errors.Errorf("unable to start session with AWS API: %w", err)
	}

	// filename is statically defined above
	//#nosec:G304
	data, err := os.ReadFile(filename) //#nosec:G304
	if err != nil {
		return errors.Errorf("could not open file: %w", err)
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return errors.Errorf("unable to upload %s to %s: %w", filename, bucket, err)
	}

	pdfFilename := strings.Replace(filename, ".html", "", -1) + ".pdf"
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(pdfFilename),
		Body:   bytes.NewReader(pdfBytes),
	})
	if err != nil {
		return errors.Errorf("unable to upload %s to %s: %w", filename, bucket, err)
	}

	slog.Info("uploaded", "filename", filename, "pdf filename", pdfFilename, "to", "bucket", bucket, "successfully")
	return nil
}

func buildPdf(data any) (string, []byte, error) {
	tmpl := template.Must(template.ParseFiles("default.html"))
	cleanupRun := func(msg string, cleanup func() error) {
		if err := cleanup(); err != nil {
			slog.Error(msg, "error", err)
		}
	}

	currentPath, err := os.Getwd()
	if err != nil {
		return "", nil, errors.Errorf("could not get current working directory: %w", err)
	}

	reportHTMLPath := filepath.Join(currentPath, "report.html")
	//#nosec: G304
	f, err := os.OpenFile(reportHTMLPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600) //#nosec: G304
	if err != nil {
		return "", nil, errors.Errorf("could not open report.html: %w", err)
	}
	defer cleanupRun("could not close file: %w", f.Close)

	if err = tmpl.Execute(f, data); err != nil {
		return "", nil, errors.Errorf("could not apply data to template: %w", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return "", nil, errors.Errorf("could not launch playwright: %w", err)
	}
	defer cleanupRun("could not stop Playwright: %w", pw.Stop)

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return "", nil, errors.Errorf("could not launch Chromium: %w", err)
	}
	context, err := browser.NewContext()
	if err != nil {
		return "", nil, errors.Errorf("could not create context: %w", err)
	}

	page, err := context.NewPage()
	if err != nil {
		return "", nil, errors.Errorf("could not create page: %w", err)
	}

	reportPage := fmt.Sprintf("file:///%s", reportHTMLPath)
	if _, err = page.Goto(reportPage); err != nil {
		return "", nil, errors.Errorf("could not goto page %s in the browser: %w", reportPage, err)
	}

	pdfBytes, err := page.PDF(playwright.PagePdfOptions{
		Path: playwright.String(reportHTMLPath),
	})
	if err != nil {
		return "", nil, errors.Errorf("could not generate pdf from page %s, err: %w", reportPage, err)

	}
	return reportHTMLPath, pdfBytes, err
}
