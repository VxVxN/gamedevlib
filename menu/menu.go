package menu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

type Menu struct {
	activeItemMenu            int
	windowWidth, windowHeight float64
	face                      *text.GoTextFace
	buttonOptions             []ButtonOptions
	options                   MenuOptions
}

type ButtonOptions struct {
	Text   string
	Action func()
}

type MenuOptions struct {
	ButtonPadding       int
	ActiveItemMenu      color.RGBA
	InactiveButtonColor color.RGBA
}

func NewMenu(windowWidth, windowHeight float64, face *text.GoTextFace, buttonOptions []ButtonOptions, options MenuOptions) (*Menu, error) {
	if isEmptyRGBA(options.ActiveItemMenu) {
		options.ActiveItemMenu = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	if isEmptyRGBA(options.InactiveButtonColor) {
		options.InactiveButtonColor = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	}
	return &Menu{
		windowWidth:   windowWidth,
		windowHeight:  windowHeight,
		face:          face,
		buttonOptions: buttonOptions,
		options:       options,
	}, nil
}

func (menu *Menu) NextMenuItem() {
	menu.activeItemMenu++
	if menu.activeItemMenu > len(menu.buttonOptions)-1 {
		menu.activeItemMenu = len(menu.buttonOptions) - 1
	}
}

func (menu *Menu) BeforeMenuItem() {
	menu.activeItemMenu--
	if menu.activeItemMenu < 0 {
		menu.activeItemMenu = 0
	}
}

func (menu *Menu) ClickActiveButton() {
	menu.buttonOptions[menu.activeItemMenu].Action()
}

func (menu *Menu) Draw(screen *ebiten.Image) {
	for i, buttonOp := range menu.buttonOptions {
		buttonColor := menu.options.InactiveButtonColor
		if menu.activeItemMenu == i {
			buttonColor = menu.options.ActiveItemMenu
		}
		op := &text.DrawOptions{}
		op.GeoM.Translate(menu.windowWidth/2, menu.windowHeight/2+float64(i)*menu.face.Size+float64(i*menu.options.ButtonPadding))
		op.ColorScale.ScaleWithColor(buttonColor)
		op.LayoutOptions.PrimaryAlign = text.AlignCenter
		text.Draw(screen, buttonOp.Text, menu.face, op)
	}
}

func isEmptyRGBA(c color.RGBA) bool {
	return c == color.RGBA{R: 0, G: 0, B: 0, A: 0}
}
