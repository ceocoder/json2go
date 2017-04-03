# json2go

[![Travis CI](https://travis-ci.org/ceocoder/json2go.svg?branch=master)](https://travis-ci.org/ceocoder/json2go)

BigQuery JSON Schema to Go schema converter

```console
json2go --help

     | |/ ____|/ __ \| \ | |      |__ \        |  _ \ / __ \
     | | (___ | |  | |  \| |  ______ ) |_____  | |_) | |  | |
 _   | |\___ \| |  | | . | |  ______/ /______| |  _ <| |  | |
| |__| |____) | |__| | |\  |       / /_        | |_) | |__| |
 \____/|_____/ \____/|_| \_|      |____|       |____/ \___\_\


BigQuery JSON schema to Go schema converter

Version v0.0.1

json2go provides support for converting BigQuery schema from JSON to Go structs

    json2go -schemaFile <file>

$json2go --help

  -packageName string
        name of package (default "main")
  -schemaFile string
        path to file with json schema
```

It expects schemaFile to contain JSON in format below

```json
[
    {
        "type": "RECORD",
        "name": "data",
        "fields": [
            {
                "type": "BOOLEAN",
                "name": "active"
            },
            {
    
                "name": "ic"
            },
            {
                "type": "INTEGER",
                "name": "id"
            },
            {
                "type": "STRING",
                "name": "name"
            }
        ]
    },
    {
        "type": "INTEGER",
        "name": "numbers"
    },
    {
        "type": "STRING",
        "name": "currency_code"
    }
]
```

It generate Go struct definition like below,

```go
package main

import "cloud.google.com/go/bigquery"

var schema = bigquery.Schema{
        &bigquery.FieldSchema{
                Name: "data",
                Type: bigquery.RecordFieldType,
                Schema: &bigquery.FieldSchema{
                        Name: "active",
                        Type: bigquery.BooleanFieldType,
                },
                &bigquery.FieldSchema{
                        Name: "ic",
                        Type: bigquery.StringFieldType,
                },
                &bigquery.FieldSchema{
                        Name: "id",
                        Type: bigquery.IntegerFieldType,
                },
                &bigquery.FieldSchema{
                        Name: "name",
                        Type: bigquery.StringFieldType,
                },
        },
        &bigquery.FieldSchema{
                Name: "numbers",
                Type: bigquery.IntegerFieldType,
        },
        &bigquery.FieldSchema{
                Name: "currency_code",
                Type: bigquery.StringFieldType,
        },
}
```



