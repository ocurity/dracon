package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ocurity/dracon/components/consumers"
	playwright "github.com/playwright-community/playwright-go"
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

	result := buildPdf(responses)
	sendToS3(result, bucket, region)
}

func sendToS3(filename, bucket, region string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	},
	)
	if err != nil {
		log.Fatalf("Unable to acquire AWS session in region %s, check your credentials", region)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}

func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func buildPdf(data any) string {
	tmpl := template.Must(template.ParseFiles("default.html"))

	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(filepath.Join(currentPath, "report.html"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tmpl.Execute(f, data)

	pw, err := playwright.Run()
	assertErrorToNilf("could not launch playwright: %w", err)

	browser, err := pw.Chromium.Launch()
	assertErrorToNilf("could not launch Chromium: %w", err)

	context, err := browser.NewContext()
	assertErrorToNilf("could not create context: %w", err)

	page, err := context.NewPage()
	assertErrorToNilf("could not create page: %w", err)

	_, err = page.Goto(fmt.Sprintf("file:///%s", filepath.Join(currentPath, "report.html")))
	assertErrorToNilf("could not goto: %w", err)

	_, err = page.PDF(playwright.PagePdfOptions{
		Path: playwright.String(filepath.Join(currentPath, "report.pdf")),
	})
	assertErrorToNilf("could not create PDF: %w", err)
	assertErrorToNilf("could not close browser: %w", browser.Close())
	assertErrorToNilf("could not stop Playwright: %w", pw.Stop())

	return filepath.Join(currentPath, "report.pdf")
}