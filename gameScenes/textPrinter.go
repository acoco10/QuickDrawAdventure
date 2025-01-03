package gameScenes

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/ebitenui/ebitenui/widget"
	"log"
	"math"
	"strings"
)

type TextPrinter struct {
	TextInput         []string
	lines             []string //slice of segmented strings for line output on eachline
	Counter           int      //TPS counter for printing text at certain intervals in the gameScenes loop
	CounterOn         bool     //all this does is make letters print every other frame, maybe a better solution
	stringPosition    int      //letter that has been printed so far in current textinput[lineindex]
	MessageIndex      int      //current index in text input slice
	charactersPerLine int      //number of characters per line
	StatusText        [3]*widget.TextInput
	NextMessage       bool //if true we should keep printing(after each complete message in status text we set to false then take playerBattleSprite input to keep printing)
	lineCounter       int
	countDown         int
}

func NewTextPrinter(initialText []string) *TextPrinter {
	return &TextPrinter{
		TextInput:         initialText,
		Counter:           0,
		CounterOn:         false,
		stringPosition:    1,
		charactersPerLine: 75,
		NextMessage:       true,
		MessageIndex:      0,
		countDown:         0,
	}
}

func (t *TextPrinter) MessageLoop(g *BattleScene) {

	//first run setup
	if t.MessageIndex == 0 && len(t.lines) == 0 {
		println("configuring new lines at first check\n")
		t.configureLines()
		t.lineCounter = 0
		for _, line := range t.lines {
			fmt.Printf("%s\n", line)
		}
	}

	t.StatusText[t.lineCounter].SetText(t.printText())

	if t.stringPosition <= len(t.lines[t.lineCounter]) {
		t.stringPosition++
	}

	if t.stringPosition == len(t.lines[t.lineCounter])+1 {

		t.stringPosition = 1

		if t.MessageIndex+1 == len(t.TextInput) && t.lineCounter+1 == len(t.lines) {
			fmt.Printf("resetting text\n")
			t.TextInput = []string{}
			t.lines = []string{}
			t.NextMessage = false
			t.lineCounter = 0
			t.MessageIndex = 0
		}

		t.lineCounter++
	}

	if t.lineCounter == len(t.lines) {
		if t.MessageIndex < len(t.TextInput)-1 {
			t.MessageIndex++
			t.configureLines()
			for _, line := range t.lines {
				fmt.Printf("%s\n", line)
			}
		} else {
			t.MessageIndex = 0
		}

		t.lineCounter = 0
		if g.battle.GetPhase() != battle.Shooting {
			t.NextMessage = false
		}
		if g.battle.GetPhase() == battle.Shooting && t.MessageIndex == g.battle.GetTurn().PlayerStartIndex+1 {
			if t.countDown == 0 {
				t.delayedTrigger(20)
			}
		}
	}
}

func (t *TextPrinter) DialogueMessageLoop() {
	//first run setup
	if t.MessageIndex == 0 && len(t.lines) == 0 {
		println("configuring new lines at first check\n")
		t.configureLines()
		t.lineCounter = 0
		for _, line := range t.lines {
			fmt.Printf("%s\n", line)
		}
	}

	t.StatusText[t.lineCounter].SetText(t.printText())

	if t.stringPosition <= len(t.lines[t.lineCounter]) {
		t.stringPosition++
	}

	if t.stringPosition == len(t.lines[t.lineCounter])+1 {

		t.stringPosition = 1

		if t.MessageIndex+1 == len(t.TextInput) && t.lineCounter+1 == len(t.lines) {
			fmt.Printf("resetting text\n")
			t.TextInput = []string{}
			t.lines = []string{}
			t.NextMessage = false
			t.lineCounter = 0
			t.MessageIndex = 0
		}

		t.lineCounter++
	}

	if t.lineCounter == len(t.lines) {
		if t.MessageIndex < len(t.TextInput)-1 {
			t.MessageIndex++
			t.configureLines()
			for _, line := range t.lines {
				fmt.Printf("%s\n", line)
			}
		} else {
			t.MessageIndex = 0
		}

		t.lineCounter = 0
		t.NextMessage = false
	}
}

func (t *TextPrinter) configureLines() {
	t.lines = t.MessageWrapperToLines(t.TextInput[t.MessageIndex])
}

func (t *TextPrinter) textWrapper(text string) string {
	if t.charactersPerLine < 10 {
		log.Fatalf(`line length too small: %d\n`, t.charactersPerLine)
	}
	if len(text) <= t.charactersPerLine+1 {
		return text
	}
	letters := strings.Split(text, "")
	letter := letters[t.charactersPerLine+1]
	i := t.charactersPerLine + 1
	for letter != " " {
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
		output = append(output, strings.TrimSpace(wrappedLine))
		//len [9:]
		text = text[len(wrappedLine):]
	}

	return output
}

func (t *TextPrinter) printText() (output string) {
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
	t.MessageIndex = 0
	t.StatusText[0].SetText("")
	t.StatusText[1].SetText("")
	t.StatusText[2].SetText("")
	t.TextInput = []string{}
	t.lines = []string{}
	t.lineCounter = 0
}

func (t *TextPrinter) ResetTPMessageTriggerNext() {
	t.stringPosition = 1

	t.StatusText[0].SetText("")
	t.StatusText[1].SetText("")
	t.StatusText[2].SetText("")

	//if there are more lines of the message trigger the printer again

	t.NextMessage = true
}
