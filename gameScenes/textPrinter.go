package gameScenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"log"
	"math"
	"strings"
)

type TextPrinter struct {
	TextInput         string
	lines             []string //slice of segmented strings for line output on eachline
	Counter           int      //TPS counter for printing text at certain intervals in the gameScenes loop
	CounterOn         bool     //all this does is make letters print every other frame, maybe a better solution
	stringPosition    int      //letter that has been printed so far in current textinput[lineindex]
	charactersPerLine int      //number of characters per line
	StatusText        [3]*widget.TextInput
	NextMessage       bool //if true we should keep printing(after each complete message in status text we set to false then take playerBattleSprite input to keep printing)
	lineCounter       int
	countDown         int
	state             TPState
	autoPlayer        bool
}

type TPState uint8

const (
	Printing TPState = iota
	NotPrinting
)

func NewTextPrinter() *TextPrinter {
	return &TextPrinter{
		TextInput:         "",
		Counter:           0,
		CounterOn:         false,
		stringPosition:    1,
		charactersPerLine: 75,
		NextMessage:       true,
		countDown:         0,
		autoPlayer:        false,
	}
}

func (t *TextPrinter) MessageLoop() {
	//first run setup
	if len(t.TextInput) <= 0 {
		t.NextMessage = false
	}

	if len(t.lines) == 0 {
		t.configureLines()
		t.lineCounter = 0
		for _, line := range t.lines {
			println(line)
		}
	}

	t.StatusText[t.lineCounter].SetText(t.PrintText())

	if t.stringPosition <= len(t.lines[t.lineCounter]) {
		t.stringPosition++
	}

	if t.stringPosition > len(t.lines[t.lineCounter]) {
		t.stringPosition = 1
		t.lineCounter++
	}

	if t.lineCounter == len(t.lines) {
		t.NextMessage = false
		t.lineCounter = 0

	}
}

// function for use in non battle scene

func (t *TextPrinter) configureLines() {
	t.lines = t.MessageWrapperToLines(t.TextInput)
}

func (t *TextPrinter) textWrapper(text string) string {
	if t.charactersPerLine < 10 {
		log.Fatalf(`line length too small: %d\n`, t.charactersPerLine)
	}
	if len(text) <= t.charactersPerLine {
		return text
	}
	letters := strings.Split(text, "")
	letter := letters[t.charactersPerLine]
	i := t.charactersPerLine
	for letter != " " && letter != "." {
		letter = letters[i]
		i--
	}
	return strings.Join(letters[:i+1], "")
}

func (t *TextPrinter) MessageWrapperToLines(text string) (output []string) {

	characters := strings.Split(text, "")
	linesF := float64(len(characters)) / float64(t.charactersPerLine)
	nLines := int(math.Ceil(linesF))

	for _ = range nLines {
		wrappedLine := t.textWrapper(text)
		//[0:10]
		if wrappedLine != "" {
			output = append(output, strings.TrimSpace(wrappedLine))
		}
		//len [9:]
		text = text[len(wrappedLine):]
	}
	return output
}

func (t *TextPrinter) PrintText() (output string) {
	characters := strings.Split(t.lines[t.lineCounter], "")
	if t.stringPosition > len(characters) {
	}

	return strings.Join(characters[:t.stringPosition], "")
}

func (t *TextPrinter) resetTextCounter() {
	t.CounterOn = false
	t.Counter = 0
	t.stringPosition = 0
}

func (t *TextPrinter) UpdateCounter() {
	t.Counter++
}

func (t *TextPrinter) updateCounterOn() {
	t.CounterOn = true
}

func (t *TextPrinter) delayedTrigger(countDown int) {
	t.countDown = countDown
	t.NextMessage = false
}

func (t *TextPrinter) countDownUpdate() {
	if t.countDown > 0 {
		t.countDown--
	}
	if t.countDown == 1 {
		t.NextMessage = true
	}
}

func (t *TextPrinter) ResetTP() {
	t.stringPosition = 1
	t.StatusText[0].SetText("")
	t.StatusText[1].SetText("")
	t.StatusText[2].SetText("")
	t.TextInput = ""
	t.lines = []string{}
	t.lineCounter = 0
	t.NextMessage = true
}

func (t *TextPrinter) UpdateTPState() {
	if !t.NextMessage {
		t.state = NotPrinting
	}

	if t.NextMessage {
		t.state = Printing
	}
}
