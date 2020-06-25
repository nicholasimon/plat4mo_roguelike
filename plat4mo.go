package main

import (
	"math/rand"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ( // MARK: var

	// sound
	track       rl.Music
	choosetrack bool
	// imgs
	backcolor     rl.Color
	backfade      = float32(0.7)
	backfadeon    bool
	backimg       rl.Texture2D
	imgs          rl.Texture2D
	teleportimg   = rl.NewRectangle(0, 169, 16, 16)
	growimg       = rl.NewRectangle(16, 169, 16, 16)
	gunimg        = rl.NewRectangle(32, 169, 16, 16)
	hp2img        = rl.NewRectangle(0, 35, 14, 12)
	hpimg         = rl.NewRectangle(46, 169, 16, 16)
	bombimg       = rl.NewRectangle(62, 168, 16, 16)
	ladderimg     = rl.NewRectangle(77, 169, 16, 16)
	key           = rl.NewRectangle(16, 37, 14, 8)
	ghost         = rl.NewRectangle(0, 444, 44, 30)
	mushroom      = rl.NewRectangle(0, 480, 32, 32)
	chicken       = rl.NewRectangle(0, 399, 32, 34)
	radish        = rl.NewRectangle(0, 352, 30, 38)
	slime         = rl.NewRectangle(0, 311, 44, 30)
	flower        = rl.NewRectangle(0, 256, 44, 42)
	coin          = rl.NewRectangle(49, 0, 16, 16)
	coinstatic    = rl.NewRectangle(49, 0, 16, 16)
	arrowkeys     = rl.NewRectangle(0, 0, 49, 35)
	runr          = rl.NewRectangle(0, 63, 16, 16)
	runl          = rl.NewRectangle(0, 47, 16, 16)
	idlel         = rl.NewRectangle(0, 79, 16, 16)
	idler         = rl.NewRectangle(0, 95, 16, 16)
	blockspecial  = rl.NewRectangle(0, 143, 16, 16)
	blockspecial2 = rl.NewRectangle(16, 127, 16, 16)
	block1        = rl.NewRectangle(0, 111, 16, 16)
	block2        = rl.NewRectangle(16, 111, 16, 16)
	block3        = rl.NewRectangle(32, 111, 16, 16)
	block4        = rl.NewRectangle(48, 111, 16, 16)
	block5        = rl.NewRectangle(64, 111, 16, 16)
	block6        = rl.NewRectangle(80, 111, 16, 16)
	ground1       = rl.NewRectangle(0, 127, 16, 16)
	// intro
	intro, fadeon, flash4 bool
	fadeintro             = float32(0.0)
	flashtimer            = rInt(10, 40)
	// level end
	levelend bool
	// extras
	extrasmap   = make([]string, levela)
	extrascount = 0
	// enemies
	enemiesmap   = make([]int, levela)
	enemymovemap = make([]bool, levela)
	// use special
	usespecial, grow bool
	bulletmap        = make([]string, levela)
	// coins
	coinslevel, coinstotal int
	// specialblock
	specialblockmap = make([]int, levela)
	activespecial   int
	// player
	followplayer, jump, playerdir, playerdied, collisionpause, run bool
	player, playerh, playerv, playerx, playery, jumpcount          int
	playerhp                                                       = 3
	collisiontimer                                                 = 3
	// draw map
	pause                                                    bool
	drawblock, drawblocknext, drawblocknexth, drawblocknextv int
	// level
	levelw   = 2000
	levelh   = 120
	levela   = levelh * levelw
	levelmap = make([]string, levela)
	tilesmap = make([]int, levela)
	// core
	lockview                             bool
	screenw, screenh, screena            int
	monh32, monw32                       int32
	monitorh, monitorw, monitornum       int
	grid16on, grid4on, debugon, lrg, sml bool
	framecount                           int
	mousepos                             rl.Vector2
	camera                               rl.Camera2D
	cameraintro                          rl.Camera2D
	camerainfobar                        rl.Camera2D
)

// MARK: notes
/*
1920X1080 / 16px blocks
width 128 blocks (127+1)
height 68 blocks (67+1)

*/
func updateall() { // MARK: updateall()
	if followplayer {
		updatedrawblockplayer()
	}
	if grid16on {
		grid16()
	}
	if grid4on {
		grid4()
	}

	gravity()
	getpositions()
	if lockview == false {
		camera.Zoom = 4.0
		camera.Target.Y = float32(playery - (monitorh / 5))
		camera.Target.X = float32(playerx - (monitorw / 8))
		lockview = true
		followplayer = true
	}
	if jump {
		jumpmove()
	}
	if usespecial {
		specialactive()
	}
	moveblocks()
	enemies()
	timers()
	animations()
	if levelend {
		nextlevel()
	}
}
func animations() { // MARK: 	animations()
	if framecount%3 == 0 {
		idlel.X += 16
		idler.X += 16
		runr.X += 16
		runl.X += 16
		coin.X += 16
		mushroom.X += 32
		ghost.X += 44
		flower.X += 44
		slime.X += 44
		chicken.X += 32
		radish.X += 30
	}
	if radish.X > 340 {
		radish.X = 0
	}
	if flower.X > 450 {
		flower.X = 0
	}
	if chicken.X > 420 {
		chicken.X = 0
	}
	if slime.X > 400 {
		slime.X = 0
	}
	if ghost.X > 400 {
		ghost.X = 0
	}
	if mushroom.X > 490 {
		mushroom.X = 0
	}
	if coin.X >= 120 {
		coin.X = 49
	}
	if idlel.X >= 82 {
		idlel.X = 0
	}
	if idler.X >= 82 {
		idler.X = 0
	}
	if runl.X >= 82 {
		runl.X = 0
	}
	if runr.X >= 82 {
		runr.X = 0
	}
}
func nextlevel() { // MARK: nextlevel()
	pause = true
	rl.StopMusicStream(track)
	for a := 0; a < levela; a++ {
		enemiesmap[a] = 0
		enemymovemap[a] = false
		specialblockmap[a] = 0
		extrasmap[a] = ""
	}
	createlevel()
	playerhp = 3
	extrascount = 0
	activespecial = 0
	coinslevel = 0
	drawblock = 0
	drawblocknext = 0
	choosebackground()
	choosetrack = flipcoin()
	musicrestart()
	levelend = false
	pause = false
}
func timers() { // MARK: timers()
	if collisionpause {
		if framecount%60 == 0 {
			collisiontimer--
			if collisiontimer == 0 {
				collisionpause = false
				collisiontimer = 3
			}
		}
	}
}
func enemies() { // MARK: enemies()
	if rolldice()+rolldice() >= 11 {
		block := player - levelw*12
		block += rInt((0 - (screenw/2 - 1)), (screenw/2 - 1))
		enemiesmap[block] = rolldice()
		enemymovemap[block] = flipcoin()
	}
}
func moveblocks() { // MARK: moveblocks()
	// left up
	for a := 0; a < levela; a++ {
		checkbullet := bulletmap[a]
		checkenemy := enemiesmap[a]
		if checkenemy != 0 {
			if player == a {
				playercollision()
			}
			if levelmap[a+levelw] == "#" || levelmap[a+levelw] == "_" {
				if enemymovemap[a] == true {
					holder := enemiesmap[a]
					holder2 := enemymovemap[a]
					enemyh := a / levelw
					enemyv := a - (enemyh * levelw)
					if checkenemymove(a, "l") {
						if rolldice()+rolldice() >= 11 {
							enemiesmap[a] = 0
							enemymovemap[a] = false
							enemiesmap[a-(levelw*2)] = holder
							enemymovemap[a-(levelw*2)] = holder2
						} else {
							if framecount%4 == 0 {
								if enemyv > 1 {
									enemiesmap[a] = 0
									enemymovemap[a] = false
									enemiesmap[a-1] = holder
									enemymovemap[a-1] = holder2
								} else {
									enemiesmap[a] = 0
									enemymovemap[a] = false
								}
							}
						}
					} else {
						enemiesmap[a] = 0
						enemymovemap[a] = false
						enemiesmap[a+1] = holder
						enemymovemap[a+1] = false
						enemymovemap[a-1] = true
					}
				}
			}
		}
		if checkbullet != "" {
			bulleth := a / levelw
			bulletv := a - (bulleth * levelw)
			if checkbullet == "<" {
				bulletmap[a] = ""
				if bulletv > 1 {
					bulletmap[a-1] = "<"
				}
			}
		}
	}
	// right down
	for a := levela - 1; a > 0; a-- {
		checkbullet := bulletmap[a]
		checkenemy := enemiesmap[a]
		if checkenemy != 0 {
			if player == a {
				playercollision()
			}
			if levelmap[a+levelw] != "_" {
				if levelmap[a+levelw] != "#" {
					holder := enemiesmap[a]
					holder2 := enemymovemap[a]
					enemiesmap[a] = 0
					enemymovemap[a] = false
					enemiesmap[a+levelw] = holder
					enemymovemap[a+levelw] = holder2
				}
			}
			if levelmap[a+levelw] == "#" || levelmap[a+levelw] == "_" {
				if enemymovemap[a] == false {
					holder := enemiesmap[a]
					holder2 := enemymovemap[a]
					enemyh := a / levelw
					enemyv := a - (enemyh * levelw)
					if checkenemymove(a, "l") {
						if framecount%4 == 0 {
							if enemyv < (levelw - 1) {
								enemiesmap[a] = 0
								enemymovemap[a] = false
								enemiesmap[a+1] = holder
								enemymovemap[a+1] = holder2
							} else {
								enemiesmap[a] = 0
								enemymovemap[a] = false
							}
						}
					} else {
						enemiesmap[a] = 0
						enemymovemap[a] = false
						enemiesmap[a-1] = holder
						enemymovemap[a-1] = true
						enemymovemap[a-1] = false
					}
				}
			}
		}
		if checkbullet != "" {
			bulleth := a / levelw
			bulletv := a - (bulleth * levelw)
			if checkbullet == ">" {
				bulletmap[a] = ""
				if bulletv < levelw-1 {
					bulletmap[a+1] = ">"
				}
			}
		}
	}
}
func playercollision() { // MARK: playercollision()
	if !collisionpause {
		playerhp--
		if playerhp <= 0 {
			playerdied = true
			playerhp = 0
		}
		collisionpause = true
	}
}
func checkenemymove(block int, dir string) bool {
	free := true
	switch dir {
	case "l":
		if enemiesmap[block-1] != 0 {
			free = false
		}
	case "r":
		if enemiesmap[block+1] != 0 {
			free = false
		}
	}
	return free
}
func specialblock(position int) { // MARK: specialblock()
	levelmap[position] = "^"
	specialblockmap[position-levelw] = rolldice()
}
func specialactive() { // MARK: specialblock()

	switch activespecial {
	case 1:
		// ladder

		if playerh > (levelh/2)+20 {
			block := player - levelw*3
			block--
			for b := 0; b < 3; b++ {
				for a := 0; a < 3; a++ {
					levelmap[block] = "#"
					tiletype := rolldice()
					tilesmap[block] = tiletype
					block++
				}
				block -= 3
				block -= levelw * 3
			}
		}
	case 2:
		// bomb

		block := player - levelw*4
		block -= 2
		for b := 0; b < 9; b++ {
			for a := 0; a < 5; a++ {
				if levelmap[block] == "#" {
					levelmap[block] = "."
					tilesmap[block] = 0
				}
				block++
			}
			block -= 5
			block += levelw
		}
	case 3:
		// teleport
		block := rInt(screenw, levelw-screenw)
		change := rInt(6, 18)
		change = levelh - change
		block += levelw * change
		if levelmap[block] == "#" {
			block -= levelw
		}
		player = block
	case 4:
		// shoot
		block := player
		if playerdir {
			bulletmap[block] = "<"
		} else {
			bulletmap[block] = ">"
		}
	case 5:
		// grow
		if grow {
			grow = false
		} else {
			grow = true
		}
	case 6:
		// playerhp
		if playerhp < 3 {
			playerhp++
		}
	}
	usespecial = false
	activespecial = 0
}
func jumpmove() { // MARK: jumpmove()
	if jumpcount != 0 {
		if levelmap[player-levelw] != "*" {
			player -= levelw
			jumpcount--
		} else {
			jumpcount = 0
			jump = false
			specialblock(player - levelw)
		}
	} else {
		jump = false
	}
}
func gravity() { // MARK: gravity()
	if !jump {
		if levelmap[player+levelw] == "." || levelmap[player+levelw] == "$" {
			player += levelw
		}
	}
}
func updatedrawblockplayer() { // MARK: updatedrawblockplayer()
	if playerh <= levelh-6 && playerh > levelh-9 {
		drawblocknext = player - (levelw * (screenh - 5))
		drawblocknext -= ((screenw / 2) - 1)
	} else if playerh <= levelh-9 {
		drawblocknext = player - (levelw * (screenh - 8))
		drawblocknext -= ((screenw / 2) - 1)
	}
}
func getpositions() { // MARK: getpositions()
	drawblocknexth = drawblocknext / levelw
	drawblocknextv = drawblocknext - (drawblocknexth * levelw)
	playerh = player / levelw
	playerv = player - (playerh * levelw)
}
func createlevel() { // MARK: createlevel()

	for a := 0; a < levela; a++ {
		levelmap[a] = "."
	}
	block := levela - levelw
	block -= levelw * 4
	for a := 0; a < (levelw * 4); a++ {
		levelmap[block] = "_"
		block++
	}
	block = levela - levelw
	block -= levelw * 8
	count := block
	for a := 0; a < (levelw); a++ {
		if rolldice()+rolldice() >= 11 {
			block = count
			platlen := rInt(8, 15)
			tiletype := rolldice()
			for b := 0; b < platlen; b++ {
				levelmap[block] = "#"
				tilesmap[block] = tiletype
				if rolldice()+rolldice() >= 11 {
					levelmap[block] = "*"
					tilesmap[block] = 0
				}
				block++
			}
		}
		count++
	}
	// coins
	block = levela - levelw
	block -= levelw * 10
	for a := 0; a < levelw; a++ {
		if rolldice()+rolldice() == 12 {
			levelmap[block] = "$"
			block += 2
			a += 2
		}
		block++
	}
	block = levela - levelw
	block -= levelw * 12
	count = block
	for a := 0; a < (levelw); a++ {
		if rolldice()+rolldice() >= 11 {
			block = count
			platlen := rInt(8, 15)
			tiletype := rolldice()
			for b := 0; b < platlen; b++ {
				levelmap[block] = "#"
				tilesmap[block] = tiletype
				if rolldice()+rolldice() >= 11 {
					levelmap[block] = "*"
					tilesmap[block] = 0
				}
				block++
			}
		}
		count++
	}
	// coins
	block = levela - levelw
	block -= levelw * 14
	for a := 0; a < levelw; a++ {
		if rolldice()+rolldice()+rolldice() >= 17 {
			levelmap[block] = "$"
			block += 2
			a += 2
		}
		block++
	}
	block = levela - levelw
	block -= levelw * 16
	count = block
	for a := 0; a < (levelw); a++ {
		if rolldice()+rolldice() >= 11 {
			block = count
			platlen := rInt(8, 15)
			tiletype := rolldice()
			for b := 0; b < platlen; b++ {
				levelmap[block] = "#"
				tilesmap[block] = tiletype
				if rolldice()+rolldice() >= 11 {
					levelmap[block] = "*"
					tilesmap[block] = 0
				}
				block++
			}
		}
		count++
	}
	// coins
	block = levela - levelw
	block -= levelw * 18
	for a := 0; a < levelw; a++ {
		if rolldice()+rolldice()+rolldice() >= 17 {
			levelmap[block] = "$"
			block += 2
			a += 2
		}
		block++
	}
	block = levela - levelw
	block -= levelw * 20
	count = block
	for a := 0; a < (levelw); a++ {
		if rolldice()+rolldice()+rolldice() >= 17 {
			block = count
			platlen := rInt(6, 11)
			tiletype := rolldice()
			for b := 0; b < platlen; b++ {
				levelmap[block] = "#"
				tilesmap[block] = tiletype
				if rolldice()+rolldice() >= 11 {
					levelmap[block] = "*"
					tilesmap[block] = 0
				}
				block++
			}
		}
		count++
	}
	block = levela - levelw
	block -= levelw * 24
	count = block
	for a := 0; a < (levelw); a++ {
		if rolldice()+rolldice()+rolldice() >= 17 {
			block = count
			platlen := rInt(6, 11)
			tiletype := rolldice()
			for b := 0; b < platlen; b++ {
				levelmap[block] = "#"
				tilesmap[block] = tiletype
				if rolldice()+rolldice() >= 11 {
					levelmap[block] = "*"
					tilesmap[block] = 0
				}
				block++
			}
		}
		count++
	}

	block = levela - levelw
	block -= levelw * 30
	count = block
	for a := 0; a < (levelw); a++ {
		if rolldice()+rolldice()+rolldice() >= 16 {
			block = count
			roomlen := rInt(8, 15)
			rooma := roomlen * roomlen
			count2 := 0
			tiletype := rolldice()
			for b := 0; b < rooma; b++ {
				levelmap[block] = "#"
				tilesmap[block] = tiletype
				block++
				count2++
				if count2 == roomlen {
					count2 = 0
					block -= levelw + roomlen
				}
			}
			block = count + 1
			block -= levelw
			roomlen -= 2
			rooma = roomlen * roomlen
			count2 = 0
			for b := 0; b < rooma; b++ {
				levelmap[block] = "."
				tilesmap[block] = 0
				block++
				count2++
				if count2 == roomlen {
					count2 = 0
					block -= levelw + roomlen
				}
			}
			if rolldice()+rolldice() >= 8 {
				block += roomlen / 2
				block += (roomlen - 2) * levelw
				extrasmap[block] = "!"
			}
			a += roomlen + 2
			count += roomlen + 2
		}
		count++
	}

}
func main() { // MARK: main()
	rand.Seed(time.Now().UnixNano()) // random numbers
	rl.SetTraceLog(rl.LogError)      // hides INFO window
	startsettings()
	createlevel()

	raylib()
}
func checkcoin(block int) { // MARK: checkcoin()
	if levelmap[block] == "$" {
		levelmap[block] = "."
		coinslevel++
		coinstotal++
	}
}
func checkpowerup(block int) { // MARK: checkpowerup()
	if specialblockmap[block] != 0 {
		activespecial = specialblockmap[block]
		specialblockmap[block] = 0
	}
}
func collectextra(block int) { // MARK: collectextra()
	extrasmap[block] = ""
	extrascount++
	if extrascount == 3 {
		levelend = true
	}
}
func startplayer() { // MARK: startplayer()
	player = levela - levelw
	player -= levelw * 5
	player += screenw / 2

	drawblocknext = player
	drawblocknext -= levelw * (screenh - 5)
	drawblocknext -= ((screenw / 2) - 1)
}
func flashtime() {
	flashtimer = rInt(10, 40)
}
func music() {
	rl.InitAudioDevice()
	if choosetrack {
		track = rl.LoadMusicStream("audio1.mp3")
	} else {
		track = rl.LoadMusicStream("audio2.mp3")
	}
	rl.PlayMusicStream(track)
}
func musicrestart() {
	if choosetrack {
		track = rl.LoadMusicStream("audio1.mp3")
	} else {
		track = rl.LoadMusicStream("audio2.mp3")
	}
	rl.PlayMusicStream(track)
}
func introscreen() { // MARK: 	introscreen()
	rl.DrawRectangle(0, 0, monw32, monh32, rl.White)
	rl.DrawText("plat4mo", 103, 103, 200, rl.Black)
	rl.DrawText("plat4mo", 102, 102, 200, rl.White)
	rl.DrawText("plat4mo", 100, 100, 200, rl.Black)
	if flash4 {
	} else {
		rl.DrawText("4", 500, 100, 200, rl.Fade(rl.Red, 0.7))
	}
	rl.DrawText("arrow keys", 280, 340, 40, rl.Black)
	rl.DrawText("move", 340, 610, 40, rl.Black)

	rl.DrawRectangle(600, 510, 380, 50, rl.Black)
	rl.DrawRectangle(608, 518, 364, 34, rl.White)
	rl.DrawText("space bar", 710, 518, 30, rl.Black)
	rl.DrawText("use special", 680, 575, 40, rl.Black)
	rl.DrawText("press space bar to start", monw32/3, monh32-140, 50, rl.Black)
	rl.DrawRectangle(monw32/4, monh32-150, monw32/2, 60, rl.Fade(rl.White, fadeintro))

	rl.DrawText("get powerups", monw32/2, 160, 40, rl.Black)
	rl.DrawText("collect coins", monw32/2, 200, 40, rl.Black)
	rl.DrawText("find 3 keys", monw32/2, 240, 40, rl.Black)
	rl.DrawText("find door", monw32/2, 280, 40, rl.Black)
	rl.DrawText("new level", monw32/2, 320, 40, rl.Black)
	rl.DrawText("repeat", monw32/2, 360, 40, rl.Black)

	rl.BeginMode2D(cameraintro)
	// arrow keys
	keysv2 := rl.NewVector2(40, 65)
	rl.DrawTextureRec(imgs, arrowkeys, keysv2, rl.White)
	rl.EndMode2D()
	//flashtimer
	if framecount%flashtimer == 0 {
		if flash4 {
			flashtime()
			flash4 = false
		} else {
			flashtime()
			flash4 = true
		}
	}
	if fadeon {
		fadeintro -= 0.1
		if fadeintro <= 0.0 {
			fadeon = false
		}
	} else {
		fadeintro += 0.1
		if fadeintro >= 1.0 {
			fadeon = true
		}
	}
}
func choosebackground() { // MARK: 	choosebackground()

	color := rolldice()

	switch color {
	case 1:
		backcolor = rl.Red
	case 2:
		backcolor = rl.DarkBlue
	case 3:
		backcolor = rl.Magenta
	case 4:
		backcolor = rl.Green
	case 5:
		backcolor = rl.DarkPurple
	case 6:
		backcolor = rl.Pink

	}

}
func raylib() { // MARK: raylib()
	rl.InitWindow(monw32, monh32, "plat4mo")
	setscreen()
	rl.CloseWindow()
	rl.InitWindow(monw32, monh32, "plat4mo")
	rl.SetExitKey(rl.KeyEnd)          // key to end the game and close window
	imgs = rl.LoadTexture("imgs.png") // load images
	music()
	perlinNoise := rl.GenImagePerlinNoise(monitorw, monitorh, 50, 50, 4.0)
	backimg = rl.LoadTextureFromImage(perlinNoise)
	rl.SetTargetFPS(30)
	for !rl.WindowShouldClose() { // MARK: WindowShouldClose
		rl.UpdateMusicStream(track)
		//mousepos = rl.GetMousePosition()
		framecount++

		rl.ClearBackground(rl.Black)
		rl.DrawTexture(backimg, 0, 0, rl.Fade(backcolor, backfade)) // MARK: draw backimg

		if backfadeon {
			if framecount%3 == 0 {
				backfade += 0.1
				if backfade >= 0.7 {
					backfadeon = false
				}
			}
		} else {
			if framecount%3 == 0 {
				backfade -= 0.1
				if backfade <= 0.1 {
					backfadeon = true
				}
			}
		}

		rl.BeginMode2D(camera)

		// MARK: draw map layer 1 / left up
		if !pause {

			count := 0
			drawx := int32(0)
			drawy := int32(0)
			drawblock = drawblocknext
			for a := 0; a < screena; a++ {
				checklevel := levelmap[drawblock]
				checkspecial := specialblockmap[drawblock]
				checkbullet := bulletmap[drawblock]
				checkenemy := enemiesmap[drawblock]
				checkextras := extrasmap[drawblock]
				checktile := tilesmap[drawblock]

				if checkextras != "" {
					if player == drawblock {
						collectextra(drawblock)
					}
				}

				switch checklevel {
				case "#":
					// rl.DrawRectangle(drawx, drawy, 15, 15, rl.Red)
				//	blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
				//	rl.DrawTextureRec(imgs, block1, blockv2, rl.Green)
				case "*":
					//	rl.DrawRectangleLines(drawx, drawy, 15, 15, rl.Red)
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, blockspecial, blockv2, rl.Orange)
				case "_":
					//	rl.DrawRectangle(drawx, drawy, 15, 15, rl.Brown)
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, ground1, blockv2, rl.White)
				case "$":
					//		rl.DrawRectangle(drawx, drawy, 15, 15, rl.Yellow)
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, coin, blockv2, rl.White)
				case "^":
					// rl.DrawRectangle(drawx, drawy, 15, 15, rl.Green)
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, blockspecial, blockv2, rl.DarkGray)
				}
				switch checktile {
				case 1:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, block1, blockv2, rl.White)
				case 2:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, block2, blockv2, rl.White)
				case 3:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, block3, blockv2, rl.White)
				case 4:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, block4, blockv2, rl.White)
				case 5:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, block5, blockv2, rl.White)
				case 6:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, block6, blockv2, rl.White)
				}

				switch checkspecial {
				case 1:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, ladderimg, blockv2, rl.White)
				case 2:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, bombimg, blockv2, rl.White)
				case 3:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, teleportimg, blockv2, rl.White)
				case 4:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, gunimg, blockv2, rl.White)
				case 5:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, growimg, blockv2, rl.White)
				case 6:
					blockv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, hpimg, blockv2, rl.White)
				}
				switch checkextras {
				case "!":
					//	rl.DrawCircleLines(drawx+8, drawy+8, 7, rl.Yellow)
					keyv2 := rl.NewVector2(float32(drawx), float32(drawy))
					rl.DrawTextureRec(imgs, key, keyv2, rl.White)
				}
				switch checkenemy {
				case 1:
					//	rl.DrawRectangleLines(drawx, drawy, 15, 15, rl.White)
					enemyv2 := rl.NewVector2(float32(drawx-8), float32(drawy-16))
					rl.DrawTextureRec(imgs, mushroom, enemyv2, rl.White)
				case 2:
					//	rl.DrawRectangleLines(drawx, drawy, 15, 15, rl.White)
					enemyv2 := rl.NewVector2(float32(drawx-8), float32(drawy-16))
					rl.DrawTextureRec(imgs, ghost, enemyv2, rl.White)
				case 3:
					//	rl.DrawRectangleLines(drawx, drawy, 15, 15, rl.White)
					enemyv2 := rl.NewVector2(float32(drawx-8), float32(drawy-16))
					rl.DrawTextureRec(imgs, slime, enemyv2, rl.White)
				case 4:
					//	rl.DrawRectangleLines(drawx, drawy, 15, 15, rl.White)
					enemyv2 := rl.NewVector2(float32(drawx-8), float32(drawy-20))
					rl.DrawTextureRec(imgs, radish, enemyv2, rl.White)
				case 5:
					//	rl.DrawRectangleLines(drawx, drawy, 15, 15, rl.White)
					enemyv2 := rl.NewVector2(float32(drawx-8), float32(drawy-16))
					rl.DrawTextureRec(imgs, chicken, enemyv2, rl.White)
				case 6:
					//	rl.DrawRectangleLines(drawx, drawy, 15, 15, rl.White)
					enemyv2 := rl.NewVector2(float32(drawx-8), float32(drawy-26))
					rl.DrawTextureRec(imgs, flower, enemyv2, rl.White)

				}
				switch checkbullet {
				case ">", "<":
					rl.DrawRectangle(drawx+4, drawy+4, 4, 4, rl.Orange)
				}
				if player == drawblock {
					if grow {
						rl.DrawRectangle(drawx-16, drawy-16, 31, 31, rl.Blue)
					} else {
						//rl.DrawRectangle(drawx, drawy, 15, 15, rl.Blue)
						if playerdir {
							playerv2 := rl.NewVector2(float32(drawx), float32(drawy))
							if run {
								rl.DrawTextureRec(imgs, runr, playerv2, rl.White)
							} else {
								rl.DrawTextureRec(imgs, idler, playerv2, rl.White)
							}
						} else {
							playerv2 := rl.NewVector2(float32(drawx), float32(drawy))
							if run {
								rl.DrawTextureRec(imgs, runl, playerv2, rl.White)
							} else {
								rl.DrawTextureRec(imgs, idlel, playerv2, rl.White)
							}
						}
					}
					playerx = int(drawx)
					playery = int(drawy)
					checkcoin(drawblock)
					checkpowerup(drawblock)
				}
				count++
				drawx += 16
				drawblock++
				if count == screenw {
					drawblock += levelw - screenw
					count = 0
					drawy += 16
					drawx = 0
				}
			}

			rl.EndMode2D() // MARK: draw no camera

			// draw info bar
			coinslevelTEXT := strconv.Itoa(coinslevel)
			playerhpTEXT := strconv.Itoa(playerhp)
			extrascountTEXT := strconv.Itoa(extrascount)

			rl.DrawRectangle(0, 0, monw32, 80, rl.Fade(rl.Black, 0.7)) // top bar rectangle

			rl.DrawText("x", 80, 28, 20, rl.White)
			rl.DrawText(coinslevelTEXT, 100, 30, 20, rl.White)
			rl.DrawText("x", 180, 28, 20, rl.White)
			rl.DrawText(playerhpTEXT, 200, 30, 20, rl.White)
			if extrascount != 0 {
				rl.DrawText("x", 350, 28, 20, rl.White)
				rl.DrawText(extrascountTEXT, 370, 30, 20, rl.White)
			}
			rl.BeginMode2D(camerainfobar) //  info bar zoom
			coinv2 := rl.NewVector2(25, 13)
			rl.DrawTextureRec(imgs, coinstatic, coinv2, rl.White)
			specialv2 := rl.NewVector2(135, 14)
			switch activespecial {
			case 1:
				rl.DrawTextureRec(imgs, ladderimg, specialv2, rl.White)
			case 2:
				rl.DrawTextureRec(imgs, bombimg, specialv2, rl.White)
			case 3:
				rl.DrawTextureRec(imgs, teleportimg, specialv2, rl.White)
			case 4:
				rl.DrawTextureRec(imgs, gunimg, specialv2, rl.White)
			case 5:
				rl.DrawTextureRec(imgs, growimg, specialv2, rl.White)
			}
			hpv2 := rl.NewVector2(80, 13)
			rl.DrawTextureRec(imgs, hpimg, hpv2, rl.White)

			if extrascount != 0 {
				keyv2 := rl.NewVector2(172, 17)
				rl.DrawTextureRec(imgs, key, keyv2, rl.White)
			}

			rl.EndMode2D() // end info bar zoom

		}
		if debugon {
			debug()
		}
		rl.BeginDrawing()
		if !intro {
			pause = true
			introscreen()
		}
		rl.EndDrawing()
		input()
		updateall()

	}
	rl.CloseWindow()
}
func setscreen() { // MARK: setscreen()
	monitornum = rl.GetMonitorCount()
	monitorh = rl.GetScreenHeight()
	monitorw = rl.GetScreenWidth()
	monh32 = int32(monitorh)
	monw32 = int32(monitorw)
	rl.SetWindowSize(monitorw, monitorh)
	setsizes()
}
func setsizes() { // MARK: setsizes()
	if monitorw >= 1600 {
		lrg = true
		sml = false
	} else if monitorw < 1600 && monitorw >= 1280 {
		lrg = false
		sml = true
	}
	screenw = (monitorw / 16) + 1
	screenh = (monitorh / 16) + 1
	screena = screenw * screenh
	startplayer()
}
func startsettings() { // MARK: start
	camera.Zoom = 1.0
	camera.Target.X = 0.0
	camera.Target.Y = 0.0
	cameraintro.Zoom = 6.0
	cameraintro.Target.X = 0.0
	cameraintro.Target.Y = 0.0
	camerainfobar.Zoom = 1.8
	cameraintro.Target.X = 0.0
	cameraintro.Target.Y = 0.0
	choosetrack = flipcoin()
	choosebackground()
	//debugon = true
	//	grid16on = true
	//selectedmenuon = true
}
func debug() { // MARK: debug
	rl.DrawRectangle(monw32-300, 0, 500, monw32, rl.Fade(rl.Black, 0.7))
	rl.DrawFPS(monw32-290, monh32-100)

	screenhTEXT := strconv.Itoa(screenh)
	screenwTEXT := strconv.Itoa(screenw)
	playerxTEXT := strconv.Itoa(playerx)
	playeryTEXT := strconv.Itoa(playery)
	playerhTEXT := strconv.Itoa(playerh)
	playervTEXT := strconv.Itoa(playerv)
	extrascountTEXT := strconv.Itoa(extrascount)

	rl.DrawText(screenwTEXT, monw32-290, 10, 10, rl.White)
	rl.DrawText("screenw", monw32-200, 10, 10, rl.White)
	rl.DrawText(screenhTEXT, monw32-290, 20, 10, rl.White)
	rl.DrawText("screenh", monw32-200, 20, 10, rl.White)
	rl.DrawText(playerxTEXT, monw32-290, 30, 10, rl.White)
	rl.DrawText("playerx", monw32-200, 30, 10, rl.White)
	rl.DrawText(playeryTEXT, monw32-290, 40, 10, rl.White)
	rl.DrawText("playery", monw32-200, 40, 10, rl.White)
	rl.DrawText(playerhTEXT, monw32-290, 50, 10, rl.White)
	rl.DrawText("playerh", monw32-200, 50, 10, rl.White)
	rl.DrawText(playervTEXT, monw32-290, 60, 10, rl.White)
	rl.DrawText("playerv", monw32-200, 60, 10, rl.White)
	rl.DrawText(extrascountTEXT, monw32-290, 70, 10, rl.White)
	rl.DrawText("extrascount", monw32-200, 70, 10, rl.White)

}
func input() { // MARK: keys input
	if rl.IsKeyPressed(rl.KeyF4) {
		if !levelend {
			levelend = true
		}
	}
	if rl.IsKeyPressed(rl.KeySpace) {
		if !intro {
			intro = true
			pause = false
		} else {
			if !usespecial {
				usespecial = true
			}
		}
	}
	if rl.IsKeyPressed(rl.KeyUp) {
		if !jump {
			jumpcount = 5
			jump = true
		}
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		if levelmap[player+levelw] == "#" {
			player += (levelw * 2)
		}
	}

	if rl.IsKeyDown(rl.KeyRight) {
		run = true
		if playerv < levelw-((screenw/2)+1) {
			player++
			playerdir = false
		}
	}
	if rl.IsKeyReleased(rl.KeyRight) {
		run = false
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		run = true
		if playerv > screenw/2+1 {
			player--
			playerdir = true
		}
	}
	if rl.IsKeyReleased(rl.KeyLeft) {
		run = false
	}
	if rl.IsKeyPressed(rl.KeyF3) {
		if followplayer {
			followplayer = false
		} else {
			followplayer = true
		}
	}
	if rl.IsKeyPressed(rl.KeyKpAdd) {
		if camera.Zoom == 1.0 {
			camera.Zoom = 2.0
			camera.Target.Y = float32(float32(playery) - (float32(monitorh) / 2.3))
			camera.Target.X = float32(playerx - (monitorw / 4))
		} else if camera.Zoom == 2.0 {
			camera.Zoom = 4.0
			camera.Target.Y = float32(playery - (monitorh / 5))
			camera.Target.X = float32(playerx - (monitorw / 8))
		}
	}
	if rl.IsKeyPressed(rl.KeyKpSubtract) {
		if camera.Zoom == 2.0 {
			camera.Zoom = 1.0
			camera.Target.Y = 0
			camera.Target.X = 0
		} else if camera.Zoom == 4.0 {
			camera.Zoom = 2.0
			camera.Target.Y = float32(float32(playery) - (float32(monitorh) / 2.3))
			camera.Target.X = float32(playerx - (monitorw / 4))
		}
	}
	if rl.IsKeyDown(rl.KeyKp8) {
		if drawblocknexth > 0 {
			drawblocknext -= levelw
		}
	}
	if rl.IsKeyDown(rl.KeyKp2) {
		if drawblocknexth < levelh-(screenh+1) {
			drawblocknext += levelw
		}
	}
	if rl.IsKeyDown(rl.KeyKp6) {
		if drawblocknextv < levelw-(screenw+1) {
			drawblocknext++
		}
	}
	if rl.IsKeyDown(rl.KeyKp4) {
		if drawblocknextv > 0 {
			drawblocknext--
		}
	}

	if rl.IsKeyPressed(rl.KeyF1) {
		if grid16on {
			grid16on = false
		} else {
			grid16on = true
		}
	}
	if rl.IsKeyPressed(rl.KeyF2) {
		if grid4on {
			grid4on = false
		} else {
			grid4on = true
		}
	}
	if rl.IsKeyPressed(rl.KeyKpDecimal) {
		if debugon {
			debugon = false
		} else {
			debugon = true
		}
	}

}
func grid16() { // MARK: grid16()
	for a := 0; a < monitorw; a += 16 {
		a32 := int32(a)
		rl.DrawLine(a32, 0, a32, monh32, rl.Fade(rl.Green, 0.1))
	}
	for a := 0; a < monitorh; a += 16 {
		a32 := int32(a)
		rl.DrawLine(0, a32, monw32, a32, rl.Fade(rl.Green, 0.1))
	}
}
func grid4() { // MARK: grid4()
	for a := 0; a < monitorw; a += 4 {
		a32 := int32(a)
		rl.DrawLine(a32, 0, a32, monh32, rl.Fade(rl.DarkGreen, 0.1))
	}
	for a := 0; a < monitorh; a += 4 {
		a32 := int32(a)
		rl.DrawLine(0, a32, monw32, a32, rl.Fade(rl.DarkGreen, 0.1))
	}
}

// random numbers
func rInt(min, max int) int {
	return rand.Intn(max-min) + min
}
func rInt32(min, max int) int32 {
	a := int32(rand.Intn(max-min) + min)
	return a
}
func rFloat32(min, max int) float32 {
	a := float32(rand.Intn(max-min) + min)
	return a
}
func flipcoin() bool {
	var b bool
	a := rInt(0, 10001)
	if a < 5000 {
		b = true
	}
	return b
}
func rolldice() int {
	a := rInt(1, 7)
	return a
}
