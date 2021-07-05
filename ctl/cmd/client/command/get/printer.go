package get

import (
	"encoding/json"
	"fmt"
	"os"

	yamljsontool "github.com/ghodss/yaml"
	"github.com/megaease/easemeshctl/cmd/client/resource"
	"github.com/megaease/easemeshctl/cmd/common"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

type (
	printer struct {
		outputFormat string
	}
)

func newPrinter(outputFormat string) *printer {
	return &printer{outputFormat: outputFormat}
}

func (p *printer) printObjects(objects []resource.MeshObject) {
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

func (p *printer) printTable(objects []resource.MeshObject) {
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

func (p *printer) printJSON(objects []resource.MeshObject) {
	yamlBuff, err := yaml.Marshal(objects)
	if err != nil {
		common.ExitWithErrorf("marshal %#v to yaml failed: %v", objects, err)
	}

	jsonBuff, err := yamljsontool.YAMLToJSON(yamlBuff)
	if err != nil {
		common.ExitWithErrorf("transform yaml %s to json failed: %v", yamlBuff, err)
	}

	var v interface{}
	err = json.Unmarshal(jsonBuff, &v)
	if err != nil {
		common.ExitWithErrorf("unmarshal %s to json failed: %v", jsonBuff, err)
	}

	prettyJSONBuff, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		common.ExitWithErrorf("unmarshal %#v to json failed: %v", v, err)
	}

	fmt.Printf("%s\n", prettyJSONBuff)
}

func (p *printer) printYAML(objects []resource.MeshObject) {
	yamlBuff, err := yaml.Marshal(objects)
	if err != nil {
		common.ExitWithErrorf("marshal %#v to yaml failed: %v", objects, err)
	}

	fmt.Printf("%s", yamlBuff)
}
