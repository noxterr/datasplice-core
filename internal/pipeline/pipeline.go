// Package pipeline runs a flow described by config.Main.
//
// ponytail: hardcoded csv-only input/output, executed in-process — no
// packages, no gRPC. This proves the row model end-to-end before paying
// for the package system (docs/roadmap.md phase 4). Transform steps (the
// ones in the middle of `steps[]`) are parsed but not executed: that needs
// a package to actually run `via`/`on` against, which is phase 6. Upgrade
// path: phase 6 replaces readCSV/writeCSV below with real calls over
// docs/proto/datasplice.proto, dialed via hashicorp/go-plugin.
package pipeline

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/datasplice-labs/datasplice/internal/config"
)

// Row mirrors the Row message in docs/proto/datasplice.proto — same shape
// so swapping this hardcoded pipeline for real packages later doesn't
// change how a row looks.
type Row = map[string]string

func stepPath(step config.Step) (string, bool) {
	p, ok := step.With["path"].(string)

	return p, ok && p != ""
}

func readCSV(path string) (header []string, rows []Row, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, nil, err
	}
	if len(records) == 0 {
		return nil, nil, nil
	}

	header = records[0]
	rows = make([]Row, 0, len(records)-1)
	for _, rec := range records[1:] {
		row := make(Row, len(header))
		for i, col := range header {
			if i < len(rec) {
				row[col] = rec[i]
			}
		}
		rows = append(rows, row)
	}
	return header, rows, nil
}

// writeCSV writes to stdout when path is empty.
func writeCSV(path string, header []string, rows []Row) error {
	out := os.Stdout
	if path != "" {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	w := csv.NewWriter(out)
	if err := w.Write(header); err != nil {
		return err
	}
	for _, row := range rows {
		rec := make([]string, len(header))
		for i, col := range header {
			rec[i] = row[col]
		}
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// Run reads the input step's csv, passes rows through unchanged, and
// writes the output step's csv (or stdout, if the output step has no
// `with.path`).
func Run(m *config.Main) error {
	if m.Config.Input.Type != "csv" || m.Config.Output.Type != "csv" {
		return fmt.Errorf("only csv input/output is wired up until the package system lands (docs/roadmap.md phase 6); got input=%q output=%q", m.Config.Input.Type, m.Config.Output.Type)
	}

	inPath, ok := stepPath(m.Steps[0])
	if !ok {
		return fmt.Errorf("input step %q needs `with.path`", m.Steps[0].Uses)
	}
	header, rows, err := readCSV(inPath)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	outPath, _ := stepPath(m.Steps[len(m.Steps)-1]) // empty path -> stdout
	if err := writeCSV(outPath, header, rows); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	return nil
}

// Describe renders the one-line-per-step macro report `datasplice plan`
// prints (docs/schema.md#plan-output).
//
// ponytail: this reads the yaml directly to guess a description. Once
// packages exist (phase 6), each package supplies its own
// ConfigureResponse.plan_description instead — see
// docs/proto/datasplice.proto.
func Describe(m *config.Main) []string {
	lines := make([]string, 0, len(m.Steps))
	for i, step := range m.Steps {
		role := "step"
		switch i {
		case 0:
			role = "input"
		case len(m.Steps) - 1:
			role = "output"
		}

		detail := ""
		switch {
		case len(step.Via) > 0:
			detail = fmt.Sprintf("via %v on columns %v", step.Via, step.On)
		default:
			if p, ok := stepPath(step); ok {
				detail = p
			}
		}

		lines = append(lines, fmt.Sprintf("%-6s %-45s %s", role, step.Uses, detail))
	}
	return lines
}
