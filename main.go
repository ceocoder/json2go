package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = `

  ____  _        ____                               _  _____  ____  _   _            ___              _____  ____
 |  _ \(_)      / __ \                             | |/ ____|/ __ \| \ | |          |__ \            / ____|/ __ \
 | |_) |_  __ _| |  | |_   _  ___ _ __ _   _       | | (___ | |  | |  \| |  ______     ) |  ______  | |  __| |  | |
 |  _ <| |/ _' | |  | | | | |/ _ \ '__| | | |  _   | |\___ \| |  | | . ' | |______|   / /  |______| | | |_ | |  | |
 | |_) | | (_| | |__| | |_| |  __/ |  | |_| | | |__| |____) | |__| | |\  |           / /_           | |__| | |__| |
 |____/|_|\__, |\___\_\\__,_|\___|_|   \__, |  \____/|_____/ \____/|_| \_|          |____|           \_____|\____/
           __/ |                        __/ |
          |___/                        |___/

BigQuery JSON schema to Go schema converter

Version %s

json2bq provides support for converting BigQuery schema from JSON to Go structs

    json2bq [flags=value]

Sample input file -

[
    {
        "type": "INTEGER",
        "name": "numbers"
    },
    {
        "type": "STRING",
        "name": "currency_code"
    }
]

Sample Output -

package main

import "cloud.google.com/go/bigquery"

var Schema = bigquery.Schema{
        &bigquery.FieldSchema{
                Name: "numbers",
                Type: bigquery.IntegerFieldType,
        },
        &bigquery.FieldSchema{
                Name: "currency_code",
                Type: bigquery.StringFieldType,
        },
}
`
	// VERSION of the tool
	VERSION            = "v0.0.1"
	defaultPackageName = "main"
)

var (
	schemaFile  string
	packageName string
)

type jsonField struct {
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Fields []jsonField `json:"fields,omitempty"`
	Mode   string      `json:"mode,omitempty"`
}

func init() {
	flag.StringVar(&schemaFile, "schemaFile", "", "path to file with json schema")
	flag.StringVar(&packageName, "packageName", defaultPackageName, "name of package")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	dat, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		log.Fatalf("Failed to read file %v", err)
	}
	schema, err := parseSchema(dat)

	if err != nil {
		log.Fatalf("Failed to build Go var for %s - reason %v", schemaFile, err)
	}
	bt, err := buildGoStructs(schema, 0)
	if err != nil {
		log.Fatalf("unable to build Go stuct %v", err)
	}
	if b, err := format.Source(bt); err != nil {
		log.Fatalf("failed - %s - %s", err, (string(bt)))
	} else {
		fmt.Println(string(b))
	}
}

func parseSchema(data []byte) ([]jsonField, error) {
	schema := []jsonField{}
	err := json.Unmarshal(data, &schema)
	if err == nil {
		return schema, nil
	}
	return nil, fmt.Errorf("failed parse schema %s", err)
}

func buildGoStructs(schema []jsonField, depth int) ([]byte, error) {

	var buf bytes.Buffer
	if depth == 0 {
		buf.WriteString(fmt.Sprintf("package %s\n", packageName))
		buf.WriteString("import \"cloud.google.com/go/bigquery\"\n")
		buf.WriteString("var Schema = bigquery.Schema{\n")
	}
	for _, field := range schema {
		buf.WriteString(fmt.Sprintf("&bigquery.FieldSchema{\n"))
		buf.WriteString(fmt.Sprintf("Name: \"%s\",\n", field.Name))

		if ft, err := buildTypeString(field.Type); err == nil {
			buf.WriteString(ft)
		} else {
			return nil, err
		}

		buf.WriteString(fmt.Sprintf(buildModeString(field.Mode)))
		if field.Type == "RECORD" {
			sub, err := buildGoStructs(field.Fields, depth+1)
			if err != nil {
				return nil, err
			}
			buf.WriteString(fmt.Sprintf("Schema: bigquery.Schema{\n %s", string(sub)))
			buf.WriteString("},\n")
		}
		buf.WriteString("},\n")
	}
	if depth == 0 {
		buf.WriteString("}\n")
	}
	return buf.Bytes(), nil
}

//
func buildModeString(mode string) string {
	switch uMode := strings.ToUpper(mode); uMode {
	case "REPEATED":
		return `Repeated: True,`
	case "REQUIRED":
		return `Repeated: True,`
	default:
		// it is OK to not have a mode
		return ""
	}
}

// buildTypeString takes field type and returns Go type
func buildTypeString(ftype string) (string, error) {
	switch ucFtype := strings.ToUpper(ftype); ucFtype {
	case "BYTES":
		return "Type: bigquery.BytesFieldType,\n", nil
	case "BOOLEAN":
		return "Type: bigquery.BooleanFieldType,\n", nil
	case "INTEGER":
		return "Type: bigquery.IntegerFieldType,\n", nil
	case "RECORD":
		return "Type: bigquery.RecordFieldType,\n", nil
	case "STRING":
		return "Type: bigquery.StringFieldType,\n", nil
	case "FLOAT":
		return "Type: bigquery.FloatFieldType,\n", nil
	case "TIMESTAMP":
		return "Type: bigquery.TimestampFieldType,\n", nil
	case "DATE":
		return "Type: bigquery.DateFieldType,\n", nil
	case "TIME":
		return "Type: bigquery.TimeFieldType,\n", nil
	case "DATETIME":
		return "Type: bigquery.DateTimeFieldType,\n", nil
	default:
		// it is NOT ok to have field without type
		return "", fmt.Errorf("unknown field type %s - %s", schemaFile, ftype)
	}
}
