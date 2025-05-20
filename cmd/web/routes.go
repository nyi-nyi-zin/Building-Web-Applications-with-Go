package main // main package ဖြစ်ကြောင်း ကြေငြာထားတာပါ

import (
	"net/http" // HTTP server အတွက်

	"github.com/go-chi/chi/v5" // Chi router အတွက်
) 

func (app *application) routes() http.Handler { // routes တွေကို handle လုပ်တဲ့ function
	mux := chi.NewRouter() // Chi router ကို create လုပ်ထားတာပါ

	mux.Get("/virtual-terminal", app.VirtualTerminal)
	return mux  // router ကို return လုပ်ပါ
}
