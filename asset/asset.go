package asset

import (
	"embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2/audio/wav"

	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/hajimehoshi/ebiten/v2"
)

const sampleRate = 44100

//go:embed image map sound
var FS embed.FS

var ImgWhiteSquare = ebiten.NewImage(8, 8)
var ImgBlackSquare = ebiten.NewImage(8, 8)

var (
	ImgTitle1 = LoadImage("image/title1.png")
	ImgTitle2 = LoadImage("image/title2.png")
	ImgTitle3 = LoadImage("image/title3.png")

	ImgBat        = LoadImage("image/bat.png")
	ImgBatBroken1 = LoadImage("image/bat_broken1.png")
	ImgBatBroken2 = LoadImage("image/bat_broken2.png")

	ImgSnowflake = LoadImage("image/snowflake.png")
)

var (
	SoundBatHit1 *audio.Player
	SoundBatHit2 *audio.Player
	SoundBatHit3 *audio.Player
	SoundBatHit4 *audio.Player
)

func init() {
	ImgWhiteSquare.Fill(color.White)
	ImgBlackSquare.Fill(color.Black)
}

func LoadSounds(ctx *audio.Context) {
	SoundBatHit1 = LoadWAV(ctx, "sound/bat_hit/hit1.wav")
	SoundBatHit2 = LoadWAV(ctx, "sound/bat_hit/hit2.wav")
	SoundBatHit3 = LoadWAV(ctx, "sound/bat_hit/hit3.wav")
	SoundBatHit4 = LoadWAV(ctx, "sound/bonk.wav")
}

func LoadImage(p string) *ebiten.Image {
	f, err := FS.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	baseImg, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(baseImg)
}

func LoadBytes(p string) []byte {
	b, err := FS.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return b
}

func LoadWAV(context *audio.Context, p string) *audio.Player {
	f, err := FS.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	stream, err := wav.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		panic(err)
	}

	player, err := context.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	// Workaround to prevent delays when playing for the first time.
	player.SetVolume(0)
	player.Play()
	player.Pause()
	player.Rewind()
	player.SetVolume(1)

	return player
}
