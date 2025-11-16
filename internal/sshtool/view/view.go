package sshtool

import (
	"netools/internal/db"
	"strconv"

	"github.com/rivo/tview"
)

const (
	pageHome         = "home"
	pageRegister     = "register"
	pageServerList   = "server-list"
	pageServerDetail = "server-detail"
	pageSession      = "session"
	pageTunnel       = "tunneling"

	keyRegister = "register"
)

type Profile struct {
	ID       int
	Name     string
	Host     string
	User     string
	Password string
	Port     int
	Desc     string
}

type View struct {
	db *db.DB

	app   *tview.Application
	Root  *tview.Flex
	pages *tview.Pages

	homeList     *tview.List
	registerForm *tview.Form
	serverList   *tview.List
	sessionView  *tview.TextView
	tunnelView   *tview.TextView
	status       *tview.TextView

	profiles  []Profile
	nextID    int
	currentID int

	onExit          func()
	onStartSession  func(p Profile) error
	onOpenTunneling func(p Profile) error
}

func NewView(app *tview.Application, db *db.DB) *View {
	v := &View{
		app:      app,
		pages:    tview.NewPages(),
		profiles: make([]Profile, 0),
		nextID:   1,
	}

	v.homeList = tview.NewList()
	v.homeList.ShowSecondaryText(true)
	v.homeList.SetBorder(true)
	v.homeList.SetTitle(" SSH Profiles ")

	v.registerForm = tview.NewForm()
	v.registerForm.SetBorder(true)

	v.serverList = tview.NewList()
	v.serverList.SetBorder(true)
	v.serverList.SetTitle(" Server Menu ")

	v.tunnelView = tview.NewTextView()
	v.tunnelView.SetDynamicColors(true)
	v.tunnelView.SetBorder(true)
	v.tunnelView.SetTitle(" Tunneling ")

	v.status = tview.NewTextView()
	v.status.SetDynamicColors(true)

	v.pages.AddPage(pageHome, v.homeList, true, true)
	v.pages.AddPage(pageRegister, v.registerForm, true, false)
	v.pages.AddPage(pageServerList, v.serverList, true, false)
	v.pages.AddPage(pageSession, v.sessionView, true, false)
	v.pages.AddPage(pageTunnel, v.tunnelView, true, false)

	v.Root = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.pages, 0, 1, true).
		AddItem(v.status, 1, 0, false)

	return v
}

func (v *View) refreshHomeList() {
	if v.homeList == nil {
		return
	}

	v.homeList.Clear()

	v.homeList.AddItem("[green::b]server register", "Add new SSH server", 'r', func() {
		v.currentID = 0
	})

}

func (v *View) showRegister(p *Profile) {
	if v.registerForm == nil {
		return
	}
	v.registerForm.Clear(true)
	v.registerForm.SetTitle(" " + keyRegister + " profile")

	var (
		name     string
		host     string
		user     string
		password string
		portStr  = "22"
		desc     string
	)

	if p != nil {
		name = p.Name
		host = p.Host
		user = p.User
		password = p.Password
		if p.Port != 0 {
			portStr = strconv.Itoa(p.Port)
		}
		desc = p.Desc
	}

	v.registerForm.AddInputField("Name", name, 20, nil, func(text string) {
		name = text
	})
	v.registerForm.AddInputField("Host", host, 20, nil, func(text string) {
		host = text
	})
	v.registerForm.AddInputField("User", user, 20, nil, func(text string) {
		user = text
	})
	v.registerForm.AddPasswordField("Password", password, 20, '*', func(text string) {
		password = text
	})
	v.registerForm.AddInputField("Port", portStr, 6, nil, func(text string) {
		portStr = text
	})
	v.registerForm.AddInputField("Desc", desc, 30, nil, func(text string) {
		desc = text
	})

	v.registerForm.AddButton("Save", func() {
		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 {
			port = 22
		}

		if p != nil && v.currentID != 0 {

		}

	})

}
