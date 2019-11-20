package list

import (
	"os"
	"time"

	"github.com/cloud-annotations/training/cacli/date"
	"github.com/cloud-annotations/training/cacli/mockwml"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"golang.org/x/sys/unix"
)

type tableItem struct {
	Name   string
	Status string
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func getWindowWidth() int {
	defaultWidowWidth := 80
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return defaultWidowWidth
	}
	return int(ws.Col)
}

func render(runs mockwml.TrainingRuns) {

	maxTableWidth := min(getWindowWidth(), 90)
	modelIDWidth := 14
	statusWidth := 11
	submittedWidth := 14
	otherPadding := 11 // left & right space + 3 columns 3 space padding.
	nameWidth := maxTableWidth - modelIDWidth - statusWidth - submittedWidth - otherPadding

	statusTransformer := text.Transformer(func(val interface{}) string {
		// a little on the hacky side...
		ti := val.(tableItem)
		if ti.Name == ti.Status {
			switch ti.Status {
			case "completed":
				return text.Colors{text.FgGreen}.Sprint(ti.Name)
			case "error":
				return text.Colors{text.FgRed}.Sprint(ti.Name)
			case "canceled":
				return text.Colors{text.Faint}.Sprint(ti.Name)
			default:
				return text.Colors{}.Sprint(ti.Name)
			}
		}
		switch ti.Status {
		case "completed":
			return text.Colors{}.Sprint(ti.Name)
		case "error":
			return text.Colors{text.Faint}.Sprint(ti.Name)
		case "canceled":
			return text.Colors{text.Faint}.Sprint(ti.Name)
		default:
			return text.Colors{}.Sprint(ti.Name)
		}
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"name", "model id", "status", "submitted"})
	t.SetStyle(table.Style{
		Box: table.BoxStyle{
			BottomLeft:       "─",
			BottomRight:      "─",
			BottomSeparator:  "───",
			Left:             " ",
			LeftSeparator:    "",
			MiddleHorizontal: "─",
			MiddleSeparator:  "───",
			MiddleVertical:   "   ",
			PaddingLeft:      "",
			PaddingRight:     "",
			Right:            " ",
			RightSeparator:   "",
			TopLeft:          "─",
			TopRight:         "─",
			TopSeparator:     "───",
		},
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateHeader:  false,
			SeparateRows:    false,
		},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:         "name",
			Align:        text.AlignLeft,
			AlignHeader:  text.AlignLeft,
			ColorsHeader: text.Colors{text.Bold},
			Transformer:  statusTransformer,
			WidthMin:     nameWidth,
			WidthMax:     nameWidth,
		}, {
			Name:         "model id",
			Align:        text.AlignCenter,
			AlignHeader:  text.AlignCenter,
			ColorsHeader: text.Colors{text.Bold},
			Transformer:  statusTransformer,
			WidthMin:     modelIDWidth,
			WidthMax:     modelIDWidth,
		}, {
			Name:         "status",
			Align:        text.AlignCenter,
			AlignHeader:  text.AlignCenter,
			ColorsHeader: text.Colors{text.Bold},
			Transformer:  statusTransformer,
			WidthMin:     statusWidth,
			WidthMax:     statusWidth,
		}, {
			Name:         "submitted",
			Align:        text.AlignRight,
			AlignHeader:  text.AlignRight,
			ColorsHeader: text.Colors{text.Bold},
			Transformer:  statusTransformer,
			WidthMin:     submittedWidth,
			WidthMax:     submittedWidth,
		},
	})

	for _, run := range runs.Resources {
		name := run.Entity.ModelDefinition.Name
		guid := run.Metadata.GUID
		status := run.Entity.Status.State
		submitted := run.Entity.Status.SubmittedAt

		t.AppendRow([]interface{}{tableItem{Name: name, Status: status}, tableItem{Name: guid, Status: status}, tableItem{Name: status, Status: status}, tableItem{Name: date.TimeElapsed(time.Now(), submitted, false), Status: status}})
	}
	t.Render()
}
