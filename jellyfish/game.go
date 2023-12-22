package jellyfish

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

//go:embed resources/image/jellyfish.png
var byteJellyFishImage []byte

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	Title        = "jellyfish"

	// クラゲの落下速度
	baseSpeed = 3

	// デフォルト増加数
	defaultIncrease = 0
)

var (
	JellyFishImage *ebiten.Image
	ImgX           int
	ImgY           int
	Sec            int // 現在の秒
)

func init() {
	var err error
	var img image.Image
	img, _, err = image.Decode(bytes.NewReader(byteJellyFishImage))
	if err != nil {
		panic(fmt.Sprintf("failed to decode image: %v", err))
	}
	JellyFishImage = ebiten.NewImageFromImage(img)
	ImgX = JellyFishImage.Bounds().Size().X
	ImgY = JellyFishImage.Bounds().Size().Y
	Sec = time.Now().Second()
}

type jellyfish struct {
	// 現在の位置
	x float64
	y float64
	// 落下速度
	speed float64
	// クラゲのサイズ倍率
	scale float64
	// 削除フラグ
	deleted bool
}

type Flag struct {
	// クリック判定
	autoClick bool
}

type Game struct {
	score float64
	// クラゲ
	jellyfishes []*jellyfish
	// クリック判定
	clicked bool
	// Flag
	flag Flag
	// 増加数
	increase int
	// debug
	debug      bool
	debugCount int
	debugX     float64
	debugY     float64
}

func NewGame() *Game {
	game := &Game{}
	game.init()
	return game
}

func (g *Game) init() {
	g.score = 0
	// クラゲ
	gf := &jellyfish{
		x:       0,
		y:       0,
		speed:   baseSpeed,
		deleted: true,
	}
	g.jellyfishes = []*jellyfish{}
	g.jellyfishes = append(g.jellyfishes, gf)
	// クリック判定
	g.clicked = false
	// Flag
	g.flag.autoClick = false

	// 増加数
	g.increase = defaultIncrease

	// 環境変数を取得
	g.debug = os.Getenv("DEBUG_JELLYFISH") == "true"
	if g.debug {
		g.debugCount = 0
		g.debugX = 0
		g.debugY = 0
	}

}

func (g *Game) Update() error {
	// フラグ
	g.setFlag()

	// クリックしたらクラゲを追加
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.clicked {
			g.addJellyFish()
			g.clicked = true
		}
	} else {
		g.clicked = false
	}

	// 1秒に1回クラゲを追加
	if g.flag.autoClick {
		current_sec := time.Now().Second()
		if Sec != current_sec {
			for i := 0; i < g.increase; i++ {
				g.addJellyFish()
			}
			Sec = current_sec
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 背景
	screen.Fill(color.RGBA{255, 245, 228, 0xff})

	// デバッグ用
	if g.debug {
		g.debugCount++
		ebitenutil.DebugPrint(screen, fmt.Sprintf("score: %0.0f\n", g.score))
	}

	// クラゲの数を表示
	// text.Draw(screen, fmt.Sprintf("score: %0.0f\n", g.score), DefaultFont, 10, 20, White)

	// クラゲ
	for i, j := range g.jellyfishes {
		if j.deleted {
			// 配列から削除
			if i == 0 {
				g.jellyfishes = g.jellyfishes[1:]
			}
			continue
		}
		// クラゲが画面外に出たら削除
		if j.y > ScreenHeight {
			j.deleted = true
		}
		op := &ebiten.DrawImageOptions{}
		j.y += j.speed // 落下
		op.GeoM.Scale(j.scale, j.scale)
		motoSizeX := float64(ImgX) * j.scale
		motoSizeY := float64(ImgY) * j.scale
		op.GeoM.Translate(-motoSizeX/2, -motoSizeY/2)
		op.GeoM.Translate(j.x, j.y)
		screen.DrawImage(JellyFishImage, op)
	}

	// スコアの表示
	face := basicfont.Face7x13
	scoreText := fmt.Sprintf("Score: %0.0f", g.score)
	text.Draw(screen, scoreText, face, 5, ScreenHeight-5, color.Black)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// JEYLLYFISH
func (g *Game) addJellyFish() {
	x := rand.Intn(ScreenWidth) // 0~ScreenWidthの間でランダムにx座標を決める
	y := 0
	speed := baseSpeed + rand.Float64()*5
	// 0.1~0.3の間でランダムにサイズを決める
	scale := 0.1 + rand.Float64()*0.2
	gf := &jellyfish{
		x:       float64(x),
		y:       float64(y),
		scale:   scale,
		speed:   float64(speed),
		deleted: false,
	}
	g.score++
	g.jellyfishes = append(g.jellyfishes, gf)
}

// FLAG
// 条件を満たしたらフラグを立てる
func (g *Game) setFlag() {
	if !g.flag.autoClick && g.score > 100 {
		g.flag.autoClick = true
		g.increase += 1
	}
}
