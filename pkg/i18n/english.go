/*

Todo list when making a new translation
- Copy this file and rename it to the language you want to translate to like someLanguage.go
- Change the EnglishTranslationSet() name to the language you want to translate to like SomeLanguageTranslationSet()
- Add an entry of someLanguage in GetTranslationSets()
- Remove this todo and the about section

*/

package i18n

type TranslationSet struct {
	NotEnoughSpace                        string
	DiffTitle                             string
	FilesTitle                            string
	BranchesTitle                         string
	CommitsTitle                          string
	StashTitle                            string
	SnakeTitle                            string
	EasterEgg                             string
	UnstagedChanges                       string
	StagedChanges                         string
	StagingTitle                          string
	MergingTitle                          string
	NormalTitle                           string
	LogTitle                              string
	LogXOfYTitle                          string
	CommitSummary                         string
	CredentialsUsername                   string
	CredentialsPassword                   string
	CredentialsPassphrase                 string
	CredentialsPIN                        string
	CredentialsToken                      string
	PassUnameWrong                        string
	Commit                                string
	CommitTooltip                         string
	AmendLastCommit                       string
	AmendLastCommitTitle                  string
	SureToAmend                           string
	NoCommitToAmend                       string
	CommitChangesWithEditor               string
	FindBaseCommitForFixup                string
	FindBaseCommitForFixupTooltip         string
	NoBaseCommitsFound                    string
	MultipleBaseCommitsFoundStaged        string
	MultipleBaseCommitsFoundUnstaged      string
	BaseCommitIsAlreadyOnMainBranch       string
	BaseCommitIsNotInCurrentView          string
	HunksWithOnlyAddedLinesWarning        string
	StatusTitle                           string
	GlobalTitle                           string
	Execute                               string
	Stage                                 string
	StageTooltip                          string
	ToggleStagedAll                       string
	ToggleStagedAllTooltip                string
	ToggleTreeView                        string
	ToggleTreeViewTooltip                 string
	OpenDiffTool                          string
	OpenMergeTool                         string
	Refresh                               string
	RefreshTooltip                        string
	Push                                  string
	Pull                                  string
	PushTooltip                           string
	PullTooltip                           string
	FileFilter                            string
	CopyToClipboardMenu                   string
	CopyFileName                          string
	CopyRelativeFilePath                  string
	CopyAbsoluteFilePath                  string
	CopyFileDiffTooltip                   string
	CopySelectedDiff                      string
	CopyAllFilesDiff                      string
	CopyFileContent                       string
	NoContentToCopyError                  string
	FileNameCopiedToast                   string
	FilePathCopiedToast                   string
	FileDiffCopiedToast                   string
	AllFilesDiffCopiedToast               string
	FileContentCopiedToast                string
	FilterStagedFiles                     string
	FilterUnstagedFiles                   string
	FilterTrackedFiles                    string
	FilterUntrackedFiles                  string
	NoFilter                              string
	FilterLabelStagedFiles                string
	FilterLabelUnstagedFiles              string
	FilterLabelTrackedFiles               string
	FilterLabelUntrackedFiles             string
	FilterLabelConflictingFiles           string
	MergeConflictsTitle                   string
	MergeConflictDescription_DD           string
	MergeConflictDescription_AU           string
	MergeConflictDescription_UA           string
	MergeConflictDescription_DU           string
	MergeConflictDescription_UD           string
	MergeConflictIncomingDiff             string
	MergeConflictCurrentDiff              string
	MergeConflictPressEnterToResolve      string
	MergeConflictKeepFile                 string
	MergeConflictDeleteFile               string
	Checkout                              string
	CheckoutTooltip                       string
	CantCheckoutBranchWhilePulling        string
	TagCheckoutTooltip                    string
	RemoteBranchCheckoutTooltip           string
	CantPullOrPushSameBranchTwice         string
	NoChangedFiles                        string
	SoftReset                             string
	AlreadyCheckedOutBranch               string
	SureForceCheckout                     string
	ForceCheckoutBranch                   string
	BranchName                            string
	NewBranchNameBranchOff                string
	CantDeleteCheckOutBranch              string
	DeleteBranchTitle                     string
	DeleteBranchesTitle                   string
	DeleteLocalBranch                     string
	DeleteLocalBranches                   string
	DeleteRemoteBranchPrompt              string
	DeleteRemoteBranchesPrompt            string
	DeleteLocalAndRemoteBranchPrompt      string
	DeleteLocalAndRemoteBranchesPrompt    string
	ForceDeleteBranchTitle                string
	ForceDeleteBranchMessage              string
	ForceDeleteBranchesMessage            string
	RebaseBranch                          string
	RebaseBranchTooltip                   string
	CantRebaseOntoSelf                    string
	CantMergeBranchIntoItself             string
	ForceCheckout                         string
	ForceCheckoutTooltip                  string
	CheckoutByName                        string
	CheckoutByNameTooltip                 string
	CheckoutPreviousBranch                string
	RemoteBranchCheckoutTitle             string
	RemoteBranchCheckoutPrompt            string
	CheckoutTypeNewBranch                 string
	CheckoutTypeNewBranchTooltip          string
	CheckoutTypeDetachedHead              string
	CheckoutTypeDetachedHeadTooltip       string
	NewBranch                             string
	NewBranchFromStashTooltip             string
	MoveCommitsToNewBranch                string
	MoveCommitsToNewBranchTooltip         string
	MoveCommitsToNewBranchFromMainPrompt  string
	MoveCommitsToNewBranchMenuPrompt      string
	MoveCommitsToNewBranchFromBaseItem    string
	MoveCommitsToNewBranchStackedItem     string
	CannotMoveCommitsFromDetachedHead     string
	CannotMoveCommitsNoUpstream           string
	CannotMoveCommitsBehindUpstream       string
	CannotMoveCommitsNoUnpushedCommits    string
	NoBranchesThisRepo                    string
	CommitWithoutMessageErr               string
	Close                                 string
	CloseCancel                           string
	Confirm                               string
	Quit                                  string
	SquashTooltip                         string
	CannotSquashOrFixupFirstCommit        string
	CannotSquashOrFixupMergeCommit        string
	Fixup                                 string
	FixupTooltip                          string
	FixupKeepMessage                      string
	FixupKeepMessageTooltip               string
	SetFixupMessage                       string
	SetFixupMessageTooltip                string
	FixupDiscardMessage                   string
	FixupDiscardMessageTooltip            string
	SureSquashThisCommit                  string
	Squash                                string
	PickCommitTooltip                     string
	Pick                                  string
	Edit                                  string
	Revert                                string
	RevertCommitTooltip                   string
	Reword                                string
	CommitRewordTooltip                   string
	DropCommit                            string
	DropCommitTooltip                     string
	MoveDownCommit                        string
	MoveUpCommit                          string
	CannotMoveAnyFurther                  string
	CannotMoveMergeCommit                 string
	EditCommit                            string
	EditCommitTooltip                     string
	AmendCommitTooltip                    string
	Amend                                 string
	ResetAuthor                           string
	ResetAuthorTooltip                    string
	SetAuthor                             string
	SetAuthorTooltip                      string
	AddCoAuthor                           string
	AmendCommitAttribute                  string
	AmendCommitAttributeTooltip           string
	SetAuthorPromptTitle                  string
	AddCoAuthorPromptTitle                string
	AddCoAuthorTooltip                    string
	RewordCommitEditor                    string
	NoCommitsThisBranch                   string
	UpdateRefHere                         string
	ExecCommandHere                       string
	Error                                 string
	Undo                                  string
	UndoReflog                            string
	RedoReflog                            string
	UndoTooltip                           string
	RedoTooltip                           string
	UndoMergeResolveTooltip               string
	DiscardAllTooltip                     string
	DiscardUnstagedTooltip                string
	DiscardUnstagedDisabled               string
	Pop                                   string
	StashPopTooltip                       string
	Drop                                  string
	StashDropTooltip                      string
	Apply                                 string
	StashApplyTooltip                     string
	NoStashEntries                        string
	StashDrop                             string
	SureDropStashEntry                    string
	StashPop                              string
	SurePopStashEntry                     string
	StashApply                            string
	SureApplyStashEntry                   string
	NoTrackedStagedFilesStash             string
	NoFilesToStash                        string
	StashChanges                          string
	RenameStash                           string
	RenameStashPrompt                     string
	OpenConfig                            string
	EditConfig                            string
	ForcePush                             string
	ForcePushPrompt                       string
	ForcePushDisabled                     string
	UpdatesRejected                       string
	UpdatesRejectedAndForcePushDisabled   string
	CheckForUpdate                        string
	CheckingForUpdates                    string
	UpdateAvailableTitle                  string
	UpdateAvailable                       string
	UpdateInProgressWaitingStatus         string
	UpdateCompletedTitle                  string
	UpdateCompleted                       string
	FailedToRetrieveLatestVersionErr      string
	OnLatestVersionErr                    string
	MajorVersionErr                       string
	CouldNotFindBinaryErr                 string
	UpdateFailedErr                       string
	ConfirmQuitDuringUpdateTitle          string
	ConfirmQuitDuringUpdate               string
	IntroPopupMessage                     string
	NonReloadableConfigWarningTitle       string
	NonReloadableConfigWarning            string
	GitconfigParseErr                     string
	EditFile                              string
	EditFileTooltip                       string
	OpenFile                              string
	OpenFileTooltip                       string
	OpenInEditor                          string
	IgnoreFile                            string
	ExcludeFile                           string
	RefreshFiles                          string
	FocusMainView                         string
	Merge                                 string
	MergeBranchTooltip                    string
	RegularMergeFastForward               string
	RegularMergeFastForwardTooltip        string
	CannotFastForwardMerge                string
	RegularMergeNonFastForward            string
	RegularMergeNonFastForwardTooltip     string
	SquashMergeUncommitted                string
	SquashMergeUncommittedTooltip         string
	SquashMergeCommitted                  string
	SquashMergeCommittedTooltip           string
	ConfirmQuit                           string
	SwitchRepo                            string
	AllBranchesLogGraph                   string
	UnsupportedGitService                 string
	CopyPullRequestURL                    string
	NoBranchOnRemote                      string
	Fetch                                 string
	FetchTooltip                          string
	CollapseAll                           string
	CollapseAllTooltip                    string
	ExpandAll                             string
	ExpandAllTooltip                      string
	DisabledInFlatView                    string
	FileEnter                             string
	FileEnterTooltip                      string
	StageSelectionTooltip                 string
	DiscardSelection                      string
	DiscardSelectionTooltip               string
	ToggleSelectHunk                      string
	SelectHunk                            string
	SelectLineByLine                      string
	ToggleSelectHunkTooltip               string
	HunkStagingHint                       string
	ToggleSelectionForPatch               string
	EditHunk                              string
	EditHunkTooltip                       string
	ToggleStagingView                     string
	ToggleStagingViewTooltip              string
	ReturnToFilesPanel                    string
	FastForward                           string
	FastForwardTooltip                    string
	FastForwarding                        string
	FoundConflictsTitle                   string
	ViewConflictsMenuItem                 string
	AbortMenuItem                         string
	PickHunk                              string
	PickAllHunks                          string
	ViewMergeRebaseOptions                string
	ViewMergeRebaseOptionsTooltip         string
	ViewMergeOptions                      string
	ViewRebaseOptions                     string
	ViewCherryPickOptions                 string
	ViewRevertOptions                     string
	NotMergingOrRebasing                  string
	AlreadyRebasing                       string
	NotMidRebase                          string
	MustSelectFixupCommit                 string
	RecentRepos                           string
	MergeOptionsTitle                     string
	RebaseOptionsTitle                    string
	CherryPickOptionsTitle                string
	RevertOptionsTitle                    string
	CommitSummaryTitle                    string
	CommitDescriptionTitle                string
	CommitDescriptionSubTitle             string
	CommitDescriptionFooter               string
	CommitDescriptionFooterTwoBindings    string
	CommitHooksDisabledSubTitle           string
	LocalBranchesTitle                    string
	SearchTitle                           string
	TagsTitle                             string
	MenuTitle                             string
	CommitMenuTitle                       string
	RemotesTitle                          string
	RemoteBranchesTitle                   string
	PatchBuildingTitle                    string
	InformationTitle                      string
	SecondaryTitle                        string
	ReflogCommitsTitle                    string
	ConflictsResolved                     string
	Continue                              string
	UnstagedFilesAfterConflictsResolved   string
	RebasingTitle                         string
	RebasingFromBaseCommitTitle           string
	SimpleRebase                          string
	InteractiveRebase                     string
	RebaseOntoBaseBranch                  string
	InteractiveRebaseTooltip              string
	RebaseOntoBaseBranchTooltip           string
	MustSelectTodoCommits                 string
	FwdNoUpstream                         string
	FwdNoLocalUpstream                    string
	FwdCommitsToPush                      string
	PullRequestNoUpstream                 string
	ErrorOccurred                         string
	ConflictLabel                         string
	PendingRebaseTodosSectionHeader       string
	PendingCherryPicksSectionHeader       string
	PendingRevertsSectionHeader           string
	CommitsSectionHeader                  string
	YouDied                               string
	RewordNotSupported                    string
	ChangingThisActionIsNotAllowed        string
	NotAllowedMidCherryPickOrRevert       string
	PickIsOnlyAllowedDuringRebase         string
	DroppingMergeRequiresSingleSelection  string
	CherryPickCopy                        string
	CherryPickCopyTooltip                 string
	PasteCommits                          string
	SureCherryPick                        string
	CherryPick                            string
	CannotCherryPickNonCommit             string
	PrevHunk                              string
	NextHunk                              string
	PrevConflict                          string
	NextConflict                          string
	SelectPrevHunk                        string
	SelectNextHunk                        string
	ScrollDown                            string
	ScrollUp                              string
	ScrollUpMainWindow                    string
	ScrollDownMainWindow                  string
	SuspendApp                            string
	CannotSuspendApp                      string
	AmendCommitTitle                      string
	AmendCommitPrompt                     string
	AmendCommitWithConflictsMenuPrompt    string
	AmendCommitWithConflictsContinue      string
	AmendCommitWithConflictsAmend         string
	DropCommitTitle                       string
	DropCommitPrompt                      string
	DropUpdateRefPrompt                   string
	DropMergeCommitPrompt                 string
	PullingStatus                         string
	PushingStatus                         string
	FetchingStatus                        string
	SquashingStatus                       string
	FixingStatus                          string
	DeletingStatus                        string
	AIGeneratingStatus                    string
	AIGenerateCommitMessage               string
	AINotEnabled                          string
	AINoStagedChanges                     string
	AIError                               string
	AISettings                            string
	AISettingsEnable                      string
	AISettingsDisable                     string
	AISettingsSetAPIKey                   string
	AISettingsAPIKeyPrompt                string
	AISettingsSetProvider                 string
	AISettingsSetModel                    string
	AISettingsModelPrompt                 string
	AISettingsSetEndpoint                 string
	AISettingsEndpointPrompt              string
	AISettingsSaved                       string
	AISettingsActiveProfile               string
	AISettingsSwitchProfile               string
	AISettingsEditProfile                 string
	AISettingsAddProfile                  string
	AISettingsNoProfiles                  string
	AISettingsProfileName                 string
	AISettingsProfileNamePrompt           string
	AISettingsNewProfileNamePrompt        string
	AISettingsMaxTokens                   string
	AISettingsMaxTokensPrompt             string
	AISettingsTimeout                     string
	AISettingsTimeoutPrompt               string
	AISettingsDeleteProfile               string
	AISettingsDeleteProfileTitle          string
	AISettingsDeleteProfilePrompt         string
	AISettingsCannotDeleteLastProfile     string
	AIAssistant                           string
	AIAssistantTitle                      string
	AIAssistantPrompt                     string
	AIAssistantStatus                     string
	AIAssistantConfirmExecute             string
	AIAssistantNoCommands                 string
	AIAssistantSilentNoCommands           string
	AIAssistantConfirmSilentExecute       string
	AIAssistantExecuting                  string
	AIAssistantExecuteError               string
	AICodeReview                          string
	AICodeReviewTitle                     string
	AIGenerationCancelled                 string
	AICommitPromptTruncated               string
	AICommitEmptyResponse                 string
	AIPromptCurrentBranch                 string
	AIPromptRecentCommits                 string
	AIPromptProjectType                   string
	AIPromptRole                          string
	AIPromptTask                          string
	AIPromptRepoContext                   string
	AIPromptCodeChanges                   string
	AIPromptOutputRules                   string
	AIPromptTypeGuide                     string
	AIPromptScopeGuide                    string
	AIPromptSubjectRequirements           string
	AIPromptBodyOptional                  string
	AIPromptScenarioGuide                 string
	AIPromptOutputFormat                  string
	AIPromptScenarioSmall                 string
	AIPromptScenarioLarge                 string
	AIPromptScenarioBugfix                string
	AIPromptScenarioRefactor              string
	AIPromptScenarioDocs                  string
	AIPromptScenarioTest                  string
	AIPromptScenarioDefault               string
	AICodeReviewConfirmTitle              string
	AICodeReviewConfirmPrompt             string
	AICodeReviewStatus                    string
	AICodeReviewNoDiff                    string
	AICodeReviewCopiedToClipboard         string
	AICodeReviewToggleZoom                string
	NewBranchWithAI                       string
	NewBranchWithAITooltip                string
	AIGeneratingBranchNameStatus          string
	CreatePRWithAIDescription             string
	CreatePRWithAIDescriptionTooltip      string
	CreatePRDirectly                      string
	AIGeneratingPRDescriptionStatus       string
	PRDescriptionCopiedToClipboard        string
	// AI Diff Filter
	AIDiffSkipped                         string
	AIDiffBinaryFile                      string
	AIDiffLockOrGeneratedFile             string
	AIDiffChangeStats                     string
	AIDiffFilesCount                      string
	AIDiffFileTypes                       string
	AIDiffChangeScale                     string
	AIDiffMajorChanges                    string
	AIDiffNewFile                         string
	AIDiffDeletedFile                     string
	AIDiffRenamedFile                     string
	AIDiffModifiedFile                    string
	AIDiffTruncated                       string
	AIDiffSmartTruncated                  string
	// AI Helper Messages
	AICodeReviewCancelled                 string
	AINotEnabledPleaseConfig              string
	AITestingConnection                   string
	AIConnectionTestFailed                string
	AIEmptyResponse                       string
	AIConnectionTestSuccess               string
	AIConnectionTestSuccessDetail         string
	APIKeyCannotBeEmpty                   string
	AIConfigComplete                      string
	AIConfigCompletePrompt                string
	AIWelcomeWizardTitle                  string
	UseDeepSeekRecommended                string
	UseOpenAI                             string
	UseAnthropicClaude                    string
	UseOllamaLocal                        string
	ConfigureLater                        string
	TestCurrentProfile                    string
	SetupProviderAPIKey                   string
	// AI Diff Filter Additional
	AIDiffSkippedFilesNote                string
	// AI Assistant Prompts
	AIAssistantSystemPrompt               string
	AIAssistantRules                      string
	AIAssistantRepoState                  string
	AIAssistantUserRequest                string
	// AI Branch Naming Prompts
	AIBranchNameSystemPrompt              string
	AIBranchNameTask                      string
	AIBranchNameRules                     string
	AIBranchNameChanges                   string
	AIBranchNameDiffSummary               string
	AIBranchNameRequirements              string
	// AI PR Description Prompts
	AIPRDescSystemPrompt                  string
	AIPRDescTask                          string
	AIPRDescCommitHistory                 string
	AIPRDescCodeChanges                   string
	AIPRDescFormatRequirements            string
	AIPRDescSummarySection                string
	AIPRDescChangesSection                string
	AIPRDescTechDetailsSection            string
	AIPRDescTestingSection                string
	AIPRDescOutputRequirements            string
	AIPRDescBranchInfo                    string
	AIPRDescDiffUnavailable               string
	AIPRDescDiffTruncated                 string
	AIPRDescMoreCommits                   string
	// AI Code Review Prompts
	AICodeReviewSystemPrompt              string
	AICodeReviewFile                      string
	AICodeReviewCorePrinciples            string
	AICodeReviewConservative              string
	AICodeReviewRespectContext            string
	AICodeReviewFocusNewLines             string
	AICodeReviewRejectFalsePositives      string
	AICodeReviewSeverityLevels            string
	AICodeReviewCritical                  string
	AICodeReviewMajor                     string
	AICodeReviewMinor                     string
	AICodeReviewNit                       string
	AICodeReviewOutputFormat              string
	AICodeReviewSummarySection            string
	AICodeReviewIssuesSection            string
	AICodeReviewIssueFormat               string
	AICodeReviewNoIssues                  string
	AICodeReviewConclusionSection         string
	AICodeReviewConclusionLGTM            string
	AICodeReviewConclusionHasIssues       string
	AICodeReviewDiffSection               string
	AICodeReviewLanguageHint              string
	AICodeReviewLanguageChecks            string
	// AI Error Messages
	AIRequestTimeout                      string
	APIKeyInvalid                         string
	APIRateLimitExceeded                  string
	NetworkConnectionFailed               string
	ModelNotAvailable                     string
	APIQuotaExhausted                     string
	InputTooLong                          string
	// AI Context Messages
	CurrentBranch                         string
	TrackingRemoteBranchAheadBehind       string
	TrackingRemoteBranchAhead             string
	TrackingRemoteBranchBehind            string
	TrackingRemoteBranchSynced            string
	NotTrackingRemoteBranch               string
	WorkingTreeState                      string
	ChangeStats                           string
	RecentCommits                         string
	ChangedFiles                          string
	MoreFiles                             string
	StashList                             string
	MoreStashes                           string
	// AI Branch and PR Messages
	AINotEnabledConfigFirst               string
	NoChangesForBranchName                string
	AIBranchNameCancelled                 string
	NoCommitsForPRDescription             string
	AIPRDescriptionCancelled              string
	AIGenericError                        string
	ChangedFilesLabel                     string
	StagedFilesLabel                      string
	UnstagedFilesLabel                    string
	DiffTruncatedNote                     string
	AICommitMessageCancelled              string
	DroppingStatus                        string
	MovingStatus                          string
	RebasingStatus                        string
	MergingStatus                         string
	LowercaseRebasingStatus               string
	LowercaseMergingStatus                string
	LowercaseCherryPickingStatus          string
	LowercaseRevertingStatus              string
	AmendingStatus                        string
	CherryPickingStatus                   string
	UndoingStatus                         string
	RedoingStatus                         string
	CheckingOutStatus                     string
	CommittingStatus                      string
	RewordingStatus                       string
	RevertingStatus                       string
	CreatingFixupCommitStatus             string
	MovingCommitsToNewBranchStatus        string
	CommitFiles                           string
	SubCommitsDynamicTitle                string
	CommitFilesDynamicTitle               string
	RemoteBranchesDynamicTitle            string
	ViewItemFiles                         string
	CommitFilesTitle                      string
	CheckoutCommitFileTooltip             string
	CannotCheckoutWithModifiedFilesErr    string
	CanOnlyDiscardFromLocalCommits        string
	Remove                                string
	DiscardOldFileChangeTooltip           string
	DiscardFileChangesTitle               string
	DiscardFileChangesPrompt              string
	DisabledForGPG                        string
	CreateRepo                            string
	BareRepo                              string
	InitialBranch                         string
	NoRecentRepositories                  string
	IncorrectNotARepository               string
	AutoStashTitle                        string
	AutoStashPrompt                       string
	AutoStashForUndo                      string
	AutoStashForCheckout                  string
	AutoStashForNewBranch                 string
	AutoStashForMovingPatchToIndex        string
	AutoStashForCherryPicking             string
	AutoStashForReverting                 string
	Discard                               string
	DiscardChangesTitle                   string
	DiscardFileChangesTooltip             string
	Cancel                                string
	DiscardAllChanges                     string
	DiscardUnstagedChanges                string
	DiscardAllChangesToAllFiles           string
	DiscardAnyUnstagedChanges             string
	DiscardUntrackedFiles                 string
	DiscardStagedChanges                  string
	HardReset                             string
	BranchDeleteTooltip                   string
	TagDeleteTooltip                      string
	Delete                                string
	Reset                                 string
	ResetTooltip                          string
	ViewResetOptions                      string
	FileResetOptionsTooltip               string
	CreateFixupCommit                     string
	CreateFixupCommitTooltip              string
	CreateAmendCommit                     string
	FixupMenu_Fixup                       string
	FixupMenu_FixupTooltip                string
	FixupMenu_AmendWithChanges            string
	FixupMenu_AmendWithChangesTooltip     string
	FixupMenu_AmendWithoutChanges         string
	FixupMenu_AmendWithoutChangesTooltip  string
	SquashAboveCommitsTooltip             string
	SquashCommitsAboveSelectedTooltip     string
	SquashCommitsInCurrentBranchTooltip   string
	SquashAboveCommits                    string
	SquashCommitsInCurrentBranch          string
	SquashCommitsAboveSelectedCommit      string
	CannotSquashCommitsInCurrentBranch    string
	ExecuteShellCommand                   string
	ExecuteShellCommandTooltip            string
	ShellCommand                          string
	ShellCommandAIMode                    string
	ShellCommandDangerousWarning          string
	CommitChangesWithoutHook              string
	ResetTo                               string
	ResetSoftTooltip                      string
	ResetMixedTooltip                     string
	ResetHardTooltip                      string
	ResetHardConfirmation                 string
	PressEnterToReturn                    string
	ViewStashOptions                      string
	ViewStashOptionsTooltip               string
	Stash                                 string
	StashTooltip                          string
	StashAllChanges                       string
	StashStagedChanges                    string
	StashAllChangesKeepIndex              string
	StashUnstagedChanges                  string
	StashIncludeUntrackedChanges          string
	StashOptions                          string
	NotARepository                        string
	WorkingDirectoryDoesNotExist          string
	ScrollLeft                            string
	ScrollRight                           string
	DiscardPatch                          string
	DiscardPatchConfirm                   string
	CantPatchWhileRebasingError           string
	ToggleAddToPatch                      string
	ToggleAddToPatchTooltip               string
	ToggleAllInPatch                      string
	ToggleAllInPatchTooltip               string
	UpdatingPatch                         string
	ViewPatchOptions                      string
	PatchOptionsTitle                     string
	NoPatchError                          string
	EmptyPatchError                       string
	EnterCommitFile                       string
	EnterCommitFileTooltip                string
	ExitCustomPatchBuilder                string
	ExitFocusedMainView                   string
	EnterUpstream                         string
	InvalidUpstream                       string
	NewRemote                             string
	NewRemoteName                         string
	NewRemoteUrl                          string
	AddForkRemote                         string
	AddForkRemoteUsername                 string
	AddForkRemoteTooltip                  string
	IncompatibleForkAlreadyExistsError    string
	NoOriginRemote                        string
	ViewBranches                          string
	EditRemoteName                        string
	EditRemoteUrl                         string
	RemoveRemote                          string
	RemoveRemoteTooltip                   string
	RemoveRemotePrompt                    string
	DeleteRemoteBranch                    string
	DeleteRemoteBranches                  string
	DeleteRemoteBranchTooltip             string
	DeleteLocalAndRemoteBranch            string
	DeleteLocalAndRemoteBranches          string
	SetAsUpstream                         string
	SetAsUpstreamTooltip                  string
	SetUpstream                           string
	UnsetUpstream                         string
	ViewDivergenceFromUpstream            string
	ViewDivergenceFromBaseBranch          string
	CouldNotDetermineBaseBranch           string
	DivergenceSectionHeaderLocal          string
	DivergenceSectionHeaderRemote         string
	ViewUpstreamResetOptions              string
	ViewUpstreamResetOptionsTooltip       string
	ViewUpstreamRebaseOptions             string
	ViewUpstreamRebaseOptionsTooltip      string
	UpstreamGenericName                   string
	SetUpstreamTitle                      string
	SetUpstreamMessage                    string
	EditRemoteTooltip                     string
	TagCommit                             string
	TagCommitTooltip                      string
	TagNameTitle                          string
	TagMessageTitle                       string
	LightweightTag                        string
	AnnotatedTag                          string
	DeleteTagTitle                        string
	DeleteLocalTag                        string
	DeleteRemoteTag                       string
	DeleteLocalAndRemoteTag               string
	SelectRemoteTagUpstream               string
	DeleteRemoteTagPrompt                 string
	DeleteLocalAndRemoteTagPrompt         string
	RemoteTagDeletedMessage               string
	PushTagTitle                          string
	PushTag                               string
	PushTagTooltip                        string
	NewTag                                string
	NewTagTooltip                         string
	CreatingTag                           string
	ForceTag                              string
	ForceTagPrompt                        string
	FetchRemoteTooltip                    string
	CheckoutCommitTooltip                 string
	NoBranchesFoundAtCommitTooltip        string
	GitFlowOptions                        string
	NotAGitFlowBranch                     string
	NewBranchNamePrompt                   string
	IgnoreTracked                         string
	ExcludeTracked                        string
	IgnoreTrackedPrompt                   string
	ExcludeTrackedPrompt                  string
	ViewResetToUpstreamOptions            string
	NextScreenMode                        string
	PrevScreenMode                        string
	CyclePagers                           string
	CyclePagersTooltip                    string
	CyclePagersDisabledReason             string
	StartSearch                           string
	StartFilter                           string
	Keybindings                           string
	KeybindingsLegend                     string
	KeybindingsMenuSectionLocal           string
	KeybindingsMenuSectionGlobal          string
	KeybindingsMenuSectionNavigation      string
	RenameBranch                          string
	Upstream                              string
	BranchUpstreamOptionsTitle            string
	ViewBranchUpstreamOptions             string
	ViewBranchUpstreamOptionsTooltip      string
	UpstreamNotSetError                   string
	UpstreamsNotSetError                  string
	NewGitFlowBranchPrompt                string
	RenameBranchWarning                   string
	OpenKeybindingsMenu                   string
	ResetCherryPick                       string
	ResetCherryPickShort                  string
	NextTab                               string
	PrevTab                               string
	CantUndoWhileRebasing                 string
	CantRedoWhileRebasing                 string
	MustStashWarning                      string
	MustStashTitle                        string
	ConfirmationTitle                     string
	PromptTitle                           string
	PromptInputCannotBeEmptyToast         string
	PrevPage                              string
	NextPage                              string
	GotoTop                               string
	GotoBottom                            string
	FilteringBy                           string
	ResetInParentheses                    string
	OpenFilteringMenu                     string
	OpenFilteringMenuTooltip              string
	FilterBy                              string
	ExitFilterMode                        string
	FilterPathOption                      string
	FilterAuthorOption                    string
	EnterFileName                         string
	EnterAuthor                           string
	FilteringMenuTitle                    string
	WillCancelExistingFilterTooltip       string
	MustExitFilterModeTitle               string
	MustExitFilterModePrompt              string
	Diff                                  string
	EnterRefToDiff                        string
	EnterRefName                          string
	ExitDiffMode                          string
	DiffingMenuTitle                      string
	SwapDiff                              string
	ViewDiffingOptions                    string
	ViewDiffingOptionsTooltip             string
	CancelDiffingMode                     string
	OpenCommandLogMenu                    string
	OpenCommandLogMenuTooltip             string
	OpenAIAssistant                       string
	ShowingGitDiff                        string
	ShowingDiffForRange                   string
	CommitDiff                            string
	CopyCommitHashToClipboard             string
	CommitHash                            string
	CommitURL                             string
	PasteCommitMessageFromClipboard       string
	SurePasteCommitMessage                string
	CommitMessage                         string
	CommitMessageBody                     string
	CommitSubject                         string
	CommitAuthor                          string
	CommitTags                            string
	CopyCommitAttributeToClipboard        string
	CopyCommitAttributeToClipboardTooltip string
	CopyBranchNameToClipboard             string
	CopyTagToClipboard                    string
	CopyPathToClipboard                   string
	CommitPrefixPatternError              string
	CopySelectedTextToClipboard           string
	NoFilesStagedTitle                    string
	NoFilesStagedPrompt                   string
	BranchNotFoundTitle                   string
	BranchNotFoundPrompt                  string
	BranchUnknown                         string
	DiscardChangeTitle                    string
	DiscardChangePrompt                   string
	CreateNewBranchFromCommit             string
	BuildingPatch                         string
	ViewCommits                           string
	MinGitVersionError                    string
	RunningCustomCommandStatus            string
	SubmoduleStashAndReset                string
	AndResetSubmodules                    string
	EnterSubmoduleTooltip                 string
	BackToParentRepo                      string
	Enter                                 string
	CopySubmoduleNameToClipboard          string
	RemoveSubmodule                       string
	RemoveSubmoduleTooltip                string
	RemoveSubmodulePrompt                 string
	ResettingSubmoduleStatus              string
	NewSubmoduleName                      string
	NewSubmoduleUrl                       string
	NewSubmodulePath                      string
	NewSubmodule                          string
	AddingSubmoduleStatus                 string
	UpdateSubmoduleUrl                    string
	UpdatingSubmoduleUrlStatus            string
	EditSubmoduleUrl                      string
	InitializingSubmoduleStatus           string
	InitSubmoduleTooltip                  string
	Update                                string
	Initialize                            string
	SubmoduleUpdateTooltip                string
	UpdatingSubmoduleStatus               string
	BulkInitSubmodules                    string
	BulkUpdateSubmodules                  string
	BulkDeinitSubmodules                  string
	BulkUpdateRecursiveSubmodules         string
	ViewBulkSubmoduleOptions              string
	BulkSubmoduleOptions                  string
	RunningCommand                        string
	SubCommitsTitle                       string
	ExitSubview                           string
	SubmodulesTitle                       string
	NavigationTitle                       string
	SuggestionsCheatsheetTitle            string
	// Unlike the cheatsheet title above, the real suggestions title has a little message saying press tab to focus
	SuggestionsTitle                         string
	SuggestionsSubtitle                      string
	ExtrasTitle                              string
	PullRequestURLCopiedToClipboard          string
	CommitDiffCopiedToClipboard              string
	CommitURLCopiedToClipboard               string
	CommitMessageCopiedToClipboard           string
	CommitMessageBodyCopiedToClipboard       string
	CommitSubjectCopiedToClipboard           string
	CommitAuthorCopiedToClipboard            string
	CommitTagsCopiedToClipboard              string
	CommitHasNoTags                          string
	CommitHasNoMessageBody                   string
	PatchCopiedToClipboard                   string
	MessageCopiedToClipboard                 string
	CopiedToClipboard                        string
	ErrCannotEditDirectory                   string
	ErrCannotCopyContentOfDirectory          string
	ErrStageDirWithInlineMergeConflicts      string
	ErrRepositoryMovedOrDeleted              string
	ErrWorktreeMovedOrRemoved                string
	CommandLog                               string
	ToggleShowCommandLog                     string
	FocusCommandLog                          string
	CopyCommandLog                           string
	CommandLogCopiedToClipboard             string
	CommandLogHeader                         string
	RandomTip                                string
	ToggleWhitespaceInDiffView               string
	ToggleWhitespaceInDiffViewTooltip        string
	IgnoreWhitespaceDiffViewSubTitle         string
	IgnoreWhitespaceNotSupportedHere         string
	IncreaseContextInDiffView                string
	IncreaseContextInDiffViewTooltip         string
	DecreaseContextInDiffView                string
	DecreaseContextInDiffViewTooltip         string
	DiffContextSizeChanged                   string
	IncreaseRenameSimilarityThreshold        string
	IncreaseRenameSimilarityThresholdTooltip string
	DecreaseRenameSimilarityThreshold        string
	DecreaseRenameSimilarityThresholdTooltip string
	RenameSimilarityThresholdChanged         string
	CreatePullRequestOptions                 string
	DefaultBranch                            string
	SelectBranch                             string
	SelectTargetRemote                       string
	NoValidRemoteName                        string
	CreatePullRequest                        string
	SelectConfigFile                         string
	NoConfigFileFoundErr                     string
	LoadingFileSuggestions                   string
	LoadingCommits                           string
	MustSpecifyOriginError                   string
	GitOutput                                string
	GitCommandFailed                         string
	AbortTitle                               string
	AbortPrompt                              string
	OpenLogMenu                              string
	OpenLogMenuTooltip                       string
	LogMenuTitle                             string
	ToggleShowGitGraphAll                    string
	ShowGitGraph                             string
	ShowGitGraphTooltip                      string
	SortOrder                                string
	SortOrderPromptLocalBranches             string
	SortOrderPromptRemoteBranches            string
	SortAlphabetical                         string
	SortByDate                               string
	SortByRecency                            string
	SortBasedOnReflog                        string
	SortOrderPrompt                          string
	SortCommits                              string
	SortCommitsTooltip                       string
	CantChangeContextSizeError               string
	OpenCommitInBrowser                      string
	ViewBisectOptions                        string
	ViewBranchesContainingCommit             string
	ViewBranchesContainingCommitTooltip      string
	NoBranchesContainingCommit               string
	EnterCommitHashToFindBranches            string
	ConfirmRevertCommit                      string
	ConfirmRevertCommitRange                 string
	RewordInEditorTitle                      string
	RewordInEditorPrompt                     string
	CheckoutAutostashPrompt                  string
	HardResetAutostashPrompt                 string
	SoftResetPrompt                          string
	UpstreamGone                             string
	NukeDescription                          string
	NukeTreeConfirmation                     string
	DiscardStagedChangesDescription          string
	EmptyOutput                              string
	Patch                                    string
	CustomPatch                              string
	CommitsCopied                            string
	CommitCopied                             string
	ResetPatch                               string
	ResetPatchTooltip                        string
	ApplyPatch                               string
	ApplyPatchTooltip                        string
	ApplyPatchInReverse                      string
	ApplyPatchInReverseTooltip               string
	RemovePatchFromOriginalCommit            string
	RemovePatchFromOriginalCommitTooltip     string
	MovePatchOutIntoIndex                    string
	MovePatchOutIntoIndexTooltip             string
	MovePatchIntoNewCommit                   string
	MovePatchIntoNewCommitTooltip            string
	MovePatchIntoNewCommitBefore             string
	MovePatchIntoNewCommitBeforeTooltip      string
	MovePatchToSelectedCommit                string
	MovePatchToSelectedCommitTooltip         string
	CopyPatchToClipboard                     string
	MustStageFilesAffectedByPatchTitle       string
	MustStageFilesAffectedByPatchWarning     string
	NoMatchesFor                             string
	MatchesFor                               string
	SearchKeybindings                        string
	SearchPrefix                             string
	FilterPrefix                             string
	FilterPrefixMenu                         string
	ExitSearchMode                           string
	ExitTextFilterMode                       string
	Switch                                   string
	SwitchToWorktree                         string
	SwitchToWorktreeTooltip                  string
	AlreadyCheckedOutByWorktree              string
	BranchCheckedOutByWorktree               string
	SomeBranchesCheckedOutByWorktreeError    string
	DetachWorktreeTooltip                    string
	Switching                                string
	RemoveWorktree                           string
	RemoveWorktreeTitle                      string
	DetachWorktree                           string
	DetachingWorktree                        string
	WorktreesTitle                           string
	WorktreeTitle                            string
	RemoveWorktreePrompt                     string
	ForceRemoveWorktreePrompt                string
	RemovingWorktree                         string
	AddingWorktree                           string
	CantDeleteCurrentWorktree                string
	AlreadyInWorktree                        string
	CantDeleteMainWorktree                   string
	NoWorktreesThisRepo                      string
	MissingWorktree                          string
	MainWorktree                             string
	NewWorktree                              string
	NewWorktreePath                          string
	NewWorktreeBase                          string
	RemoveWorktreeTooltip                    string
	NewBranchName                            string
	NewBranchNameLeaveBlank                  string
	ViewWorktreeOptions                      string
	CreateWorktreeFrom                       string
	CreateWorktreeFromDetached               string
	LcWorktree                               string
	ChangingDirectoryTo                      string
	Name                                     string
	Branch                                   string
	Path                                     string
	MarkedBaseCommitStatus                   string
	MarkAsBaseCommit                         string
	MarkAsBaseCommitTooltip                  string
	CancelMarkedBaseCommit                   string
	MarkedCommitMarker                       string
	FailedToOpenURL                          string
	InvalidLazygitEditURL                    string
	NoCopiedCommits                          string
	DisabledMenuItemPrefix                   string
	QuickStartInteractiveRebase              string
	QuickStartInteractiveRebaseTooltip       string
	CannotQuickStartInteractiveRebase        string
	ToggleRangeSelect                        string
	DismissRangeSelect                       string
	RangeSelectUp                            string
	RangeSelectDown                          string
	RangeSelectNotSupported                  string
	NoItemSelected                           string
	SelectedItemIsNotABranch                 string
	SelectedItemDoesNotHaveFiles             string
	MultiSelectNotSupportedForSubmodules     string
	CommandDoesNotSupportOpeningInEditor     string
	CustomCommands                           string
	NoApplicableCommandsInThisContext        string
	SelectCommitsOfCurrentBranch             string
	Actions                                  Actions
	Bisect                                   Bisect
	Log                                      Log
	BreakingChangesTitle                     string
	BreakingChangesMessage                   string
	BreakingChangesByVersion                 map[string]string
	ViewMergeConflictOptions                 string
	ViewMergeConflictOptionsTooltip          string
	NoFilesWithMergeConflicts                string
	MergeConflictOptionsTitle                string
	UseCurrentChanges                        string
	UseIncomingChanges                       string
	UseBothChanges                           string

	// AI Common
	AICancel                                 string
	AIOK                                     string
	AIConfirm                                string
	AIYes                                    string
	AINo                                     string
	AISuccess                                string
	AIFailed                                 string
	AIWarning                                string
	AIUnknown                                string
	AIExecuting                              string
	AIThinking                               string
	AIIdle                                   string
	AICancelled                              string
	AIThinkingInProgress                     string

	// AI Agent
	AIAgentToolNotAllowedInPlanning          string
	AIAgentCriticalStepFailed                string
	AIAgentStepTimeout                       string
	AIAgentUserRejectedTool                  string
	AIAgentResolveConflictManually           string
	AIAgentSetUpstreamBranch                 string
	AIAgentConflict                          string
	AIAgentToolName                          string
	AIAgentStageFilesFirst                   string
	AIAgentPossibleReasons                   string
	AIAgentExampleCommitMsg                  string
	AIAgentDont                              string
	AIAgentRepoStatusAndUserInstruction      string
	AIAgentUnknownTool                       string
	AIAgentUserRejectedExecution             string
	AIAgentMaxStepsReached                   string
	AIAgentToolLabel                         string
	AIAgentDescriptionLabel                  string
	AIAgentPermissionLabel                   string
	AIAgentParamsLabel                       string

	// AI Tools
	AIToolMissingParam                       string
	AIToolMissingNameParam                   string
	AIToolMissingPathParam                   string
	AIToolMissingMessageParam                string
	AIToolMissingHashParam                   string
	AIToolFilePath                           string
	AIToolBranchName                         string
	AIToolTagName                            string
	AIToolCommitMessage                      string
	AIToolNoChanges                          string
	AIToolWorkingDir                         string
	AIToolStagingArea                        string
	AIToolTargetRefOrHash                    string
	AIToolResetSteps                         string
	AIToolStashIndex                         string
	AIToolMaxLines                           string
	AIToolOffset                             string
	AIToolTargetRef                          string
	AIToolPushConfigError                    string
	AIToolRebasedTo                          string
	AIToolRenameFailed                       string
	AIToolDiscardChangesFailed               string
	AIToolParam                              string
	AIToolValue                              string

	// AI Tools - Schema descriptions
	AIToolGetStatusDesc              string
	AIToolGetStagedDiffDesc          string
	AIToolGetDiffDesc                string
	AIToolGetFileDiffDesc            string
	AIToolGetFileDiffStagedParam     string
	AIToolGetLogDesc                 string
	AIToolGetLogCountParam           string
	AIToolGetBranchesDesc            string
	AIToolGetStashListDesc           string
	AIToolGetRemotesDesc             string
	AIToolGetTagsDesc                string
	AIToolGetStashDiffDesc           string
	AIToolGetStashDiffIndexParam     string
	AIToolGetCommitDiffDesc          string
	AIToolGetCommitDiffHashParam     string
	AIToolGetBranchDiffDesc          string
	AIToolGetBranchDiffBaseParam     string
	AIToolGetBranchDiffTargetParam   string
	AIToolGetBranchDiffEmpty         string
	AIToolStageAllDesc               string
	AIToolStageFileDesc              string
	AIToolUnstageAllDesc             string
	AIToolUnstageFileDesc            string
	AIToolDiscardFileDesc            string
	AIToolCommitDesc                 string
	AIToolCommitMsgParam             string
	AIToolAmendHeadDesc              string
	AIToolAmendMsgParam              string
	AIToolRevertCommitDesc           string
	AIToolRevertHashParam            string
	AIToolResetSoftDesc              string
	AIToolResetMixedDesc             string
	AIToolResetHardDesc              string
	AIToolCherryPickDesc             string
	AIToolCherryPickHashParam        string
	AIToolCheckoutDesc               string
	AIToolCheckoutNameParam          string
	AIToolCreateBranchDesc           string
	AIToolCreateBranchNameParam      string
	AIToolCreateBranchBaseParam      string
	AIToolCreateBranchCheckoutParam  string
	AIToolDeleteBranchDesc           string
	AIToolDeleteBranchForceParam     string
	AIToolRenameBranchDesc           string
	AIToolRenameBranchOldParam       string
	AIToolMergeBranchDesc            string
	AIToolMergeBranchNameParam       string
	AIToolRebaseBranchDesc           string
	AIToolRebaseBranchTargetParam    string
	AIToolStashDesc                  string
	AIToolStashMsgParam              string
	AIToolStashPopDesc               string
	AIToolStashApplyDesc             string
	AIToolStashDropDesc              string
	AIToolCreateTagDesc              string
	AIToolDeleteTagDesc              string
	AIToolPullDesc                   string
	AIToolPullRemoteParam            string
	AIToolPullBranchParam            string
	AIToolFetchDesc                  string
	AIToolPushDesc                   string
	AIToolPushForceDesc              string
	AIToolAbortOperationDesc         string
	AIToolAbortOperationTypeParam    string
	AIToolContinueOperationDesc      string
	AIToolContinueOperationTypeParam string

	// AI Tools - Output messages (success)
	AIToolStagedDiffEmpty               string
	AIToolUnstagedDiffEmpty             string
	AIToolNoStashEntries                string
	AIToolNoRemotes                     string
	AIToolNoTags                        string
	AIToolFileNotInWorkdir              string
	AIToolStashEntryEmpty               string
	AIToolStatusFiles                   string
	AIToolStatusClean                   string
	AIToolStatusInProgress              string
	AIToolStageAllSuccess               string
	AIToolStageFileSuccess              string
	AIToolUnstageAllSuccess             string
	AIToolUnstageFileSuccess            string
	AIToolCommitSuccess                 string
	AIToolAmendSuccess                  string
	AIToolRevertSuccess                 string
	AIToolResetSoftSuccess              string
	AIToolResetMixedSuccess             string
	AIToolResetHardSuccess              string
	AIToolCherryPickSuccess             string
	AIToolCheckoutSuccess               string
	AIToolCreateBranchSuccess           string
	AIToolCreateBranchNoCheckoutSuccess string
	AIToolDeleteBranchSuccess           string
	AIToolRenameBranchSuccess           string
	AIToolMergeBranchSuccess            string
	AIToolStashSuccess                  string
	AIToolStashPopSuccess               string
	AIToolStashApplySuccess             string
	AIToolStashDropSuccess              string
	AIToolCreateTagSuccess              string
	AIToolDeleteTagSuccess              string
	AIToolPullSuccess                   string
	AIToolFetchSuccess                  string
	AIToolPushSuccess                   string
	AIToolPushForceSuccess              string
	AIToolAbortSuccess                  string
	AIToolContinueSuccess               string
	AIToolTruncated                     string

	// AI Tools - Error messages
	AIToolGetStagedDiffFailed   string
	AIToolGetDiffFailed         string
	AIToolGetStashDiffFailed    string
	AIToolGetCommitDiffFailed   string
	AIToolGetBranchDiffFailed   string
	AIToolStageAllFailed        string
	AIToolStageFileFailed       string
	AIToolUnstageAllFailed      string
	AIToolUnstageFileFailed     string
	AIToolCommitFailed          string
	AIToolAmendFailed           string
	AIToolRevertFailed          string
	AIToolResetSoftFailed       string
	AIToolResetMixedFailed      string
	AIToolResetHardFailed       string
	AIToolCherryPickFailed      string
	AIToolCheckoutFailed        string
	AIToolCreateBranchFailed    string
	AIToolDeleteBranchFailed    string
	AIToolMergeBranchFailed     string
	AIToolRebaseBranchFailed    string
	AIToolStashFailed           string
	AIToolStashPopFailed        string
	AIToolStashApplyFailed      string
	AIToolStashDropFailed       string
	AIToolCreateTagFailed       string
	AIToolDeleteTagFailed       string
	AIToolPullFailed            string
	AIToolFetchFailed           string
	AIToolPushFailed            string
	AIToolPushForceFailed       string
	AIToolAbortFailed           string
	AIToolContinueFailed        string
	AIToolMissingTargetParam    string
	AIToolMissingOldOrNameParam string
	AIToolUnknownOperationType  string

	// AI Skills
	AISkillCurrentBranch                     string
	AISkillBranchNameOnly                    string
	AISkillBranchNameFormat                  string
	AISkillWindowsGitBash                    string
	AISkillRuntime                           string
	AISkillOutputJSONArray                   string
	AISkillExplanation                       string
	AISkillCommitSubject                     string
	AISkillTestScenario                      string
	AISkillOutputCommitMsg                   string
	AISkillRefactorScenario                  string
	AISkillGeneratePRDesc                    string
	AISkillPRSummary                         string
	AISkillPRTesting                         string
	AISkillCodeChanges                       string
	AISkillDiffSummary                       string
	AISkillRepoContext                       string
	AISkillCodeChangesTitle                  string
	AISkillBranchInfo                        string
	AISkillCommitHistory                     string

	// AI Chat (GUI)
	AIChatNotEnabled                         string
	AIChatCanInputNext                       string
	AIChatGeneratingPlan                     string
	AIChatTemplateBranchName                 string
	AIChatTemplateTagName                    string
	AIChatTemplateMessage                    string
	AIChatPushingToRemote                    string
	AIChatAbortMerge                         string
	AIChatResolveConflict                    string
	AIChatConflictFiles                      string
	AIChatMergeConflict                      string
	AIChatUncommittedChanges                 string
	AIChatDeleteSuccess                      string
	AIChatConfirmSuffix                      string

	// AI Repository Context
	AIMoreItems                              string
	AIRepoWorkingDirClean                    string
	AIRepoInProgress                         string
	AIRepoRemoteSynced                       string
	AIRepoChanges                            string
	AIRepoRemoteAheadBehind                  string
	AIRepoBranch                             string
	AIRepoRecentCommits                      string
	AIRepoStashCount                         string

	// AI Manager
	AIManagerGenerateBranchName              string
	AIManagerParam                           string
	AIManagerValue                           string
	AIManagerStagedDiff                      string
	AIManagerFeatureDesc                     string
	AIManagerGenerateCommitMsg               string

	// AI Analyze Tool
	AIAnalyzeToolDescription                 string
	AIAnalyzeToolStagedParam                 string
	AIAnalyzeToolFocusParam                  string
	AIAnalyzeWorkingDirClean                 string
	AIAnalyzeNoChanges                       string
	AIAnalyzeCancelled                       string
	AIAnalyzeFailed                          string
	AIAnalyzeReportTitle                     string
	AIAnalyzeReportTitleWithFocus            string
	AIAnalyzeFileCount                       string
	AIAnalyzeTotalLines                      string
	AIAnalyzeDetailedAnalysis                string
	AIAnalyzeAnalysisFailed                  string
	AIAnalyzeNoChangesInfo                   string
	AIAnalyzeOverallSuggestions              string
	AIAnalyzeSuggestion1                     string
	AIAnalyzeSuggestion2                     string
	AIAnalyzeSuggestion3                     string
	AIAnalyzeCodeReviewExpert                string
	AIAnalyzeFileLabel                       string
	AIAnalyzePromptIntro                     string
	AIAnalyzeFocusLabel                      string
	AIAnalyzeMainChanges                     string
	AIAnalyzePotentialIssues                 string
	AIAnalyzeImprovementSuggestions          string

	// Command Completion
	CompletionBranch                         string
	CompletionRemote                         string
	CompletionCommitRef                      string
	CompletionTag                            string
	CompletionGitDesc                        string
	CompletionCdDesc                         string
	CompletionLsDesc                         string
	CompletionPwdDesc                        string
	CompletionCatDesc                        string
	CompletionGrepDesc                       string
	CompletionFindDesc                       string
	CompletionGitAddDesc                     string
	CompletionGitCommitDesc                  string
	CompletionGitPushDesc                    string
	CompletionGitPullDesc                    string
	CompletionGitCheckoutDesc                string
	CompletionGitSwitchDesc                  string
	CompletionGitBranchDesc                  string
	CompletionGitMergeDesc                   string
	CompletionGitRebaseDesc                  string
	CompletionGitResetDesc                   string
	CompletionGitRevertDesc                  string
	CompletionGitStashDesc                   string
	CompletionGitLogDesc                     string
	CompletionGitDiffDesc                    string
	CompletionGitStatusDesc                  string
	CompletionGitTagDesc                     string
	CompletionGitFetchDesc                   string
	CompletionGitCloneDesc                   string
	CompletionGitInitDesc                    string
	CompletionGitCleanDesc                   string
	CompletionGitCherryPickDesc              string
	CompletionGitShowDesc                    string
	CompletionGitRmDesc                      string
	CompletionGitMvDesc                      string
	CompletionGitGrepDesc                    string
	CompletionGitBisectDesc                  string
	CompletionFlagAmendDesc                  string
	CompletionFlagNoEditDesc                 string
	CompletionFlagMDesc                      string
	CompletionFlagADesc                      string
	CompletionFlagAllDesc                    string
	CompletionFlagFixupDesc                  string
	CompletionFlagSignoffDesc                string
	CompletionFlagSDesc                      string
	CompletionFlagNoVerifyDesc               string
	CompletionFlagAllowEmptyDesc             string
	CompletionFlagForceDesc                  string
	CompletionFlagForceWithLeaseDesc         string
	CompletionFlagSetUpstreamDesc            string
	CompletionFlagUDesc                      string
	CompletionFlagTagsDesc                   string
	CompletionFlagDeleteDesc                 string
	CompletionFlagDryRunDesc                 string
	CompletionFlagAllBranchesDesc            string
	CompletionFlagSoftDesc                   string
	CompletionFlagMixedDesc                  string
	CompletionFlagHardDesc                   string
	CompletionStatusConflicted               string
	CompletionStatusPartiallyStaged          string
	CompletionStatusStaged                   string
	CompletionStatusModified                 string
	CompletionStatusUntracked                string
	CompletionStatusTracked                  string

	// AI Command Helper
	AICommandNotEnabled                      string
	AICommandGenerationCancelled             string
	AICommandInvalidFormat                   string
	AICommandExplainPrompt                   string
	AICommandExplainCancelled                string
	AICommandRiskHardReset                   string
	AICommandRiskCleanFdx                    string
	AICommandRiskCleanFd                     string
	AICommandRiskForcePush1                  string
	AICommandRiskForcePush2                  string
	AICommandRiskReflogExpire                string
	AICommandRiskRmRf                        string
	AICommandRiskBranchD                     string
	AICommandRiskRebaseI                     string
	AICommandRiskGcAggressive                string
	AICommandSuggestionHardReset             string
	AICommandSuggestionCleanFdx              string
	AICommandSuggestionCleanFd               string
	AICommandSuggestionForcePush1            string
	AICommandSuggestionForcePush2            string
	AICommandSuggestionBranchD               string

	// AI Skills - Commit Message
	AISkillCommitMsgSystemPrompt             string
	AISkillCommitMsgRepoBackground           string
	AISkillCommitMsgCodeChanges              string
	AISkillCommitMsgOutputRules              string
	AISkillCommitMsgFormatExample            string
	AISkillCommitMsgTypeList                 string
	AISkillCommitMsgSubjectRules             string
	AISkillCommitMsgScopeOptional            string
	AISkillCommitMsgBodyRequired             string
	AISkillCommitMsgScenarioBugfix           string
	AISkillCommitMsgScenarioRefactor         string
	AISkillCommitMsgScenarioDocs             string
	AISkillCommitMsgScenarioTest             string
	AISkillCommitMsgScenarioDefault          string
	AISkillCommitMsgScenarioLarge            string
	AISkillCommitMsgProjectType              string

	// AI Skills - Branch Name
	AISkillBranchNamePromptIntro             string
	AISkillBranchNameStagedFiles             string
	AISkillBranchNameUnstagedFiles           string
	AISkillBranchNameMoreFiles               string
	AISkillBranchNameDiffSummaryTitle        string
	AISkillBranchNameRules                   string
	AISkillBranchNameFormatRule              string
	AISkillBranchNameTypeRule                string
	AISkillBranchNameDescRule                string
	AISkillBranchNameOutputRule              string
	AISkillBranchNameDescriptionHint        string
	AISkillBranchNameSystemPrompt            string

	// AI Skills - PR Description
	AISkillPRDescSystemPrompt                string
	AISkillPRDescBranchInfo                  string
	AISkillPRDescCommitHistory               string
	AISkillPRDescCodeChangesSection          string
	AISkillPRDescGeneratePrompt              string
	AISkillPRDescSummarySection              string
	AISkillPRDescChangesSection              string
	AISkillPRDescBreakingSection             string
	AISkillPRDescTestingSection              string
	AISkillPRDescChecklistSection            string

	// AI Skills - Shell Command
	AISkillShellCmdSystemPrompt              string
	AISkillShellCmdRuntime                   string
	AISkillShellCmdRepoStatus                string
	AISkillShellCmdUserIntent                string
	AISkillShellCmdOutputFormat              string
	AISkillShellCmdCommandField              string
	AISkillShellCmdExplanationField          string
	AISkillShellCmdRiskLevelField            string
	AISkillShellCmdAlternativesField         string
	AISkillShellCmdOutputNote                string
	AISkillShellCmdWindowsHint               string
	AISkillShellCmdMacOSHint                 string
	AISkillShellCmdLinuxHint                 string

	// AI Skills - Code Review
	AISkillCodeReviewSystemPrompt            string

	// AI Skills - Explain Diff
	AISkillExplainDiffSystemPrompt           string

	// AI Skills - Release Notes
	AISkillReleaseNotesSystemPrompt          string

	// AI Skills - Stash Name
	AISkillStashNameSystemPrompt             string

	// AI Chat Helper
	AIChatWelcomeSystem                      string
	AIChatWelcomeMessage                     string
	AIChatConfigPrompt                       string
	AIChatPreviousContext                    string
	AIChatNoContentToCopy                    string
	AIChatNoExecutableReply                  string
	AIChatConfirmExecution                   string
	AIChatExecutionPlan                      string
	AIChatNotInitialized                     string
	AIChatRequestFailed                      string
	AIChatCopyFailed                         string
	AIChatCopiedToClipboard                  string
	AIChatNoCommandsFound                    string
	AIChatClearHistoryTitle                  string
	AIChatClearHistoryPrompt                 string
	AIChatHistoryCleared                     string
	AIChatHowCanIHelp                        string
	AIChatGenerationStopped                  string
	AIChatCompleted                          string
	AIChatWaitingConfirm                     string
	AIChatConfirmPrompt                      string
	AIChatExecutingPlan                      string
	AIChatGeneratingReply                    string
	AIChatCallingTool                        string
	AIChatToolCompleted                      string
	AIChatToolFailed                         string
	AIChatPlanGenerated                      string
	AIChatStatusLabel                        string
	AIChatActionLabel                        string
	AIChatGreeting                           string
	AIChatCapabilities                       string
	AIChatInputPrompt                        string
	AIChatStoppedGeneration                  string
	AIChatCallingToolPrefix                  string
	AIChatToolCompletedPrefix                string
	AIChatToolFailedPrefix                   string

	// AI Two Phase Agent
	AITwoPhaseAgentSystemPromptIntro         string
	AITwoPhaseAgentWorkflowTitle             string
	AITwoPhaseAgentWorkflowStep1             string
	AITwoPhaseAgentWorkflowStep2             string
	AITwoPhaseAgentWorkflowStep2Sub1         string
	AITwoPhaseAgentWorkflowStep2Sub2         string
	AITwoPhaseAgentWorkflowStep2Sub3         string
	AITwoPhaseAgentWorkflowStep3             string
	AITwoPhaseAgentWorkflowStep3Sub1         string
	AITwoPhaseAgentWorkflowStep3Sub2         string
	AITwoPhaseAgentWorkflowStep4             string
	AITwoPhaseAgentWorkflowStep5             string
	AITwoPhaseAgentWorkflowStep6             string
	AITwoPhaseAgentToolNameTitle             string
	AITwoPhaseAgentToolNameIntro             string
	AITwoPhaseAgentToolNameStageFile         string
	AITwoPhaseAgentToolNameDontUseAdd        string
	AITwoPhaseAgentToolNameCommit            string
	AITwoPhaseAgentToolNameDontUseGitCommit  string
	AITwoPhaseAgentToolNameCheckout          string
	AITwoPhaseAgentToolNameDontUseSwitch     string
	AITwoPhaseAgentToolNameCreateBranch      string
	AITwoPhaseAgentToolNameDontUseBranch     string
	AITwoPhaseAgentSpecialToolTitle          string
	AITwoPhaseAgentSpecialToolIntro          string
	AITwoPhaseAgentSpecialToolUsage1         string
	AITwoPhaseAgentSpecialToolUsage2         string
	AITwoPhaseAgentSpecialToolUsage3         string
	AITwoPhaseAgentSpecialToolExample        string
	AITwoPhaseAgentSpecialToolExampleReturn  string
	AITwoPhaseAgentSpecialToolExamplePlan    string
	AITwoPhaseAgentPlanFormatTitle           string
	AITwoPhaseAgentPlanFormatExample         string
	AITwoPhaseAgentNotesTitle                string
	AITwoPhaseAgentNotesParam                string
	AITwoPhaseAgentNotesCriticalTrue         string
	AITwoPhaseAgentNotesCriticalFalse        string
	AITwoPhaseAgentNotesMinimal              string
	AITwoPhaseAgentExecuting                 string
	AITwoPhaseAgentRepoStatusTitle           string
	AITwoPhaseAgentUserInstructionTitle      string
	AITwoPhaseAgentPlanAdjustment            string
	AITwoPhaseAgentExecutionCancelled        string
	AITwoPhaseAgentPlanValidationFailed      string
	AITwoPhaseAgentPlanErrorsIntro           string
	AITwoPhaseAgentPlanRegeneratePrompt      string
	AITwoPhaseAgentContinueAnalysis          string
	AITwoPhaseAgentToolCallWarning           string
	AITwoPhaseAgentSystemPrefix              string
	AITwoPhaseAgentUserFeedbackPrompt        string
	AITwoPhaseAgentToolResultPrefix          string
	AITwoPhaseAgentMaxStepsExceeded          string
	AITwoPhaseAgentEmptyResponseError        string
}

type Bisect struct {
	MarkStart                   string
	ResetTitle                  string
	ResetPrompt                 string
	ResetOption                 string
	ChooseTerms                 string
	OldTermPrompt               string
	NewTermPrompt               string
	BisectMenuTitle             string
	Mark                        string
	SkipCurrent                 string
	SkipSelected                string
	CompleteTitle               string
	CompletePrompt              string
	CompletePromptIndeterminate string
	Bisecting                   string
}

type Log struct {
	EditRebase               string
	HandleUndo               string
	RemoveFile               string
	CopyToClipboard          string
	Remove                   string
	CreateFileWithContent    string
	AppendingLineToFile      string
	EditRebaseFromBaseCommit string
}

type Actions struct {
	CheckoutCommit                   string
	CheckoutBranchAtCommit           string
	CheckoutCommitAsDetachedHead     string
	CheckoutTag                      string
	CheckoutBranch                   string
	CheckoutBranchOrCommit           string
	ForceCheckoutBranch              string
	DeleteLocalBranch                string
	Merge                            string
	SquashMerge                      string
	RebaseBranch                     string
	RenameBranch                     string
	CreateBranch                     string
	FastForwardBranch                string
	AutoForwardBranches              string
	CherryPick                       string
	CheckoutFile                     string
	SquashCommitDown                 string
	FixupCommit                      string
	FixupCommitKeepMessage           string
	RewordCommit                     string
	DropCommit                       string
	EditCommit                       string
	AmendCommit                      string
	ResetCommitAuthor                string
	SetCommitAuthor                  string
	AddCommitCoAuthor                string
	RevertCommit                     string
	CreateFixupCommit                string
	SquashAllAboveFixupCommits       string
	MoveCommitUp                     string
	MoveCommitDown                   string
	CopyCommitMessageToClipboard     string
	CopyCommitMessageBodyToClipboard string
	CopyCommitSubjectToClipboard     string
	CopyCommitDiffToClipboard        string
	CopyCommitHashToClipboard        string
	CopyCommitURLToClipboard         string
	CopyCommitAuthorToClipboard      string
	CopyCommitAttributeToClipboard   string
	CopyCommitTagsToClipboard        string
	CopyPatchToClipboard             string
	CustomCommand                    string
	DiscardAllChangesInFile          string
	DiscardAllUnstagedChangesInFile  string
	StageFile                        string
	StageResolvedFiles               string
	UnstageFile                      string
	UnstageAllFiles                  string
	StageAllFiles                    string
	ResolveConflictByKeepingFile     string
	ResolveConflictByDeletingFile    string
	NotEnoughContextToStage          string
	NotEnoughContextToDiscard        string
	NotEnoughContextForCustomPatch   string
	IgnoreExcludeFile                string
	IgnoreFileErr                    string
	ExcludeFile                      string
	ExcludeGitIgnoreErr              string
	Commit                           string
	Push                             string
	Pull                             string
	OpenFile                         string
	StashAllChanges                  string
	StashAllChangesKeepIndex         string
	StashStagedChanges               string
	StashUnstagedChanges             string
	StashIncludeUntrackedChanges     string
	GitFlowFinish                    string
	GitFlowStart                     string
	CopyToClipboard                  string
	CopySelectedTextToClipboard      string
	RemovePatchFromCommit            string
	MovePatchToSelectedCommit        string
	MovePatchIntoIndex               string
	MovePatchIntoNewCommit           string
	DeleteRemoteBranch               string
	SetBranchUpstream                string
	AddRemote                        string
	AddForkRemote                    string
	RemoveRemote                     string
	UpdateRemote                     string
	ApplyPatch                       string
	Stash                            string
	PopStash                         string
	ApplyStash                       string
	DropStash                        string
	RenameStash                      string
	RemoveSubmodule                  string
	ResetSubmodule                   string
	AddSubmodule                     string
	UpdateSubmoduleUrl               string
	InitialiseSubmodule              string
	BulkInitialiseSubmodules         string
	BulkUpdateSubmodules             string
	BulkDeinitialiseSubmodules       string
	BulkUpdateRecursiveSubmodules    string
	UpdateSubmodule                  string
	CreateLightweightTag             string
	CreateAnnotatedTag               string
	DeleteLocalTag                   string
	DeleteRemoteTag                  string
	PushTag                          string
	NukeWorkingTree                  string
	DiscardUnstagedFileChanges       string
	RemoveUntrackedFiles             string
	RemoveStagedFiles                string
	SoftReset                        string
	MixedReset                       string
	HardReset                        string
	Undo                             string
	Redo                             string
	CopyPullRequestURL               string
	OpenMergeTool                    string
	OpenCommitInBrowser              string
	OpenPullRequest                  string
	StartBisect                      string
	ResetBisect                      string
	BisectSkip                       string
	BisectMark                       string
	AddWorktree                      string
}

const englishIntroPopupMessage = `
Thanks for using lazygit! Seriously you rock. Three things to share with you:

 1) If you want to learn about lazygit's features, watch this vid:
      https://youtu.be/CPLdltN7wgE

 2) Be sure to read the latest release notes at:
      https://github.com/dswcpp/lazygit/releases

 3) If you're using git, that makes you a programmer! With your help we can make
    lazygit better, so consider becoming a contributor and joining the fun at
      https://github.com/dswcpp/lazygit
    Or even just star the repo to share the love!

 4) If lazygit has made your life easier, you can say thanks by clicking the
    donate button at the bottom right. Donation does not grant priority support,
    but it is much appreciated.

Press {{confirmationKey}} to get started.
`

const englishNonReloadableConfigWarning = `The following config settings were changed, but the change doesn't take effect immediately. Please quit and restart lazygit for changes to take effect:

{{configs}}`

const englishHunkStagingHint = `Hunk selection mode is now the default for staging. If you want to stage individual lines, press '%s' to switch to line-by-line mode.

If you prefer to use line-by-line mode by default (like in earlier lazygit versions), add

gui:
  useHunkModeInStagingView: false

to your lazygit config.`

// exporting this so we can use it in tests
func EnglishTranslationSet() *TranslationSet {
	return &TranslationSet{
		NotEnoughSpace:                       "Not enough space to render panels",
		DiffTitle:                            "Diff",
		FilesTitle:                           "Files",
		BranchesTitle:                        "Branches",
		CommitsTitle:                         "Commits",
		StashTitle:                           "Stash",
		SnakeTitle:                           "Snake",
		EasterEgg:                            "Easter egg",
		UnstagedChanges:                      "Unstaged changes",
		StagedChanges:                        "Staged changes",
		StagingTitle:                         "Main panel (staging)",
		MergingTitle:                         "Main panel (merging)",
		NormalTitle:                          "Main panel (normal)",
		LogTitle:                             "Log",
		LogXOfYTitle:                         "Log (%d of %d)",
		CommitSummary:                        "Commit summary",
		CredentialsUsername:                  "Username",
		CredentialsPassword:                  "Password",
		CredentialsPassphrase:                "Enter passphrase for SSH key",
		CredentialsPIN:                       "Enter PIN for SSH key",
		CredentialsToken:                     "Enter Token for SSH key",
		PassUnameWrong:                       "Password, passphrase and/or username wrong",
		Commit:                               "Commit",
		CommitTooltip:                        "Commit staged changes.",
		AmendLastCommit:                      "Amend last commit",
		AmendLastCommitTitle:                 "Amend last commit",
		SureToAmend:                          "Are you sure you want to amend last commit? Afterwards, you can change the commit message from the commits panel.",
		NoCommitToAmend:                      "There's no commit to amend.",
		CommitChangesWithEditor:              "Commit changes using git editor",
		FindBaseCommitForFixup:               "Find base commit for fixup",
		FindBaseCommitForFixupTooltip:        "Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/dswcpp/lazygit/tree/master/docs/Fixup_Commits.md>",
		NoBaseCommitsFound:                   "No base commits found",
		MultipleBaseCommitsFoundStaged:       "Multiple base commits found. (Try staging fewer changes at once)",
		MultipleBaseCommitsFoundUnstaged:     "Multiple base commits found. (Try staging some of the changes)",
		BaseCommitIsAlreadyOnMainBranch:      "The base commit for this change is already on the main branch",
		BaseCommitIsNotInCurrentView:         "Base commit is not in current view",
		HunksWithOnlyAddedLinesWarning:       "There are ranges of only added lines in the diff; be careful to check that these belong in the found base commit.\n\nProceed?",
		StatusTitle:                          "Status",
		Execute:                              "Execute",
		Stage:                                "Stage",
		StageTooltip:                         "Toggle staged for selected file.",
		ToggleStagedAll:                      "Stage all",
		ToggleStagedAllTooltip:               "Toggle staged/unstaged for all files in working tree.",
		ToggleTreeView:                       "Toggle file tree view",
		ToggleTreeViewTooltip:                "Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.\n\nThe default can be changed in the config file with the key 'gui.showFileTree'.",
		OpenDiffTool:                         "Open external diff tool (git difftool)",
		OpenMergeTool:                        "Open external merge tool",
		Refresh:                              "Refresh",
		RefreshTooltip:                       "Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`.",
		Push:                                 "Push",
		PushTooltip:                          "Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch.",
		Pull:                                 "Pull",
		PullTooltip:                          "Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch.",
		MergeConflictsTitle:                  "Merge conflicts",
		MergeConflictDescription_DD:          "Conflict: this file was moved or renamed both in the current and the incoming changes, but to different destinations. I don't know which ones, but they should both show up as conflicts too (marked 'AU' and 'UA', respectively). The most likely resolution is to delete this file, and pick one of the destinations and delete the other.",
		MergeConflictDescription_AU:          "Conflict: this file is the destination of a move or rename in the current changes, but was moved or renamed to a different destination in the incoming changes. That other destination should also show up as a conflict (marked 'UA'), as well as the file that both were renamed from (marked 'DD').",
		MergeConflictDescription_UA:          "Conflict: this file is the destination of a move or rename in the incoming changes, but was moved or renamed to a different destination in the current changes. That other destination should also show up as a conflict (marked 'AU'), as well as the file that both were renamed from (marked 'DD').",
		MergeConflictDescription_DU:          "Conflict: this file was deleted in the current changes and modified in the incoming changes.\n\nThe most likely resolution is to delete the file after applying the incoming modifications manually to some other place in the code.",
		MergeConflictDescription_UD:          "Conflict: this file was modified in the current changes and deleted in incoming changes.\n\nThe most likely resolution is to delete the file after applying the current modifications manually to some other place in the code.",
		MergeConflictIncomingDiff:            "Incoming changes:",
		MergeConflictCurrentDiff:             "Current changes:",
		MergeConflictPressEnterToResolve:     "Press %s to resolve.",
		MergeConflictKeepFile:                "Keep file",
		MergeConflictDeleteFile:              "Delete file",
		Checkout:                             "Checkout",
		CheckoutTooltip:                      "Checkout selected item.",
		CantCheckoutBranchWhilePulling:       "You cannot checkout another branch while pulling the current branch",
		TagCheckoutTooltip:                   "Checkout the selected tag as a detached HEAD.",
		RemoteBranchCheckoutTooltip:          "Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head.",
		CantPullOrPushSameBranchTwice:        "You cannot push or pull a branch while it is already being pushed or pulled",
		FileFilter:                           "Filter files by status",
		CopyToClipboardMenu:                  "Copy to clipboard",
		CopyFileName:                         "File name",
		CopyRelativeFilePath:                 "Relative path",
		CopyAbsoluteFilePath:                 "Absolute path",
		CopyFileDiffTooltip:                  "If there are staged items, this command considers only them. Otherwise, it considers all the unstaged ones.",
		CopySelectedDiff:                     "Diff of selected file",
		CopyAllFilesDiff:                     "Diff of all files",
		CopyFileContent:                      "Content of selected file",
		NoContentToCopyError:                 "Nothing to copy",
		FileNameCopiedToast:                  "File name copied to clipboard",
		FilePathCopiedToast:                  "File path copied to clipboard",
		FileDiffCopiedToast:                  "File diff copied to clipboard",
		AllFilesDiffCopiedToast:              "All files diff copied to clipboard",
		FileContentCopiedToast:               "File content copied to clipboard",
		FilterStagedFiles:                    "Show only staged files",
		FilterUnstagedFiles:                  "Show only unstaged files",
		FilterTrackedFiles:                   "Show only tracked files",
		FilterUntrackedFiles:                 "Show only untracked files",
		NoFilter:                             "No filter",
		FilterLabelStagedFiles:               "(only staged)",
		FilterLabelUnstagedFiles:             "(only unstaged)",
		FilterLabelTrackedFiles:              "(only tracked)",
		FilterLabelUntrackedFiles:            "(only untracked)",
		FilterLabelConflictingFiles:          "(only conflicting)",
		NoChangedFiles:                       "No changed files",
		SoftReset:                            "Soft reset",
		AlreadyCheckedOutBranch:              "You have already checked out this branch",
		SureForceCheckout:                    "Are you sure you want force checkout? You will lose all local changes",
		ForceCheckoutBranch:                  "Force checkout branch",
		BranchName:                           "Branch name",
		NewBranchNameBranchOff:               "New branch name (branch is off of '{{.branchName}}')",
		CantDeleteCheckOutBranch:             "You cannot delete the checked out branch!",
		DeleteBranchTitle:                    "Delete branch '{{.selectedBranchName}}'?",
		DeleteBranchesTitle:                  "Delete selected branches?",
		DeleteLocalBranch:                    "Delete local branch",
		DeleteLocalBranches:                  "Delete local branches",
		DeleteRemoteBranchPrompt:             "Are you sure you want to delete the remote branch '{{.selectedBranchName}}' from '{{.upstream}}'?",
		DeleteRemoteBranchesPrompt:           "Are you sure you want to delete the remote branches of the selected branches from their respective remotes?",
		DeleteLocalAndRemoteBranchPrompt:     "Are you sure you want to delete both '{{.localBranchName}}' from your machine, and '{{.remoteBranchName}}' from '{{.remoteName}}'?",
		DeleteLocalAndRemoteBranchesPrompt:   "Are you sure you want to delete both the selected branches from your machine, and their remote branches from their respective remotes?",
		ForceDeleteBranchTitle:               "Force delete branch",
		ForceDeleteBranchMessage:             "'{{.selectedBranchName}}' is not fully merged. Are you sure you want to delete it?",
		ForceDeleteBranchesMessage:           "Some of the selected branches are not fully merged. Are you sure you want to delete them?",
		RebaseBranch:                         "Rebase",
		RebaseBranchTooltip:                  "Rebase the checked-out branch onto the selected branch.",
		CantRebaseOntoSelf:                   "You cannot rebase a branch onto itself",
		CantMergeBranchIntoItself:            "You cannot merge a branch into itself",
		ForceCheckout:                        "Force checkout",
		ForceCheckoutTooltip:                 "Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch.",
		CheckoutByName:                       "Checkout by name",
		CheckoutByNameTooltip:                "Checkout by name. In the input box you can enter '-' to switch to the previous branch.",
		CheckoutPreviousBranch:               "Checkout previous branch",
		RemoteBranchCheckoutTitle:            "Checkout {{.branchName}}",
		RemoteBranchCheckoutPrompt:           "How would you like to check out this branch?",
		CheckoutTypeNewBranch:                "New local branch",
		CheckoutTypeNewBranchTooltip:         "Checkout the remote branch as a local branch, tracking the remote branch.",
		CheckoutTypeDetachedHead:             "Detached head",
		CheckoutTypeDetachedHeadTooltip:      "Checkout the remote branch as a detached head, which can be useful if you just want to test the branch but not work on it yourself. You can still create a local branch from it later.",
		NewBranch:                            "New branch",
		NewBranchFromStashTooltip:            "Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit.",
		MoveCommitsToNewBranch:               "Move commits to new branch",
		MoveCommitsToNewBranchTooltip:        "Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.\n\nNote that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which).",
		MoveCommitsToNewBranchFromMainPrompt: "This will take all unpushed commits and move them to a new branch (off of {{.baseBranchName}}). It will then hard-reset the current branch to its upstream branch. Do you want to continue?",
		MoveCommitsToNewBranchMenuPrompt:     "This will take all unpushed commits and move them to a new branch. This new branch can either be created from the main branch ({{.baseBranchName}}) or stacked on top of the current branch. Which of these would you like to do?",
		MoveCommitsToNewBranchFromBaseItem:   "New branch from base branch (%s)",
		MoveCommitsToNewBranchStackedItem:    "New branch stacked on current branch (%s)",
		CannotMoveCommitsFromDetachedHead:    "Cannot move commits from a detached head",
		CannotMoveCommitsNoUpstream:          "Cannot move commits from a branch that has no upstream branch",
		CannotMoveCommitsBehindUpstream:      "Cannot move commits from a branch that is behind its upstream branch",
		CannotMoveCommitsNoUnpushedCommits:   "There are no unpushed commits to move to a new branch",
		NoBranchesThisRepo:                   "No branches for this repo",
		CommitWithoutMessageErr:              "You cannot commit without a commit message",
		Close:                                "Close",
		CloseCancel:                          "Close/Cancel",
		Confirm:                              "Confirm",
		Quit:                                 "Quit",
		SquashTooltip:                        "Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it.",
		NoCommitsThisBranch:                  "No commits for this branch",
		UpdateRefHere:                        "Update branch '{{.ref}}' here",
		ExecCommandHere:                      "Execute the following command here:",
		CannotSquashOrFixupFirstCommit:       "There's no commit below to squash into",
		CannotSquashOrFixupMergeCommit:       "Cannot squash or fixup a merge commit",
		Fixup:                                "Fixup",
		FixupKeepMessage:                     "Fixup and use this commit's message",
		FixupKeepMessageTooltip:              "Squash the selected commit into the commit below, using this commit's message, discarding the message of the commit below.",
		SetFixupMessage:                      "Set fixup message",
		SetFixupMessageTooltip:               "Set the message option for the fixup commit. The -C option means to use this commit's message instead of the target commit's message.",
		FixupDiscardMessage:                  "Fixup and discard this commit's message",
		FixupDiscardMessageTooltip:           "Squash the selected commit into the commit below, discarding this commit's message.",
		SureSquashThisCommit:                 "Are you sure you want to squash the selected commit(s) into the commit below?",
		Squash:                               "Squash",
		PickCommitTooltip:                    "Mark the selected commit to be picked (when mid-rebase). This means that the commit will be retained upon continuing the rebase.",
		Pick:                                 "Pick",
		Edit:                                 "Edit",
		Revert:                               "Revert",
		RevertCommitTooltip:                  "Create a revert commit for the selected commit, which applies the selected commit's changes in reverse.",
		Reword:                               "Reword",
		CommitRewordTooltip:                  "Reword the selected commit's message.",
		DropCommit:                           "Drop",
		DropCommitTooltip:                    "Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts.",
		MoveDownCommit:                       "Move commit down one",
		MoveUpCommit:                         "Move commit up one",
		CannotMoveAnyFurther:                 "Cannot move any further",
		CannotMoveMergeCommit:                "Cannot move a merge commit",
		EditCommit:                           "Edit (start interactive rebase)",
		EditCommitTooltip:                    "Edit the selected commit. Use this to start an interactive rebase from the selected commit. When already mid-rebase, this will mark the selected commit for editing, which means that upon continuing the rebase, the rebase will pause at the selected commit to allow you to make changes.",
		AmendCommitTooltip:                   "Amend commit with staged changes. If the selected commit is the HEAD commit, this will perform `git commit --amend`. Otherwise the commit will be amended via a rebase.",
		Amend:                                "Amend",
		ResetAuthor:                          "Reset author",
		ResetAuthorTooltip:                   "Reset the commit's author to the currently configured user. This will also renew the author timestamp",
		SetAuthor:                            "Set author",
		SetAuthorTooltip:                     "Set the author based on a prompt",
		AddCoAuthor:                          "Add co-author",
		AmendCommitAttribute:                 "Amend commit attribute",
		AmendCommitAttributeTooltip:          "Set/Reset commit author or set co-author.",
		SetAuthorPromptTitle:                 "Set author (must look like 'Name <Email>')",
		AddCoAuthorPromptTitle:               "Add co-author (must look like 'Name <Email>')",
		AddCoAuthorTooltip:                   "Add co-author using the Github/Gitlab metadata Co-authored-by.",
		RewordCommitEditor:                   "Reword with editor",
		Error:                                "Error",
		PickHunk:                             "Pick hunk",
		PickAllHunks:                         "Pick all hunks",
		Undo:                                 "Undo",
		UndoReflog:                           "Undo",
		RedoReflog:                           "Redo",
		UndoTooltip:                          "The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration.",
		RedoTooltip:                          "The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration.",
		UndoMergeResolveTooltip:              "Undo last merge conflict resolution.",
		DiscardAllTooltip:                    "Discard both staged and unstaged changes in '{{.path}}'.",
		DiscardUnstagedTooltip:               "Discard unstaged changes in '{{.path}}'.",
		DiscardUnstagedDisabled:              "The selected items don't have both staged and unstaged changes.",
		Pop:                                  "Pop",
		StashPopTooltip:                      "Apply the stash entry to your working directory and remove the stash entry.",
		Drop:                                 "Drop",
		StashDropTooltip:                     "Remove the stash entry from the stash list.",
		Apply:                                "Apply",
		StashApplyTooltip:                    "Apply the stash entry to your working directory.",
		NoStashEntries:                       "No stash entries",
		StashDrop:                            "Stash drop",
		SureDropStashEntry:                   "Are you sure you want to drop the selected stash entry(ies)?",
		StashPop:                             "Stash pop",
		SurePopStashEntry:                    "Are you sure you want to pop this stash entry?",
		StashApply:                           "Stash apply",
		SureApplyStashEntry:                  "Are you sure you want to apply this stash entry?",
		NoTrackedStagedFilesStash:            "You have no tracked/staged files to stash",
		NoFilesToStash:                       "You have no files to stash",
		StashChanges:                         "Stash changes",
		RenameStash:                          "Rename stash",
		RenameStashPrompt:                    "Rename stash: {{.stashName}}",
		OpenConfig:                           "Open config file",
		EditConfig:                           "Edit config file",
		ForcePush:                            "Force push",
		ForcePushPrompt:                      "Your branch has diverged from the remote branch. Press {{.cancelKey}} to cancel, or {{.confirmKey}} to force push.",
		ForcePushDisabled:                    "Your branch has diverged from the remote branch and you've disabled force pushing",
		UpdatesRejected:                      "Updates were rejected. Please fetch and examine the remote changes before pushing again.",
		UpdatesRejectedAndForcePushDisabled:  "Updates were rejected and you have disabled force pushing",
		CheckForUpdate:                       "Check for update",
		CheckingForUpdates:                   "Checking for updates...",
		UpdateAvailableTitle:                 "Update available!",
		UpdateAvailable:                      "Download and install version {{.newVersion}}?",
		UpdateInProgressWaitingStatus:        "Updating",
		UpdateCompletedTitle:                 "Update completed!",
		UpdateCompleted:                      "Update has been installed successfully. Restart lazygit for it to take effect.",
		FailedToRetrieveLatestVersionErr:     "Failed to retrieve version information",
		OnLatestVersionErr:                   "You already have the latest version",
		MajorVersionErr:                      "New version ({{.newVersion}}) has non-backwards compatible changes compared to the current version ({{.currentVersion}})",
		CouldNotFindBinaryErr:                "Could not find any binary at {{.url}}",
		UpdateFailedErr:                      "Update failed: {{.errMessage}}",
		ConfirmQuitDuringUpdateTitle:         "Currently updating",
		ConfirmQuitDuringUpdate:              "An update is in progress. Are you sure you want to quit?",
		IntroPopupMessage:                    englishIntroPopupMessage,
		NonReloadableConfigWarningTitle:      "Config changed",
		NonReloadableConfigWarning:           englishNonReloadableConfigWarning,
		GitconfigParseErr:                    `Gogit failed to parse your gitconfig file due to the presence of unquoted '\' characters. Removing these should fix the issue.`,
		EditFile:                             `Edit file`,
		EditFileTooltip:                      "Open file in external editor.",
		OpenFile:                             `Open file`,
		OpenFileTooltip:                      "Open file in default application.",
		OpenInEditor:                         "Open in editor",
		IgnoreFile:                           `Add to .gitignore`,
		ExcludeFile:                          `Add to .git/info/exclude`,
		RefreshFiles:                         `Refresh files`,
		FocusMainView:                        "Focus main view",
		Merge:                                `Merge`,
		MergeBranchTooltip:                   "View options for merging the selected item into the current branch (regular merge, squash merge)",
		RegularMergeFastForward:              "Regular merge (fast-forward)",
		RegularMergeFastForwardTooltip:       "Fast-forward '{{.checkedOutBranch}}' to '{{.selectedBranch}}' without creating a merge commit.",
		CannotFastForwardMerge:               "Cannot fast-forward '{{.checkedOutBranch}}' to '{{.selectedBranch}}'",
		RegularMergeNonFastForward:           "Regular merge (with merge commit)",
		RegularMergeNonFastForwardTooltip:    "Merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}', creating a merge commit.",
		SquashMergeUncommitted:               "Squash merge and leave uncommitted",
		SquashMergeUncommittedTooltip:        "Squash merge '{{.selectedBranch}}' into the working tree.",
		SquashMergeCommitted:                 "Squash merge and commit",
		SquashMergeCommittedTooltip:          "Squash merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}' as a single commit.",
		ConfirmQuit:                          `Are you sure you want to quit?`,
		SwitchRepo:                           `Switch to a recent repo`,
		AllBranchesLogGraph:                  `Show/cycle all branch logs`,
		UnsupportedGitService:                `Unsupported git service`,
		CreatePullRequest:                    `Create pull request`,
		CopyPullRequestURL:                   `Copy pull request URL to clipboard`,
		NoBranchOnRemote:                     `This branch doesn't exist on remote. You need to push it to remote first.`,
		Fetch:                                `Fetch`,
		FetchTooltip:                         "Fetch changes from remote.",
		CollapseAll:                          "Collapse all files",
		CollapseAllTooltip:                   "Collapse all directories in the files tree",
		ExpandAll:                            "Expand all files",
		ExpandAllTooltip:                     "Expand all directories in the file tree",
		DisabledInFlatView:                   "Not available in flat view",
		FileEnter:                            `Stage lines / Collapse directory`,
		FileEnterTooltip:                     "If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it.",
		StageSelectionTooltip:                `Toggle selection staged / unstaged.`,
		DiscardSelection:                     `Discard`,
		DiscardSelectionTooltip:              "When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change.",
		ToggleRangeSelect:                    "Toggle range select",
		DismissRangeSelect:                   "Dismiss range select",
		ToggleSelectHunk:                     "Toggle hunk selection",
		SelectHunk:                           "Select hunks",
		SelectLineByLine:                     "Select line-by-line",
		ToggleSelectHunkTooltip:              "Toggle line-by-line vs. hunk selection mode.",
		HunkStagingHint:                      englishHunkStagingHint,
		ToggleSelectionForPatch:              `Toggle lines in patch`,
		EditHunk:                             `Edit hunk`,
		EditHunkTooltip:                      "Edit selected hunk in external editor.",
		ToggleStagingView:                    "Switch view",
		ToggleStagingViewTooltip:             "Switch to other view (staged/unstaged changes).",
		ReturnToFilesPanel:                   `Return to files panel`,
		FastForward:                          `Fast-forward`,
		FastForwardTooltip:                   "Fast-forward selected branch from its upstream.",
		FastForwarding:                       "Fast-forwarding",
		FoundConflictsTitle:                  "Conflicts!",
		ViewConflictsMenuItem:                "View conflicts",
		AbortMenuItem:                        "Abort the %s",
		ViewMergeRebaseOptions:               "View merge/rebase options",
		ViewMergeRebaseOptionsTooltip:        "View options to abort/continue/skip the current merge/rebase.",
		ViewMergeOptions:                     "View merge options",
		ViewRebaseOptions:                    "View rebase options",
		ViewCherryPickOptions:                "View cherry-pick options",
		ViewRevertOptions:                    "View revert options",
		NotMergingOrRebasing:                 "You are currently neither rebasing nor merging",
		AlreadyRebasing:                      "Can't perform this action during a rebase",
		NotMidRebase:                         "This action only works during an interactive rebase",
		MustSelectFixupCommit:                "This action only works on fixup commits",
		RecentRepos:                          "Recent repositories",
		MergeOptionsTitle:                    "Merge options",
		RebaseOptionsTitle:                   "Rebase options",
		CherryPickOptionsTitle:               "Cherry-pick options",
		RevertOptionsTitle:                   "Revert options",
		CommitSummaryTitle:                   "Commit summary",
		CommitDescriptionTitle:               "Commit description",
		CommitDescriptionSubTitle:            "Press {{.togglePanelKeyBinding}} to toggle focus, {{.commitMenuKeybinding}} to open menu",
		CommitDescriptionFooter:              "Press {{.confirmInEditorKeybinding}} to submit",
		CommitDescriptionFooterTwoBindings:   "Press {{.confirmInEditorKeybinding1}} or {{.confirmInEditorKeybinding2}} to submit",
		CommitHooksDisabledSubTitle:          "(hooks disabled)",
		LocalBranchesTitle:                   "Local branches",
		SearchTitle:                          "Search",
		TagsTitle:                            "Tags",
		MenuTitle:                            "Menu",
		CommitMenuTitle:                      "Commit Menu",
		RemotesTitle:                         "Remotes",
		RemoteBranchesTitle:                  "Remote branches",
		PatchBuildingTitle:                   "Main panel (patch building)",
		InformationTitle:                     "Information",
		SecondaryTitle:                       "Secondary",
		ReflogCommitsTitle:                   "Reflog",
		GlobalTitle:                          "Global keybindings",
		ConflictsResolved:                    "All merge conflicts resolved. Continue the %s?",
		Continue:                             "Continue",
		UnstagedFilesAfterConflictsResolved:  "Files have been modified since conflicts were resolved. Auto-stage them and continue?",
		Keybindings:                          "Keybindings",
		KeybindingsMenuSectionLocal:          "Local",
		KeybindingsMenuSectionGlobal:         "Global",
		KeybindingsMenuSectionNavigation:     "Navigation",
		RebasingTitle:                        "Rebase '{{.checkedOutBranch}}'",
		RebasingFromBaseCommitTitle:          "Rebase '{{.checkedOutBranch}}' from marked base",
		SimpleRebase:                         "Simple rebase onto '{{.ref}}'",
		InteractiveRebase:                    "Interactive rebase onto '{{.ref}}'",
		RebaseOntoBaseBranch:                 "Rebase onto base branch ({{.baseBranch}})",
		InteractiveRebaseTooltip:             "Begin an interactive rebase with a break at the start, so you can update the TODO commits before continuing.",
		RebaseOntoBaseBranchTooltip:          "Rebase the checked out branch onto its base branch (i.e. the closest main branch).",
		MustSelectTodoCommits:                "When rebasing, this action only works on a selection of TODO commits.",
		FwdNoUpstream:                        "Cannot fast-forward a branch with no upstream",
		FwdNoLocalUpstream:                   "Cannot fast-forward a branch whose remote is not registered locally",
		FwdCommitsToPush:                     "Cannot fast-forward a branch with commits to push",
		PullRequestNoUpstream:                "Cannot open a pull request for a branch with no upstream",
		ErrorOccurred:                        "An error occurred! Please create an issue at",
		ConflictLabel:                        "CONFLICT",
		PendingRebaseTodosSectionHeader:      "Pending rebase todos",
		PendingCherryPicksSectionHeader:      "Pending cherry-picks",
		PendingRevertsSectionHeader:          "Pending reverts",
		CommitsSectionHeader:                 "Commits",
		YouDied:                              "YOU DIED!",
		RewordNotSupported:                   "Rewording commits while interactively rebasing is not currently supported",
		ChangingThisActionIsNotAllowed:       "Changing this kind of rebase todo entry is not allowed",
		NotAllowedMidCherryPickOrRevert:      "This action is not allowed while cherry-picking or reverting",
		PickIsOnlyAllowedDuringRebase:        "This action is only allowed while rebasing",
		DroppingMergeRequiresSingleSelection: "Dropping a merge commit requires a single selected item",
		CherryPickCopy:                       "Copy (cherry-pick)",
		CherryPickCopyTooltip:                "Mark commit as copied. Then, within the local commits view, you can press `{{.paste}}` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `{{.escape}}` to cancel the selection.",
		PasteCommits:                         "Paste (cherry-pick)",
		SureCherryPick:                       "Are you sure you want to cherry-pick the {{.numCommits}} copied commit(s) onto this branch?",
		CherryPick:                           "Cherry-pick",
		CannotCherryPickNonCommit:            "Cannot cherry-pick this kind of todo item",
		PrevHunk:                             "Go to previous hunk",
		NextHunk:                             "Go to next hunk",
		PrevConflict:                         "Previous conflict",
		NextConflict:                         "Next conflict",
		SelectPrevHunk:                       "Previous hunk",
		SelectNextHunk:                       "Next hunk",
		ScrollDown:                           "Scroll down",
		ScrollUp:                             "Scroll up",
		ScrollUpMainWindow:                   "Scroll up main window",
		ScrollDownMainWindow:                 "Scroll down main window",
		SuspendApp:                           "Suspend the application",
		CannotSuspendApp:                     "Suspending the application is not supported on Windows",
		AmendCommitTitle:                     "Amend commit",
		AmendCommitPrompt:                    "Are you sure you want to amend this commit with your staged files?",
		AmendCommitWithConflictsMenuPrompt:   "WARNING: you are about to amend the last finished commit with your resolved conflicts. This is very unlikely to be what you want at this point. More likely, you simply want to continue the rebase instead.\n\nDo you still want to amend the previous commit?",
		AmendCommitWithConflictsContinue:     "No, continue rebase",
		AmendCommitWithConflictsAmend:        "Yes, amend previous commit",
		DropCommitTitle:                      "Drop commit",
		DropCommitPrompt:                     "Are you sure you want to drop the selected commit(s)?",
		DropMergeCommitPrompt:                "Are you sure you want to drop the selected merge commit? Note that it will also drop all the commits that were merged in by it.",
		DropUpdateRefPrompt:                  "Are you sure you want to delete the selected update-ref todo(s)? This is irreversible except by aborting the rebase.",
		PullingStatus:                        "Pulling",
		PushingStatus:                        "Pushing",
		FetchingStatus:                       "Fetching",
		SquashingStatus:                      "Squashing",
		FixingStatus:                         "Fixing up",
		DeletingStatus:                       "Deleting",
		AIGeneratingStatus:                   "AI generating commit message...",
		AIGenerateCommitMessage:              "Generate commit message with AI",
		AINotEnabled:                         "AI is not enabled. Set ai.enabled: true in your config",
		AINoStagedChanges:                    "No staged changes to generate a commit message from",
		AIError:                              "AI error: %s",
		AISettings:                           "AI Settings",
		AISettingsEnable:                     "Enable AI",
		AISettingsDisable:                    "Disable AI",
		AISettingsSetAPIKey:                  "Set API Key",
		AISettingsAPIKeyPrompt:               "Enter API key (or env var ref like ${DEEPSEEK_API_KEY})",
		AISettingsSetProvider:                "Set Provider",
		AISettingsSetModel:                   "Set Model",
		AISettingsModelPrompt:                "Enter model name (e.g. deepseek-reasoner, gpt-4o-mini)",
		AISettingsSetEndpoint:                "Set Endpoint",
		AISettingsEndpointPrompt:             "Enter API endpoint URL (e.g. http://localhost:11434/v1)",
		AISettingsSaved:                      "AI settings saved",
		AISettingsActiveProfile:              "Active profile",
		AISettingsSwitchProfile:              "Switch active profile",
		AISettingsEditProfile:                "Edit active profile",
		AISettingsAddProfile:                 "Add new profile",
		AISettingsNoProfiles:                 "No AI profiles configured",
		AISettingsProfileName:                "Name",
		AISettingsProfileNamePrompt:          "Enter profile name",
		AISettingsNewProfileNamePrompt:       "Enter name for new profile",
		AISettingsMaxTokens:                  "Max tokens",
		AISettingsMaxTokensPrompt:            "Enter max tokens (e.g. 8000)",
		AISettingsTimeout:                    "Timeout (s)",
		AISettingsTimeoutPrompt:              "Enter timeout in seconds (e.g. 60)",
		AISettingsDeleteProfile:              "Delete this profile",
		AISettingsDeleteProfileTitle:         "Delete profile",
		AISettingsDeleteProfilePrompt:        "Delete profile '%s'?",
		AISettingsCannotDeleteLastProfile:    "Cannot delete the last profile",
		AIAssistant:                          "Open AI git assistant",
		AIAssistantTitle:                     "AI Git Assistant",
		AIAssistantPrompt:                    "Describe your git task (e.g. 'squash last 3 commits')",
		AIAssistantStatus:                    "AI generating commands...",
		AIAssistantConfirmExecute:            "Execute these commands?",
		AIAssistantNoCommands:                "AI did not generate any commands",
		AIAssistantSilentNoCommands:          "No executable commands found",
		AIAssistantConfirmSilentExecute:      "Silently execute these commands?",
		AIAssistantExecuting:                 "Executing commands...",
		AIAssistantExecuteError:              "Command execution failed",
		AICodeReview:                         "AI code review",
		AICodeReviewTitle:                    "AI Code Review",
		AIGenerationCancelled:                "AI commit message generation cancelled",
		AICommitPromptTruncated:              "\n[Content too large, truncated at 120000 characters]",
		AICommitEmptyResponse:                "AI returned an empty response. Please try again or check AI settings (Ctrl+A)",
		AIPromptCurrentBranch:                "Current branch: %s\n",
		AIPromptRecentCommits:                "Recent commits:\n",
		AIPromptProjectType:                  "Project type: %s\n\n",
		AIPromptRole:                         "You are a professional git commit message generator, proficient in Conventional Commits specification.\n\n",
		AIPromptTask:                         "## Task\nGenerate a standard commit message based on the following code changes.\n\n",
		AIPromptRepoContext:                  "## Repository Context\n",
		AIPromptCodeChanges:                  "## Code Changes\n",
		AIPromptOutputRules:                  "## Output Rules\n1. Format: <type>(<scope>): <subject>\n\n",
		AIPromptTypeGuide:                    "2. Type Selection Guide:\n   - feat: New feature (user-visible new functionality)\n   - fix: Bug fix (fixing user-encountered issues)\n   - refactor: Refactoring (code improvement without changing external behavior)\n   - perf: Performance optimization\n   - docs: Documentation update\n   - test: Test-related\n   - chore: Build/tool/dependency updates\n   - style: Code formatting (no logic changes)\n   - ci: CI/CD configuration\n\n",
		AIPromptScopeGuide:                   "3. Scope Selection:\n   - Use the most relevant module/component name\n   - If multiple modules involved, use the core one\n   - Examples: (auth), (api), (ui), (db)\n\n",
		AIPromptSubjectRequirements:          "4. Subject Requirements:\n   - Use imperative mood (\"add\" not \"added\")\n   - No more than 72 characters\n   - No period at the end\n   - Must use English\n\n",
		AIPromptBodyOptional:                 "5. Body (optional):\n   - Add body if changes are complex or need explanation\n   - Start after a blank line\n   - Explain \"why\" not \"what\"\n\n",
		AIPromptScenarioGuide:                "6. Scenario Guidance:\n",
		AIPromptOutputFormat:                 "7. Output Format:\n   - Output only the commit message itself\n   - Do not add markdown code blocks (like ```)\n   - Do not add any explanatory text\n",
		AIPromptScenarioSmall:                "   - Small change detected, generate a concise single-line commit message\n",
		AIPromptScenarioLarge:                "   - Large-scale change detected, summarize core goal in subject\n   - List main changes in body (3-5 items)\n",
		AIPromptScenarioBugfix:               "   - Bug fix detected, suggest explaining problem and solution in body\n",
		AIPromptScenarioRefactor:             "   - Refactoring detected, suggest explaining refactoring purpose and improvements in body\n",
		AIPromptScenarioDocs:                 "   - Documentation update detected, use docs type\n",
		AIPromptScenarioTest:                 "   - Test-related detected, use test type\n",
		AIPromptScenarioDefault:              "   - Choose appropriate type and scope based on changes\n",
		AICodeReviewConfirmTitle:             "AI Code Review",
		AICodeReviewConfirmPrompt:            "Start AI code review for the following file?\n\n%s",
		AICodeReviewStatus:                   "AI reviewing, please wait...",
		AICodeReviewNoDiff:                   "No diff to review for this file",
		AICodeReviewCopiedToClipboard:        "AI code review copied to clipboard",
		AICodeReviewToggleZoom:               "Toggle zoom",
		NewBranchWithAI:                      "New branch with AI suggestion",
		NewBranchWithAITooltip:               "Use AI to suggest branch name based on working tree changes (format: feature/add-user-auth)",
		AIGeneratingBranchNameStatus:         "AI generating branch name...",
		CreatePRWithAIDescription:            "Generate AI PR description",
		CreatePRWithAIDescriptionTooltip:     "Use AI to generate a professional PR description based on commit history and code changes, copied to clipboard",
		CreatePRDirectly:                     "Open PR directly",
		AIGeneratingPRDescriptionStatus:      "AI generating PR description...",
		PRDescriptionCopiedToClipboard:       "PR description copied to clipboard! You can paste it when the PR page opens.",
		// AI Diff Filter
		AIDiffSkipped:                        "Skipped %s: %s",
		AIDiffBinaryFile:                     "binary file",
		AIDiffLockOrGeneratedFile:            "lock/generated file",
		AIDiffChangeStats:                    "# Change Statistics",
		AIDiffFilesCount:                     "- Files: %d",
		AIDiffFileTypes:                      "- File types: ",
		AIDiffChangeScale:                    "- Change scale: +%d/-%d lines",
		AIDiffMajorChanges:                   "- Major changes: ",
		AIDiffNewFile:                        "[New]",
		AIDiffDeletedFile:                    "[Deleted]",
		AIDiffRenamedFile:                    "[Renamed]",
		AIDiffModifiedFile:                   "[Modified]",
		AIDiffTruncated:                      "[%s diff is large, %d lines total, truncated to first %d lines]",
		AIDiffSmartTruncated:                 "[%s smart truncated: preserved %d/%d lines of key code (function signatures, important comments, etc.)]",
		// AI Helper Messages
		AICodeReviewCancelled:                "AI code review cancelled",
		AINotEnabledPleaseConfig:             "AI is not enabled. Please enable AI and configure a profile first",
		AITestingConnection:                  "Testing AI connection...",
		AIConnectionTestFailed:               "AI connection test failed",
		AIEmptyResponse:                      "AI returned an empty response",
		AIConnectionTestSuccess:              "✓ AI connection test successful!",
		AIConnectionTestSuccessDetail:        "✓ AI connection test successful!\nProfile: %s\nResponse: %s",
		APIKeyCannotBeEmpty:                  "API Key cannot be empty",
		AIConfigComplete:                     "AI Configuration Complete",
		AIConfigCompletePrompt:               "✓ %s Profile created and activated!\n\nTest connection now?",
		AIWelcomeWizardTitle:                 "Welcome to AI Features - First Time Setup Wizard",
		UseDeepSeekRecommended:               "Use DeepSeek (Recommended)",
		UseOpenAI:                             "Use OpenAI",
		UseAnthropicClaude:                   "Use Anthropic Claude",
		UseOllamaLocal:                       "Use Ollama (Local Model)",
		ConfigureLater:                       "Configure Later (Go to AI Settings)",
		TestCurrentProfile:                   "Test Current Profile Connection",
		SetupProviderAPIKey:                  "Setup %s API Key",
		// AI Diff Filter Additional
		AIDiffSkippedFilesNote:               " (plus %d lock/binary/generated files skipped)",
		// AI Assistant Prompts
		AIAssistantSystemPrompt:              "You are a git command generator. Based on user requirements and repository status, generate shell/git commands that need to be executed.\n\n",
		AIAssistantRules:                     "Rules:\n- Output only directly executable commands, one per line\n- Do not output any explanations, comments (starting with #), or markdown\n- Commands should be in execution order\n- If the requirement cannot be safely completed with git commands, output on the first line: CANNOT_EXECUTE: <reason>\n\n",
		AIAssistantRepoState:                 "Current repository status:\n%s\n",
		AIAssistantUserRequest:               "User requirement: %s",
		// AI Branch Naming Prompts
		AIBranchNameSystemPrompt:             "You are a git branch naming assistant.\n\n",
		AIBranchNameTask:                     "Task: Based on the following changes, suggest a standard git branch name.\n\n",
		AIBranchNameRules:                    "Branch naming rules:\n- Format: <type>/<short-description>\n- Type: feature (new feature), bugfix (fix bug), refactor (refactoring), docs (documentation), test (testing), chore (miscellaneous)\n- Description: use kebab-case (lowercase letters and hyphens), concise and clear, no more than 50 characters\n- Examples: feature/user-authentication, bugfix/login-crash, refactor/api-client\n\n",
		AIBranchNameChanges:                  "Changes:\n%s\n",
		AIBranchNameDiffSummary:              "Diff summary:\n```diff\n%s\n```\n\n",
		AIBranchNameRequirements:             "Requirements:\n- Output only one branch name (do not include any explanations or quotes)\n- Must conform to the above format and rules\n- Use English naming\n",
		// AI PR Description Prompts
		AIPRDescSystemPrompt:                 "You are a Pull Request description generation assistant.\n\n",
		AIPRDescTask:                         "Task: Based on the following commit history and code changes, generate a professional PR description.\n\n",
		AIPRDescCommitHistory:                "Commit history:\n",
		AIPRDescCodeChanges:                  "Code changes:\n```diff\n%s\n```\n\n",
		AIPRDescFormatRequirements:           "PR description format requirements:\n",
		AIPRDescSummarySection:               "## Summary\nSummarize the main purpose and value of this PR in one sentence.\n\n",
		AIPRDescChangesSection:               "## Changes\n- List main feature changes, bug fixes, or refactoring (3-5 points)\n\n",
		AIPRDescTechDetailsSection:           "## Technical Details (optional)\n- If there are important technical implementation details or architectural changes, briefly explain\n\n",
		AIPRDescTestingSection:               "## Testing\n- [ ] Unit tests passed\n- [ ] Manual testing completed\n- [ ] Code review ready\n\n",
		AIPRDescOutputRequirements:           "Output requirements:\n- Use Simplified Chinese (except for code and technical terms)\n- Use Markdown format\n- Concise and professional, highlight key points\n- Do not include title (# Pull Request), start directly from summary\n",
		AIPRDescBranchInfo:                   "Source branch: %s/%s\nTarget branch: %s\n",
		AIPRDescDiffUnavailable:              "[Unable to get complete diff: %v]\nNumber of commits: %d\n",
		AIPRDescDiffTruncated:                "\n[diff truncated, showing only first 15000 characters]",
		AIPRDescMoreCommits:                  "\n... %d more commits\n",
		// AI Code Review Prompts
		AICodeReviewSystemPrompt:             "You are a senior software engineer conducting a code review on the following git diff.\n\n",
		AICodeReviewFile:                     "**File:** %s\n\n",
		AICodeReviewCorePrinciples:           "## Core Principles\n",
		AICodeReviewConservative:             "- **Conservative review**: Only report issues you are **certain** exist. When uncertain, prefer not to report rather than guess.\n",
		AICodeReviewRespectContext:           "- **Respect context limitations**: You can only see the diff, not the complete file. If an issue requires full file context to judge (such as whether an error is already handled elsewhere), skip it and do not assume.\n",
		AICodeReviewFocusNewLines:            "- **Focus on new lines**: Focus on reviewing new lines starting with `+`; `-` deleted lines and context lines are only for understanding intent, do not comment on them.\n",
		AICodeReviewRejectFalsePositives:     "- **Reject false positives**: Do not flag correct idiomatic code as issues; do not mark code as problematic just because it's \"not how you would write it\".\n",
		AICodeReviewSeverityLevels:           "\n## Severity Levels (only use when confirmed)\n",
		AICodeReviewCritical:                 "- **CRITICAL**: Bugs that will cause crashes, data corruption, security vulnerabilities, or clear logic errors.\n",
		AICodeReviewMajor:                    "- **MAJOR**: Resource leaks, clear missing error handling (visible in diff), API usage errors.\n",
		AICodeReviewMinor:                    "- **MINOR**: Edge cases that might cause problems, code that could be more robust but currently works.\n",
		AICodeReviewNit:                      "- **NIT**: Pure style issues, only report when it truly affects readability.\n",
		AICodeReviewOutputFormat:             "\n## Output Format (output in Simplified Chinese, keep code snippets in original language)\n\n",
		AICodeReviewSummarySection:           "### Summary\nOne sentence explaining the purpose of this change and whether it is overall correct.\n\n",
		AICodeReviewIssuesSection:            "### Issue List\nUse the following format for each issue, with blank lines between issues:\n\n",
		AICodeReviewIssueFormat:              "**[Level] Category — Title**\nCode: `<problematic code snippet>`\nIssue: <issue description and impact>\nSuggestion: <specific fix or code>\n\n",
		AICodeReviewNoIssues:                 "If no issues, write directly: No issues\n\n",
		AICodeReviewConclusionSection:        "### Conclusion\n",
		AICodeReviewConclusionLGTM:           "No issues: LGTM, one sentence explaining it can be merged.\n",
		AICodeReviewConclusionHasIssues:      "Has issues: List CRITICAL/MAJOR items that must be fixed; MINOR/NIT can be summarized in one sentence.\n\n",
		AICodeReviewDiffSection:              "---\n\n## Diff\n",
		AICodeReviewLanguageHint:             " (%s)",
		AICodeReviewLanguageChecks:           "\n## Language-Specific Checks%s\n%s\n",
		// AI Error Messages
		AIRequestTimeout:                     "AI request timed out. Please try again later or adjust timeout in AI settings (Ctrl+A → Edit Profile → Timeout)",
		APIKeyInvalid:                        "API key is invalid. Please check the key in AI settings (Ctrl+A → Edit Profile → API Key)",
		APIRateLimitExceeded:                 "API rate limit exceeded. Please try again later or consider switching providers (Ctrl+A → Switch Profile)",
		NetworkConnectionFailed:              "Network connection failed. Please check your network connection or Endpoint configuration in AI settings (Ctrl+A → Edit Profile → Endpoint)",
		ModelNotAvailable:                    "Model is not available. Please select another model in AI settings (Ctrl+A → Edit Profile → Model)",
		APIQuotaExhausted:                    "API quota exhausted. Please check your account balance or switch providers (Ctrl+A → Switch Profile)",
		InputTooLong:                         "Input is too long and exceeds model limits. Please reduce the number of staged files or increase MaxTokens in AI settings (Ctrl+A → Edit Profile → MaxTokens)",
		// AI Context Messages
		CurrentBranch:                        "Current branch: %s",
		TrackingRemoteBranchAheadBehind:      "  Tracking remote branch: %s (local ahead %s commits, behind %s commits)",
		TrackingRemoteBranchAhead:            "  Tracking remote branch: %s (local ahead %s commits, not pushed)",
		TrackingRemoteBranchBehind:           "  Tracking remote branch: %s (local behind %s commits, need to pull)",
		TrackingRemoteBranchSynced:           "  Tracking remote branch: %s (synced)",
		NotTrackingRemoteBranch:              "  Not tracking remote branch (local branch)",
		WorkingTreeState:                     "\nWorking tree state: %s",
		ChangeStats:                          "\nChange stats: staged %d, unstaged %d, untracked %d",
		RecentCommits:                        "\nRecent commits:",
		ChangedFiles:                         "\nChanged files:",
		MoreFiles:                            "  ... (%d more files)",
		StashList:                            "\nStash list: %d stashes",
		MoreStashes:                          "  ... (%d more stashes)",
		// AI Branch and PR Messages
		AINotEnabledConfigFirst:              "AI is not enabled. Press Ctrl+A to configure AI settings",
		NoChangesForBranchName:               "No changes in working tree. Please make some changes before using AI branch name suggestion",
		AIBranchNameCancelled:                "AI branch name suggestion cancelled",
		NoCommitsForPRDescription:            "Current branch has no commits available for generating PR description",
		AIPRDescriptionCancelled:             "AI PR description generation cancelled",
		AIGenericError:                       "AI request failed: %v\n\nTip: Press Ctrl+A to check AI settings or switch provider",
		ChangedFilesLabel:                    "Changed files:",
		StagedFilesLabel:                     "Staged:",
		UnstagedFilesLabel:                   "Unstaged:",
		DiffTruncatedNote:                    "\n[diff truncated, showing first 8000 characters only]",
		DroppingStatus:                       "Dropping",
		MovingStatus:                         "Moving",
		RebasingStatus:                       "Rebasing",
		MergingStatus:                        "Merging",
		LowercaseRebasingStatus:              "rebasing",       // lowercase because it shows up in parentheses
		LowercaseMergingStatus:               "merging",        // lowercase because it shows up in parentheses
		LowercaseCherryPickingStatus:         "cherry-picking", // lowercase because it shows up in parentheses
		LowercaseRevertingStatus:             "reverting",      // lowercase because it shows up in parentheses
		AmendingStatus:                       "Amending",
		CherryPickingStatus:                  "Cherry-picking",
		UndoingStatus:                        "Undoing",
		RedoingStatus:                        "Redoing",
		CheckingOutStatus:                    "Checking out",
		CommittingStatus:                     "Committing",
		RewordingStatus:                      "Rewording",
		RevertingStatus:                      "Reverting",
		CreatingFixupCommitStatus:            "Creating fixup commit",
		MovingCommitsToNewBranchStatus:       "Moving commits to new branch",
		CommitFiles:                          "Commit files",
		SubCommitsDynamicTitle:               "Commits (%s)",
		CommitFilesDynamicTitle:              "Diff files (%s)",
		RemoteBranchesDynamicTitle:           "Remote branches (%s)",
		ViewItemFiles:                        "View files",
		CommitFilesTitle:                     "Commit files",
		CheckoutCommitFileTooltip:            "Checkout file. This replaces the file in your working tree with the version from the selected commit.",
		CannotCheckoutWithModifiedFilesErr:   "You have local modifications for the file(s) you are trying to check out. You need to stash or discard these first.",
		CanOnlyDiscardFromLocalCommits:       "Changes can only be discarded from local commits",
		Remove:                               "Remove",
		DiscardOldFileChangeTooltip:          "Discard this commit's changes to this file. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes this file.",
		DiscardFileChangesTitle:              "Discard file changes",
		DiscardFileChangesPrompt:             "Are you sure you want to remove changes to the selected file(s) from this commit?\n\nThis action will start a rebase, reverting these file changes. Be aware that if subsequent commits depend on these changes, you may need to resolve conflicts.\nNote: This will also reset any active custom patches.",
		DisabledForGPG:                       "Feature not available for users using GPG.\n\nIf you are using a passphrase agent (e.g. gpg-agent) so that you don't have to type your passphrase when signing, you can enable this feature by adding\n\ngit:\n  overrideGpg: true\n\nto your lazygit config file.",
		CreateRepo:                           "Not in a git repository. Create a new git repository? (y/N): ",
		BareRepo:                             "You've attempted to open Lazygit in a bare repo but Lazygit does not yet support bare repos. Open most recent repo? (y/n) ",
		InitialBranch:                        "Branch name? (leave empty for git's default): ",
		NoRecentRepositories:                 "Must open lazygit in a git repository. No valid recent repositories. Exiting.",
		IncorrectNotARepository:              "The value of 'notARepository' is incorrect. It should be one of 'prompt', 'create', 'skip', or 'quit'.",
		AutoStashTitle:                       "Autostash?",
		AutoStashPrompt:                      "You must stash and pop your changes to bring them across. Do this automatically? (enter/esc)",
		AutoStashForUndo:                     "Auto-stashing changes for undoing to %s",
		AutoStashForCheckout:                 "Auto-stashing changes for checking out %s",
		AutoStashForNewBranch:                "Auto-stashing changes for creating new branch %s",
		AutoStashForMovingPatchToIndex:       "Auto-stashing changes for moving custom patch to index from %s",
		AutoStashForCherryPicking:            "Auto-stashing changes for cherry-picking commits",
		AutoStashForReverting:                "Auto-stashing changes for reverting commits",
		Discard:                              "Discard",
		DiscardFileChangesTooltip:            "View options for discarding changes to the selected file.",
		DiscardChangesTitle:                  "Discard changes",
		Cancel:                               "Cancel",
		DiscardAllChanges:                    "Discard all changes",
		DiscardUnstagedChanges:               "Discard unstaged changes",
		DiscardAllChangesToAllFiles:          "Nuke working tree",
		DiscardAnyUnstagedChanges:            "Discard unstaged changes",
		DiscardUntrackedFiles:                "Discard untracked files",
		DiscardStagedChanges:                 "Discard staged changes",
		HardReset:                            "Hard reset",
		BranchDeleteTooltip:                  "View delete options for local/remote branch.",
		TagDeleteTooltip:                     "View delete options for local/remote tag.",
		Delete:                               "Delete",
		Reset:                                "Reset",
		ResetTooltip:                         "View reset options (soft/mixed/hard) for resetting onto selected item.",
		ResetSoftTooltip:                     "Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as staged changes.",
		ResetMixedTooltip:                    "Reset HEAD to the chosen commit, and keep the changes between the current and chosen commit as unstaged changes.",
		ResetHardTooltip:                     "Reset HEAD to the chosen commit, and discard all changes between the current and chosen commit, as well as all current modifications in the working tree.",
		ResetHardConfirmation:                "Are you sure you want to do a hard reset? This will discard all uncommitted changes (both staged and unstaged), which is not undoable.",
		ViewResetOptions:                     `Reset`,
		FileResetOptionsTooltip:              "View reset options for working tree (e.g. nuking the working tree).",
		FixupTooltip:                         "Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded.",
		CreateFixupCommit:                    "Create fixup commit",
		CreateFixupCommitTooltip:             "Create 'fixup!' commit for the selected commit. Later on, you can press `{{.squashAbove}}` on this same commit to apply all above fixup commits.",
		CreateAmendCommit:                    `Create "amend!" commit`,
		FixupMenu_Fixup:                      "fixup! commit",
		FixupMenu_FixupTooltip:               "Lets you fixup another commit and keep the original commit's message.",
		FixupMenu_AmendWithChanges:           "amend! commit with changes",
		FixupMenu_AmendWithChangesTooltip:    "Lets you fixup another commit and also change its commit message.",
		FixupMenu_AmendWithoutChanges:        "amend! commit without changes (pure reword)",
		FixupMenu_AmendWithoutChangesTooltip: "Lets you change the commit message of another commit without changing its content.",
		SquashAboveCommits:                   "Apply fixup commits",
		SquashAboveCommitsTooltip:            `Squash all 'fixup!' commits, either above the selected commit, or all in current branch (autosquash).`,
		SquashCommitsAboveSelectedTooltip:    `Squash all 'fixup!' commits above the selected commit (autosquash).`,
		SquashCommitsInCurrentBranchTooltip:  `Squash all 'fixup!' commits in the current branch (autosquash).`,
		SquashCommitsInCurrentBranch:         "In current branch",
		SquashCommitsAboveSelectedCommit:     "Above the selected commit",
		CannotSquashCommitsInCurrentBranch:   "Cannot squash commits in current branch: the HEAD commit is a merge commit or is present on the main branch.",
		ExecuteShellCommand:                  "Execute shell command",
		ExecuteShellCommandTooltip:           "Bring up a prompt where you can enter a shell command to execute.",
		ShellCommand:                         "Shell command:",
		ShellCommandAIMode:                   "Shell command (AI mode):",
		ShellCommandDangerousWarning:         "Dangerous Command Warning",
		CommitChangesWithoutHook:             "Commit changes without pre-commit hook",
		ResetTo:                              `Reset to`,
		PressEnterToReturn:                   "Press enter to return to lazygit",
		ViewStashOptions:                     "View stash options",
		ViewStashOptionsTooltip:              "View stash options (e.g. stash all, stash staged, stash unstaged).",
		Stash:                                "Stash",
		StashTooltip:                         "Stash all changes. For other variations of stashing, use the view stash options keybinding.",
		StashAllChanges:                      "Stash all changes",
		StashStagedChanges:                   "Stash staged changes",
		StashAllChangesKeepIndex:             "Stash all changes and keep index",
		StashUnstagedChanges:                 "Stash unstaged changes",
		StashIncludeUntrackedChanges:         "Stash all changes including untracked files",
		StashOptions:                         "Stash options",
		NotARepository:                       "Error: must be run inside a git repository",
		WorkingDirectoryDoesNotExist:         "Error: the current working directory does not exist",
		ScrollLeft:                           "Scroll left",
		ScrollRight:                          "Scroll right",
		DiscardPatch:                         "Discard patch",
		DiscardPatchConfirm:                  "You can only build a patch from one commit/stash-entry at a time. Discard current patch?",
		CantPatchWhileRebasingError:          "You cannot build a patch or run patch commands while in a merging or rebasing state",
		ToggleAddToPatch:                     "Toggle file included in patch",
		ToggleAddToPatchTooltip:              "Toggle whether the file is included in the custom patch. See {{.doc}}.",
		ToggleAllInPatch:                     "Toggle all files",
		ToggleAllInPatchTooltip:              "Add/remove all commit's files to custom patch. See {{.doc}}.",
		UpdatingPatch:                        "Updating patch",
		ViewPatchOptions:                     "View custom patch options",
		PatchOptionsTitle:                    "Patch options",
		NoPatchError:                         "No patch created yet. To start building a patch, use 'space' on a commit file or enter to add specific lines",
		EmptyPatchError:                      "Patch is still empty. Add some files or lines to your patch first.",
		EnterCommitFile:                      "Enter file / Toggle directory collapsed",
		EnterCommitFileTooltip:               "If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory.",
		ExitCustomPatchBuilder:               `Exit custom patch builder`,
		ExitFocusedMainView:                  "Exit back to side panel",
		EnterUpstream:                        `Enter upstream as '<remote> <branchname>'`,
		InvalidUpstream:                      "Invalid upstream. Must be in the format '<remote> <branchname>'",
		NewRemote:                            `New remote`,
		NewRemoteName:                        `New remote name:`,
		NewRemoteUrl:                         `New remote url:`,
		AddForkRemoteUsername:                `Fork owner (username/org). Use username:branch to check out a branch`,
		AddForkRemote:                        `Add fork remote`,
		AddForkRemoteTooltip:                 `Quickly add a fork remote by replacing the owner in the origin URL and optionally check out a branch from new remote.`,
		IncompatibleForkAlreadyExistsError:   `Remote {{.remoteName}} already exists and has different URL`,
		NoOriginRemote:                       "Action needs 'origin' remote",
		ViewBranches:                         "View branches",
		EditRemoteName:                       `Enter updated remote name for {{.remoteName}}:`,
		EditRemoteUrl:                        `Enter updated remote url for {{.remoteName}}:`,
		RemoveRemote:                         `Remove remote`,
		RemoveRemoteTooltip:                  `Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected.`,
		RemoveRemotePrompt:                   "Are you sure you want to remove remote?",
		DeleteRemoteBranch:                   "Delete remote branch",
		DeleteRemoteBranches:                 "Delete remote branches",
		DeleteRemoteBranchTooltip:            "Delete the remote branch from the remote.",
		DeleteLocalAndRemoteBranch:           "Delete local and remote branch",
		DeleteLocalAndRemoteBranches:         "Delete local and remote branches",
		SetAsUpstream:                        "Set as upstream",
		SetAsUpstreamTooltip:                 "Set the selected remote branch as the upstream of the checked-out branch.",
		SetUpstream:                          "Set upstream of selected branch",
		UnsetUpstream:                        "Unset upstream of selected branch",
		ViewDivergenceFromUpstream:           "View divergence from upstream",
		ViewDivergenceFromBaseBranch:         "View divergence from base branch ({{.baseBranch}})",
		CouldNotDetermineBaseBranch:          "Couldn't determine base branch",
		DivergenceSectionHeaderLocal:         "Local",
		DivergenceSectionHeaderRemote:        "Remote",
		ViewUpstreamResetOptions:             "Reset checked-out branch onto {{.upstream}}",
		ViewUpstreamResetOptionsTooltip:      "View options for resetting the checked-out branch onto {{upstream}}. Note: this will not reset the selected branch onto the upstream, it will reset the checked-out branch onto the upstream.",
		ViewUpstreamRebaseOptions:            "Rebase checked-out branch onto {{.upstream}}",
		ViewUpstreamRebaseOptionsTooltip:     "View options for rebasing the checked-out branch onto {{upstream}}. Note: this will not rebase the selected branch onto the upstream, it will rebase the checked-out branch onto the upstream.",
		UpstreamGenericName:                  "upstream of selected branch",
		SetUpstreamTitle:                     "Set upstream branch",
		SetUpstreamMessage:                   "Are you sure you want to set the upstream branch of '{{.checkedOut}}' to '{{.selected}}'?",
		EditRemoteTooltip:                    "Edit the selected remote's name or URL.",
		TagCommit:                            "Tag commit",
		TagCommitTooltip:                     "Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description.",
		TagNameTitle:                         "Tag name",
		TagMessageTitle:                      "Tag description",
		AnnotatedTag:                         "Annotated tag",
		LightweightTag:                       "Lightweight tag",
		DeleteTagTitle:                       "Delete tag '{{.tagName}}'?",
		DeleteLocalTag:                       "Delete local tag",
		DeleteRemoteTag:                      "Delete remote tag",
		DeleteLocalAndRemoteTag:              "Delete local and remote tag",
		RemoteTagDeletedMessage:              "Remote tag deleted",
		SelectRemoteTagUpstream:              "Remote from which to remove tag '{{.tagName}}':",
		DeleteRemoteTagPrompt:                "Are you sure you want to delete the remote tag '{{.tagName}}' from '{{.upstream}}'?",
		DeleteLocalAndRemoteTagPrompt:        "Are you sure you want to delete '{{.tagName}}' from both your machine and from '{{.upstream}}'?",
		PushTagTitle:                         "Remote to push tag '{{.tagName}}' to:",
		// Using 'push tag' rather than just 'push' to disambiguate from a global push
		PushTag:                        "Push tag",
		PushTagTooltip:                 "Push the selected tag to a remote. You'll be prompted to select a remote.",
		NewTag:                         "New tag",
		NewTagTooltip:                  "Create new tag from current commit. You'll be prompted to enter a tag name and optional description.",
		CreatingTag:                    "Creating tag",
		ForceTag:                       "Force Tag",
		ForceTagPrompt:                 "The tag '{{.tagName}}' exists already. Press {{.cancelKey}} to cancel, or {{.confirmKey}} to overwrite.",
		FetchRemoteTooltip:             "Fetch updates from the remote repository. This retrieves new commits and branches without merging them into your local branches.",
		CheckoutCommitTooltip:          "Checkout the selected commit as a detached HEAD.",
		NoBranchesFoundAtCommitTooltip: "No branches found at selected commit.",
		GitFlowOptions:                 "Show git-flow options",
		NotAGitFlowBranch:              "This does not seem to be a git flow branch",
		NewGitFlowBranchPrompt:         "New {{.branchType}} name:",

		IgnoreTracked:                    "Ignore tracked file",
		IgnoreTrackedPrompt:              "Are you sure you want to ignore a tracked file?",
		ExcludeTracked:                   "Exclude tracked file",
		ExcludeTrackedPrompt:             "Are you sure you want to exclude a tracked file?",
		ViewResetToUpstreamOptions:       "View upstream reset options",
		NextScreenMode:                   "Next screen mode (normal/half/fullscreen)",
		PrevScreenMode:                   "Prev screen mode",
		CyclePagers:                      "Cycle pagers",
		CyclePagersTooltip:               "Choose the next pager in the list of configured pagers",
		CyclePagersDisabledReason:        "No other pagers configured",
		StartSearch:                      "Search the current view by text",
		StartFilter:                      "Filter the current view by text",
		KeybindingsLegend:                "Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b",
		RenameBranch:                     "Rename branch",
		BranchUpstreamOptionsTitle:       "Upstream options",
		ViewBranchUpstreamOptions:        "View upstream options",
		ViewBranchUpstreamOptionsTooltip: "View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream.",
		UpstreamNotSetError:              "The selected branch has no upstream (or the upstream is not stored locally)",
		UpstreamsNotSetError:             "Some of the selected branches have no upstream (or the upstream is not stored locally)",
		Upstream:                         "Upstream",
		NewBranchNamePrompt:              "Enter new branch name for branch",
		RenameBranchWarning:              "This branch is tracking a remote. This action will only rename the local branch name, not the name of the remote branch. Continue?",
		OpenKeybindingsMenu:              "Open keybindings menu",
		ResetCherryPick:                  "Reset copied (cherry-picked) commits selection",
		ResetCherryPickShort:             "Reset copied commits",
		NextTab:                          "Next tab",
		PrevTab:                          "Previous tab",
		CantUndoWhileRebasing:            "Can't undo while rebasing",
		CantRedoWhileRebasing:            "Can't redo while rebasing",
		MustStashWarning:                 "Pulling a patch out into the index requires stashing and unstashing your changes. If something goes wrong, you'll be able to access your files from the stash. Continue?",
		MustStashTitle:                   "Must stash",
		ConfirmationTitle:                "Confirmation panel",
		PromptTitle:                      "Input prompt",
		PromptInputCannotBeEmptyToast:    "Empty input is not allowed",
		PrevPage:                         "Previous page",
		NextPage:                         "Next page",
		GotoTop:                          "Scroll to top",
		GotoBottom:                       "Scroll to bottom",
		FilteringBy:                      "Filtering by",
		ResetInParentheses:               "(Reset)",
		OpenFilteringMenu:                "View filter options",
		OpenFilteringMenuTooltip:         "View options for filtering the commit log, so that only commits matching the filter are shown.",
		FilterBy:                         "Filter by",
		ExitFilterMode:                   "Stop filtering",
		FilterPathOption:                 "Enter path to filter by",
		FilterAuthorOption:               "Enter author to filter by",
		EnterFileName:                    "Enter path:",
		EnterAuthor:                      "Enter author:",
		FilteringMenuTitle:               "Filtering",
		WillCancelExistingFilterTooltip:  "Note: this will cancel the existing filter",
		MustExitFilterModeTitle:          "Command not available",
		MustExitFilterModePrompt:         "Command not available in filter-by-path mode. Exit filter-by-path mode?",
		Diff:                             "Diff",
		EnterRefToDiff:                   "Enter ref to diff",
		EnterRefName:                     "Enter ref:",
		ExitDiffMode:                     "Exit diff mode",
		DiffingMenuTitle:                 "Diffing",
		SwapDiff:                         "Reverse diff direction",
		ViewDiffingOptions:               "View diffing options",
		ViewDiffingOptionsTooltip:        "View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction.",
		CancelDiffingMode:                "Cancel diffing mode",
		// the actual view is the extras view which I intend to give more tabs in future but for now we'll only mention the command log part
		OpenCommandLogMenu:                       "View command log options",
		OpenCommandLogMenuTooltip:                "View options for the command log e.g. show/hide the command log and focus the command log.",
		OpenAIAssistant:                          "Open AI git assistant",
		ShowingGitDiff:                           "Showing output for:",
		ShowingDiffForRange:                      "Showing diff for range",
		CommitDiff:                               "Commit diff",
		CopyCommitHashToClipboard:                "Copy commit hash to clipboard",
		CommitHash:                               "Commit hash",
		CommitURL:                                "Commit URL",
		PasteCommitMessageFromClipboard:          "Paste commit message from clipboard",
		SurePasteCommitMessage:                   "Pasting will overwrite the current commit message, continue?",
		CommitMessage:                            "Commit message (subject and body)",
		CommitMessageBody:                        "Commit message body",
		CommitSubject:                            "Commit subject",
		CommitAuthor:                             "Commit author",
		CommitTags:                               "Commit tags",
		CopyCommitAttributeToClipboard:           "Copy commit attribute to clipboard",
		CopyCommitAttributeToClipboardTooltip:    "Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author).",
		CopyBranchNameToClipboard:                "Copy branch name to clipboard",
		CopyTagToClipboard:                       "Copy tag to clipboard",
		CopyPathToClipboard:                      "Copy path to clipboard",
		CopySelectedTextToClipboard:              "Copy selected text to clipboard",
		CommitPrefixPatternError:                 "Error in commitPrefix pattern",
		NoFilesStagedTitle:                       "No files staged",
		NoFilesStagedPrompt:                      "You have not staged any files. Commit all files?",
		BranchNotFoundTitle:                      "Branch not found",
		BranchNotFoundPrompt:                     "Branch not found. Create a new branch named",
		BranchUnknown:                            "Branch unknown",
		DiscardChangeTitle:                       "Discard change",
		DiscardChangePrompt:                      "Are you sure you want to discard this change (git reset)? It is irreversible.\nTo disable this dialogue set the config key of 'gui.skipDiscardChangeWarning' to true",
		CreateNewBranchFromCommit:                "Create new branch off of commit",
		BuildingPatch:                            "Building patch",
		ViewCommits:                              "View commits",
		MinGitVersionError:                       "Git version must be at least %s. Please upgrade your git version.",
		RunningCustomCommandStatus:               "Running custom command",
		SubmoduleStashAndReset:                   "Stash uncommitted submodule changes and update",
		AndResetSubmodules:                       "And reset submodules",
		Enter:                                    "Enter",
		EnterSubmoduleTooltip:                    "Enter submodule. After entering the submodule, you can press `{{.escape}}` to escape back to the parent repo.",
		BackToParentRepo:                         "Back to parent repo",
		CopySubmoduleNameToClipboard:             "Copy submodule name to clipboard",
		RemoveSubmodule:                          "Remove submodule",
		RemoveSubmodulePrompt:                    "Are you sure you want to remove submodule '%s' and its corresponding directory? This is irreversible.",
		RemoveSubmoduleTooltip:                   "Remove the selected submodule and its corresponding directory.",
		ResettingSubmoduleStatus:                 "Resetting submodule",
		NewSubmoduleName:                         "New submodule name:",
		NewSubmoduleUrl:                          "New submodule URL:",
		NewSubmodulePath:                         "New submodule path:",
		NewSubmodule:                             "New submodule",
		AddingSubmoduleStatus:                    "Adding submodule",
		UpdateSubmoduleUrl:                       "Update URL for submodule '%s'",
		UpdatingSubmoduleUrlStatus:               "Updating URL",
		EditSubmoduleUrl:                         "Update submodule URL",
		InitializingSubmoduleStatus:              "Initializing submodule",
		InitSubmoduleTooltip:                     "Initialize the selected submodule to prepare for fetching. You probably want to follow this up by invoking the 'update' action to fetch the submodule.",
		Update:                                   "Update",
		Initialize:                               "Initialize",
		SubmoduleUpdateTooltip:                   "Update selected submodule.",
		UpdatingSubmoduleStatus:                  "Updating submodule",
		BulkInitSubmodules:                       "Bulk init submodules",
		BulkUpdateSubmodules:                     "Bulk update submodules",
		BulkDeinitSubmodules:                     "Bulk deinit submodules",
		BulkUpdateRecursiveSubmodules:            "Bulk init and update submodules recursively",
		ViewBulkSubmoduleOptions:                 "View bulk submodule options",
		BulkSubmoduleOptions:                     "Bulk submodule options",
		RunningCommand:                           "Running command",
		SubCommitsTitle:                          "Sub-commits",
		ExitSubview:                              "Exit subview",
		SubmodulesTitle:                          "Submodules",
		NavigationTitle:                          "List panel navigation",
		SuggestionsCheatsheetTitle:               "Suggestions",
		SuggestionsTitle:                         "Suggestions (press %s to focus)",
		SuggestionsSubtitle:                      "(press %s to delete, %s to edit)",
		ExtrasTitle:                              "Command log",
		PullRequestURLCopiedToClipboard:          "Pull request URL copied to clipboard",
		CommitDiffCopiedToClipboard:              "Commit diff copied to clipboard",
		CommitURLCopiedToClipboard:               "Commit URL copied to clipboard",
		CommitMessageCopiedToClipboard:           "Commit message copied to clipboard",
		CommitMessageBodyCopiedToClipboard:       "Commit message body copied to clipboard",
		CommitSubjectCopiedToClipboard:           "Commit subject copied to clipboard",
		CommitAuthorCopiedToClipboard:            "Commit author copied to clipboard",
		CommitTagsCopiedToClipboard:              "Commit tags copied to clipboard",
		CommitHasNoTags:                          "Commit has no tags",
		CommitHasNoMessageBody:                   "Commit has no message body",
		PatchCopiedToClipboard:                   "Patch copied to clipboard",
		MessageCopiedToClipboard:                 "Message copied to clipboard",
		CopiedToClipboard:                        "copied to clipboard",
		ErrCannotEditDirectory:                   "Cannot edit directories: you can only edit individual files",
		ErrCannotCopyContentOfDirectory:          "Cannot copy content of directories: you can only copy content of individual files",
		ErrStageDirWithInlineMergeConflicts:      "Cannot stage/unstage directory containing files with inline merge conflicts. Please fix up the merge conflicts first",
		ErrRepositoryMovedOrDeleted:              "Cannot find repo. It might have been moved or deleted ¯\\_(ツ)_/¯",
		CommandLog:                               "Command log",
		ErrWorktreeMovedOrRemoved:                "Cannot find worktree. It might have been moved or removed ¯\\_(ツ)_/¯",
		ToggleShowCommandLog:                     "Toggle show/hide command log",
		FocusCommandLog:                          "Focus command log",
		CopyCommandLog:                           "Copy command log to clipboard",
		CommandLogCopiedToClipboard:             "Command log copied to clipboard",
		CommandLogHeader:                         "You can hide/focus this panel by pressing '%s'\n",
		RandomTip:                                "Random tip",
		ToggleWhitespaceInDiffView:               "Toggle whitespace",
		ToggleWhitespaceInDiffViewTooltip:        "Toggle whether or not whitespace changes are shown in the diff view.\n\nThe default can be changed in the config file with the key 'git.ignoreWhitespaceInDiffView'.",
		IgnoreWhitespaceDiffViewSubTitle:         "(ignoring whitespace)",
		IgnoreWhitespaceNotSupportedHere:         "Ignoring whitespace is not supported in this view",
		IncreaseContextInDiffView:                "Increase diff context size",
		IncreaseContextInDiffViewTooltip:         "Increase the amount of the context shown around changes in the diff view.\n\nThe default can be changed in the config file with the key 'git.diffContextSize'.",
		DecreaseContextInDiffView:                "Decrease diff context size",
		DecreaseContextInDiffViewTooltip:         "Decrease the amount of the context shown around changes in the diff view.\n\nThe default can be changed in the config file with the key 'git.diffContextSize'.",
		DiffContextSizeChanged:                   "Changed diff context size to %d",
		IncreaseRenameSimilarityThresholdTooltip: "Increase the similarity threshold for a deletion and addition pair to be treated as a rename.\n\nThe default can be changed in the config file with the key 'git.renameSimilarityThreshold'.",
		IncreaseRenameSimilarityThreshold:        "Increase rename similarity threshold",
		DecreaseRenameSimilarityThresholdTooltip: "Decrease the similarity threshold for a deletion and addition pair to be treated as a rename.\n\nThe default can be changed in the config file with the key 'git.renameSimilarityThreshold'.",
		DecreaseRenameSimilarityThreshold:        "Decrease rename similarity threshold",
		RenameSimilarityThresholdChanged:         "Changed rename similarity threshold to %d%%",
		CreatePullRequestOptions:                 "View create pull request options",
		DefaultBranch:                            "Default branch",
		SelectBranch:                             "Select branch",
		SelectTargetRemote:                       "Select target remote",
		NoValidRemoteName:                        "A remote named '%s' does not exist",
		SelectConfigFile:                         "Select config file",
		NoConfigFileFoundErr:                     "No config file found",
		LoadingFileSuggestions:                   "Loading file suggestions",
		LoadingCommits:                           "Loading commits",
		MustSpecifyOriginError:                   "Must specify a remote if specifying a branch",
		GitOutput:                                "Git output:",
		GitCommandFailed:                         "Git command failed. Check command log for details (open with %s)",
		AbortTitle:                               "Abort %s",
		AbortPrompt:                              "Are you sure you want to abort the current %s?",
		OpenLogMenu:                              "View log options",
		OpenLogMenuTooltip:                       "View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph.",
		LogMenuTitle:                             "Commit Log Options",
		ToggleShowGitGraphAll:                    "Toggle show whole git graph (pass the `--all` flag to `git log`)",
		ShowGitGraph:                             "Show git graph",
		ShowGitGraphTooltip:                      "Show or hide the git graph in the commit log.\n\nThe default can be changed in the config file with the key 'git.log.showGraph'.",
		SortOrder:                                "Sort order",
		SortOrderPromptLocalBranches:             "The default sort order for local branches can be set in the config file with the key 'git.localBranchSortOrder'.",
		SortOrderPromptRemoteBranches:            "The default sort order for remote branches can be set in the config file with the key 'git.remoteBranchSortOrder'.",
		SortAlphabetical:                         "Alphabetical",
		SortByDate:                               "Date",
		SortByRecency:                            "Recency",
		SortBasedOnReflog:                        "(based on reflog)",
		SortCommits:                              "Commit sort order",
		SortCommitsTooltip:                       "Change the sort order of the commits in the commit log.\n\nThe default can be changed in the config file with the key 'git.log.sortOrder'.",
		CantChangeContextSizeError:               "Cannot change context while in patch building mode because we were too lazy to support it when releasing the feature. If you really want it, please let us know!",
		OpenCommitInBrowser:                      "Open commit in browser",
		ViewBisectOptions:                        "View bisect options",
		ViewBranchesContainingCommit:             "View branches containing this commit",
		ViewBranchesContainingCommitTooltip:      "Show all local and remote branches that contain this commit.",
		NoBranchesContainingCommit:               "No branches contain this commit",
		EnterCommitHashToFindBranches:            "Enter commit hash:",
		ConfirmRevertCommit:                      "Are you sure you want to revert {{.selectedCommit}}?",
		ConfirmRevertCommitRange:                 "Are you sure you want to revert the selected commits?",
		RewordInEditorTitle:                      "Reword in editor",
		RewordInEditorPrompt:                     "Are you sure you want to reword this commit in your editor?",
		HardResetAutostashPrompt:                 "Are you sure you want to hard reset to '%s'? An auto-stash will be performed if necessary.",
		SoftResetPrompt:                          "Are you sure you want to soft reset to '%s'?",
		CheckoutAutostashPrompt:                  "Are you sure you want to checkout '%s'? An auto-stash will be performed if necessary.",
		UpstreamGone:                             "(upstream gone)",
		NukeDescription:                          "If you want to make all the changes in the worktree go away, this is the way to do it. If there are dirty submodule changes this will stash those changes in the submodule(s).",
		NukeTreeConfirmation:                     "Are you sure you want to nuke the working tree? This will discard all changes in the worktree (staged, unstaged and untracked), which is not undoable.",
		DiscardStagedChangesDescription:          "This will create a new stash entry containing only staged files and then drop it, so that the working tree is left with only unstaged changes",
		EmptyOutput:                              "<Empty output>",
		Patch:                                    "Patch",
		CustomPatch:                              "Custom patch",
		CommitsCopied:                            "commits copied", // lowercase because it's used in a sentence
		CommitCopied:                             "commit copied",  // lowercase because it's used in a sentence
		ResetPatch:                               "Reset patch",
		ResetPatchTooltip:                        "Clear the current patch.",
		ApplyPatch:                               "Apply patch",
		ApplyPatchTooltip:                        "Apply the current patch to the working tree.",
		ApplyPatchInReverse:                      "Apply patch in reverse",
		ApplyPatchInReverseTooltip:               "Apply the current patch in reverse to the working tree.",
		RemovePatchFromOriginalCommit:            "Remove patch from original commit (%s)",
		RemovePatchFromOriginalCommitTooltip:     "Remove the current patch from its commit. This is achieved by starting an interactive rebase at the commit, applying the patch in reverse, and then continuing the rebase. If later commits depend on the patch, you may need to resolve conflicts.",
		MovePatchOutIntoIndex:                    "Move patch out into index",
		MovePatchOutIntoIndexTooltip:             "Move the patch out of its commit and into the index. This is achieved by starting an interactive rebase at the commit, applying the patch in reverse, continuing the rebase to completion, and then applying the patch to the index. If later commits depend on the patch, you may need to resolve conflicts.",
		MovePatchIntoNewCommit:                   "Move patch into new commit after the original commit",
		MovePatchIntoNewCommitTooltip:            "Move the patch out of its commit and into a new commit sitting on top of the original commit. This is achieved by starting an interactive rebase at the original commit, applying the patch in reverse, then applying the patch to the index and committing it as a new commit, before continuing the rebase to completion. If later commits depend on the patch, you may need to resolve conflicts.",
		MovePatchIntoNewCommitBefore:             "Move patch into new commit before the original commit",
		MovePatchIntoNewCommitBeforeTooltip:      "Move the patch out of its commit and into a new commit before the original commit. This works best when the custom patch contains only entire hunks or even entire files; if it contains partial hunks, you are likely to get conflicts.",
		MovePatchToSelectedCommit:                "Move patch to selected commit (%s)",
		MovePatchToSelectedCommitTooltip:         "Move the patch out of its original commit and into the selected commit. This is achieved by starting an interactive rebase at the original commit, applying the patch in reverse, then continuing the rebase up to the selected commit, before applying the patch forward and amending the selected commit. The rebase is then continued to completion. If commits between the source and destination commit depend on the patch, you may need to resolve conflicts.",
		CopyPatchToClipboard:                     "Copy patch to clipboard",
		MustStageFilesAffectedByPatchTitle:       "Must stage files",
		MustStageFilesAffectedByPatchWarning:     "Applying a patch to the index requires staging the unstaged files that are affected by the patch. Note that you might get conflicts when applying the patch. Continue?",
		NoMatchesFor:                             "No matches for '%s' %s",
		ExitSearchMode:                           "%s: Exit search mode",
		ExitTextFilterMode:                       "%s: Exit filter mode",
		MatchesFor:                               "matches for '%s' (%d of %d) %s", // lowercase because it's after other text
		SearchKeybindings:                        "%s: Next match, %s: Previous match, %s: Exit search mode",
		SearchPrefix:                             "Search: ",
		FilterPrefix:                             "Filter: ",
		FilterPrefixMenu:                         "Filter (prepend '@' to filter keybindings): ",
		WorktreesTitle:                           "Worktrees",
		WorktreeTitle:                            "Worktree",
		Switch:                                   "Switch",
		SwitchToWorktree:                         "Switch to worktree",
		SwitchToWorktreeTooltip:                  "Switch to the selected worktree.",
		AlreadyCheckedOutByWorktree:              "This branch is checked out by worktree {{.worktreeName}}. Do you want to switch to that worktree?",
		BranchCheckedOutByWorktree:               "Branch {{.branchName}} is checked out by worktree {{.worktreeName}}",
		SomeBranchesCheckedOutByWorktreeError:    "Some of the selected branches are checked out by other worktrees. Select them one by one to delete them.",
		DetachWorktreeTooltip:                    "This will run `git checkout --detach` on the worktree so that it stops hogging the branch, but the worktree's working tree will be left alone.",
		Switching:                                "Switching",
		RemoveWorktree:                           "Remove worktree",
		RemoveWorktreeTitle:                      "Remove worktree",
		RemoveWorktreePrompt:                     "Are you sure you want to remove worktree '{{.worktreeName}}'?",
		ForceRemoveWorktreePrompt:                "'{{.worktreeName}}' contains modified or untracked files, or submodules (or all of these). Are you sure you want to remove it?",
		RemovingWorktree:                         "Deleting worktree",
		DetachWorktree:                           "Detach worktree",
		DetachingWorktree:                        "Detaching worktree",
		AddingWorktree:                           "Adding worktree",
		CantDeleteCurrentWorktree:                "You cannot remove the current worktree!",
		AlreadyInWorktree:                        "You are already in the selected worktree",
		CantDeleteMainWorktree:                   "You cannot remove the main worktree!",
		NoWorktreesThisRepo:                      "No worktrees",
		MissingWorktree:                          "(missing)",
		MainWorktree:                             "(main)",
		NewWorktree:                              "New worktree",
		NewWorktreePath:                          "New worktree path",
		NewWorktreeBase:                          "New worktree base ref",
		RemoveWorktreeTooltip:                    "Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory.",
		NewBranchName:                            "New branch name",
		NewBranchNameLeaveBlank:                  "New branch name (leave blank to checkout {{.default}})",
		ViewWorktreeOptions:                      "View worktree options",
		CreateWorktreeFrom:                       "Create worktree from {{.ref}}",
		CreateWorktreeFromDetached:               "Create worktree from {{.ref}} (detached)",
		LcWorktree:                               "worktree",
		ChangingDirectoryTo:                      "Changing directory to {{.path}}",
		Name:                                     "Name",
		Branch:                                   "Branch",
		Path:                                     "Path",
		MarkedBaseCommitStatus:                   "Marked a base commit for rebase",
		MarkAsBaseCommit:                         "Mark as base commit for rebase",
		MarkAsBaseCommitTooltip:                  "Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command.",
		CancelMarkedBaseCommit:                   "Cancel marked base commit",
		MarkedCommitMarker:                       "↑↑↑ Will rebase from here ↑↑↑",
		FailedToOpenURL:                          "Failed to open URL %s\n\nError: %v",
		InvalidLazygitEditURL:                    "Invalid lazygit-edit URL format: %s",
		DisabledMenuItemPrefix:                   "Disabled: ",
		NoCopiedCommits:                          "No copied commits",
		QuickStartInteractiveRebase:              "Start interactive rebase",
		QuickStartInteractiveRebaseTooltip:       "Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.\nIf you would instead like to start an interactive rebase from the selected commit, press `{{.editKey}}`.",
		CannotQuickStartInteractiveRebase:        "Cannot start interactive rebase: the HEAD commit is a merge commit or is present on the main branch, so there is no appropriate base commit to start the rebase from. You can start an interactive rebase from a specific commit by selecting the commit and pressing `{{.editKey}}`.",
		RangeSelectUp:                            "Range select up",
		RangeSelectDown:                          "Range select down",
		RangeSelectNotSupported:                  "Action does not support range selection, please select a single item",
		NoItemSelected:                           "No item selected",
		SelectedItemIsNotABranch:                 "Selected item is not a branch",
		SelectedItemDoesNotHaveFiles:             "Selected item does not have files to view",
		MultiSelectNotSupportedForSubmodules:     "Multiselection not supported for submodules",
		CommandDoesNotSupportOpeningInEditor:     "This command doesn't support switching to the editor",
		CustomCommands:                           "Custom commands",
		NoApplicableCommandsInThisContext:        "(No applicable commands in this context)",
		SelectCommitsOfCurrentBranch:             "Select commits of current branch",
		ViewMergeConflictOptions:                 "View merge conflict options",
		ViewMergeConflictOptionsTooltip:          "View options for resolving merge conflicts.",
		NoFilesWithMergeConflicts:                "There are no files with merge conflicts.",
		MergeConflictOptionsTitle:                "Resolve merge conflicts",
		UseCurrentChanges:                        "Use current changes",
		UseIncomingChanges:                       "Use incoming changes",
		UseBothChanges:                           "Use both",

		Actions: Actions{
			// TODO: combine this with the original keybinding descriptions (those are all in lowercase atm)
			CheckoutCommit:                   "Checkout commit",
			CheckoutBranchAtCommit:           "Checkout branch '%s'",
			CheckoutCommitAsDetachedHead:     "Checkout commit %s as detached head",
			CheckoutTag:                      "Checkout tag",
			CheckoutBranch:                   "Checkout branch",
			ForceCheckoutBranch:              "Force checkout branch",
			CheckoutBranchOrCommit:           "Checkout branch or commit",
			DeleteLocalBranch:                "Delete local branch",
			Merge:                            "Merge",
			SquashMerge:                      "Squash merge",
			RebaseBranch:                     "Rebase branch",
			RenameBranch:                     "Rename branch",
			CreateBranch:                     "Create branch",
			CherryPick:                       "(Cherry-pick) paste commits",
			CheckoutFile:                     "Checkout file",
			SquashCommitDown:                 "Squash commit down",
			FixupCommit:                      "Fixup commit",
			FixupCommitKeepMessage:           "Fixup commit (keep message)",
			RewordCommit:                     "Reword commit",
			DropCommit:                       "Drop commit",
			EditCommit:                       "Edit commit",
			AmendCommit:                      "Amend commit",
			ResetCommitAuthor:                "Reset commit author",
			SetCommitAuthor:                  "Set commit author",
			AddCommitCoAuthor:                "Add commit co-author",
			RevertCommit:                     "Revert commit",
			CreateFixupCommit:                "Create fixup commit",
			SquashAllAboveFixupCommits:       "Squash all above fixup commits",
			CreateLightweightTag:             "Create lightweight tag",
			CreateAnnotatedTag:               "Create annotated tag",
			CopyCommitMessageToClipboard:     "Copy commit message to clipboard",
			CopyCommitMessageBodyToClipboard: "Copy commit message body to clipboard",
			CopyCommitSubjectToClipboard:     "Copy commit subject to clipboard",
			CopyCommitTagsToClipboard:        "Copy commit tags to clipboard",
			CopyCommitDiffToClipboard:        "Copy commit diff to clipboard",
			CopyCommitHashToClipboard:        "Copy full commit hash to clipboard",
			CopyCommitURLToClipboard:         "Copy commit URL to clipboard",
			CopyCommitAuthorToClipboard:      "Copy commit author to clipboard",
			CopyCommitAttributeToClipboard:   "Copy to clipboard",
			CopyPatchToClipboard:             "Copy patch to clipboard",
			MoveCommitUp:                     "Move commit up",
			MoveCommitDown:                   "Move commit down",
			CustomCommand:                    "Custom command",
			DiscardAllChangesInFile:          "Discard all changes in selected file(s)",
			DiscardAllUnstagedChangesInFile:  "Discard all unstaged changes selected file(s)",
			StageFile:                        "Stage file",
			StageResolvedFiles:               "Stage files whose merge conflicts were resolved",
			UnstageFile:                      "Unstage file",
			UnstageAllFiles:                  "Unstage all files",
			StageAllFiles:                    "Stage all files",
			ResolveConflictByKeepingFile:     "Resolve by keeping file",
			ResolveConflictByDeletingFile:    "Resolve by deleting file",
			NotEnoughContextToStage:          "Staging or unstaging changes is not possible with a diff context size of 0. Increase the context using '%s'.",
			NotEnoughContextToDiscard:        "Discarding changes is not possible with a diff context size of 0. Increase the context using '%s'.",
			NotEnoughContextForCustomPatch:   "Creating custom patches is not possible with a diff context size of 0. Increase the context using '%s'.",
			IgnoreExcludeFile:                "Ignore or exclude file",
			IgnoreFileErr:                    "Cannot ignore .gitignore",
			ExcludeFile:                      "Exclude file",
			ExcludeGitIgnoreErr:              "Cannot exclude .gitignore",
			Commit:                           "Commit",
			Push:                             "Push",
			Pull:                             "Pull",
			OpenFile:                         "Open file",
			StashAllChanges:                  "Stash all changes",
			StashAllChangesKeepIndex:         "Stash all changes and keep index",
			StashStagedChanges:               "Stash staged changes",
			StashUnstagedChanges:             "Stash unstaged changes",
			StashIncludeUntrackedChanges:     "Stash all changes including untracked files",
			GitFlowFinish:                    "git flow finish",
			GitFlowStart:                     "git flow start",
			CopyToClipboard:                  "Copy to clipboard",
			CopySelectedTextToClipboard:      "Copy selected text to clipboard",
			RemovePatchFromCommit:            "Remove patch from commit",
			MovePatchToSelectedCommit:        "Move patch to selected commit",
			MovePatchIntoIndex:               "Move patch into index",
			MovePatchIntoNewCommit:           "Move patch into new commit",
			DeleteRemoteBranch:               "Delete remote branch",
			SetBranchUpstream:                "Set branch upstream",
			AddRemote:                        "Add remote",
			AddForkRemote:                    "Add fork remote",
			RemoveRemote:                     "Remove remote",
			UpdateRemote:                     "Update remote",
			ApplyPatch:                       "Apply patch",
			Stash:                            "Stash",
			PopStash:                         "Pop stash",
			ApplyStash:                       "Apply stash",
			DropStash:                        "Drop stash",
			RenameStash:                      "Rename stash",
			RemoveSubmodule:                  "Remove submodule",
			ResetSubmodule:                   "Reset submodule",
			AddSubmodule:                     "Add submodule",
			UpdateSubmoduleUrl:               "Update submodule URL",
			InitialiseSubmodule:              "Initialise submodule",
			BulkInitialiseSubmodules:         "Bulk initialise submodules",
			BulkUpdateSubmodules:             "Bulk update submodules",
			BulkDeinitialiseSubmodules:       "Bulk deinitialise submodules",
			BulkUpdateRecursiveSubmodules:    "Bulk initialise and update submodules recursively",
			UpdateSubmodule:                  "Update submodule",
			DeleteLocalTag:                   "Delete local tag",
			DeleteRemoteTag:                  "Delete remote tag",
			PushTag:                          "Push tag",
			NukeWorkingTree:                  "Nuke working tree",
			DiscardUnstagedFileChanges:       "Discard unstaged file changes",
			RemoveUntrackedFiles:             "Remove untracked files",
			RemoveStagedFiles:                "Remove staged files",
			SoftReset:                        "Soft reset",
			MixedReset:                       "Mixed reset",
			HardReset:                        "Hard reset",
			FastForwardBranch:                "Fast forward branch",
			AutoForwardBranches:              "Auto-forward branches",
			Undo:                             "Undo",
			Redo:                             "Redo",
			CopyPullRequestURL:               "Copy pull request URL",
			OpenMergeTool:                    "Open merge tool",
			OpenCommitInBrowser:              "Open commit in browser",
			OpenPullRequest:                  "Open pull request in browser",
			StartBisect:                      "Start bisect",
			ResetBisect:                      "Reset bisect",
			BisectSkip:                       "Bisect skip",
			BisectMark:                       "Bisect mark",
			AddWorktree:                      "Add worktree",
		},
		Bisect: Bisect{
			Mark:                        "Mark current commit (%s) as %s",
			MarkStart:                   "Mark %s as %s (start bisect)",
			SkipCurrent:                 "Skip current commit (%s)",
			SkipSelected:                "Skip selected commit (%s)",
			ResetTitle:                  "Reset 'git bisect'",
			ResetPrompt:                 "Are you sure you want to reset 'git bisect'?",
			ResetOption:                 "Reset bisect",
			ChooseTerms:                 "Choose bisect terms",
			OldTermPrompt:               "Term for old/good commit:",
			NewTermPrompt:               "Term for new/bad commit:",
			BisectMenuTitle:             "Bisect",
			CompleteTitle:               "Bisect complete",
			CompletePrompt:              "Bisect complete! The following commit introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
			CompletePromptIndeterminate: "Bisect complete! Some commits were skipped, so any of the following commits may have introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
			Bisecting:                   "Bisecting",
		},
		Log: Log{
			EditRebase:               "Beginning interactive rebase at '{{.ref}}'",
			HandleUndo:               "Undoing last conflict resolution",
			RemoveFile:               "Deleting path '{{.path}}'",
			CopyToClipboard:          "Copying '{{.str}}' to clipboard",
			Remove:                   "Removing '{{.filename}}'",
			CreateFileWithContent:    "Creating file '{{.path}}'",
			AppendingLineToFile:      "Appending '{{.line}}' to file '{{.filename}}'",
			EditRebaseFromBaseCommit: "Beginning interactive rebase from '{{.baseCommit}}' onto '{{.targetBranchName}}",
		},
		BreakingChangesTitle: "Breaking Changes",
		BreakingChangesMessage: `You are updating to a new version of lazygit which contains breaking changes. Please review the notes below and update your configuration if necessary.
For more information, see the full release notes at <https://github.com/dswcpp/lazygit/releases>.`,
		BreakingChangesByVersion: map[string]string{
			"0.41.0": `- When you press 'g' to bring up the git reset menu, the 'mixed' option is now the first and default, rather than 'soft'. This is because 'mixed' is the most commonly used option.
- The commit message panel now automatically hard-wraps by default (i.e. it adds newline characters when you reach the margin). You can adjust the config like so:

git:
  commit:
    autoWrapCommitMessage: true
    autoWrapWidth: 72

- The 'v' key was already being used in the staging view to start a range select, but now you can use it to start a range select in any view. Unfortunately this clashes with the 'v' keybinding for pasting commits (cherry-pick), so now pasting commits is done via 'shift+V' and for the sake of consistency, copying commits is now done via 'shift+C' instead of just 'c'. Note that the 'v' keybinding is only one way to start a range-select: you can use shift+up/down arrow instead. So, if you want to configure the cherry-pick keybindings to get the old behaviour, set the following in your config:

keybinding:
  universal:
      toggleRangeSelect: <something other than v>
    commits:
      cherryPickCopy: 'c'
      pasteCommits: 'v'

- Squashing fixups using 'shift-S' now brings up a menu, with the default option being to squash all fixup commits in the branch. The original behaviour of only squashing fixup commits above the selected commit is still available as the second option in that menu.
- Push/pull/fetch loading statuses are now shown against the branch rather than in a popup. This allows you to e.g. fetch multiple branches in parallel and see the status for each branch.
- The git log graph in the commits view is now always shown by default (previously it was only shown when the view was maximised). If you find this too noisy, you can change it back via ctrl+L -> 'Show git graph' -> 'when maximised'
- Pressing space on a remote branch used to show a prompt for entering a name for a new local branch to check out from the remote branch. Now it just checks out the remote branch directly, letting you choose between a new local branch with the same name, or a detached head. The old behavior is still available via the 'n' keybinding.
- Filtering (e.g. when pressing '/') is less fuzzy by default; it only matches substrings now. Multiple substrings can be matched by separating them with spaces. If you want to revert to the old behavior, set the following in your config:

gui:
  filterMode: 'fuzzy'
`,
			"0.44.0": `- The gui.branchColors config option is deprecated; it will be removed in a future version. Please use gui.branchColorPatterns instead.
- The automatic coloring of branches starting with "feature/", "bugfix/", or "hotfix/" has been removed; if you want this, it's easy to set up using the new gui.branchColorPatterns option.`,
			"0.49.0": `- Executing shell commands (with the ':' prompt) no longer uses an interactive shell, which means that if you want to use your shell aliases in this prompt, you need to do a little bit of setup work. See https://github.com/dswcpp/lazygit/blob/master/docs/Config.md#using-aliases-or-functions-in-shell-commands for details.`,
			"0.50.0": `- After fetching, main branches now get auto-forwarded to their upstream if they fall behind. This is useful for keeping your main or master branch up to date automatically. If you don't want this, you can disable it by setting the following in your config:

git:
  autoForwardBranches: none

If, on the other hand, you want this even for feature branches, you can set it to 'allBranches' instead.`,
			"0.51.0": `- The 'subprocess', 'stream', and 'showOutput' fields of custom commands have been replaced by a single 'output' field. This should be transparent, if you used these in your config file it should have been automatically updated for you. There's one notable change though: the 'stream' field used to mean both that the command's output would be streamed to the command log, and that the command would be run in a pseudo terminal (pty). We converted this to 'output: log', which means that the command's output will be streamed to the command log, but not use a pty, on the assumption that this is what most people wanted. If you do actually want to run a command in a pty, you can change this to 'output: logWithPty' instead.`,
			"0.54.0": `- The default sort order for local and remote branches has changed: it used to be 'recency' (based on reflog) for local branches, and 'alphabetical' for remote branches. Both of these have been changed to 'date' (which means committerdate). If you do liked the old defaults better, you can revert to them with the following config:

git:
  localBranchSortOrder: recency
  remoteBranchSortOrder: alphabetical

- The default selection mode in the staging and custom patch building views has been changed to hunk mode. This is the more useful mode in most cases, as it usually saves a lot of keystrokes. If you want to switch back to the old line mode default, you can do so by adding the following to your config:

gui:
  useHunkModeInStagingView: false
`,
			"0.55.0": `- The 'redo' command, which used to be bound to ctrl-z, is now bound to shift-Z instead. This is because ctrl-z is now used for suspending the application; it is a commonly known keybinding for that in the Linux world. If you want to revert this change, you can do so by adding the following to your config:

keybinding:
  universal:
    suspendApp: <disabled>
    redo: <c-z>

- The 'git.paging.useConfig' option has been removed. If you were relying on it to configure your pager, you'll have to explicitly set the pager again using the 'git.paging.pager' option.
`,
		},

		// AI Common
		AICancel:                            "Cancel",
		AIOK:                                "OK",
		AIConfirm:                           "Confirm",
		AIYes:                               "Yes",
		AINo:                                "No",
		AISuccess:                           "Success",
		AIFailed:                            "Failed",
		AIWarning:                           "Warning",
		AIUnknown:                           "Unknown",
		AIExecuting:                         "Executing",
		AIThinking:                          "Thinking",
		AIIdle:                              "Idle",
		AICancelled:                         "Cancelled",
		AIThinkingInProgress:                "Thinking...",

		// AI Agent
		AIAgentToolNotAllowedInPlanning:     "Tool not allowed in planning phase: %s",
		AIAgentCriticalStepFailed:           "Critical step failed: %s — %s",
		AIAgentStepTimeout:                  "⏱️ Step execution timeout (%v): %s",
		AIAgentUserRejectedTool:             "[User rejected] Tool %s was not executed, please adjust subsequent operations.",
		AIAgentResolveConflictManually:      "Resolve conflict manually and continue",
		AIAgentSetUpstreamBranch:            "Set upstream branch",
		AIAgentConflict:                     "Conflict",
		AIAgentToolName:                     "Tool name",
		AIAgentStageFilesFirst:              "Stage files first (stage_all or stage_file)",
		AIAgentPossibleReasons:              "\n\n💡 Possible reasons:",
		AIAgentExampleCommitMsg:             "feat: add user login feature",
		AIAgentDont:                         "Don't",
		AIAgentRepoStatusAndUserInstruction: "## Current Repository Status\n\n%s\n\n## User Instruction\n\n%s",
		AIAgentUnknownTool:                  "Unknown tool: %s",
		AIAgentUserRejectedExecution:        "User rejected execution: %s",
		AIAgentMaxStepsReached:              "Maximum steps (%d) reached, stopping execution.",
		AIAgentToolLabel:                    "Tool: %s\nDescription: %s\nPermission: %s",
		AIAgentDescriptionLabel:             "Description",
		AIAgentPermissionLabel:              "Permission",
		AIAgentParamsLabel:                  "Parameters",

		// AI Tools
		AIToolMissingParam:          "Missing %s parameter",
		AIToolMissingNameParam:      "Missing name parameter",
		AIToolMissingPathParam:      "Missing path parameter",
		AIToolMissingMessageParam:   "Missing message parameter",
		AIToolMissingHashParam:      "Missing hash parameter",
		AIToolFilePath:              "File path",
		AIToolBranchName:            "Branch name",
		AIToolTagName:               "Tag name",
		AIToolCommitMessage:         "Commit message",
		AIToolNoChanges:             "No changes",
		AIToolWorkingDir:            "Working directory",
		AIToolStagingArea:           "Staging area",
		AIToolTargetRefOrHash:       "Target ref or hash (preferred)",
		AIToolResetSteps:            "Reset steps (used when ref is empty, default 1)",
		AIToolStashIndex:            "Stash index, default 0",
		AIToolMaxLines:              "Maximum lines to return (default 300, 0 for unlimited)",
		AIToolOffset:                "Starting line offset for pagination (0-based, default 0)",
		AIToolTargetRef:             "Target ref (default HEAD)",
		AIToolPushConfigError:       "Push configuration error: %v",
		AIToolRebasedTo:             "Rebased current branch to %s",
		AIToolRenameFailed:          "Rename failed: %v",
		AIToolDiscardChangesFailed:  "Discard changes failed: %v",
		AIToolParam:                 "Parameter",
		AIToolValue:                 "Value",

		// AI Tools - Schema descriptions
		AIToolGetStatusDesc:              "Get current repository status (branch, working tree files, rebase/merge progress)",
		AIToolGetStagedDiffDesc:          "Get the staged diff",
		AIToolGetDiffDesc:                "Get unstaged working tree diff",
		AIToolGetFileDiffDesc:            "Get diff for a specific file (staged or unstaged)",
		AIToolGetFileDiffStagedParam:     "true = staged diff, false = unstaged diff (default false)",
		AIToolGetLogDesc:                 "Get recent commit history",
		AIToolGetLogCountParam:           "Number of entries to return (default 15, max 50)",
		AIToolGetBranchesDesc:            "List local branches (current branch marked with *)",
		AIToolGetStashListDesc:           "List all stash entries",
		AIToolGetRemotesDesc:             "List all configured remotes",
		AIToolGetTagsDesc:                "List all tags",
		AIToolGetStashDiffDesc:           "Show diff of a stash entry (useful for previewing before apply)",
		AIToolGetStashDiffIndexParam:     "Stash index, default 0 (most recent stash)",
		AIToolGetCommitDiffDesc:          "Get diff for a specific commit (default HEAD)",
		AIToolGetCommitDiffHashParam:     "Commit hash; leave empty for HEAD",
		AIToolGetBranchDiffDesc:        "Get diff between two branches or commits. Uses three-dot syntax (A...B): shows changes on target since it diverged from base",
		AIToolGetBranchDiffBaseParam:   "Base ref: branch name, tag, or commit hash (e.g. main, v1.0)",
		AIToolGetBranchDiffTargetParam: "Target ref to diff against base (default HEAD)",
		AIToolGetBranchDiffEmpty:       "No differences between %s and %s",
		AIToolStageAllDesc:               "Stage all working tree changes",
		AIToolStageFileDesc:              "Stage a specific file",
		AIToolUnstageAllDesc:             "Unstage all staged changes (git reset HEAD)",
		AIToolUnstageFileDesc:            "Unstage a specific file",
		AIToolDiscardFileDesc:            "Discard all changes in a file (restore to HEAD)",
		AIToolCommitDesc:                 "Create a commit from staged changes. Call get_staged_diff first; generate the commit message yourself — do not ask the user",
		AIToolCommitMsgParam:             "Conventional Commits message (AI-generated from diff, e.g. feat: add login page)",
		AIToolAmendHeadDesc:              "Rewrite the most recent commit message (git commit --amend)",
		AIToolAmendMsgParam:              "New commit message",
		AIToolRevertCommitDesc:           "Revert a commit by creating a new reverse commit",
		AIToolRevertHashParam:            "Hash of the commit to revert",
		AIToolResetSoftDesc:              "git reset --soft (keep changes in the staging area)",
		AIToolResetMixedDesc:             "git reset --mixed (keep changes in the working tree, unstaged)",
		AIToolResetHardDesc:              "git reset --hard (discard ALL uncommitted changes — irreversible, use with caution)",
		AIToolCherryPickDesc:             "Cherry-pick a specific commit onto the current branch",
		AIToolCherryPickHashParam:        "Hash of the commit to cherry-pick",
		AIToolCheckoutDesc:               "Switch to a branch, tag, or commit hash (hash enters detached HEAD state)",
		AIToolCheckoutNameParam:          "Branch name, tag, or commit hash",
		AIToolCreateBranchDesc:           "Create a new branch. AI generates the name from the user description (kebab-case, e.g. feature/user-login); do not ask the user. checkout=true (default) switches to the new branch",
		AIToolCreateBranchNameParam:      "Branch name in type/description format (kebab-case), e.g. feature/user-login",
		AIToolCreateBranchBaseParam:      "Base ref (default HEAD)",
		AIToolCreateBranchCheckoutParam:  "Switch to the new branch after creation (default true)",
		AIToolDeleteBranchDesc:           "Delete a local branch (must be fully merged unless force=true)",
		AIToolDeleteBranchForceParam:     "Force delete even if unmerged (default false)",
		AIToolRenameBranchDesc:           "Rename a local branch",
		AIToolRenameBranchOldParam:       "Current branch name",
		AIToolMergeBranchDesc:            "Merge a branch into the current branch",
		AIToolMergeBranchNameParam:       "Name of the branch to merge",
		AIToolRebaseBranchDesc:           "Rebase the current branch onto the target branch (git rebase <target>)",
		AIToolRebaseBranchTargetParam:    "Target branch name or ref",
		AIToolStashDesc:                  "Save working tree changes to the stash",
		AIToolStashMsgParam:              "Stash description (optional)",
		AIToolStashPopDesc:               "Restore a stash entry and remove it from the stash list",
		AIToolStashApplyDesc:             "Apply a stash entry without removing it",
		AIToolStashDropDesc:              "Delete a stash entry",
		AIToolCreateTagDesc:              "Create a lightweight tag",
		AIToolDeleteTagDesc:              "Delete a local tag",
		AIToolPullDesc:                   "Pull from remote and merge into the current branch (git pull)",
		AIToolPullRemoteParam:            "Remote name (default: current tracking remote)",
		AIToolPullBranchParam:            "Remote branch name (default: current tracking branch)",
		AIToolFetchDesc:                  "Fetch latest refs from all remotes (git fetch)",
		AIToolPushDesc:                   "Push the current branch to its remote (git push). Use push_force for force-push",
		AIToolPushForceDesc:              "Force-push using --force-with-lease (safer than --force; aborts if remote has unpulled commits). Still rewrites remote history — use with caution",
		AIToolAbortOperationDesc:         "Abort an in-progress rebase, merge, or cherry-pick",
		AIToolAbortOperationTypeParam:    `Operation type: "rebase" | "merge" | "cherry-pick"`,
		AIToolContinueOperationDesc:      "Continue a paused rebase or merge after resolving conflicts",
		AIToolContinueOperationTypeParam: `Operation type: "rebase" | "merge" (default rebase)`,

		// AI Tools - Output messages (success)
		AIToolStagedDiffEmpty:               "No staged changes",
		AIToolUnstagedDiffEmpty:             "No unstaged changes",
		AIToolNoStashEntries:                "No stash entries",
		AIToolNoRemotes:                     "No remotes configured",
		AIToolNoTags:                        "No tags",
		AIToolFileNotInWorkdir:              "File not in working tree: %s",
		AIToolStashEntryEmpty:               "stash[%d] is empty",
		AIToolStatusFiles:                   "Changed files: %d (%d staged, %d unstaged, %d untracked)",
		AIToolStatusClean:                   "Working tree: clean",
		AIToolStatusInProgress:              "In progress: %s",
		AIToolStageAllSuccess:               "Staged all changes",
		AIToolStageFileSuccess:              "Staged: %s",
		AIToolUnstageAllSuccess:             "Unstaged all changes",
		AIToolUnstageFileSuccess:            "Unstaged: %s",
		AIToolCommitSuccess:                 "Committed: \"%s\"",
		AIToolAmendSuccess:                  "Amended commit message to: \"%s\"",
		AIToolRevertSuccess:                 "Reverted commit: %s",
		AIToolResetSoftSuccess:              "reset --soft to %s; changes are in the staging area",
		AIToolResetMixedSuccess:             "reset --mixed to %s; changes are in the working tree",
		AIToolResetHardSuccess:              "reset --hard to %s; all uncommitted changes discarded",
		AIToolCherryPickSuccess:             "Cherry-picked: %s",
		AIToolCheckoutSuccess:               "Switched to: %s",
		AIToolCreateBranchSuccess:           "Created and switched to branch %s (from %s)",
		AIToolCreateBranchNoCheckoutSuccess: "Created branch %s (from %s, not switched)",
		AIToolDeleteBranchSuccess:           "Deleted local branch: %s",
		AIToolRenameBranchSuccess:           "Renamed branch %s to %s",
		AIToolMergeBranchSuccess:            "Merged %s into current branch",
		AIToolStashSuccess:                  "Stashed changes: %s",
		AIToolStashPopSuccess:               "Restored stash[%d]",
		AIToolStashApplySuccess:             "Applied stash[%d] (entry kept)",
		AIToolStashDropSuccess:              "Deleted stash[%d]",
		AIToolCreateTagSuccess:              "Created tag %s (pointing to %s)",
		AIToolDeleteTagSuccess:              "Deleted local tag: %s",
		AIToolPullSuccess:                   "Pull successful",
		AIToolFetchSuccess:                  "Fetch complete",
		AIToolPushSuccess:                   "Push successful",
		AIToolPushForceSuccess:              "Force push (--force-with-lease) successful",
		AIToolAbortSuccess:                  "Aborted %s",
		AIToolContinueSuccess:               "Resumed %s",
		AIToolTruncated:                     "\n... (truncated, %d lines total, showing first %d)",

		// AI Tools - Error messages
		AIToolGetStagedDiffFailed:   "Failed to get staged diff: %v",
		AIToolGetDiffFailed:         "Failed to get diff: %v",
		AIToolGetStashDiffFailed:    "Failed to get stash[%d] diff: %v",
		AIToolGetCommitDiffFailed:   "Failed to get commit diff: %v",
		AIToolGetBranchDiffFailed:   "Failed to get branch diff (%s...%s): %v",
		AIToolStageAllFailed:        "Failed to stage all: %v",
		AIToolStageFileFailed:       "Failed to stage file: %v",
		AIToolUnstageAllFailed:      "Failed to unstage all: %v",
		AIToolUnstageFileFailed:     "Failed to unstage: %v",
		AIToolCommitFailed:          "Commit failed: %v",
		AIToolAmendFailed:           "Failed to amend commit: %v",
		AIToolRevertFailed:          "Revert failed: %v",
		AIToolResetSoftFailed:       "reset --soft failed: %v",
		AIToolResetMixedFailed:      "reset --mixed failed: %v",
		AIToolResetHardFailed:       "reset --hard failed: %v",
		AIToolCherryPickFailed:      "Cherry-pick failed: %v",
		AIToolCheckoutFailed:        "Checkout failed: %v",
		AIToolCreateBranchFailed:    "Failed to create branch: %v",
		AIToolDeleteBranchFailed:    "Failed to delete branch: %v (ensure branch is fully merged or use force=true)",
		AIToolMergeBranchFailed:     "Merge failed: %v",
		AIToolRebaseBranchFailed:    "Rebase failed: %v (resolve conflicts manually then continue)",
		AIToolStashFailed:           "Failed to stash: %v",
		AIToolStashPopFailed:        "stash pop failed: %v",
		AIToolStashApplyFailed:      "stash apply failed: %v",
		AIToolStashDropFailed:       "stash drop failed: %v",
		AIToolCreateTagFailed:       "Failed to create tag: %v",
		AIToolDeleteTagFailed:       "Failed to delete tag: %v",
		AIToolPullFailed:            "Pull failed: %v",
		AIToolFetchFailed:           "Fetch failed: %v",
		AIToolPushFailed:            "Push failed: %v (check remote configuration and authentication)",
		AIToolPushForceFailed:       "Force push failed: %v",
		AIToolAbortFailed:           "Failed to abort %s: %v",
		AIToolContinueFailed:        "Failed to resume %s: %v (ensure all conflicts are resolved and staged)",
		AIToolMissingTargetParam:    "Missing target parameter",
		AIToolMissingOldOrNameParam: "Missing old or name parameter",
		AIToolUnknownOperationType:  `Unknown operation type: %q, supported: %s`,

		// AI Skills
		AISkillCurrentBranch:     "Current branch: %s\n",
		AISkillBranchNameOnly:    "- Output branch name only, no explanation\n",
		AISkillBranchNameFormat:  "- description: lowercase kebab-case, 2-5 words\n",
		AISkillWindowsGitBash:    "Windows + Git Bash, use && to connect commands",
		AISkillRuntime:           "Runtime environment: %s\n\n",
		AISkillOutputJSONArray:   "Output JSON array, each element contains:\n",
		AISkillExplanation:       "- explanation: Chinese explanation (1-2 sentences)\n",
		AISkillCommitSubject:     "- subject: Chinese, verb-first, imperative, max 72 chars\n",
		AISkillTestScenario:      "Scenario hint: This is test-related change, use test type.\n",
		AISkillOutputCommitMsg:   "\nPlease output commit message directly:",
		AISkillRefactorScenario:  "Scenario hint: This is refactoring, prefer refactor type.\n",
		AISkillGeneratePRDesc:    "## Please generate PR description with the following sections\n",
		AISkillPRSummary:         "### Summary\nOne sentence describing the purpose of this PR.\n\n",
		AISkillPRTesting:         "### Testing\n- Explain how to verify these changes\n",
		AISkillCodeChanges:       "## Code Changes\n```diff\n",
		AISkillDiffSummary:       "\nDiff summary:\n```diff\n",
		AISkillRepoContext:       "## Repository Context\n",
		AISkillCodeChangesTitle:  "## Code Changes\n",
		AISkillBranchInfo:        "## Branch Info\nMerging from `%s` to `%s`\n\n",
		AISkillCommitHistory:     "## Commit History\n",

		// AI Chat (GUI)
		AIChatNotEnabled:         "AI not enabled",
		AIChatCanInputNext:       "You can input the next command",
		AIChatGeneratingPlan:     "Analyzing and generating execution plan",
		AIChatTemplateBranchName: "branch-name: branch name",
		AIChatTemplateTagName:    "tag-name: tag name",
		AIChatTemplateMessage:    "message: commit message",
		AIChatPushingToRemote:    "Pushing to remote repository...",
		AIChatAbortMerge:         "Abort merge",
		AIChatResolveConflict:    "Resolve conflict",
		AIChatConflictFiles:      "Conflict files:\n",
		AIChatMergeConflict:      "Merge conflict",
		AIChatUncommittedChanges: "Uncommitted changes detected, how to handle?",
		AIChatDeleteSuccess:      "Delete successful",
		AIChatConfirmSuffix:      "?",

		// AI Repository Context
		AIMoreItems:             "... %d more\n",
		AIRepoWorkingDirClean:   "Working directory: clean\n",
		AIRepoInProgress:        "⚠ In progress: %s\n",
		AIRepoRemoteSynced:      "Remote: %s [synced]\n",
		AIRepoChanges:           "Changes: %d (staged %d, unstaged %d, untracked %d)\n",
		AIRepoRemoteAheadBehind: "Remote: %s [↑%s ↓%s]\n",
		AIRepoBranch:            "Branch: %s\n",
		AIRepoRecentCommits:     "Recent commits:\n",
		AIRepoStashCount:        "Stash: %d entries\n",

		// AI Manager
		AIManagerGenerateBranchName: "Generate appropriate Git branch name based on feature description (kebab-case with type prefix)",
		AIManagerParam:              "Parameter",
		AIManagerValue:              "Value",
		AIManagerStagedDiff:         "Output of git diff --staged",
		AIManagerFeatureDesc:        "Feature or purpose of the branch",
		AIManagerGenerateCommitMsg:  "Generate commit message following Conventional Commits specification based on staged diff",

		// AI Analyze Tool
		AIAnalyzeToolDescription:          "Intelligently analyze current changes: analyze diff file by file and integrate results (suitable for large change scenarios)",
		AIAnalyzeToolStagedParam:          "true=analyze staging area, false=analyze working directory (default false)",
		AIAnalyzeToolFocusParam:           "Analysis focus (e.g.: security issues, performance optimization, code quality, etc.), leave empty for comprehensive analysis",
		AIAnalyzeWorkingDirClean:          "Working directory is clean, no changed files",
		AIAnalyzeNoChanges:                "No changes in %s",
		AIAnalyzeCancelled:                "Analysis cancelled",
		AIAnalyzeFailed:                   "AI analysis failed: %w",
		AIAnalyzeReportTitle:              "# Change Analysis Report\n\n",
		AIAnalyzeReportTitleWithFocus:     "# Change Analysis Report (Focus: %s)\n\n",
		AIAnalyzeFileCount:                "**File count**: %d files (successfully analyzed %d, failed %d)\n",
		AIAnalyzeTotalLines:               "**Total changed lines**: approximately %d lines\n\n",
		AIAnalyzeDetailedAnalysis:         "## Detailed Analysis\n\n",
		AIAnalyzeAnalysisFailed:           "❌ **Analysis failed**: %s\n\n",
		AIAnalyzeNoChangesInfo:            "ℹ️ No changes\n\n",
		AIAnalyzeOverallSuggestions:       "## Overall Suggestions\n\n",
		AIAnalyzeSuggestion1:              "1. Confirm all changes meet expectations\n",
		AIAnalyzeSuggestion2:              "2. Run tests to ensure functionality is normal\n",
		AIAnalyzeSuggestion3:              "3. Check for any missing files\n",
		AIAnalyzeCodeReviewExpert:         "You are a code review expert, skilled at analyzing code changes. Please analyze the diff content concisely and accurately.",
		AIAnalyzeFileLabel:                "## File: %s\n\n",
		AIAnalyzePromptIntro:              "Please analyze the following diff and summarize in 2-3 sentences:\n",
		AIAnalyzeFocusLabel:               "**Analysis focus**: %s\n\n",
		AIAnalyzeMainChanges:              "- Main changes\n",
		AIAnalyzePotentialIssues:          "- Potential issues (if any)\n",
		AIAnalyzeImprovementSuggestions:   "- Improvement suggestions (if any)\n\n",

		// Command Completion
		CompletionBranch:                  "Branch",
		CompletionRemote:                  "Remote",
		CompletionCommitRef:               "Commit reference",
		CompletionTag:                     "Tag",
		CompletionGitDesc:                 "Version control system",
		CompletionCdDesc:                  "Change directory",
		CompletionLsDesc:                  "List files",
		CompletionPwdDesc:                 "Print working directory",
		CompletionCatDesc:                 "Display file content",
		CompletionGrepDesc:                "Search text",
		CompletionFindDesc:                "Find files",
		CompletionGitAddDesc:              "Add files to staging area",
		CompletionGitCommitDesc:           "Commit changes",
		CompletionGitPushDesc:             "Push to remote",
		CompletionGitPullDesc:             "Pull from remote",
		CompletionGitCheckoutDesc:         "Switch branch",
		CompletionGitSwitchDesc:           "Switch branch (new)",
		CompletionGitBranchDesc:           "Manage branches",
		CompletionGitMergeDesc:            "Merge branches",
		CompletionGitRebaseDesc:           "Rebase",
		CompletionGitResetDesc:            "Reset commits",
		CompletionGitRevertDesc:           "Revert commits",
		CompletionGitStashDesc:            "Stash working directory",
		CompletionGitLogDesc:              "View commit history",
		CompletionGitDiffDesc:             "View differences",
		CompletionGitStatusDesc:           "View status",
		CompletionGitTagDesc:              "Manage tags",
		CompletionGitFetchDesc:            "Fetch remote updates",
		CompletionGitCloneDesc:            "Clone repository",
		CompletionGitInitDesc:             "Initialize repository",
		CompletionGitCleanDesc:            "Clean untracked files",
		CompletionGitCherryPickDesc:       "Cherry-pick commits",
		CompletionGitShowDesc:             "Show commit details",
		CompletionGitRmDesc:               "Remove files",
		CompletionGitMvDesc:               "Move files",
		CompletionGitGrepDesc:             "Search content",
		CompletionGitBisectDesc:           "Binary search for problematic commit",
		CompletionFlagAmendDesc:           "Amend last commit",
		CompletionFlagNoEditDesc:          "Don't edit commit message",
		CompletionFlagMDesc:               "Specify commit message",
		CompletionFlagADesc:               "Commit all tracked files",
		CompletionFlagAllDesc:             "Commit all tracked files",
		CompletionFlagFixupDesc:           "Create fixup commit",
		CompletionFlagSignoffDesc:         "Add Signed-off-by line",
		CompletionFlagSDesc:               "GPG sign",
		CompletionFlagNoVerifyDesc:        "Skip pre-commit hooks",
		CompletionFlagAllowEmptyDesc:      "Allow empty commit",
		CompletionFlagForceDesc:           "Force push",
		CompletionFlagForceWithLeaseDesc:  "Safe force push",
		CompletionFlagSetUpstreamDesc:     "Set upstream branch",
		CompletionFlagUDesc:               "Set upstream branch",
		CompletionFlagTagsDesc:            "Push tags",
		CompletionFlagDeleteDesc:          "Delete remote branch",
		CompletionFlagDryRunDesc:          "Preview push",
		CompletionFlagAllBranchesDesc:     "Push all branches",
		CompletionFlagSoftDesc:            "Keep staging area and working directory",
		CompletionFlagMixedDesc:           "Keep working directory",
		CompletionFlagHardDesc:            "Discard all changes",
		CompletionStatusConflicted:        "Conflicted",
		CompletionStatusPartiallyStaged:   "Partially staged",
		CompletionStatusStaged:            "Staged",
		CompletionStatusModified:          "Modified",
		CompletionStatusUntracked:         "Untracked",
		CompletionStatusTracked:           "Tracked",

		// AI Command Helper
		AICommandNotEnabled:               "AI feature is not enabled",
		AICommandGenerationCancelled:      "AI command generation cancelled",
		AICommandInvalidFormat:            "Invalid format returned by AI: %v\nResponse content: %s",
		AICommandExplainPrompt:            "Explain what this shell command does, in concise Chinese:\n\nCommand: %s\n\nPlease explain:\n1. What this command does (1 line)\n2. What impact it will have (1-2 lines)\n3. Whether there are risks (if any, explain the risk points, 1 line)\n4. Suggestions or notes (optional, 1 line)\n\nKeep the answer concise (3-5 lines total).",
		AICommandExplainCancelled:         "Command explanation cancelled",
		AICommandRiskHardReset:            "⚠️ Will lose all uncommitted changes",
		AICommandRiskCleanFdx:             "⚠️ Will delete all untracked and ignored files (including files in .gitignore)",
		AICommandRiskCleanFd:              "⚠️ Will delete all untracked files and directories",
		AICommandRiskForcePush1:           "⚠️ May overwrite remote branch history, affecting other collaborators",
		AICommandRiskForcePush2:           "⚠️ May overwrite remote branch history, affecting other collaborators",
		AICommandRiskReflogExpire:         "⚠️ Will permanently delete reflog records, cannot be recovered",
		AICommandRiskRmRf:                 "⚠️ Dangerous: recursive file deletion, may delete important data",
		AICommandRiskBranchD:              "⚠️ Force delete branch, even if branch is not merged",
		AICommandRiskRebaseI:              "⚠️ Rewrite history, may cause collaboration issues",
		AICommandRiskGcAggressive:         "⚠️ Aggressive garbage collection, may delete recent objects",
		AICommandSuggestionHardReset:      "Suggestion: Use 'git stash' to save changes first, or use 'git reset --soft' to keep changes",
		AICommandSuggestionCleanFdx:       "Suggestion: Use 'git clean -fdn' to preview files to be deleted first",
		AICommandSuggestionCleanFd:        "Suggestion: Use 'git clean -fdn' to preview files to be deleted first",
		AICommandSuggestionForcePush1:     "Suggestion: Use 'git push --force-with-lease' for safe force push",
		AICommandSuggestionForcePush2:     "Suggestion: Use 'git push --force-with-lease' for safe force push",
		AICommandSuggestionBranchD:        "Suggestion: Check if branch is merged first, use 'git branch -d' for safe deletion",

		// AI Skills - Commit Message
		AISkillCommitMsgSystemPrompt:      "You are an experienced software engineer specializing in writing high-quality Git commit messages.\nFollow the Conventional Commits specification (https://www.conventionalcommits.org/).\nOnly output the commit message itself, without any additional explanations, prefixes, or quotes.\nCommit messages (subject and body) must be in Chinese.",
		AISkillCommitMsgRepoBackground:    "## Repository Background\n",
		AISkillCommitMsgCodeChanges:       "## Code Changes\n",
		AISkillCommitMsgOutputRules:       "## Output Rules\n",
		AISkillCommitMsgFormatExample:     "- Format:\n  ```\n  <type>(<scope>): <subject>\n  \n  <body>\n  ```\n",
		AISkillCommitMsgTypeList:          "- type: feat | fix | refactor | docs | test | chore | perf | style | ci | revert\n",
		AISkillCommitMsgSubjectRules:      "- subject: Chinese, verb-first, imperative, max 72 characters\n",
		AISkillCommitMsgScopeOptional:     "- scope is optional\n",
		AISkillCommitMsgBodyRequired:      "- body: Required, leave one blank line between subject and body, explain the reason and main content of this change in Chinese (1-4 lines)\n\n",
		AISkillCommitMsgScenarioBugfix:    "Scenario hint: This is a bug fix, prefer fix type.\n",
		AISkillCommitMsgScenarioRefactor:  "Scenario hint: This is refactoring, prefer refactor type.\n",
		AISkillCommitMsgScenarioDocs:      "Scenario hint: This is documentation update, use docs type.\n",
		AISkillCommitMsgScenarioTest:      "Scenario hint: This is test-related change, use test type.\n",
		AISkillCommitMsgScenarioDefault:   "Please output the commit message directly:",
		AISkillCommitMsgScenarioLarge:     "Scenario hint: This is a large changeset. Focus on the primary intent, mention key modules affected in the body. Be comprehensive but concise. Please output the commit message:\n",
		AISkillCommitMsgProjectType:       "Project type: %s\n",

		// AI Skills - Branch Name
		AISkillBranchNamePromptIntro:      "Recommend a branch name based on the following working directory changes.\n\n",
		AISkillBranchNameStagedFiles:      "Staged files:\n",
		AISkillBranchNameUnstagedFiles:    "Unstaged files:\n",
		AISkillBranchNameMoreFiles:        "  ... %d more\n",
		AISkillBranchNameDiffSummaryTitle: "\nDiff summary:\n```diff\n",
		AISkillBranchNameRules:            "\nNaming rules:\n",
		AISkillBranchNameFormatRule:       "- Format: <type>/<description> (e.g. feature/add-user-auth)\n",
		AISkillBranchNameTypeRule:         "- type: feature | fix | refactor | docs | test | chore\n",
		AISkillBranchNameDescRule:         "- description: lowercase kebab-case, 2-5 words\n",
		AISkillBranchNameOutputRule:       "- Only output the branch name, no explanation\n",
		AISkillBranchNameDescriptionHint:  "Intent description: %s\n\n",
		AISkillBranchNameSystemPrompt:     "You are a Git branch naming expert. Recommend concise, descriptive branch names based on change content. Only output the branch name itself.",

		// AI Skills - PR Description
		AISkillPRDescSystemPrompt:         "You are a senior software engineer responsible for writing clear, professional Pull Request descriptions.\nOutput in Markdown format, including Summary, Changes, and Testing sections.\nUse concise Chinese.",
		AISkillPRDescBranchInfo:           "## Branch Info\nMerging from `%s` to `%s`\n\n",
		AISkillPRDescCommitHistory:        "## Commit History\n",
		AISkillPRDescCodeChangesSection:   "## Code Changes\n```diff\n",
		AISkillPRDescGeneratePrompt:       "## Please generate a PR description with the following sections\n\n",
		AISkillPRDescSummarySection:       "### Summary\nOne sentence describing the purpose and motivation of this PR.\n\n",
		AISkillPRDescChangesSection:       "### Changes\n- List main changes (3-5 items, be specific)\n\n",
		AISkillPRDescBreakingSection:      "### Breaking Changes\nIf there are breaking changes, describe them and the migration steps. Otherwise write: None.\n\n",
		AISkillPRDescTestingSection:       "### Testing\n- Describe how to verify these changes (manual steps or test commands)\n\n",
		AISkillPRDescChecklistSection:     "### Checklist\n- [ ] Tests added or updated\n- [ ] Documentation updated\n- [ ] No breaking changes (or documented above)\n",

		// AI Skills - Shell Command
		AISkillShellCmdSystemPrompt:       "You are a Git command expert. Generate precise shell commands based on user intent.\n\n",
		AISkillShellCmdRuntime:            "Runtime environment: %s\n\n",
		AISkillShellCmdRepoStatus:         "Repository status:\n%s\n\n",
		AISkillShellCmdUserIntent:         "User intent: %s\n\n",
		AISkillShellCmdOutputFormat:       "Output JSON array, each element contains:\n",
		AISkillShellCmdCommandField:       "- command: Complete executable command\n",
		AISkillShellCmdExplanationField:   "- explanation: Chinese explanation (1-2 sentences)\n",
		AISkillShellCmdRiskLevelField:     "- risk_level: \"safe\" | \"medium\" | \"dangerous\"\n",
		AISkillShellCmdAlternativesField:  "- alternatives: Alternative commands (optional)\n\n",
		AISkillShellCmdOutputNote:         "Return 1-3 suggestions, sorted by recommendation. Only output JSON, no other content.",
		AISkillShellCmdWindowsHint:        "Windows + Git Bash, use && to connect commands",
		AISkillShellCmdMacOSHint:          "macOS + zsh/bash, use && to connect commands",
		AISkillShellCmdLinuxHint:          "Linux + bash, use && to connect commands",

		// AI Skills - Code Review
		AISkillCodeReviewSystemPrompt:     "You are a senior software engineer conducting a thorough code review.\nFocus on correctness, security, and maintainability.\nBe conservative: only report issues you are certain exist.\nOutput in Simplified Chinese.",

		// AI Skills - Explain Diff
		AISkillExplainDiffSystemPrompt:    "You are a technical writer who explains code changes clearly.\nDescribe what changed and why it matters — do not judge or assign severity.\nUse Simplified Chinese. Be concise: 2–4 sentences per section.",

		// AI Skills - Release Notes
		AISkillReleaseNotesSystemPrompt:   "You are a release manager writing a professional changelog.\nGroup commits by type, rewrite them as user-facing descriptions.\nUse Simplified Chinese. Output clean Markdown only.",

		// AI Skills - Stash Name
		AISkillStashNameSystemPrompt:      "You are a Git expert generating descriptive stash messages.\nOutput exactly one line: a concise stash message in Simplified Chinese.\nFormat: `<type>: <description>` — no quotes, no extra text.",

		// AI Chat Helper
		AIChatWelcomeSystem:               "Welcome to AI Assistant!",
		AIChatWelcomeMessage:              "Hello! I'm your Git Agent\n\nI can directly help you operate the repository, for example:\n  • \"Help me commit current changes\"\n  • \"Create a feature/login branch\"\n  • \"View recent commit history\"\n  • \"Stash these changes and switch to main branch\"\n\nTell me what you want to do, and I'll execute it.",
		AIChatConfigPrompt:                "Please enable and configure AI features in settings first.\nTip: Press 'o' to open settings menu",
		AIChatPreviousContext:             "─── The following content is from the previous AI analysis, you can continue asking ───",
		AIChatNoContentToCopy:             "No content to copy",
		AIChatNoExecutableReply:           "No executable AI reply",
		AIChatConfirmExecution:            "Confirm execution",
		AIChatExecutionPlan:               "Execution plan",
		AIChatNotInitialized:              "AI not initialized, please configure AI features first.",
		AIChatRequestFailed:               "AI request failed, you can enter the next instruction",
		AIChatCopyFailed:                  "Copy failed",
		AIChatCopiedToClipboard:           "Copied to clipboard",
		AIChatNoCommandsFound:             "No executable commands found in the last AI reply",
		AIChatClearHistoryTitle:           "Confirm clear",
		AIChatClearHistoryPrompt:          "Are you sure you want to clear the conversation history? This action cannot be undone.",
		AIChatHistoryCleared:              "Conversation history cleared",
		AIChatHowCanIHelp:                 "How can I help you?",
		AIChatGenerationStopped:           "Generation stopped, you can enter the next instruction",
		AIChatCompleted:                   "Completed",
		AIChatWaitingConfirm:              "Waiting for confirmation",
		AIChatConfirmPrompt:               "Enter Y to execute, N to cancel, or enter additional instructions to adjust the plan",
		AIChatExecutingPlan:               "Executing plan",
		AIChatGeneratingReply:             "Generating reply",
		AIChatCallingTool:                 "Calling",
		AIChatToolCompleted:               "Completed tool",
		AIChatToolFailed:                  "Tool failed",
		AIChatPlanGenerated:               "Execution plan generated, waiting for confirmation",
		AIChatStatusLabel:                 "Status:",
		AIChatActionLabel:                 "Action:",
		AIChatGreeting:                    "Welcome to AI Assistant!",
		AIChatCapabilities:                "I can directly help you operate the repository",
		AIChatInputPrompt:                 "▶ Enter Y to confirm execution, N to cancel, or enter additional instructions to adjust the plan",
		AIChatStoppedGeneration:           "Generation stopped",
		AIChatCallingToolPrefix:           "Calling",
		AIChatToolCompletedPrefix:         "Completed tool",
		AIChatToolFailedPrefix:            "Tool",

		// AI Two Phase Agent
		AITwoPhaseAgentSystemPromptIntro:        "You are the built-in AI for lazygit, responsible for analyzing user requirements and creating Git operation plans.\n\n",
		AITwoPhaseAgentWorkflowTitle:            "## Workflow\n\n",
		AITwoPhaseAgentWorkflowStep1:            "1. Call read-only tools (get_status, get_diff, etc.) to collect necessary information\n",
		AITwoPhaseAgentWorkflowStep2:            "2. **If commit message generation is needed**:\n",
		AITwoPhaseAgentWorkflowStep2Sub1:        "   - First call get_staged_diff to get staged changes\n",
		AITwoPhaseAgentWorkflowStep2Sub2:        "   - Then call commit_msg tool to generate commit message (returned content is used directly as commit message parameter)\n",
		AITwoPhaseAgentWorkflowStep2Sub3:        "   - **Important**: commit_msg can only be called during planning phase, not in execution plan\n",
		AITwoPhaseAgentWorkflowStep3:            "3. **If branch name generation is needed**:\n",
		AITwoPhaseAgentWorkflowStep3Sub1:        "   - Call branch_name tool to generate branch name\n",
		AITwoPhaseAgentWorkflowStep3Sub2:        "   - **Important**: branch_name can only be called during planning phase, not in execution plan\n",
		AITwoPhaseAgentWorkflowStep4:            "4. After information collection, output a ```plan``` block containing the complete execution plan\n",
		AITwoPhaseAgentWorkflowStep5:            "5. After the ```plan``` block, add a brief natural language explanation, prompting the user to enter Y to confirm, N to cancel, or provide additional instructions\n",
		AITwoPhaseAgentWorkflowStep6:            "6. Strictly prohibit calling any write operation tools during the planning phase\n\n",
		AITwoPhaseAgentToolNameTitle:            "## Important: Tool Name Conventions\n\n",
		AITwoPhaseAgentToolNameIntro:            "**Must use the exact tool names from the tool list below**, do not use git command names:\n",
		AITwoPhaseAgentToolNameStageFile:        "- ✅ Stage files: stage_all (stage all) or stage_file (stage single file)\n",
		AITwoPhaseAgentToolNameDontUseAdd:       "- ❌ Do not use: add, git_add\n",
		AITwoPhaseAgentToolNameCommit:           "- ✅ Commit: commit (parameter message)\n",
		AITwoPhaseAgentToolNameDontUseGitCommit: "- ❌ Do not use: git_commit\n",
		AITwoPhaseAgentToolNameCheckout:         "- ✅ Switch branch: checkout\n",
		AITwoPhaseAgentToolNameDontUseSwitch:    "- ❌ Do not use: switch\n",
		AITwoPhaseAgentToolNameCreateBranch:     "- ✅ Create branch: create_branch\n",
		AITwoPhaseAgentToolNameDontUseBranch:    "- ❌ Do not use: branch\n\n",
		AITwoPhaseAgentSpecialToolTitle:         "## Special Tool Instructions\n\n",
		AITwoPhaseAgentSpecialToolIntro:         "**commit_msg and branch_name are helper tools that can only be called during the planning phase**:\n",
		AITwoPhaseAgentSpecialToolUsage1:        "- Call commit_msg during planning phase to get commit message\n",
		AITwoPhaseAgentSpecialToolUsage2:        "- Use the returned commit message as the message parameter for the commit tool\n",
		AITwoPhaseAgentSpecialToolUsage3:        "- **Do not** put commit_msg in the execution plan steps\n\n",
		AITwoPhaseAgentSpecialToolExample:       "Example:\n```tool\n{\"name\": \"commit_msg\", \"params\": {\"diff\": \"...\"}}\n```\n",
		AITwoPhaseAgentSpecialToolExampleReturn: "Returns: \"feat: add user login functionality\"\n\n",
		AITwoPhaseAgentSpecialToolExamplePlan:   "Then in the execution plan:\n```plan\n{\n  \"steps\": [\n    {\"tool\": \"commit\", \"params\": {\"message\": \"feat: add user login functionality\"}}\n  ]\n}\n```\n\n",
		AITwoPhaseAgentPlanFormatTitle:          "## Plan Format\n\n",
		AITwoPhaseAgentPlanFormatExample:        "```plan\n{\n  \"summary\": \"Overall description (one sentence)\",\n  \"steps\": [\n    {\n      \"id\": \"1\",\n      \"description\": \"Human-readable step description\",\n      \"tool\": \"tool name\",\n      \"params\": {\"param name\": \"specific value\"},\n      \"critical\": true\n    }\n  ]\n}\n```\n\n",
		AITwoPhaseAgentNotesTitle:               "## Notes\n\n",
		AITwoPhaseAgentNotesParam:               "- All step params must be specific values, no placeholders\n",
		AITwoPhaseAgentNotesCriticalTrue:        "- critical=true means abort entire execution if this step fails\n",
		AITwoPhaseAgentNotesCriticalFalse:       "- critical=false means skip and continue if this step fails\n",
		AITwoPhaseAgentNotesMinimal:             "- Only include necessary steps",
		AITwoPhaseAgentExecuting:                "Executing, please wait...",
		AITwoPhaseAgentRepoStatusTitle:          "## Current Repository Status\n\n",
		AITwoPhaseAgentUserInstructionTitle:     "## User Instruction\n\n",
		AITwoPhaseAgentPlanAdjustment:           "Plan adjustment feedback",
		AITwoPhaseAgentExecutionCancelled:       "Execution cancelled",
		AITwoPhaseAgentPlanValidationFailed:     "Plan validation failed",
		AITwoPhaseAgentPlanErrorsIntro:          "❌ The plan contains the following errors, please correct them:\n\n",
		AITwoPhaseAgentPlanRegeneratePrompt:     "\nPlease regenerate a correct execution plan.",
		AITwoPhaseAgentContinueAnalysis:         "Please continue analysis. After collecting enough information, output a ```plan block.",
		AITwoPhaseAgentToolCallWarning:          "⚠️ Warning: Tool %s has been called %d times (with same parameters).\nPlease avoid repeatedly calling the same tool. If enough information has been collected, please output a ```plan block directly.",
		AITwoPhaseAgentSystemPrefix:             "[System] ",
		AITwoPhaseAgentUserFeedbackPrompt:       "User has the following feedback on the above plan, please adjust the plan according to the feedback and output a new ```plan block:\n\n%s",
		AITwoPhaseAgentToolResultPrefix:         "[Tool result %s]\n%s",
		AITwoPhaseAgentMaxStepsExceeded:         "Planning phase exceeded maximum steps (%d), failed to generate execution plan",
		AITwoPhaseAgentEmptyResponseError:       "AI produced no valid output (no tool calls or execution plan) 3 times in a row. Please rephrase your task or check AI configuration",
	}
}
