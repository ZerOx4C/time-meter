package ui

import (
	"syscall"
	"time"
	"time-meter/logic"
	"time-meter/setting"
	"time-meter/textmap"
	winapi2 "time-meter/winapi"
	"time-meter/wrapped"

	"github.com/cwchiu/go-winapi"
)

type Controller interface {
	SetTextMap(textMap textmap.TextMap)
	SetSettings(settings *setting.Settings)
	SetTasks(tasks []logic.Task)
	SetErrorMessage(message string)
	OnPopupMenuCommand(handler PopupMenuCommandHandler)
	ShowErrorMessageBox(message string)
	Initialize() error
	Run()
	Quit()
	Finalize() error
}

type PopupMenuCommandHandler func(menuId MenuId)

type controller struct {
	textMap                 textmap.TextMap
	settings                *setting.Settings
	tasks                   []logic.Task
	popupMenuCommandHandler PopupMenuCommandHandler
	meterWindow             *MeterWindow
	tipWindow               *TipWindow
	meterRenderer           *MeterRenderer
	tipRenderer             *TipRenderer
	contextMenu             *PopupMenu
}

type MenuId int16

const (
	MID_ZERO MenuId = iota
	MID_EDIT_SCHEDULE
	MID_QUIT
)

func NewController() Controller {
	ret := new(controller)
	ret.meterWindow = new(MeterWindow)
	ret.tipWindow = new(TipWindow)
	ret.meterRenderer = new(MeterRenderer)
	ret.tipRenderer = new(TipRenderer)
	ret.contextMenu = new(PopupMenu)
	return ret
}

func (c *controller) SetTextMap(textMap textmap.TextMap) {
	c.textMap = textMap
	c.tipRenderer.textMap = textMap
}

func (c *controller) SetSettings(settings *setting.Settings) {
	c.settings = settings
	c.meterWindow.settings = settings
	c.tipWindow.settings = settings
	c.meterRenderer.settings = settings
	c.tipRenderer.settings = settings
}

func (c *controller) SetTasks(tasks []logic.Task) {
	c.tasks = []logic.Task{}
	c.tasks = append(c.tasks, tasks...)

	c.meterRenderer.tasks = c.tasks
}

func (c *controller) SetErrorMessage(message string) {
	c.tipRenderer.errorMessage = message
}

func (c *controller) OnPopupMenuCommand(handler PopupMenuCommandHandler) {
	c.popupMenuCommandHandler = handler
}

func (c *controller) ShowErrorMessageBox(message string) {
	caption := c.textMap.Of("NOUN_TIME_METER").String()
	captionPtr, _ := syscall.UTF16PtrFromString(caption)
	messagePtr, _ := syscall.UTF16PtrFromString(message)
	winapi.MessageBox(c.meterWindow.hWnd, messagePtr, captionPtr, winapi.MB_ICONERROR|winapi2.MB_TOPMOST)
}

func (c *controller) Initialize() error {
	if err := c.meterWindow.Initialize(); err != nil {
		return err
	}

	if err := c.tipWindow.Initialize(); err != nil {
		return err
	}

	if err := c.meterRenderer.Initialize(); err != nil {
		return err
	}

	if err := c.tipRenderer.Initialize(); err != nil {
		return err
	}

	if err := c.contextMenu.Initialize(); err != nil {
		return err
	}

	c.contextMenu.AppendStringItem(MID_EDIT_SCHEDULE, c.textMap.Of("VERB_EDIT_SCHEDULE").String())
	c.contextMenu.AppendStringItem(MID_QUIT, c.textMap.Of("VERB_QUIT").String())

	c.meterWindow.onPaint = func() {
		c.meterRenderer.width = c.meterWindow.bound.Width()
		c.meterRenderer.height = c.meterWindow.bound.Height()
		c.meterRenderer.Draw(c.meterWindow.hWnd)
	}

	c.meterWindow.onMouseMove = func() {
		var cursorPos wrapped.POINT
		winapi.GetCursorPos(cursorPos.Unwrap())

		focusRatio := 1 - float64(cursorPos.Y-c.meterWindow.bound.Top)/float64(c.meterWindow.bound.Height())
		totalDuration := c.settings.FutureDuration + c.settings.PastDuration
		focusAt := time.Now().Add(-c.settings.PastDuration + time.Duration(focusRatio*float64(totalDuration)))

		var focusTasks []logic.Task
		for _, task := range c.tasks {
			if task.OverlapWith(focusAt, focusAt) {
				focusTasks = append(focusTasks, task)
			}
		}

		c.tipRenderer.tasks = focusTasks

		if c.tipRenderer.errorMessage != "" {
			// NOTE: workaround.
			c.tipWindow.Show()

		} else if 0 < len(focusTasks) {
			c.tipWindow.Show()

		} else {
			c.tipWindow.Hide()
		}

		c.tipWindow.boundLeft = c.meterWindow.bound.Right
		c.tipWindow.Update()
	}

	c.meterWindow.onMouseEnter = func() {
		c.tipWindow.Show()
	}

	c.meterWindow.onMouseLeave = func() {
		c.tipWindow.Hide()
	}

	c.meterWindow.onMouseRightClick = func() {
		c.contextMenu.Popup(c.meterWindow.hWnd)
	}

	c.meterWindow.onPopupMenuCommand = func() {
		if c.popupMenuCommandHandler != nil {
			c.popupMenuCommandHandler(c.meterWindow.lastMenuId)
		}
	}

	c.tipWindow.onPaint = func() {
		c.tipRenderer.Draw(c.tipWindow.hWnd)
	}

	return nil
}

func (c *controller) Run() {
	c.meterWindow.Show()

	var msg winapi.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}
}

func (c *controller) Quit() {
	winapi.SendMessage(c.meterWindow.hWnd, winapi.WM_CLOSE, 0, 0)
}

func (c *controller) Finalize() error {
	c.contextMenu.Finalize()
	c.tipRenderer.Finalize()
	c.meterRenderer.Finalize()
	c.tipWindow.Finalize()
	c.meterWindow.Finalize()

	return nil
}
