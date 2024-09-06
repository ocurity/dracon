package enricher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/atom"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/conf"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/fs"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/logging"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/search"
)

type (
	enricher struct {
		cfg        *conf.Conf
		atomReader *atom.Reader
		readWriter *fs.ReadWriter
	}
)

// NewEnricher returns a new reachability enricher.
func NewEnricher(
	cfg *conf.Conf,
	atomReader *atom.Reader,
	readWriter *fs.ReadWriter,
) (*enricher, error) {
	switch {
	case cfg == nil:
		return nil, errors.New("invalid nil configuration provided")
	case atomReader == nil:
		return nil, errors.New("invalid nil atom reader provided")
	case readWriter == nil:
		return nil, errors.New("invalid nil read writer provided")
	}

	return &enricher{
		cfg:        cfg,
		atomReader: atomReader,
		readWriter: readWriter,
	}, nil
}

// Enrich looks for untagged inputs and processes them outputting if any of them is reachable.
// The reachability checks leverage atom - https://github.com/AppThreat/atom.
func (r *enricher) Enrich(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	taggedRes, err := r.readWriter.ReadTaggedResponse()
	if err != nil {
		return fmt.Errorf("could not read tagged response: %w", err)
	}

	reachablesRes, err := r.atomReader.Read()
	if err != nil {
		return fmt.Errorf("could not read atom reachables from path %s: %w", r.cfg.ATOMFilePath, err)
	}

	reachablePurls, err := r.atomReader.ReachablePurls(ctx, reachablesRes)
	if err != nil {
		return fmt.Errorf("could not get reachable purls: %w", err)
	}

	searcher, err := search.NewSearcher(reachablesRes.Reachables, reachablePurls)
	if err != nil {
		return fmt.Errorf("could not create searcher: %w", err)
	}

	for _, taggedEntry := range taggedRes {
		var (
			enrichedIssues []*v1.EnrichedIssue
			issues         = taggedEntry.GetIssues()
		)
		for _, issue := range issues {
			// Search.
			ok, err := searcher.Search(issue.Target)
			if err != nil {
				logger.Error(
					"could not search target. Continuing...",
					slog.String("target", issue.Target),
					slog.String("err", err.Error()),
				)
				continue
			}

			// Enrich.
			enrichedIssues = append(enrichedIssues, &v1.EnrichedIssue{
				RawIssue: issue,
				Annotations: map[string]string{
					"reachable": strconv.FormatBool(ok),
				},
			})
		}

		for _, ei := range enrichedIssues {
			v, ok := ei.Annotations["reachable"]
			if ok && v == "true" {
				_ = v
			}
		}

		// Write results.
		if err := r.readWriter.WriteEnrichedResults(taggedEntry, enrichedIssues); err != nil {
			logger.Error(
				"could not write enriched results. Continuing...",
				slog.String("err", err.Error()),
			)
			continue
		}

		if err := r.readWriter.WriteRawResults(taggedEntry); err != nil {
			logger.Error(
				"could not write raw results. Continuing...",
				slog.String("err", err.Error()),
			)
			continue
		}
	}

	return nil
}
