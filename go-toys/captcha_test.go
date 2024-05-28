package main_test

import "testing"
import "github.com/steambap/captcha"

func TestCaptchaGen(t *testing.T) {
	t.Log("TestCaptchaGen")
	// create a captcha of 150x50px
	data, _ := captcha.New(150, 50)

	// session come from other library such as gorilla/sessions
	session.Values["captcha"] = data.Text
	session.Save(r, w)
	// send image data to client
	data.WriteImage(w)
}
