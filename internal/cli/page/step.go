package page

type Step int

const (
	StepLoginUser Step = iota
	StepLoginPass
	StepRootMenu
	StepNotesMenu
	StepPasswordsMenu
	StepRemindersMenu
	StepListLoading
	StepListView
	StepNoteDetailLoading
	StepNoteDetailView
	StepAddTitle
	StepAddText
	StepAddTags
	StepUpdateID
	StepUpdateTitle
	StepUpdateText
	StepUpdateTags
	StepDeleteID
	StepInfo
)
