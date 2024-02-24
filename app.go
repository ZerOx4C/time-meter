package main

type App struct {
}

func (a *App) initialize() error {
	meter := new(MeterWindow)

	if err := meter.initialize(); err != nil {
		return err
	}

	meter.show()

	meter.finalize()

	return nil
}
