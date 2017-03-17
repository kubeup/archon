package common // import "go.pedge.io/pb/etc/cmd/common"

import (
	"bytes"
	"encoding/csv"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go.pedge.io/env"
)

// GenerateEnv is the environment for Generators.
type GenerateEnv struct {
	// must be absolute path
	RepoDir string `env:"REPO_DIR"`
	// relative to RepoDir
	TmplDir string `env:"TMPL_DIR"`
	// relative to RepoDir/TmplDir, comma separated
	GoTmplFiles string `env:"GO_TMPL_FILES"`
	// relative to RepoDir/TmplDir, comma separated
	PbTmplFiles string `env:"PB_TMPL_FILES"`
}

// CSVGenerateEnv is the environment for CSVGenerators.
type CSVGenerateEnv struct {
	GenerateEnv
	// relative to RepoDir
	CSVFile string `env:"CSV_FILE"`
}

// TmplData is what is passed into the template.
type TmplData struct {
	// either proto, go, or gogo
	Type string
	// The data from the GenerateHelper
	Data interface{}
}

// GenerateHelper does stuff for Generators.
type GenerateHelper interface {
	ExtraTmplFuncs() map[string]interface{}
}

// CSVGenerateHelper does stuff for CSVGenerators.
type CSVGenerateHelper interface {
	GenerateHelper
	TmplData(records [][]string) (interface{}, error)
}

// Main is the main function for commands.
func Main(generateHelper GenerateHelper) {
	env.Main(
		func(generateEnvObj interface{}) error {
			return generate(nil, generateHelper, generateEnvObj.(*GenerateEnv))
		},
		&GenerateEnv{},
	)
}

// CSVMain is the main function for csv commands.
func CSVMain(csvGenerateHelper CSVGenerateHelper) {
	env.Main(
		func(csvGenerateEnvObj interface{}) error {
			csvGenerateEnv := csvGenerateEnvObj.(*CSVGenerateEnv)
			records, err := getCSVRecords(filepath.Join(csvGenerateEnv.GenerateEnv.RepoDir, csvGenerateEnv.CSVFile))
			if err != nil {
				return err
			}
			data, err := csvGenerateHelper.TmplData(records)
			if err != nil {
				return err
			}
			return generate(data, csvGenerateHelper, &csvGenerateEnv.GenerateEnv)
		},
		&CSVGenerateEnv{},
	)
}

func generate(data interface{}, generateHelper GenerateHelper, generateEnv *GenerateEnv) error {
	tmplFiles, err := getTmplFiles(generateEnv)
	if err != nil {
		return err
	}
	for _, tmplFile := range tmplFiles {
		if err := tmplFile.generate(data, generateHelper.ExtraTmplFuncs()); err != nil {
			return err
		}
	}
	return nil
}

type tmplFile struct {
	// absolute path
	Path string
	// absolute output path of generated file
	OutputPath string
	// either proto, go, or gogo
	Type string
}

func getTmplFiles(generateEnv *GenerateEnv) ([]*tmplFile, error) {
	var tmplFiles []*tmplFile
	if generateEnv.GoTmplFiles != "" {
		for _, goTmplFilePath := range strings.Split(generateEnv.GoTmplFiles, ",") {
			for _, t := range []string{"go", "gogo"} {
				tmplFile, err := getTmplFile(
					generateEnv.RepoDir,
					generateEnv.TmplDir,
					goTmplFilePath,
					t,
				)
				if err != nil {
					return nil, err
				}
				tmplFiles = append(tmplFiles, tmplFile)
			}
		}
	}
	if generateEnv.PbTmplFiles != "" {
		for _, pbTmplFilePath := range strings.Split(generateEnv.PbTmplFiles, ",") {
			tmplFile, err := getTmplFile(
				generateEnv.RepoDir,
				generateEnv.TmplDir,
				pbTmplFilePath,
				"proto",
			)
			if err != nil {
				return nil, err
			}
			tmplFiles = append(tmplFiles, tmplFile)
		}
	}
	return tmplFiles, nil
}

func getTmplFile(repoPath string, tmplRelPath string, relPath string, t string) (*tmplFile, error) {
	return &tmplFile{
		Path:       filepath.Join(repoPath, tmplRelPath, relPath),
		OutputPath: filepath.Join(repoPath, t, strings.TrimSuffix(relPath, ".tmpl")),
		Type:       t,
	}, nil
}

func (t *tmplFile) tmplData(data interface{}) *TmplData {
	return &TmplData{
		Type: t.Type,
		Data: data,
	}
}

func (t *tmplFile) generate(data interface{}, extraTmplFuncs map[string]interface{}) (retErr error) {
	tmpl, err := newTemplate(t.Path, extraTmplFuncs)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buffer, t.tmplData(data)); err != nil {
		return err
	}
	output := buffer.Bytes()
	if t.Type == "go" || t.Type == "gogo" {
		output, err = format.Source(output)
		if err != nil {
			return err
		}
	}
	if err = os.MkdirAll(filepath.Dir(t.OutputPath), 0744); err != nil {
		return err
	}
	file, err := os.Create(t.OutputPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	if _, err := file.Write(output); err != nil {
		return err
	}
	return nil
}

func newTemplate(file string, extraTmplFuncs map[string]interface{}) (*template.Template, error) {
	funcMap := template.FuncMap{
		"add": func(i int, j int) int {
			return i + j
		},
	}
	for key, value := range extraTmplFuncs {
		funcMap[key] = value
	}
	return template.New(filepath.Base(file)).Funcs(funcMap).ParseFiles(file)
}

func getCSVRecords(csvFilePath string) (_ [][]string, retErr error) {
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := csvFile.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	records, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}
	for i, record := range records {
		records[i] = cleanRecord(record)
	}
	return records, nil
}

func cleanRecord(record []string) []string {
	for i, elem := range record {
		record[i] = strings.TrimSpace(elem)
	}
	return record
}
