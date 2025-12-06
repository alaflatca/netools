package sshtool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"netools/internal/db"
	sshtool "netools/internal/sshtool/api"
	"strconv"

	"github.com/rivo/tview"
)

const (
	pageHome        = "home"
	pageRegister    = "register"
	pageProfileMenu = "profile-menu"
	pageSession     = "session"
	pageTunnel      = "tunneling"

	keyRegister = "register"
)

type textViewWriter struct {
	app  *tview.Application
	view *tview.TextView
}

func (w *textViewWriter) Write(p []byte) (int, error) {
	var (
		n   int
		err error
	)

	w.app.QueueUpdateDraw(func() {
		aw := tview.ANSIWriter(w.view)
		n, err = aw.Write(p)
		w.view.ScrollToEnd()
	})

	return n, err
}

type Profile struct {
	ID       int
	Name     string
	Host     string
	User     string
	Password string
	Port     int
	KeyPath  string
	Desc     string
}

func (p Profile) Addr() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

func (p Profile) AuthMethod() string {
	switch {
	case p.KeyPath != "":
		return "key"
	case p.Password != "":
		return "password"
	default:
		return "none"
	}
}

type View struct {
	db *db.DB

	app *tview.Application

	Root   *tview.Flex
	header *tview.TextView // 상단 타이틀
	body   tview.Primitive // 가운데 화면 (homeList / registerForm / ...)
	status *tview.TextView // 하단 status bar

	homeList     *tview.List
	registerForm *tview.Form
	profileMenu  *tview.List
	// sessionView  *tview.TextView
	sessionTerm *TerminalView
	tunnelView  *tview.TextView

	profiles  []Profile
	nextID    int
	currentID int

	sessionStdin  io.WriteCloser
	sessionCancel context.CancelFunc

	onExit          func()
	onStartSession  func(p Profile) error
	onOpenTunneling func(p Profile) error
}

func NewView(app *tview.Application, dbConn *db.DB) *View {
	v := &View{
		app:      app,
		db:       dbConn,
		profiles: make([]Profile, 0),
		nextID:   1,
	}

	// ----- header -----
	v.header = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	v.header.SetBorder(true)

	// ----- homeList -----
	v.homeList = tview.NewList()
	v.homeList.ShowSecondaryText(true)
	v.homeList.SetBorder(false) // 테두리는 header에서만
	v.homeList.SetWrapAround(true)

	// ----- registerForm -----
	v.registerForm = tview.NewForm()
	v.registerForm.SetBorder(false)

	// ----- profileMenu -----
	v.profileMenu = tview.NewList()
	v.profileMenu.SetBorder(false)

	// ----- sessionTerm -----
	v.sessionTerm = NewTerminalView(app)

	// ----- tunnelView -----
	v.tunnelView = tview.NewTextView()
	v.tunnelView.SetDynamicColors(true)
	v.tunnelView.SetBorder(false)

	// ----- status -----
	v.status = tview.NewTextView()
	v.status.SetDynamicColors(true)
	v.status.SetBorder(true)

	// 초기 body 는 homeList
	v.body = v.homeList

	// Root: [header] [body] [status]
	v.Root = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.header, 3, 0, false). // 타이틀 박스 높이 대충 3
		AddItem(v.body, 0, 1, true).    // 나머지 대부분
		AddItem(v.status, 3, 0, false)  // status 한 줄
	// DB에서 프로필 불러오기
	v.loadProfilesFromDB()

	v.onStartSession = func(p Profile) error {
		// 1) 화면 쪽은 바로 반응
		v.ClearSession()
		v.AppendSessionLine(
			fmt.Sprintf("connecting to %s@%s:%d (%s auth)...",
				p.User, p.Host, p.Port, p.AuthMethod()),
		)

		// stdin 파이프 준비 (키 입력 → sessionStdin.Write → pr → SSH)
		pr, pw := io.Pipe()
		v.sessionStdin = pw

		ctx, cancel := context.WithCancel(context.Background())
		v.sessionCancel = cancel

		outWriter := v.sessionTerm.StdoutWriter()

		// 2) SSH 연결 + Session 은 별도 고루틴에서 실행
		go func() {
			defer func() {
				pr.Close()
				pw.Close()
				v.sessionStdin = nil
				v.sessionCancel = nil
			}()

			// (1) config 생성
			cfg, err := sshtool.CreateSshConfig(p.User, p.Password, p.KeyPath)
			if err != nil {
				v.app.QueueUpdateDraw(func() {
					v.AppendSessionLine(fmt.Sprintf("\n[red]config error: %v", err))
				})
				return
			}

			// (2) SSH 연결
			client, err := sshtool.NewSSHClient(sshtool.SessionArgs{
				Network:   "tcp",
				Host:      p.Host,
				Port:      p.Port,
				ClientCfg: cfg,
			})
			if err != nil {
				v.app.QueueUpdateDraw(func() {
					v.AppendSessionLine(fmt.Sprintf("\n[red]connect error: %v", err))
				})
				return
			}

			// (3) 세션 실행
			err = sshtool.Session(ctx, client, pr, outWriter)

			v.app.QueueUpdateDraw(func() {
				if err != nil && !errors.Is(err, io.EOF) {
					v.AppendSessionLine(fmt.Sprintf("\n[red]session error: %v", err))
				} else {
					v.AppendSessionLine("\n[green]session closed")
				}
			})
		}()

		// 3) UI는 바로 세션 화면으로 전환하고 반환
		v.switchTo(pageSession)
		return nil
	}

	// 상태 & 초기 화면
	v.setStatus("SSH tool ready. Add a profile or select an existing one.")
	v.switchTo(pageHome)

	return v
}

func (v *View) Primitive() tview.Primitive {
	return v.Root
}

/* =============== 콜백 설정 =============== */

func (v *View) SetOnExit(fn func()) {
	v.onExit = fn
}
func (v *View) SetOnStartSession(fn func(Profile) error) {
	v.onStartSession = fn
}
func (v *View) SetOnOpenTunneling(fn func(Profile) error) {
	v.onOpenTunneling = fn
}

/* =============== 프로필 세팅 =============== */

func (v *View) SetProfiles(ps []Profile) {
	v.profiles = make([]Profile, len(ps))
	copy(v.profiles, ps)

	maxID := 0
	for _, p := range v.profiles {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	v.nextID = maxID + 1
	v.refreshHomeList()
}

func (v *View) Profiles() []Profile {
	out := make([]Profile, len(v.profiles))
	copy(out, v.profiles)
	return out
}

/* =============== 화면 전환 =============== */

func (v *View) switchTo(name string) {
	var (
		title string
		body  tview.Primitive
	)

	switch name {
	case pageHome:
		title = "[white::b] SSH Profiles "
		body = v.homeList
	case pageRegister:
		title = "[white::b] " + keyRegister + " profile "
		body = v.registerForm
	case pageProfileMenu:
		title = "[white::b] Profile menu "
		body = v.profileMenu
	case pageSession:
		title = "[white::b] Session "
		body = v.sessionTerm
	case pageTunnel:
		title = "[white::b] Tunneling "
		body = v.tunnelView
	default:
		return
	}

	v.header.SetText(title)

	if v.body != body {
		v.body = body
		// Root 레이아웃을 다시 구성
		v.Root.Clear()
		v.Root.
			SetDirection(tview.FlexRow).
			AddItem(v.header, 3, 0, false).
			AddItem(v.body, 0, 1, true).
			AddItem(v.status, 3, 0, false)
	}

	v.app.SetFocus(body)
}

func (v *View) setStatus(msg string) {
	if v.status == nil {
		return
	}
	v.status.SetText(msg)
}

/* =============== 프로필 삭제 =============== */

func (v *View) deleteProfileByID(id int) {
	for i, p := range v.profiles {
		if p.ID == id {
			v.profiles = append(v.profiles[:i], v.profiles[i+1:]...)
			break
		}
	}

	if v.db != nil {
		if err := db.DeleteSSHConfig(context.Background(), v.db, int64(id)); err != nil {
			v.setStatus(fmt.Sprintf("[red]failed to delete from db: %v", err))
		}
	}
}

/* =============== 홈 리스트 =============== */

func (v *View) refreshHomeList() {
	if v.homeList == nil {
		return
	}
	v.homeList.Clear()

	// server register
	v.homeList.AddItem("[green::b]server register", "Add new SSH server", 'r', func() {
		v.currentID = 0
		v.showRegister(nil)
	})

	for i := range v.profiles {
		p := v.profiles[i]
		main := p.Name
		if main == "" {
			main = fmt.Sprintf("%s@%s", p.User, p.Host)
		}
		secondary := fmt.Sprintf("%s@%s:%d", p.User, p.Host, p.Port)

		pCopy := p
		v.homeList.AddItem(main, secondary, 0, func() {
			v.currentID = pCopy.ID
			v.showProfileMenu(pCopy)
		})
	}
}

/* =============== 프로필 메뉴 =============== */

func (v *View) showProfileMenu(p Profile) {
	if v.profileMenu == nil {
		return
	}
	v.profileMenu.Clear()

	v.profileMenu.SetTitle(fmt.Sprintf(" %s ", p.Name))

	v.profileMenu.AddItem("session", "Open SSH session", 's', func() {
		v.setStatus(fmt.Sprintf("Starting session for %s...", p.Name))

		if v.onStartSession != nil {
			if err := v.onStartSession(p); err != nil {
				v.setStatus(fmt.Sprintf("[red]Session error: %v", err))
				return
			}
		}
		v.showSession(p)
	})

	v.profileMenu.AddItem("tunneling", "Open tunneling", 't', func() {
		v.setStatus(fmt.Sprintf("Opening tunneling for %s...", p.Name))
		if v.onOpenTunneling != nil {
			if err := v.onOpenTunneling(p); err != nil {
				v.setStatus(fmt.Sprintf("[red]Tunneling error: %v", err))
				return
			}
		}
		v.showTunnel(p)
	})

	v.profileMenu.AddItem("edit profile", "Modify this profile", 'e', func() {
		v.currentID = p.ID
		v.showRegister(&p)
	})

	v.profileMenu.AddItem("delete profile", "Remove this profile", 'd', func() {
		v.deleteProfileByID(p.ID)
		v.refreshHomeList()
		v.setStatus(fmt.Sprintf("Deleted profile %s", p.Name))
		v.switchTo(pageHome)
	})

	v.profileMenu.AddItem("back", "Back to server list", 'b', func() {
		v.switchTo(pageHome)
	})

	v.switchTo(pageProfileMenu)
}

/* =============== 프로필 등록/수정 =============== */

func (v *View) showRegister(p *Profile) {
	if v.registerForm == nil {
		return
	}
	v.registerForm.Clear(true)

	var (
		name     string
		host     string
		user     string
		password string
		portStr  = "22"
		keyPath  string
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
		keyPath = p.KeyPath
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
	v.registerForm.AddInputField("KeyPath", keyPath, 50, nil, func(text string) {
		keyPath = text
	})
	v.registerForm.AddInputField("Desc", desc, 30, nil, func(text string) {
		desc = text
	})

	v.registerForm.AddButton("Save", func() {
		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 {
			port = 22
		}
		portStrDB := strconv.Itoa(port)

		if p != nil && v.currentID != 0 {
			// 수정
			for i := range v.profiles {
				if v.profiles[i].ID == v.currentID {
					v.profiles[i].Name = name
					v.profiles[i].Host = host
					v.profiles[i].User = user
					v.profiles[i].Password = password
					v.profiles[i].Port = port
					v.profiles[i].KeyPath = keyPath
					v.profiles[i].Desc = desc
					break
				}
			}

			if v.db != nil {
				cfg := db.SSHConfig{
					ID:       int64(v.currentID),
					Name:     name,
					IP:       host,
					User:     user,
					Password: password,
					Port:     portStrDB,
					KeyPath:  keyPath,
					Desc:     desc,
				}
				if err := db.UpdateSSHConfig(context.Background(), v.db, cfg); err != nil {
					v.setStatus(fmt.Sprintf("[red]DB update error: %v", err))
				} else {
					v.setStatus(fmt.Sprintf("Updated profile %s", name))
				}
			} else {
				v.setStatus(fmt.Sprintf("Updated profile %s", name))
			}
		} else {
			// 신규
			if v.db != nil {
				cfg := db.SSHConfig{
					Name:     name,
					IP:       host,
					User:     user,
					Password: password,
					Port:     portStrDB,
					KeyPath:  keyPath,
					Desc:     desc,
				}
				newID, err := db.InsertSSHConfig(context.Background(), v.db, cfg)
				if err != nil {
					v.setStatus(fmt.Sprintf("[red]DB insert error: %v", err))
					return
				}

				np := Profile{
					ID:       int(newID),
					Name:     name,
					Host:     host,
					User:     user,
					Password: password,
					Port:     port,
					KeyPath:  keyPath,
					Desc:     desc,
				}
				v.profiles = append(v.profiles, np)
				v.setStatus(fmt.Sprintf("Added profile %s", name))
			} else {
				id := v.nextID
				v.nextID++

				np := Profile{
					ID:       id,
					Name:     name,
					Host:     host,
					User:     user,
					Password: password,
					Port:     port,
					KeyPath:  keyPath,
					Desc:     desc,
				}
				v.profiles = append(v.profiles, np)
				v.setStatus(fmt.Sprintf("Added profile %s (in-memory only)", name))
			}
		}

		v.refreshHomeList()
		v.switchTo(pageHome)
	})

	v.registerForm.AddButton("Cancel", func() {
		v.switchTo(pageHome)
	})

	v.switchTo(pageRegister)
}

/* =============== DB 로딩 =============== */

func (v *View) loadProfilesFromDB() {
	v.profiles = v.profiles[:0]
	v.nextID = 1

	if v.db != nil {
		cfgs, err := db.SelectSSHConfigs(context.Background(), v.db)
		if err != nil {
			v.setStatus(fmt.Sprintf("[red]failed to load ssh configs: %v", err))
		} else {
			for _, c := range cfgs {
				port := 22
				if c.Port != "" {
					if p, err := strconv.Atoi(c.Port); err == nil && p > 0 {
						port = p
					}
				}

				p := Profile{
					ID:       int(c.ID),
					Name:     c.Name,
					Host:     c.IP,
					User:     c.User,
					Password: c.Password,
					Port:     port,
					KeyPath:  c.KeyPath,
					Desc:     c.Desc,
				}
				v.profiles = append(v.profiles, p)
				if p.ID >= v.nextID {
					v.nextID = p.ID + 1
				}
			}
		}
	}

	v.refreshHomeList()
}

/* =============== 세션 / 터널 =============== */

func (v *View) showSession(p Profile) {
	v.sessionTerm.ClearScreen()
	v.sessionTerm.PrintLine(
		fmt.Sprintf("Session for %s (%s@%s:%d, %s auth)",
			p.Name, p.User, p.Host, p.Port, p.AuthMethod()),
	)
	v.switchTo(pageSession)
}

func (v *View) showTunnel(p Profile) {
	v.tunnelView.Clear()
	fmt.Fprintf(v.tunnelView, "Tunneling for [yellow]%s[white]\n", p.Name)
	fmt.Fprintf(v.tunnelView, "%s@%s:%d\n", p.User, p.Host, p.Port)
	if p.Desc != "" {
		fmt.Fprintf(v.tunnelView, "\n%s\n", p.Desc)
	}
	v.switchTo(pageTunnel)
}
