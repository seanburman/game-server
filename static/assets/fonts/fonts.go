package fonts

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"net/http"

	"code.rocketnine.space/tslocum/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/seanburman/game/config"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type Font int

const (
	Roboto Font = iota
)

type FontSize int

const (
	Small FontSize = iota
)

type FontFace struct {
	*sfnt.Font
	io.Closer
}

type TextInput struct {
	Text  string
	Field *messeji.InputField
	Font  Font
}

func (t *TextInput) Update() {
	t.Field.Update()
}

func (t *TextInput) Draw(screen *ebiten.Image) {
	t.Field.Draw(screen)
}

func NewTextInput(font Font) *TextInput {
	ti := TextInput{}
	ff, err := ti.GetFont(Roboto, 32, 72)
	if err != nil {
		log.Fatal("failed to get font:", err)
	}
	field := messeji.NewInputField(ff)
	field.SetHandleKeyboard(true)
	field.SetRect(image.Rect(0, 0, 200, 50))
	field.SetForegroundColor(color.Black)

	field.SetChangedFunc(ti.ChangedFunc())
	field.Write([]byte("Cool"))
	return &ti
}

func (t *TextInput) GetFont(font Font, size float64, dpi float64) (font.Face, error) {
	var url string
	switch font {
	case Roboto:
		url = fmt.Sprintf("%s:5500/assets/fonts/Roboto-Regular.ttf", config.Env().HOST)
	default:
		url = fmt.Sprintf("%s:5500/assets/fonts/Roboto-Regular.ttf", config.Env().HOST)
	}

	b, err := LoadFont(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	f, err := opentype.Parse(b)
	if err != nil {
		return nil, err
	}
	ff, err := opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: dpi})
	if err != nil {
		return nil, err
	}
	return ff, nil
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

func LoadFont(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sead body: %v", err)
	}

	return b, nil
}
