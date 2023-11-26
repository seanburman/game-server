package input

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"net/http"

	"code.rocketnine.space/tslocum/messeji"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

// var Input = messeji.NewInputField(font)

type FontFace struct {
	*sfnt.Font
	io.Closer
}

func Roboto() (*font.Face, error) {
	// b, err := os.ReadFile("assets/input/Roboto-Regular.ttf")
	// b, err := GetFont("http://127.0.0.1:5500/assets/input/Roboto-Regular.ttf")
	b, err := GetFont("http://192.168.1.154:5500/assets/input/Roboto-Regular.ttf")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	f, err := opentype.Parse(b)
	if err != nil {
		return nil, err
	}
	ff, err := opentype.NewFace(f, &opentype.FaceOptions{Size: 32, DPI: 72})
	if err != nil {
		return nil, err
	}
	return &ff, nil
}

type TextInput struct {
	Text  string
	Field *messeji.InputField
}

func NewTextInput() *TextInput {
	font, err := Roboto()
	if err != nil {
		log.Fatal(err)
	}
	field := messeji.NewInputField(*font)
	ti := TextInput{Field: field}
	field.SetHandleKeyboard(true)
	field.SetFont(*font)
	field.SetRect(image.Rect(0, 0, 200, 50))
	field.SetText("Test")
	field.SetForegroundColor(color.Black)

	field.SetChangedFunc(ti.ChangedFunc())
	field.Write([]byte("Cool"))
	return &ti
}

func (t *TextInput) SetText(text string) {
	t.Text = text
}

func (t *TextInput) ChangedFunc() func(r rune) bool {
	return func(r rune) bool {
		fmt.Println(r)
		text := string(r)
		t.SetText(text)
		return true
	}
}

func GetFont(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sead body: %v", err)
	}

	return data, nil
}
