package sshtool

import "github.com/rivo/tview"

const (
	addConfig = "add-config"
)

type SSHView struct {
	Name         string
	Root         *tview.Flex
	Mid          *tview.List
	Inner        *tview.Pages
	addForm      *tview.Form
	profilePages map[string]tview.Primitive
	profileTabs  map[string]*tview.List
	app          *tview.Application
}

func NewView(app *tview.Application) *SSHView {
	v := &SSHView{
		Name:         "ssh",
		Mid:          tview.NewList(),
		Inner:        tview.NewPages(),
		profilePages: make(map[string]tview.Primitive),
		profileTabs:  make(map[string]*tview.List),
		app:          app,
	}

	host := tview.NewInputField().SetLabel("Host: ").SetFieldWidth(24)
	user := tview.NewInputField().SetLabel("User: ").SetFieldWidth(16)
	password := tview.NewInputField().SetLabel("Password: ").SetMaskCharacter('*').SetFieldWidth(24)
	port := tview.NewInputField().SetLabel("Port: ").SetFieldWidth(6).
		SetAcceptanceFunc(func(s string, r rune) bool {
			return r >= '0' && r <= '9'
		})

	v.addForm = tview.NewForm().
		AddFormItem(host).
		AddFormItem(user).
		AddFormItem(password).
		AddFormItem(port).
		AddButton("Save", func() {
			// 저장 로직
		}).
		AddButton("Cancel", func() {
			if v.app != nil {
				v.app.SetFocus(v.Mid)
			}

		})
	v.addForm.SetBorder(true).SetTitle(addConfig)

	v.Inner.AddPage(addConfig, v.addForm, true, true)

	v.Mid.
		AddItem(addConfig, "server config register", '+', func() {
			v.Inner.SwitchToPage(addConfig)
			if v.app != nil {
				v.app.SetFocus(v.addForm)
			}
		}).
		SetBorder(true).
		SetTitle("ssh")

	// v.Root = tview.NewFlex().
	// 	AddItem(v.Mid, 0, 0, true).
	// 	AddItem(v.Inner, 0, 0, true)

	return v
}

func (v *SSHView) Primitive() tview.Primitive {
	return v.Root
}
func (v *SSHView) FocusDefault() tview.Primitive {
	return v.Mid
}
