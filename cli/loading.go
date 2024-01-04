package main

import "github.com/rivo/tview"

func (a *App) initLoadingView() {
	a.loadingView = tview.NewTextView()
	a.pages.AddPage("loading", a.loadingView, true, false)
}

func (a *App) showLoadingView(loadingMsg string) {
	a.loadingView.SetText(loadingMsg)
	a.showView(loading)
}
