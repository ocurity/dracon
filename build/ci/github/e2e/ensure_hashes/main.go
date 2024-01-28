package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbURL          string
	dbName         string
	collectionName string
	hashesLoc      string
)

func init() {
	flag.StringVar(&dbURL, "db-uri", "https://localhost:8529", "URI to connect to MongoDB.")
	flag.StringVar(&dbName, "db-name", "consumer-mongodb", "The MongoDB database name to use.")
	flag.StringVar(&collectionName, "collection-name", "consumer-mongodb", "The MongoDB collection name to use.")
	flag.StringVar(&hashesLoc, "hashLoc", "/hashes.txt", "The path to the file containing the hashes that should be in the database")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	bytesRead, err := os.ReadFile(hashesLoc)
	if err != nil {
		fmt.Println(err)
		panic("Could not read hashes file")
	}
	fileContent := string(bytesRead)
	content := strings.ReplaceAll(fileContent, "'", "")
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		log.Fatal("hash file is empty")
	}
	log.Println("started successfully")
	if err := run(ctx, lines); err != nil {
		log.Fatalf("run not successful: %s", err)
	}
}

func run(ctx context.Context, want []string) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	if err != nil {
		return err
	}
	log.Println("connected to db")

	defer func() {
		if client == nil {
			return
		}
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	coll := client.Database(dbName).Collection(collectionName)

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	var issues []*v1.EnrichedIssue
	if err := cur.All(ctx, &issues); err != nil {
		return err
	}
	if len(issues) == 0 {
		return fmt.Errorf("could not retrieve any issues from the database")
	}
	log.Println("found", len(issues), "hashes")

	var hashes []string
	for _, iss := range issues {
		hashes = append(hashes, iss.GetHash())
	}
	if err != nil {
		return err
	}
	if !compare(want, hashes) {
		return fmt.Errorf("%s", "lists are different, comparison failed")
	}
	return nil
}

func compare(want, have []string) bool {
	if len(want) == 0 && len(have) == 0 {
		return true
	}
	var extraHave []string
	var extraWant []string
	wantLen := len(want)
	haveLen := len(have)

	visited := make([]bool, haveLen)
	for i := 0; i < wantLen; i++ {
		if want[i] == "" {
			continue
		}
		found := false

		for j := 0; j < haveLen; j++ {
			if visited[j] {
				continue
			}
			if want[i] == have[j] {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraWant = append(extraWant, want[i])
		}
	}

	for j := 0; j < haveLen; j++ {
		if visited[j] {
			continue
		}
		extraHave = append(extraHave, have[j])
	}

	if len(extraWant) == 0 && len(extraHave) == 0 {
		return true
	}
	fmt.Println("Hash lists are different want(len:", len(want), ")!=have(len:", len(have), ")",
		"local list has the following", len(extraWant), "elements different")
	fmt.Println(extraWant)
	fmt.Println("Remote list has the following", len(extraHave), "extra elements")
	fmt.Println(extraHave)

	return false
}
