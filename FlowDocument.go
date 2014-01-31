//
// A FlowDocument is used for display multiple sentences with word wrapping in a confined
// amount of space. On top of that, it provides information as to how to how many lines
// are actually visible.
//

package main
import "termbox-go"
import "container/list"

type FlowDocument struct {
    width, height int
    lines *list.List
}
type flowDocumentLine struct {
    line []rune
    fg, bg termbox.Attribute
}
func NewFlowDocument(width, height int) *FlowDocument {
    return &FlowDocument{ width: width, height: height, lines : list.New() } 
}
// Adds a paragraph. If that filled up the document, it returns false.
func (f *FlowDocument) AddParagraph(paragraph string, fg, bg termbox.Attribute) bool {
    if (f.lines.Len() >= f.height) { return true }

    // loop until the entire string has been processed, or the flow document is full.
    
    paragraphRunes := []rune(paragraph)
    
    for f.lines.Len() < f.height && len(paragraphRunes) > 0 {
        workingLine := paragraphRunes
        if len(workingLine) > f.width {
            workingLine = TrimLength(workingLine, f.width)
            paragraphRunes = paragraphRunes[f.width:]//:len(paragraph) - f.width - 1]
        } else {
            paragraphRunes = []rune("")
        }
        
        f.lines.PushBack(&flowDocumentLine{line:workingLine, fg:fg, bg:bg})        
    }    
    
    return f.lines.Len() >= f.height
}
// get lines to display.

func (f *FlowDocument) Write(r *ConsoleRange, x, y int) {
    // ---- todo 
    
    i := 0
    for e := f.lines.Front(); e != nil; e = e.Next() {
        line := e.Value.(*flowDocumentLine)
        r.DisplayRunesWithColor(line.line, x , y+i, line.fg, line.bg)
        i++
    }
}

func (f *FlowDocument) IsFull() bool {
    return f.lines.Len() >= f.height
}