package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFormatKeyValueInstructions(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "legacy environment",
			input: "ENV NAME hello world\n",
			want:  "ENV\tNAME\thello world\n",
		},
		{
			name:  "modern environment",
			input: "ENV FIRST=one SECOND=\"two words\" EMPTY=\n",
			want:  "ENV\tFIRST=one\tSECOND=\"two words\"\tEMPTY=\n",
		},
		{
			name:  "legacy label",
			input: "LABEL description \"two words\"\n",
			want:  "LABEL\tdescription=\"two words\"\n",
		},
		{
			name:  "modern labels",
			input: "LABEL first=one second=\"two words\" empty=\n",
			want:  "LABEL\tfirst=one\tsecond=\"two words\"\tempty=\n",
		},
		{
			name:  "empty exec argument remains an argument",
			input: "CMD [\"\", \"two words\"]\n",
			want:  "CMD\t[\"\",\"two words\"]\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "Dockerfile")
			if err := os.WriteFile(path, []byte("FROM alpine\n"+tc.input), 0644); err != nil {
				t.Fatal(err)
			}
			cmd := formatCommand{write: true}
			if err := cmd.Run(context.Background(), []string{path}); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			want := "\t\nFROM\talpine\n" + tc.want
			if string(got) != want {
				t.Fatalf("formatted Dockerfile = %q, want %q", got, want)
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 || entries[0].Name() != "Dockerfile" {
				t.Fatalf("formatter left unexpected backup files: %v", entries)
			}
		})
	}
}

func TestFormatInvalidEnvironmentLeavesFileUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Dockerfile")
	input := "FROM alpine\nENV ONLYNAME\n"
	if err := os.WriteFile(path, []byte(input), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := formatCommand{write: true}
	if err := cmd.Run(context.Background(), []string{path}); err == nil {
		t.Fatal("formatting an environment instruction without a value succeeded")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != input {
		t.Fatalf("failed formatting changed source to %q", got)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "Dockerfile" {
		t.Fatalf("failed formatting left unexpected backup files: %v", entries)
	}
}
