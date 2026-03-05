package helpers

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/commands/git_commands"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type CommitsHelper struct {
	c *HelperCommon

	loadingHelper *LoadingHelper
	aiHelper      *AIHelper

	getCommitSummary              func() string
	setCommitSummary              func(string)
	getCommitDescription          func() string
	getUnwrappedCommitDescription func() string
	setCommitDescription          func(string)
}

func NewCommitsHelper(
	c *HelperCommon,
	loadingHelper *LoadingHelper,
	aiHelper *AIHelper,
	getCommitSummary func() string,
	setCommitSummary func(string),
	getCommitDescription func() string,
	getUnwrappedCommitDescription func() string,
	setCommitDescription func(string),
) *CommitsHelper {
	return &CommitsHelper{
		c:                             c,
		loadingHelper:                 loadingHelper,
		aiHelper:                      aiHelper,
		getCommitSummary:              getCommitSummary,
		setCommitSummary:              setCommitSummary,
		getCommitDescription:          getCommitDescription,
		getUnwrappedCommitDescription: getUnwrappedCommitDescription,
		setCommitDescription:          setCommitDescription,
	}
}

// SetAIHelper sets the AI helper after initialization
func (self *CommitsHelper) SetAIHelper(aiHelper *AIHelper) {
	self.aiHelper = aiHelper
}

func (self *CommitsHelper) SplitCommitMessageAndDescription(message string) (string, string) {
	msg, description, _ := strings.Cut(message, "\n")
	return msg, strings.TrimSpace(description)
}

func (self *CommitsHelper) SetMessageAndDescriptionInView(message string) {
	summary, description := self.SplitCommitMessageAndDescription(message)

	self.setCommitSummary(summary)
	self.setCommitDescription(description)
	self.c.Contexts().CommitMessage.RenderSubtitle()
}

func (self *CommitsHelper) JoinCommitMessageAndUnwrappedDescription() string {
	if len(self.getUnwrappedCommitDescription()) == 0 {
		return self.getCommitSummary()
	}
	return self.getCommitSummary() + "\n" + self.getUnwrappedCommitDescription()
}

func TryRemoveHardLineBreaks(message string, autoWrapWidth int) string {
	lastHardLineStart := 0
	result := message
	for i, b := range message {
		if b == '\n' {
			// Try to make this a soft linebreak by turning it into a space, and
			// checking whether it still wraps to the same result then.
			str := message[lastHardLineStart:i] + " " + message[i+1:]
			softLineBreakIndices := gocui.AutoWrapContent(str, autoWrapWidth)

			// See if auto-wrapping inserted a soft line break:
			if len(softLineBreakIndices) > 0 && softLineBreakIndices[0] == i-lastHardLineStart+1 {
				// It did, so change it to a space in the result.
				result = result[:i] + " " + result[i+1:]
			}
			lastHardLineStart = i + 1
		}
	}

	return result
}

func (self *CommitsHelper) SwitchToEditor() error {
	message := lo.Ternary(len(self.getCommitDescription()) == 0,
		self.getCommitSummary(),
		self.getCommitSummary()+"\n\n"+self.getCommitDescription())
	filepath := filepath.Join(self.c.OS().GetTempDir(), self.c.Git().RepoPaths.RepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".msg")
	err := self.c.OS().CreateFileWithContent(filepath, message)
	if err != nil {
		return err
	}

	self.CloseCommitMessagePanel()

	return self.c.Contexts().CommitMessage.SwitchToEditor(filepath)
}

func (self *CommitsHelper) UpdateCommitPanelView(message string) {
	if message != "" {
		self.SetMessageAndDescriptionInView(message)
		return
	}

	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		preservedMessage := self.c.Contexts().CommitMessage.GetPreservedMessageAndLogError()
		self.SetMessageAndDescriptionInView(preservedMessage)
		return
	}

	self.SetMessageAndDescriptionInView("")
}

type OpenCommitMessagePanelOpts struct {
	CommitIndex      int
	SummaryTitle     string
	DescriptionTitle string
	PreserveMessage  bool
	OnConfirm        func(summary string, description string) error
	OnSwitchToEditor func(string) error
	InitialMessage   string

	// The following two fields are only for the display of the "(hooks
	// disabled)" display in the commit message panel. They have no effect on
	// the actual behavior; make sure what you are passing in matches that.
	// Leave unassigned if the concept of skipping hooks doesn't make sense for
	// what you are doing, e.g. when creating a tag.
	ForceSkipHooks  bool
	SkipHooksPrefix string
}

func (self *CommitsHelper) OpenCommitMessagePanel(opts *OpenCommitMessagePanelOpts) {
	onConfirm := func(summary string, description string) error {
		self.CloseCommitMessagePanel()

		return opts.OnConfirm(summary, description)
	}

	self.c.Contexts().CommitMessage.SetPanelState(
		opts.CommitIndex,
		opts.SummaryTitle,
		opts.DescriptionTitle,
		opts.PreserveMessage,
		opts.InitialMessage,
		onConfirm,
		opts.OnSwitchToEditor,
		opts.ForceSkipHooks,
		opts.SkipHooksPrefix,
	)

	self.UpdateCommitPanelView(opts.InitialMessage)

	self.c.Context().Push(self.c.Contexts().CommitMessage, types.OnFocusOpts{})
}

func (self *CommitsHelper) ClearPreservedCommitMessage() {
	self.c.Contexts().CommitMessage.SetPreservedMessageAndLogError("")
}

func (self *CommitsHelper) HandleCommitConfirm() error {
	summary, description := self.getCommitSummary(), self.getCommitDescription()

	if summary == "" {
		return errors.New(self.c.Tr.CommitWithoutMessageErr)
	}

	err := self.c.Contexts().CommitMessage.OnConfirm(summary, description)
	if err != nil {
		return err
	}

	return nil
}

func (self *CommitsHelper) CloseCommitMessagePanel() {
	if self.c.Contexts().CommitMessage.GetPreserveMessage() {
		message := self.JoinCommitMessageAndUnwrappedDescription()
		if message != self.c.Contexts().CommitMessage.GetInitialMessage() {
			self.c.Contexts().CommitMessage.SetPreservedMessageAndLogError(message)
		}
	} else {
		self.SetMessageAndDescriptionInView("")
	}

	self.c.Contexts().CommitMessage.SetHistoryMessage("")

	self.c.Views().CommitMessage.Visible = false
	self.c.Views().CommitDescription.Visible = false

	self.c.Context().Pop()
}

func (self *CommitsHelper) OpenCommitMenu(suggestionFunc func(string) []*types.Suggestion) error {
	var disabledReasonForOpenInEditor *types.DisabledReason
	if !self.c.Contexts().CommitMessage.CanSwitchToEditor() {
		disabledReasonForOpenInEditor = &types.DisabledReason{
			Text: self.c.Tr.CommandDoesNotSupportOpeningInEditor,
		}
	}

	menuItems := []*types.MenuItem{
		{
			Label: self.c.Tr.OpenInEditor,
			OnPress: func() error {
				return self.SwitchToEditor()
			},
			Key:            'e',
			DisabledReason: disabledReasonForOpenInEditor,
		},
		{
			Label: self.c.Tr.AddCoAuthor,
			OnPress: func() error {
				return self.addCoAuthor(suggestionFunc)
			},
			Key: 'c',
		},
		{
			Label: self.c.Tr.PasteCommitMessageFromClipboard,
			OnPress: func() error {
				return self.pasteCommitMessageFromClipboard()
			},
			Key: 'p',
		},
	}
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CommitMenuTitle,
		Items: menuItems,
	})
}

func (self *CommitsHelper) addCoAuthor(suggestionFunc func(string) []*types.Suggestion) error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.AddCoAuthorPromptTitle,
		FindSuggestionsFunc: suggestionFunc,
		HandleConfirm: func(value string) error {
			commitDescription := self.getCommitDescription()
			commitDescription = git_commands.AddCoAuthorToDescription(commitDescription, value)
			self.setCommitDescription(commitDescription)
			return nil
		},
	})

	return nil
}

func (self *CommitsHelper) pasteCommitMessageFromClipboard() error {
	message, err := self.c.OS().PasteFromClipboard()
	if err != nil {
		return err
	}
	if message == "" {
		return nil
	}

	currentMessage := self.JoinCommitMessageAndUnwrappedDescription()
	return self.c.ConfirmIf(currentMessage != "", types.ConfirmOpts{
		Title:  self.c.Tr.PasteCommitMessageFromClipboard,
		Prompt: self.c.Tr.SurePasteCommitMessage,
		HandleConfirm: func() error {
			self.SetMessageAndDescriptionInView(message)
			return nil
		},
	})
}

func (self *CommitsHelper) AIGenerateCommitMessage() error {
	if self.c.AIManager == nil {
		return self.aiHelper.ShowFirstTimeWizard()
	}

	self.loadingHelper.WithCenteredLoadingStatus(self.c.Tr.AIGeneratingStatus, func(_ gocui.Task) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		rawDiff, err := self.c.Git().Diff.GetDiff(true)
		if err != nil {
			return err
		}
		if strings.TrimSpace(rawDiff) == "" {
			return errors.New(self.c.Tr.AINoStagedChanges)
		}

		diff := FilterDiffForAI(rawDiff, self.c.Tr)
		const maxDiffChars = 120_000
		safetyNote := ""
		if len(diff) > maxDiffChars {
			diff = diff[:maxDiffChars]
			safetyNote = self.c.Tr.AICommitPromptTruncated
		}

		output, err := self.c.AIManager.RunSkill(ctx, "commit_msg", map[string]any{
			"diff":         diff,
			"project_type": self.detectProjectType(),
			"safety_note":  safetyNote,
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return errors.New(self.c.Tr.AIGenerationCancelled)
			}
			return self.aiHelper.HandleAIError(err)
		}

		message := strings.TrimSpace(output.Content)
		if message == "" {
			return errors.New(self.c.Tr.AICommitEmptyResponse)
		}

		self.SetMessageAndDescriptionInView(message)
		return nil
	})
	return nil
}

// AIGenerateCommitMessageAndOpen 生成 AI 提交信息后调用 openPanel 打开提交面板。
func (self *CommitsHelper) AIGenerateCommitMessageAndOpen(openPanel func(message string) error) error {
	if self.c.AIManager == nil {
		return self.aiHelper.ShowFirstTimeWizard()
	}

	self.loadingHelper.WithCenteredLoadingStatus(self.c.Tr.AIGeneratingStatus, func(_ gocui.Task) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		rawDiff, err := self.c.Git().Diff.GetDiff(true)
		if err != nil {
			return err
		}
		if strings.TrimSpace(rawDiff) == "" {
			return errors.New(self.c.Tr.AINoStagedChanges)
		}

		diff := FilterDiffForAI(rawDiff, self.c.Tr)
		const maxDiffChars = 120_000
		safetyNote := ""
		if len(diff) > maxDiffChars {
			diff = diff[:maxDiffChars]
			safetyNote = self.c.Tr.AICommitPromptTruncated
		}

		output, err := self.c.AIManager.RunSkill(ctx, "commit_msg", map[string]any{
			"diff":         diff,
			"project_type": self.detectProjectType(),
			"safety_note":  safetyNote,
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return errors.New(self.c.Tr.AIGenerationCancelled)
			}
			return self.aiHelper.HandleAIError(err)
		}

		message := strings.TrimSpace(output.Content)
		if message == "" {
			return errors.New(self.c.Tr.AICommitEmptyResponse)
		}

		return openPanel(message)
	})
	return nil
}


// detectProjectType detects the project type based on file extensions.
func (self *CommitsHelper) detectProjectType() string {
	files := self.c.Model().Files

	counts := make(map[string]int)
	typeMap := map[string]string{
		".go":   "Go",
		".ts":   "TypeScript",
		".tsx":  "TypeScript",
		".js":   "JavaScript",
		".jsx":  "JavaScript",
		".py":   "Python",
		".java": "Java",
		".rs":   "Rust",
		".cpp":  "C++",
		".c":    "C",
		".cs":   "C#",
		".rb":   "Ruby",
		".php":  "PHP",
	}

	for _, f := range files {
		path := strings.ToLower(f.Path)
		for ext, lang := range typeMap {
			if strings.HasSuffix(path, ext) {
				counts[lang]++
				break
			}
		}
	}

	// Find the most common type
	maxCount := 0
	projectType := "Mixed"
	for lang, count := range counts {
		if count > maxCount {
			maxCount = count
			projectType = lang
		}
	}

	return projectType
}


