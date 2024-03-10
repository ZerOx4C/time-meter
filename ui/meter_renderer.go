package ui

import (
	"time"
	"time-meter/logic"
	"time-meter/setting"
	"time-meter/wrapped"

	"github.com/cwchiu/go-winapi"
)

type MeterRenderer struct {
	settings        *setting.Settings
	tasks           []logic.Task
	width           int32
	height          int32
	backgroundBrush winapi.HBRUSH
	headPen         winapi.HPEN
	hourPen         winapi.HPEN
	chartBrush      winapi.HBRUSH
}

func (mr *MeterRenderer) Initialize() error {
	mr.backgroundBrush = winapi.CreateSolidBrush(mr.settings.BackgroundColor)
	mr.headPen = winapi.CreatePen(winapi.PS_SOLID, 1, mr.settings.MainScaleColor)
	mr.hourPen = winapi.CreatePen(winapi.PS_SOLID, 1, mr.settings.SubScalesColor)
	mr.chartBrush = winapi.CreateSolidBrush(mr.settings.ChartColor)

	return nil
}

func (mr *MeterRenderer) Finalize() error {
	winapi.DeleteObject(winapi.HGDIOBJ(mr.backgroundBrush))
	winapi.DeleteObject(winapi.HGDIOBJ(mr.headPen))
	winapi.DeleteObject(winapi.HGDIOBJ(mr.hourPen))
	winapi.DeleteObject(winapi.HGDIOBJ(mr.chartBrush))

	return nil
}

func (mr *MeterRenderer) Draw(hWnd winapi.HWND) {
	var paint winapi.PAINTSTRUCT
	hdc := winapi.BeginPaint(hWnd, &paint)

	backBuffer := new(BackBuffer)
	backDc := backBuffer.begin(hWnd, hdc)

	var clientRect wrapped.RECT
	winapi.GetClientRect(hWnd, clientRect.Unwrap())
	winapi.FillRect(backDc, clientRect.Unwrap(), mr.backgroundBrush)

	mr.drawAllCharts(backDc,
		mr.tasks,
		time.Now(),
		mr.settings.FutureDuration,
		mr.settings.PastDuration,
	)

	mr.drawAllScaleLines(backDc,
		mr.settings.FutureDuration,
		mr.settings.PastDuration,
		mr.settings.ScaleInterval,
	)

	backBuffer.end()

	winapi.EndPaint(hWnd, &paint)
}

func (mr *MeterRenderer) drawAllCharts(hdc winapi.HDC, tasks []logic.Task, now time.Time, futureDuration, pastDuration time.Duration) {
	chartBeginAt := now.Add(-pastDuration)
	chartEndAt := now.Add(futureDuration)
	totalSeconds := int32((futureDuration + pastDuration) / time.Second)

	tracks := [][]logic.Task{}

	for _, task := range tasks {
		if !task.OverlapWith(chartBeginAt, chartEndAt) {
			continue
		}

		found := false

		for index := range tracks {
			if !mr.isTaskConflict(tracks[index], task) {
				tracks[index] = append(tracks[index], task)
				found = true
				break
			}
		}

		if !found {
			tracks = append(tracks, []logic.Task{task})
		}
	}

	if len(tracks) == 0 {
		return
	}

	trackWidth := int(mr.width) / len(tracks)

	for trackIndex, track := range tracks {
		for _, task := range track {
			var rect wrapped.RECT
			rect.Left = int32(trackIndex*trackWidth) + 1
			rect.Right = rect.Left + int32(trackWidth) - 2
			rect.Top = mr.height - mr.height*int32(task.EndAt.Sub(chartBeginAt)/time.Second)/totalSeconds + 1
			rect.Bottom = mr.height - mr.height*int32(task.BeginAt.Sub(chartBeginAt)/time.Second)/totalSeconds - 1
			mr.drawChart(hdc, &rect)
		}
	}
}

func (mr *MeterRenderer) isTaskConflict(tasks []logic.Task, desiredTask logic.Task) bool {
	for _, task := range tasks {
		if task.OverlapWith(desiredTask.BeginAt, desiredTask.EndAt) {
			return true
		}
	}
	return false
}

func (mr *MeterRenderer) drawChart(hdc winapi.HDC, rect *wrapped.RECT) {
	winapi.FillRect(hdc, rect.Unwrap(), mr.chartBrush)
}

func (mr *MeterRenderer) drawAllScaleLines(hdc winapi.HDC, futureDuration, pastDuration, interval time.Duration) {
	offset := futureDuration
	totalDuration := futureDuration + pastDuration
	totalSeconds := int32(totalDuration / time.Second)

	winapi.SelectObject(hdc, winapi.HGDIOBJ(mr.hourPen))

	for interval < offset {
		offset -= interval
	}

	for offset < totalDuration {
		if offset != futureDuration {
			mr.drawScaleLine(hdc, mr.height*int32(offset/time.Second)/totalSeconds)
		}
		offset += interval
	}

	winapi.SelectObject(hdc, winapi.HGDIOBJ(mr.headPen))
	mr.drawScaleLine(hdc, mr.height*int32(futureDuration/time.Second)/totalSeconds)
}

func (mr *MeterRenderer) drawScaleLine(hdc winapi.HDC, y int32) {
	winapi.MoveToEx(hdc, 0, y, nil)
	winapi.LineTo(hdc, mr.width, y)
}
