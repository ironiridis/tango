package rpc

// UIDefineCategory is sent by a process asking to create a category container in the UI
type UIDefineCategory struct {
	Label    string
	CSSClass string

	// Parent may be the zero UUID to indicate the root category
	Parent UUID
}

func (z *UIDefineCategory) T() string {
	return "UIDefineCategory"
}

// UICategoryCreated is delivered to a process when a category it requested is created
type UICategoryCreated struct {
	UUID UUID
}

func (z *UICategoryCreated) T() string {
	return "UICategoryCreated"
}

// UIRemoveCategory is sent by a process asking to remove a category from the UI
type UIRemoveCategory struct {
	UUID UUID
}

func (z *UIRemoveCategory) T() string {
	return "UIRemoveCategory"
}

// UIUserAction is delivered to a process when an element it created calls tango.UserAction, usually in response to, say, a clicked button
type UIUserAction struct {
	UUID   UUID
	Params []string
}

func (z *UIUserAction) T() string {
	return "UIUserAction"
}
