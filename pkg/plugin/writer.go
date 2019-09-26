package plugin

import (
	"github.com/olekukonko/tablewriter"
)

func (o PluginOptions) PrintTable(title string, data [][]string) {
	table := tablewriter.NewWriter(o.Output)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	//table.SetAutoWrapText(true)
	o.Printf("%s:\n", title)
	table.Render()
	o.Printf("")
}
