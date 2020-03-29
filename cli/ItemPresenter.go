package cli

import (
	"fmt"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"time"
)

type ItemPresenter interface {
	TimelineView(items entities.ItemCollection) error
	BoardView(items entities.ItemCollection) error
}

type itemPresenter struct {
	renderer Renderer
}

func getItemSubtype(item entities.Manageable) string {
	subtype := item.GetStatus()
	if subtype == "" {
		subtype = item.GetType()
	}
	return subtype
}

type itemConfig struct {
	color  COLOR
	symbol string
}

var iconConfig = map[string]itemConfig{
	"note":                      {color: BLUE, symbol: "●"},
	string(entities.ToDo):       {color: MAGENTA, symbol: "☐"},
	string(entities.InProgress): {color: BLUE, symbol: "…"},
	string(entities.Stopped):    {color: YELLOW, symbol: "☐"},
	string(entities.Completed):  {color: GREEN, symbol: "✔"},
	string(entities.Cancelled):  {color: RED, symbol: "✖"},
}

func (p *itemPresenter) getColor(subtype string) COLOR {
	config := iconConfig[subtype]
	return config.color
}

func (p *itemPresenter) getIcon(item entities.Manageable) string {
	subtype := getItemSubtype(item)
	config := iconConfig[subtype]
	if config.symbol == "" {
		return ""
	}
	return p.renderer.Colorify(config.symbol, config.color)
}

func GetColorStatus(min, max, i int) COLOR {
	switch i {
	case min:
		return GREY
	case max:
		return GREEN
	default:
		return YELLOW
	}
}

func (p *itemPresenter) PrintSummary(summary Summary, prefix string) error {
	donePercentage := summary.GetDonePercentage()
	color := GetColorStatus(0, summary.TasksCount, summary.DoneCount)

	percentageTxt := fmt.Sprintf("%v%%", donePercentage)
	_, err := fmt.Printf("\n%v%v %v\n", prefix, p.renderer.Colorify(percentageTxt, color), p.renderer.Colorify("of all tasks complete.", GREY))

	fmt.Printf(
		"%v%v %v  %v %v  %v %v  %v %v\n",
		prefix,
		p.renderer.Colorify(summary.DoneCount, p.getColor(string(entities.Completed))),
		p.renderer.Colorify(" done", GREY),
		p.renderer.Colorify(summary.InProgressCount, p.getColor(string(entities.InProgress))),
		p.renderer.Colorify(" in-progress", GREY),
		p.renderer.Colorify(summary.PendingCount, p.getColor(string(entities.ToDo))),
		p.renderer.Colorify(" pending", GREY),
		p.renderer.Colorify(summary.NoteCount, p.getColor("note")),
		p.renderer.Colorify(" notes", GREY),
	)
	return err
}

func (p *itemPresenter) renderItemForTimelineView(item entities.Manageable) string {
	tags := ""
	for _, tag := range item.GetTags() {
		tags = fmt.Sprintf("%s @%s", tags, tag)
	}
	star := ""
	if item.HasStar() {
		star = p.renderer.Colorify(" ★ ", YELLOW)
	}
	return fmt.Sprintf(
		"%v. %v %v %v%v",
		p.renderer.Colorify(item.GetId(), GREY),
		p.getIcon(item),
		item.GetTitle(),
		star,
		p.renderer.Colorify(tags, GREY),
	)
}

func (p *itemPresenter) renderItemForBoardView(item entities.Manageable) string {
	duration := GetDurationText(time.Now(), item.GetCreationDateTime())
	star := ""
	if item.HasStar() {
		star = p.renderer.Colorify(" ★ ", YELLOW)
	}
	return fmt.Sprintf(
		"%v. %v %v %v%v",
		p.renderer.Colorify(item.GetId(), GREY),
		p.getIcon(item),
		item.GetTitle(),
		star,
		p.renderer.Colorify(duration, GREY),
	)
}

func (p *itemPresenter) TimelineView(items entities.ItemCollection) error {
	timeline := GroupBy(items, GroupByCreationDate)
	summary := Summarize(items)
	for _, group := range timeline {
		group.Print(p.renderer, p.renderItemForTimelineView)
	}
	return p.PrintSummary(summary, "  ")
}

func (p *itemPresenter) BoardView(items entities.ItemCollection) error {
	flattened := FlatByTags(items)
	timeline := GroupBy(flattened, GroupByTag)
	summary := Summarize(items)
	for _, group := range timeline {
		group.Print(p.renderer, p.renderItemForBoardView)
	}
	return p.PrintSummary(summary, "  ")
}

func NewItemPresenter(renderer Renderer) ItemPresenter {
	return &itemPresenter{renderer: renderer}
}
