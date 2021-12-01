/*
 * Copyright (c) 2021, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package printer

import (
	"fmt"
	"os"
	"strings"

	"github.com/megaease/easemeshctl/cmd/client/resource/meta"
	"github.com/megaease/easemeshctl/cmd/common"
	"github.com/olekukonko/tablewriter"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

type (
	// Printer prints information about the EaseMesh objects
	Printer interface {
		PrintObjects(objects []meta.MeshObject)
	}

	printer struct {
		outputFormat string
	}
)

// New creates a Printer
func New(outputFormat string) Printer {
	return &printer{outputFormat: outputFormat}
}

func (p *printer) PrintObjects(objects []meta.MeshObject) {
	if len(objects) == 0 {
		fmt.Println("No resource")
		return
	}
	switch p.outputFormat {
	case "table":
		p.printTable(objects)
	case "json":
		p.printJSON(objects)
	case "yaml":
		p.printYAML(objects)
	default:
		common.ExitWithErrorf("unsupported output format: %s", p.outputFormat)
	}
}

func (p *printer) printTable(objects []meta.MeshObject) {
	table := tablewriter.NewWriter(os.Stdout)

	header := []string{"Kind", "Name", "Labels"}

	var headerColumns []*meta.TableColumn
	for _, object := range objects {
		if tableObject, ok := object.(meta.TableObject); ok {
			headerColumns = tableObject.Columns()
			break
		}
	}

	for _, column := range headerColumns {
		header = append(header, column.Name)
	}

	table.SetHeader(header)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetColumnSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, object := range objects {
		var labels []string
		for k, v := range object.Labels() {
			labels = append(labels, k+"="+v)
		}

		row := []string{
			object.Kind(),
			object.Name(),
			strings.Join(labels, ","),
		}

		tableObject, ok := object.(meta.TableObject)
		if ok {
			for _, column := range tableObject.Columns() {
				row = append(row, column.Value)
			}
		}

		table.Append(row)
	}

	table.Render()
}

func (p *printer) printYAML(objects []meta.MeshObject) {
	
	yamlBuff, err := yaml.Marshal(objects)
	if err != nil {
		common.ExitWithErrorf("marshal %#v to yaml failed: %v", objects, err)
	}

	fmt.Printf("%s", yamlBuff)
}

func (p *printer) printJSON(objects []meta.MeshObject) {
	yamlBuff, err := yaml.Marshal(objects)
	if err != nil {
		common.ExitWithErrorf("marshal %#v to yaml failed: %v", objects, err)
	}

	var m interface{}
	err = yaml.Unmarshal(yamlBuff, &m)
	if err != nil {
		common.ExitWithErrorf("unmarshal %#v to yaml failed: %v", objects, err)
	}

	prettyJSONBuff, err := jsoniter.MarshalIndent(m, "", "  ")
	if err != nil {
		common.ExitWithErrorf("marshal %#v to json failed: %v", m, err)
	}

	fmt.Printf("%s\n", prettyJSONBuff)
}
