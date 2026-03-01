package mergeconflicts

import (
	"bytes"

	"github.com/dswcpp/lazygit/pkg/gui/style"
	"github.com/dswcpp/lazygit/pkg/theme"
	"github.com/dswcpp/lazygit/pkg/utils"
)

func ColoredConflictFile(state *State) string {
	content := state.GetContent()
	if len(state.conflicts) == 0 {
		return content
	}
	conflict, remainingConflicts := shiftConflict(state.conflicts)
	var outputBuffer bytes.Buffer
	for i, line := range utils.SplitLines(content) {
		textStyle := theme.DefaultTextColor
		if conflict.isMarkerLine(i) {
			textStyle = style.FgRed
		}

		if i == conflict.end && len(remainingConflicts) > 0 {
			conflict, remainingConflicts = shiftConflict(remainingConflicts)
		}
		outputBuffer.WriteString(textStyle.Sprint(line) + "\n")
	}
	return outputBuffer.String()
}

func shiftConflict(conflicts []*mergeConflict) (*mergeConflict, []*mergeConflict) {
	return conflicts[0], conflicts[1:]
}
