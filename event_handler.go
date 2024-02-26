package main

type EventHandler func()

func (eh EventHandler) Invoke() {
	if eh != nil {
		eh()
	}
}
