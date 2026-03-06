package aii18n

import (
	"fmt"

	"github.com/dswcpp/lazygit/pkg/i18n"
)

// Translator provides translation functions for AI modules
type Translator struct {
	tr *i18n.TranslationSet
}

// NewTranslator creates a new translator
func NewTranslator(tr *i18n.TranslationSet) *Translator {
	return &Translator{tr: tr}
}

// Common translations
func (t *Translator) Cancel() string                { return t.tr.AICancel }
func (t *Translator) OK() string                    { return t.tr.AIOK }
func (t *Translator) Confirm() string               { return t.tr.AIConfirm }
func (t *Translator) Yes() string                   { return t.tr.AIYes }
func (t *Translator) No() string                    { return t.tr.AINo }
func (t *Translator) Success() string               { return t.tr.AISuccess }
func (t *Translator) Failed() string                { return t.tr.AIFailed }
func (t *Translator) Warning() string               { return t.tr.AIWarning }
func (t *Translator) Unknown() string               { return t.tr.AIUnknown }
func (t *Translator) Executing() string             { return t.tr.AIExecuting }
func (t *Translator) Thinking() string              { return t.tr.AIThinking }
func (t *Translator) Idle() string                  { return t.tr.AIIdle }
func (t *Translator) Cancelled() string             { return t.tr.AICancelled }
func (t *Translator) ThinkingInProgress() string    { return t.tr.AIThinkingInProgress }

// Agent translations
func (t *Translator) AgentToolNotAllowedInPlanning(tool string) string {
	return fmt.Sprintf(t.tr.AIAgentToolNotAllowedInPlanning, tool)
}

func (t *Translator) AgentCriticalStepFailed(step, reason string) string {
	return fmt.Sprintf(t.tr.AIAgentCriticalStepFailed, step, reason)
}

func (t *Translator) AgentStepTimeout(duration, step string) string {
	return fmt.Sprintf(t.tr.AIAgentStepTimeout, duration, step)
}

func (t *Translator) AgentUserRejectedTool(tool string) string {
	return fmt.Sprintf(t.tr.AIAgentUserRejectedTool, tool)
}

func (t *Translator) AgentResolveConflictManually() string { return t.tr.AIAgentResolveConflictManually }
func (t *Translator) AgentSetUpstreamBranch() string       { return t.tr.AIAgentSetUpstreamBranch }
func (t *Translator) AgentConflict() string                { return t.tr.AIAgentConflict }
func (t *Translator) AgentToolName() string                { return t.tr.AIAgentToolName }
func (t *Translator) AgentStageFilesFirst() string         { return t.tr.AIAgentStageFilesFirst }
func (t *Translator) AgentPossibleReasons() string         { return t.tr.AIAgentPossibleReasons }
func (t *Translator) AgentExampleCommitMsg() string        { return t.tr.AIAgentExampleCommitMsg }
func (t *Translator) AgentDont() string                    { return t.tr.AIAgentDont }

func (t *Translator) AgentRepoStatusAndUserInstruction(status, instruction string) string {
	return fmt.Sprintf(t.tr.AIAgentRepoStatusAndUserInstruction, status, instruction)
}

func (t *Translator) AgentUnknownTool(tool string) string {
	return fmt.Sprintf(t.tr.AIAgentUnknownTool, tool)
}

func (t *Translator) AgentUserRejectedExecution(tool string) string {
	return fmt.Sprintf(t.tr.AIAgentUserRejectedExecution, tool)
}

func (t *Translator) AgentMaxStepsReached(maxSteps int) string {
	return fmt.Sprintf(t.tr.AIAgentMaxStepsReached, maxSteps)
}

func (t *Translator) AgentToolLabel(name, desc, perm string) string {
	return fmt.Sprintf(t.tr.AIAgentToolLabel, name, desc, perm)
}

func (t *Translator) AgentDescriptionLabel() string { return t.tr.AIAgentDescriptionLabel }
func (t *Translator) AgentPermissionLabel() string  { return t.tr.AIAgentPermissionLabel }
func (t *Translator) AgentParamsLabel() string      { return t.tr.AIAgentParamsLabel }

// Tool translations
func (t *Translator) ToolMissingParam(param string) string {
	return fmt.Sprintf(t.tr.AIToolMissingParam, param)
}

func (t *Translator) ToolMissingNameParam() string    { return t.tr.AIToolMissingNameParam }
func (t *Translator) ToolMissingPathParam() string    { return t.tr.AIToolMissingPathParam }
func (t *Translator) ToolMissingMessageParam() string { return t.tr.AIToolMissingMessageParam }
func (t *Translator) ToolMissingHashParam() string    { return t.tr.AIToolMissingHashParam }
func (t *Translator) ToolFilePath() string            { return t.tr.AIToolFilePath }
func (t *Translator) ToolBranchName() string          { return t.tr.AIToolBranchName }
func (t *Translator) ToolTagName() string             { return t.tr.AIToolTagName }
func (t *Translator) ToolCommitMessage() string       { return t.tr.AIToolCommitMessage }
func (t *Translator) ToolNoChanges() string           { return t.tr.AIToolNoChanges }
func (t *Translator) ToolWorkingDir() string          { return t.tr.AIToolWorkingDir }
func (t *Translator) ToolStagingArea() string         { return t.tr.AIToolStagingArea }
func (t *Translator) ToolTargetRefOrHash() string     { return t.tr.AIToolTargetRefOrHash }
func (t *Translator) ToolResetSteps() string          { return t.tr.AIToolResetSteps }
func (t *Translator) ToolStashIndex() string          { return t.tr.AIToolStashIndex }
func (t *Translator) ToolMaxLines() string            { return t.tr.AIToolMaxLines }
func (t *Translator) ToolTargetRef() string           { return t.tr.AIToolTargetRef }

func (t *Translator) ToolPushConfigError(err error) string {
	return fmt.Sprintf(t.tr.AIToolPushConfigError, err)
}

func (t *Translator) ToolRebasedTo(target string) string {
	return fmt.Sprintf(t.tr.AIToolRebasedTo, target)
}

func (t *Translator) ToolRenameFailed(err error) string {
	return fmt.Sprintf(t.tr.AIToolRenameFailed, err)
}

func (t *Translator) ToolDiscardChangesFailed(err error) string {
	return fmt.Sprintf(t.tr.AIToolDiscardChangesFailed, err)
}

func (t *Translator) ToolParam() string { return t.tr.AIToolParam }
func (t *Translator) ToolValue() string { return t.tr.AIToolValue }

// Skill translations
func (t *Translator) SkillCurrentBranch(branch string) string {
	return fmt.Sprintf(t.tr.AISkillCurrentBranch, branch)
}

func (t *Translator) SkillBranchNameOnly() string    { return t.tr.AISkillBranchNameOnly }
func (t *Translator) SkillBranchNameFormat() string  { return t.tr.AISkillBranchNameFormat }
func (t *Translator) SkillWindowsGitBash() string    { return t.tr.AISkillWindowsGitBash }
func (t *Translator) SkillOutputJSONArray() string   { return t.tr.AISkillOutputJSONArray }
func (t *Translator) SkillExplanation() string       { return t.tr.AISkillExplanation }
func (t *Translator) SkillCommitSubject() string     { return t.tr.AISkillCommitSubject }
func (t *Translator) SkillTestScenario() string      { return t.tr.AISkillTestScenario }
func (t *Translator) SkillOutputCommitMsg() string   { return t.tr.AISkillOutputCommitMsg }
func (t *Translator) SkillRefactorScenario() string  { return t.tr.AISkillRefactorScenario }
func (t *Translator) SkillGeneratePRDesc() string    { return t.tr.AISkillGeneratePRDesc }
func (t *Translator) SkillPRSummary() string         { return t.tr.AISkillPRSummary }
func (t *Translator) SkillPRTesting() string         { return t.tr.AISkillPRTesting }
func (t *Translator) SkillCodeChanges() string       { return t.tr.AISkillCodeChanges }
func (t *Translator) SkillDiffSummary() string       { return t.tr.AISkillDiffSummary }
func (t *Translator) SkillRepoContext() string       { return t.tr.AISkillRepoContext }
func (t *Translator) SkillCodeChangesTitle() string  { return t.tr.AISkillCodeChangesTitle }
func (t *Translator) SkillCommitHistory() string     { return t.tr.AISkillCommitHistory }

func (t *Translator) SkillRuntime(runtime string) string {
	return fmt.Sprintf(t.tr.AISkillRuntime, runtime)
}

func (t *Translator) SkillBranchInfo(from, to string) string {
	return fmt.Sprintf(t.tr.AISkillBranchInfo, from, to)
}

// Repository context translations
func (t *Translator) MoreItems(count int) string {
	return fmt.Sprintf(t.tr.AIMoreItems, count)
}

func (t *Translator) RepoWorkingDirClean() string { return t.tr.AIRepoWorkingDirClean }

func (t *Translator) RepoInProgress(operation string) string {
	return fmt.Sprintf(t.tr.AIRepoInProgress, operation)
}

func (t *Translator) RepoRemoteSynced(remote string) string {
	return fmt.Sprintf(t.tr.AIRepoRemoteSynced, remote)
}

func (t *Translator) RepoChanges(total, staged, unstaged, untracked int) string {
	return fmt.Sprintf(t.tr.AIRepoChanges, total, staged, unstaged, untracked)
}

func (t *Translator) RepoRemoteAheadBehind(remote, ahead, behind string) string {
	return fmt.Sprintf(t.tr.AIRepoRemoteAheadBehind, remote, ahead, behind)
}

func (t *Translator) RepoBranch(branch string) string {
	return fmt.Sprintf(t.tr.AIRepoBranch, branch)
}

func (t *Translator) RepoRecentCommits() string { return t.tr.AIRepoRecentCommits }

func (t *Translator) RepoStashCount(count int) string {
	return fmt.Sprintf(t.tr.AIRepoStashCount, count)
}

// Manager translations
func (t *Translator) ManagerGenerateBranchName() string { return t.tr.AIManagerGenerateBranchName }
func (t *Translator) ManagerParam() string              { return t.tr.AIManagerParam }
func (t *Translator) ManagerValue() string              { return t.tr.AIManagerValue }
func (t *Translator) ManagerStagedDiff() string         { return t.tr.AIManagerStagedDiff }
func (t *Translator) ManagerFeatureDesc() string        { return t.tr.AIManagerFeatureDesc }
func (t *Translator) ManagerGenerateCommitMsg() string  { return t.tr.AIManagerGenerateCommitMsg }

// Analyze tool translations
func (t *Translator) AnalyzeToolDescription() string { return t.tr.AIAnalyzeToolDescription }
func (t *Translator) AnalyzeToolStagedParam() string { return t.tr.AIAnalyzeToolStagedParam }
func (t *Translator) AnalyzeToolFocusParam() string  { return t.tr.AIAnalyzeToolFocusParam }
func (t *Translator) AnalyzeWorkingDirClean() string { return t.tr.AIAnalyzeWorkingDirClean }

func (t *Translator) AnalyzeNoChanges(label string) string {
	return fmt.Sprintf(t.tr.AIAnalyzeNoChanges, label)
}

func (t *Translator) AnalyzeCancelled() string { return t.tr.AIAnalyzeCancelled }

func (t *Translator) AnalyzeFailed(err error) string {
	return fmt.Sprintf(t.tr.AIAnalyzeFailed, err)
}

func (t *Translator) AnalyzeReportTitle() string { return t.tr.AIAnalyzeReportTitle }

func (t *Translator) AnalyzeReportTitleWithFocus(focus string) string {
	return fmt.Sprintf(t.tr.AIAnalyzeReportTitleWithFocus, focus)
}

func (t *Translator) AnalyzeFileCount(total, success, fail int) string {
	return fmt.Sprintf(t.tr.AIAnalyzeFileCount, total, success, fail)
}

func (t *Translator) AnalyzeTotalLines(lines int) string {
	return fmt.Sprintf(t.tr.AIAnalyzeTotalLines, lines)
}

func (t *Translator) AnalyzeDetailedAnalysis() string { return t.tr.AIAnalyzeDetailedAnalysis }

func (t *Translator) AnalyzeAnalysisFailed(err string) string {
	return fmt.Sprintf(t.tr.AIAnalyzeAnalysisFailed, err)
}

func (t *Translator) AnalyzeNoChangesInfo() string       { return t.tr.AIAnalyzeNoChangesInfo }
func (t *Translator) AnalyzeOverallSuggestions() string  { return t.tr.AIAnalyzeOverallSuggestions }
func (t *Translator) AnalyzeSuggestion1() string         { return t.tr.AIAnalyzeSuggestion1 }
func (t *Translator) AnalyzeSuggestion2() string         { return t.tr.AIAnalyzeSuggestion2 }
func (t *Translator) AnalyzeSuggestion3() string         { return t.tr.AIAnalyzeSuggestion3 }
func (t *Translator) AnalyzeCodeReviewExpert() string    { return t.tr.AIAnalyzeCodeReviewExpert }
func (t *Translator) AnalyzePromptIntro() string         { return t.tr.AIAnalyzePromptIntro }
func (t *Translator) AnalyzeMainChanges() string         { return t.tr.AIAnalyzeMainChanges }
func (t *Translator) AnalyzePotentialIssues() string     { return t.tr.AIAnalyzePotentialIssues }
func (t *Translator) AnalyzeImprovementSuggestions() string { return t.tr.AIAnalyzeImprovementSuggestions }

func (t *Translator) AnalyzeFileLabel(path string) string {
	return fmt.Sprintf(t.tr.AIAnalyzeFileLabel, path)
}

func (t *Translator) AnalyzeFocusLabel(focus string) string {
	return fmt.Sprintf(t.tr.AIAnalyzeFocusLabel, focus)
}

// Commit message skill translations
func (t *Translator) SkillCommitMsgSystemPrompt() string { return t.tr.AISkillCommitMsgSystemPrompt }
func (t *Translator) SkillRepoBackground() string        { return t.tr.AISkillCommitMsgRepoBackground }
func (t *Translator) SkillCodeChangesSection() string    { return t.tr.AISkillCommitMsgCodeChanges }
func (t *Translator) SkillOutputRules() string           { return t.tr.AISkillCommitMsgOutputRules }
func (t *Translator) SkillFormatExample() string         { return t.tr.AISkillCommitMsgFormatExample }
func (t *Translator) SkillTypeList() string              { return t.tr.AISkillCommitMsgTypeList }
func (t *Translator) SkillSubjectRules() string          { return t.tr.AISkillCommitMsgSubjectRules }
func (t *Translator) SkillScopeOptional() string         { return t.tr.AISkillCommitMsgScopeOptional }
func (t *Translator) SkillBodyRequired() string          { return t.tr.AISkillCommitMsgBodyRequired }
func (t *Translator) SkillScenarioBugfix() string        { return t.tr.AISkillCommitMsgScenarioBugfix }
func (t *Translator) SkillScenarioRefactor() string      { return t.tr.AISkillCommitMsgScenarioRefactor }
func (t *Translator) SkillScenarioDocs() string          { return t.tr.AISkillCommitMsgScenarioDocs }
func (t *Translator) SkillScenarioTest() string          { return t.tr.AISkillCommitMsgScenarioTest }
func (t *Translator) SkillScenarioDefault() string       { return t.tr.AISkillCommitMsgScenarioDefault }

// Branch name skill translations
func (t *Translator) SkillBranchNamePromptIntro() string      { return t.tr.AISkillBranchNamePromptIntro }
func (t *Translator) SkillBranchNameStagedFiles() string      { return t.tr.AISkillBranchNameStagedFiles }
func (t *Translator) SkillBranchNameUnstagedFiles() string    { return t.tr.AISkillBranchNameUnstagedFiles }
func (t *Translator) SkillBranchNameDiffSummaryTitle() string { return t.tr.AISkillBranchNameDiffSummaryTitle }
func (t *Translator) SkillBranchNameRules() string            { return t.tr.AISkillBranchNameRules }
func (t *Translator) SkillBranchNameFormatRule() string       { return t.tr.AISkillBranchNameFormatRule }
func (t *Translator) SkillBranchNameTypeRule() string         { return t.tr.AISkillBranchNameTypeRule }
func (t *Translator) SkillBranchNameDescRule() string         { return t.tr.AISkillBranchNameDescRule }
func (t *Translator) SkillBranchNameOutputRule() string       { return t.tr.AISkillBranchNameOutputRule }
func (t *Translator) SkillBranchNameSystemPrompt() string     { return t.tr.AISkillBranchNameSystemPrompt }

func (t *Translator) SkillBranchNameMoreFiles(count int) string {
	return fmt.Sprintf(t.tr.AISkillBranchNameMoreFiles, count)
}

// PR description skill translations
func (t *Translator) SkillPRDescSystemPrompt() string       { return t.tr.AISkillPRDescSystemPrompt }
func (t *Translator) SkillPRDescCommitHistory() string      { return t.tr.AISkillPRDescCommitHistory }
func (t *Translator) SkillPRDescCodeChangesSection() string { return t.tr.AISkillPRDescCodeChangesSection }
func (t *Translator) SkillPRDescGeneratePrompt() string     { return t.tr.AISkillPRDescGeneratePrompt }
func (t *Translator) SkillPRDescSummarySection() string     { return t.tr.AISkillPRDescSummarySection }
func (t *Translator) SkillPRDescChangesSection() string     { return t.tr.AISkillPRDescChangesSection }
func (t *Translator) SkillPRDescTestingSection() string     { return t.tr.AISkillPRDescTestingSection }

func (t *Translator) SkillPRDescBranchInfo(from, to string) string {
	return fmt.Sprintf(t.tr.AISkillPRDescBranchInfo, from, to)
}

// Shell command skill translations
func (t *Translator) SkillShellCmdSystemPrompt() string       { return t.tr.AISkillShellCmdSystemPrompt }
func (t *Translator) SkillShellCmdOutputFormat() string       { return t.tr.AISkillShellCmdOutputFormat }
func (t *Translator) SkillShellCmdCommandField() string       { return t.tr.AISkillShellCmdCommandField }
func (t *Translator) SkillShellCmdExplanationField() string   { return t.tr.AISkillShellCmdExplanationField }
func (t *Translator) SkillShellCmdRiskLevelField() string     { return t.tr.AISkillShellCmdRiskLevelField }
func (t *Translator) SkillShellCmdAlternativesField() string  { return t.tr.AISkillShellCmdAlternativesField }
func (t *Translator) SkillShellCmdOutputNote() string         { return t.tr.AISkillShellCmdOutputNote }
func (t *Translator) SkillShellCmdWindowsHint() string        { return t.tr.AISkillShellCmdWindowsHint }
func (t *Translator) SkillShellCmdMacOSHint() string          { return t.tr.AISkillShellCmdMacOSHint }
func (t *Translator) SkillShellCmdLinuxHint() string          { return t.tr.AISkillShellCmdLinuxHint }

func (t *Translator) SkillShellCmdRuntime(runtime string) string {
	return fmt.Sprintf(t.tr.AISkillShellCmdRuntime, runtime)
}

func (t *Translator) SkillShellCmdRepoStatus(status string) string {
	return fmt.Sprintf(t.tr.AISkillShellCmdRepoStatus, status)
}

func (t *Translator) SkillShellCmdUserIntent(intent string) string {
	return fmt.Sprintf(t.tr.AISkillShellCmdUserIntent, intent)
}

// Additional Chat translations (extended from line 157-171)
func (t *Translator) ChatWelcomeSystem() string               { return t.tr.AIChatWelcomeSystem }
func (t *Translator) ChatWelcomeMessage() string              { return t.tr.AIChatWelcomeMessage }
func (t *Translator) ChatConfigPrompt() string                { return t.tr.AIChatConfigPrompt }
func (t *Translator) ChatPreviousContext() string             { return t.tr.AIChatPreviousContext }
func (t *Translator) ChatNoContentToCopy() string             { return t.tr.AIChatNoContentToCopy }
func (t *Translator) ChatNoExecutableReply() string           { return t.tr.AIChatNoExecutableReply }
func (t *Translator) ChatConfirmExecution() string            { return t.tr.AIChatConfirmExecution }
func (t *Translator) ChatExecutionPlan() string               { return t.tr.AIChatExecutionPlan }
func (t *Translator) ChatNotInitialized() string              { return t.tr.AIChatNotInitialized }
func (t *Translator) ChatRequestFailed() string               { return t.tr.AIChatRequestFailed }
func (t *Translator) ChatCopyFailed() string                  { return t.tr.AIChatCopyFailed }
func (t *Translator) ChatCopiedToClipboard() string           { return t.tr.AIChatCopiedToClipboard }
func (t *Translator) ChatNoCommandsFound() string             { return t.tr.AIChatNoCommandsFound }
func (t *Translator) ChatClearHistoryTitle() string           { return t.tr.AIChatClearHistoryTitle }
func (t *Translator) ChatClearHistoryPrompt() string          { return t.tr.AIChatClearHistoryPrompt }
func (t *Translator) ChatHistoryCleared() string              { return t.tr.AIChatHistoryCleared }
func (t *Translator) ChatHowCanIHelp() string                 { return t.tr.AIChatHowCanIHelp }
func (t *Translator) ChatGenerationStopped() string           { return t.tr.AIChatGenerationStopped }
func (t *Translator) ChatCompleted() string                   { return t.tr.AIChatCompleted }
func (t *Translator) ChatWaitingConfirm() string              { return t.tr.AIChatWaitingConfirm }
func (t *Translator) ChatConfirmPrompt() string               { return t.tr.AIChatConfirmPrompt }
func (t *Translator) ChatExecutingPlan() string               { return t.tr.AIChatExecutingPlan }
func (t *Translator) ChatGeneratingReply() string             { return t.tr.AIChatGeneratingReply }
func (t *Translator) ChatCallingTool() string                 { return t.tr.AIChatCallingTool }
func (t *Translator) ChatToolCompleted() string               { return t.tr.AIChatToolCompleted }
func (t *Translator) ChatToolFailed() string                  { return t.tr.AIChatToolFailed }
func (t *Translator) ChatPlanGenerated() string               { return t.tr.AIChatPlanGenerated }
func (t *Translator) ChatStatusLabel() string                 { return t.tr.AIChatStatusLabel }
func (t *Translator) ChatActionLabel() string                 { return t.tr.AIChatActionLabel }
func (t *Translator) ChatGreeting() string                    { return t.tr.AIChatGreeting }
func (t *Translator) ChatCapabilities() string                { return t.tr.AIChatCapabilities }
func (t *Translator) ChatInputPrompt() string                 { return t.tr.AIChatInputPrompt }
func (t *Translator) ChatStoppedGeneration() string           { return t.tr.AIChatStoppedGeneration }
func (t *Translator) ChatCallingToolPrefix() string           { return t.tr.AIChatCallingToolPrefix }
func (t *Translator) ChatToolCompletedPrefix() string         { return t.tr.AIChatToolCompletedPrefix }
func (t *Translator) ChatToolFailedPrefix() string            { return t.tr.AIChatToolFailedPrefix }

// Two Phase Agent translations
func (t *Translator) TwoPhaseAgentSystemPromptIntro() string        { return t.tr.AITwoPhaseAgentSystemPromptIntro }
func (t *Translator) TwoPhaseAgentWorkflowTitle() string            { return t.tr.AITwoPhaseAgentWorkflowTitle }
func (t *Translator) TwoPhaseAgentWorkflowStep1() string            { return t.tr.AITwoPhaseAgentWorkflowStep1 }
func (t *Translator) TwoPhaseAgentWorkflowStep2() string            { return t.tr.AITwoPhaseAgentWorkflowStep2 }
func (t *Translator) TwoPhaseAgentWorkflowStep2Sub1() string        { return t.tr.AITwoPhaseAgentWorkflowStep2Sub1 }
func (t *Translator) TwoPhaseAgentWorkflowStep2Sub2() string        { return t.tr.AITwoPhaseAgentWorkflowStep2Sub2 }
func (t *Translator) TwoPhaseAgentWorkflowStep2Sub3() string        { return t.tr.AITwoPhaseAgentWorkflowStep2Sub3 }
func (t *Translator) TwoPhaseAgentWorkflowStep3() string            { return t.tr.AITwoPhaseAgentWorkflowStep3 }
func (t *Translator) TwoPhaseAgentWorkflowStep3Sub1() string        { return t.tr.AITwoPhaseAgentWorkflowStep3Sub1 }
func (t *Translator) TwoPhaseAgentWorkflowStep3Sub2() string        { return t.tr.AITwoPhaseAgentWorkflowStep3Sub2 }
func (t *Translator) TwoPhaseAgentWorkflowStep4() string            { return t.tr.AITwoPhaseAgentWorkflowStep4 }
func (t *Translator) TwoPhaseAgentWorkflowStep5() string            { return t.tr.AITwoPhaseAgentWorkflowStep5 }
func (t *Translator) TwoPhaseAgentWorkflowStep6() string            { return t.tr.AITwoPhaseAgentWorkflowStep6 }
func (t *Translator) TwoPhaseAgentToolNameTitle() string            { return t.tr.AITwoPhaseAgentToolNameTitle }
func (t *Translator) TwoPhaseAgentToolNameIntro() string            { return t.tr.AITwoPhaseAgentToolNameIntro }
func (t *Translator) TwoPhaseAgentToolNameStageFile() string        { return t.tr.AITwoPhaseAgentToolNameStageFile }
func (t *Translator) TwoPhaseAgentToolNameDontUseAdd() string       { return t.tr.AITwoPhaseAgentToolNameDontUseAdd }
func (t *Translator) TwoPhaseAgentToolNameCommit() string           { return t.tr.AITwoPhaseAgentToolNameCommit }
func (t *Translator) TwoPhaseAgentToolNameDontUseGitCommit() string { return t.tr.AITwoPhaseAgentToolNameDontUseGitCommit }
func (t *Translator) TwoPhaseAgentToolNameCheckout() string         { return t.tr.AITwoPhaseAgentToolNameCheckout }
func (t *Translator) TwoPhaseAgentToolNameDontUseSwitch() string    { return t.tr.AITwoPhaseAgentToolNameDontUseSwitch }
func (t *Translator) TwoPhaseAgentToolNameCreateBranch() string     { return t.tr.AITwoPhaseAgentToolNameCreateBranch }
func (t *Translator) TwoPhaseAgentToolNameDontUseBranch() string    { return t.tr.AITwoPhaseAgentToolNameDontUseBranch }
func (t *Translator) TwoPhaseAgentSpecialToolTitle() string         { return t.tr.AITwoPhaseAgentSpecialToolTitle }
func (t *Translator) TwoPhaseAgentSpecialToolIntro() string         { return t.tr.AITwoPhaseAgentSpecialToolIntro }
func (t *Translator) TwoPhaseAgentSpecialToolUsage1() string        { return t.tr.AITwoPhaseAgentSpecialToolUsage1 }
func (t *Translator) TwoPhaseAgentSpecialToolUsage2() string        { return t.tr.AITwoPhaseAgentSpecialToolUsage2 }
func (t *Translator) TwoPhaseAgentSpecialToolUsage3() string        { return t.tr.AITwoPhaseAgentSpecialToolUsage3 }
func (t *Translator) TwoPhaseAgentSpecialToolExample() string       { return t.tr.AITwoPhaseAgentSpecialToolExample }
func (t *Translator) TwoPhaseAgentSpecialToolExampleReturn() string { return t.tr.AITwoPhaseAgentSpecialToolExampleReturn }
func (t *Translator) TwoPhaseAgentSpecialToolExamplePlan() string   { return t.tr.AITwoPhaseAgentSpecialToolExamplePlan }
func (t *Translator) TwoPhaseAgentPlanFormatTitle() string          { return t.tr.AITwoPhaseAgentPlanFormatTitle }
func (t *Translator) TwoPhaseAgentPlanFormatExample() string        { return t.tr.AITwoPhaseAgentPlanFormatExample }
func (t *Translator) TwoPhaseAgentNotesTitle() string               { return t.tr.AITwoPhaseAgentNotesTitle }
func (t *Translator) TwoPhaseAgentNotesParam() string               { return t.tr.AITwoPhaseAgentNotesParam }
func (t *Translator) TwoPhaseAgentNotesCriticalTrue() string        { return t.tr.AITwoPhaseAgentNotesCriticalTrue }
func (t *Translator) TwoPhaseAgentNotesCriticalFalse() string       { return t.tr.AITwoPhaseAgentNotesCriticalFalse }
func (t *Translator) TwoPhaseAgentNotesMinimal() string             { return t.tr.AITwoPhaseAgentNotesMinimal }
func (t *Translator) TwoPhaseAgentExecuting() string                { return t.tr.AITwoPhaseAgentExecuting }
func (t *Translator) TwoPhaseAgentRepoStatusTitle() string          { return t.tr.AITwoPhaseAgentRepoStatusTitle }
func (t *Translator) TwoPhaseAgentUserInstructionTitle() string     { return t.tr.AITwoPhaseAgentUserInstructionTitle }
func (t *Translator) TwoPhaseAgentPlanAdjustment() string           { return t.tr.AITwoPhaseAgentPlanAdjustment }
func (t *Translator) TwoPhaseAgentExecutionCancelled() string       { return t.tr.AITwoPhaseAgentExecutionCancelled }
func (t *Translator) TwoPhaseAgentPlanValidationFailed() string     { return t.tr.AITwoPhaseAgentPlanValidationFailed }
func (t *Translator) TwoPhaseAgentPlanErrorsIntro() string          { return t.tr.AITwoPhaseAgentPlanErrorsIntro }
func (t *Translator) TwoPhaseAgentPlanRegeneratePrompt() string     { return t.tr.AITwoPhaseAgentPlanRegeneratePrompt }
func (t *Translator) TwoPhaseAgentContinueAnalysis() string         { return t.tr.AITwoPhaseAgentContinueAnalysis }
func (t *Translator) TwoPhaseAgentSystemPrefix() string             { return t.tr.AITwoPhaseAgentSystemPrefix }

func (t *Translator) TwoPhaseAgentToolCallWarning(tool string, count int) string {
	return fmt.Sprintf(t.tr.AITwoPhaseAgentToolCallWarning, tool, count)
}

func (t *Translator) TwoPhaseAgentUserFeedbackPrompt(feedback string) string {
	return fmt.Sprintf(t.tr.AITwoPhaseAgentUserFeedbackPrompt, feedback)
}

func (t *Translator) TwoPhaseAgentToolResultPrefix(toolName, output string) string {
	return fmt.Sprintf(t.tr.AITwoPhaseAgentToolResultPrefix, toolName, output)
}

func (t *Translator) TwoPhaseAgentMaxStepsExceeded(maxSteps int) string {
	return fmt.Sprintf(t.tr.AITwoPhaseAgentMaxStepsExceeded, maxSteps)
}

// BuildPlanningSystemPrompt builds the planning phase system prompt
func (t *Translator) BuildPlanningSystemPrompt() string {
	return t.TwoPhaseAgentSystemPromptIntro() +
		t.TwoPhaseAgentWorkflowTitle() +
		t.TwoPhaseAgentWorkflowStep1() +
		t.TwoPhaseAgentWorkflowStep2() +
		t.TwoPhaseAgentWorkflowStep2Sub1() +
		t.TwoPhaseAgentWorkflowStep2Sub2() +
		t.TwoPhaseAgentWorkflowStep2Sub3() +
		t.TwoPhaseAgentWorkflowStep3() +
		t.TwoPhaseAgentWorkflowStep3Sub1() +
		t.TwoPhaseAgentWorkflowStep3Sub2() +
		t.TwoPhaseAgentWorkflowStep4() +
		t.TwoPhaseAgentWorkflowStep5() +
		t.TwoPhaseAgentWorkflowStep6() +
		t.TwoPhaseAgentToolNameTitle() +
		t.TwoPhaseAgentToolNameIntro() +
		t.TwoPhaseAgentToolNameStageFile() +
		t.TwoPhaseAgentToolNameDontUseAdd() +
		t.TwoPhaseAgentToolNameCommit() +
		t.TwoPhaseAgentToolNameDontUseGitCommit() +
		t.TwoPhaseAgentToolNameCheckout() +
		t.TwoPhaseAgentToolNameDontUseSwitch() +
		t.TwoPhaseAgentToolNameCreateBranch() +
		t.TwoPhaseAgentToolNameDontUseBranch() +
		t.TwoPhaseAgentSpecialToolTitle() +
		t.TwoPhaseAgentSpecialToolIntro() +
		t.TwoPhaseAgentSpecialToolUsage1() +
		t.TwoPhaseAgentSpecialToolUsage2() +
		t.TwoPhaseAgentSpecialToolUsage3() +
		t.TwoPhaseAgentSpecialToolExample() +
		t.TwoPhaseAgentSpecialToolExampleReturn() +
		t.TwoPhaseAgentSpecialToolExamplePlan() +
		t.TwoPhaseAgentPlanFormatTitle() +
		t.TwoPhaseAgentPlanFormatExample() +
		t.TwoPhaseAgentNotesTitle() +
		t.TwoPhaseAgentNotesParam() +
		t.TwoPhaseAgentNotesCriticalTrue() +
		t.TwoPhaseAgentNotesCriticalFalse() +
		t.TwoPhaseAgentNotesMinimal()
}
