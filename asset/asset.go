package asset

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	_ "image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
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

	SoundCreepDie *audio.Player

	SoundTitleMusic *audio.Player
	SoundLevelMusic *audio.Player
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

	SoundCreepDie = LoadOGG(ctx, "sound/creep_die/creep_die.ogg", false)

	SoundTitleMusic = LoadOGG(ctx, "sound/title_music.ogg", true)
	SoundTitleMusic.SetVolume(0.5)

	SoundLevelMusic = LoadOGG(ctx, "sound/level_music.ogg", true)
	SoundLevelMusic.SetVolume(0.4)
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

func LoadOGG(context *audio.Context, p string, loop bool) *audio.Player {
	b := LoadBytes(p)

	stream, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	var s io.Reader
	if loop {
		s = audio.NewInfiniteLoop(stream, stream.Length())
	} else {
		s = stream
	}

	player, err := context.NewPlayer(s)
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
