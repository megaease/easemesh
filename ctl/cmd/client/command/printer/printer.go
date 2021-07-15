/*
 * Copyright (c) 2017, MegaEase
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
	"encoding/json"
	"fmt"
	"os"

	yamljsontool "github.com/ghodss/yaml"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common"
	"github.com/olekukonko/tablewriter"
)

type (
	Printer struct {
		outputFormat string
	}
)

func New(outputFormat string) *Printer {
	return &Printer{outputFormat: outputFormat}
}

func (p *Printer) PrintObjects(objects []resource.MeshObject) {
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

func (p *Printer) printTable(objects []resource.MeshObject) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kind", "Name", "Labels"})

	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetColumnSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, object := range objects {
		var labels string
		for k, v := range object.Labels() {
			labels += k + "=" + v
		}
		table.Append([]string{
			object.Kind(),
			object.Name(),
			labels,
		})
	}

	table.Render()
}

func (p *Printer) printYAML(objects []resource.MeshObject) {
	jsonBuff, err := json.Marshal(objects)
	if err != nil {
		common.ExitWithErrorf("marshal %#v to json failed: %v", objects, err)
	}

	yamlBuff, err := yamljsontool.JSONToYAML(jsonBuff)
	if err != nil {
		common.ExitWithErrorf("transform yaml %s to json failed: %v", yamlBuff, err)
	}

	fmt.Printf("%s", yamlBuff)
}

func (p *Printer) printJSON(objects []resource.MeshObject) {
	prettyJSONBuff, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		common.ExitWithErrorf("unmarshal %#v to json failed: %v", objects, err)
	}

	fmt.Printf("%s\n", prettyJSONBuff)
}
