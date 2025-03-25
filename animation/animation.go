package animation

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type Animation struct {
	x, y                float64
	scaleX, scaleY      float64
	currentFrame        float64
	isRepeatable        bool
	isReverse           bool
	nowReverseAnimation bool
	start               bool
	callbackDone        bool
	images              []*ebiten.Image
	player              *audio.Player
	callback            func()
}

func NewAnimation(images []*ebiten.Image) *Animation {
	return &Animation{
		images: images,
	}
}

func (animation *Animation) Update(speed float64) {
	if !animation.start {
		return
	}

	if animation.player != nil {
		animation.player.Play()
	}
	if animation.callback != nil && int(animation.currentFrame) >= len(animation.images) && !animation.callbackDone {
		animation.callback()
		animation.callbackDone = true
	}
	if animation.nowReverseAnimation {
		animation.currentFrame -= speed
	} else {
		animation.currentFrame += speed
	}
}

func (animation *Animation) Draw(screen *ebiten.Image) {
	if !animation.start {
		return
	}
	if !animation.isRepeatable && int(animation.currentFrame) >= len(animation.images) {
		return
	}
	if animation.isReverse && int(animation.currentFrame) < 0 {
		animation.currentFrame = 0
		animation.nowReverseAnimation = false
	}
	if int(animation.currentFrame) >= len(animation.images) {
		if animation.isReverse {
			animation.nowReverseAnimation = true
			animation.currentFrame--
		} else {
			animation.currentFrame = 0
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(animation.x, animation.y)
	op.GeoM.Scale(animation.scaleX, animation.scaleY)
	screen.DrawImage(animation.images[int(animation.currentFrame)], op)
}

func (animation *Animation) SetPosition(x, y float64) {
	animation.x = x
	animation.y = y
}

func (animation *Animation) SetRepeatable(enabled bool) {
	animation.isRepeatable = enabled
}

func (animation *Animation) SetReverse(isReverse bool) {
	animation.isReverse = isReverse
}

func (animation *Animation) Start() {
	animation.start = true
}

func (animation *Animation) Stop() {
	animation.start = false
}

func (animation *Animation) Reset() {
	animation.currentFrame = 0
	animation.callbackDone = false
	animation.start = false
	if animation.player != nil {
		animation.player.Pause()
		animation.player.SetPosition(0)
	}
}

func (animation *Animation) SetCallback(callback func()) {
	animation.callback = callback
}

func (animation *Animation) SetScale(scaleX, scaleY float64) {
	animation.scaleX = scaleX
	animation.scaleY = scaleY
}

func (animation *Animation) SetSound(audioContext *audio.Context, fileName string) error {
	musicFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open music: %v", err)
	}
	mp3Stream, err := mp3.DecodeF32(musicFile)
	if err != nil {
		return fmt.Errorf("failed to decode music: %v", err)
	}

	player, err := audioContext.NewPlayerF32(mp3Stream)
	if err != nil {
		return fmt.Errorf("failed to create player: %v", err)
	}
	animation.player = player

	return nil
}

func (animation *Animation) SetVolume(volume float64) {
	animation.player.SetVolume(volume)
}
