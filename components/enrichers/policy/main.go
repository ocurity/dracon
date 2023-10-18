package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	opaclient "github.com/ocurity/dracon/components/enrichers/policy/opaClient"
	"github.com/ocurity/dracon/pkg/putil"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	policy     string
	regoServer string
	readPath   string
	writePath  string
	annotation string
)

const defaultAnnotation = "Policy Pass: "

type opaIssue struct {
	Input map[string]interface{}
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func enrichIssue(i *v1.Issue, client opaclient.Client) (*v1.EnrichedIssue, error) {
	var strFinding map[string]interface{}

	if err := json.Unmarshal([]byte(protojson.Format(i)), &strFinding); err != nil {
		log.Println("Could not marshal proto to json err:", err)
	}

	opaIssue := opaIssue{
		Input: strFinding,
	}
	marshalled, err := json.Marshal(opaIssue)
	if err != nil {
		log.Printf("Could not marshal issue err is: %v", err)
	}
	passed, err := client.Decide(marshalled)
	if err != nil {
		log.Printf("Could not make a decision on finding %v, returning false, err is: %v", i, err)
	}
	issue := &v1.EnrichedIssue{
		RawIssue:    i,
		Annotations: map[string]string{},
	}
	issue.Annotations[annotation+client.PolicyPath] = strconv.FormatBool(passed)
	log.Printf("Evaluated %s to %t on issue titled %s", client.PolicyPath, passed, i.Title)
	return issue, nil
}

func run() {
	res, err := putil.LoadTaggedToolResponse(readPath)
	if err != nil {
		log.Fatalf("could not load tool response from path %s , error:%v", readPath, err)
	}
	if annotation == "" {
		annotation = defaultAnnotation
	}
	client := opaclient.Client{
		RemoteURI: regoServer,
		Policy:    policy,
	}
	if err := client.Bootstrap(); err != nil {
		log.Fatalf("Could not bootstrap OPA policies, err: %v", err)
		return
	}
	log.Printf("successfully bootstrapped policy %s", client.PolicyPath)

	for _, r := range res {
		enrichedIssues := []*v1.EnrichedIssue{}
		for _, i := range r.GetIssues() {
			eI, err := enrichIssue(i, client)
			if err != nil {
				log.Println(err)
				continue
			}
			enrichedIssues = append(enrichedIssues, eI)
		}
		if len(enrichedIssues) > 0 {
			if err := putil.WriteEnrichedResults(r, enrichedIssues,
				filepath.Join(writePath, fmt.Sprintf("%s.policy.enriched.pb", r.GetToolName())),
			); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("no enriched issues were created for", r.GetToolName())
		}
		if len(r.GetIssues()) > 0 {
			scanStartTime := r.GetScanInfo().GetScanStartTime().AsTime()
			if err := putil.WriteResults(
				r.GetToolName(),
				r.GetIssues(),
				filepath.Join(writePath, fmt.Sprintf("%s.raw.pb", r.GetToolName())),
				r.GetScanInfo().GetScanUuid(),
				scanStartTime.Format(time.RFC3339),
			); err != nil {
				log.Fatalf("could not write results: %s", err)
			}
		}

	}
}

func main() {
	flag.StringVar(&policy, "policy", lookupEnvOrString("POLICY", ""), "base64 encoded policy")
	flag.StringVar(&annotation, "annotation", lookupEnvOrString("ANNOTATION", defaultAnnotation), "How to label the issues this binary will enrich by default `Policy Pass: <name of the policy>`")
	flag.StringVar(&regoServer, "opa_server", lookupEnvOrString("OPA_SERVER", "http://127.0.0.1:8181"), "where to find the rego server")
	flag.StringVar(&readPath, "read_path", lookupEnvOrString("READ_PATH", ""), "where to find producer results")
	flag.StringVar(&writePath, "write_path", lookupEnvOrString("WRITE_PATH", ""), "where to put enriched results")
	flag.Parse()
	run()
}
