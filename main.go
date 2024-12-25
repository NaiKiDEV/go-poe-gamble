package main

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	NORMAL = iota + 1
	MAGIC
	RARE
	UNIQUE
)

const (
	ORB_AUG    = "ORB_AUG"
	ORB_TRANS  = "ORB_TRANS"
	ORB_REGAL  = "ORB_REGAL"
	ORB_ALCH   = "ORB_ALCH"
	ORB_DIVINE = "ORB_DIVINE"
	ORB_EXALT  = "ORB_EXALT"
	ORB_CHAOS  = "ORB_CHAOS"
)

const (
	AFFIX_FIRE_RES                = "AFFIX_FIRE_RES"
	AFFIX_COLD_RES                = "AFFIX_COLD_RES"
	AFFIX_LIGHTNING_RES           = "AFFIX_LIGHTNING_RES"
	AFFIX_ALL_RES                 = "AFFIX_ALL_RES"
	AFFIX_CHAOS_RES               = "AFFIX_CHAOS_RES"
	AFFIX_EVASION_RATING          = "AFFIX_EVASION_RATING"
	AFFIX_ARMOUR                  = "AFFIX_ARMOUR"
	AFFIX_ENERGY_SHIELD           = "AFFIX_ENERGY_SHIELD"
	AFFIX_INCREASED_ENERGY_SHIELD = "AFFIX_INCREASED_ENERGY_SHIELD"
)

// Renderer consts
const (
	materialMenuOffsetX int32 = 50
	materialMenuOffsetY int32 = 50
	materialBoxSize     int32 = 54
	materialImageSize   int32 = materialBoxSize - 4
	materialFontSize    int32 = 12
	materialMenuGap     int32 = 4
)

type Material struct {
	Name        string
	Type        string
	Description string
	Rarity      int8
	Amount      int32
}

type AffixPoolItem struct {
	Name string
	Min  int32
	Max  int32
	Tier int32
}

type ItemAffix struct {
	PoolItem    AffixPoolItem
	Description string
	MinValue    int32
	MaxValue    int32
}

type Item struct {
	Name    string
	Affixes []ItemAffix
}

type GameState struct {
	Materials   map[string]Material
	Item        *Item
	SelectedOrb string
}

func getSizeGapPerElementOffset(size, gap, element int32) int32 {
	return element*size + element*gap
}

var (
	augmentationTex  rl.Texture2D
	transmutationTex rl.Texture2D
	regalTex         rl.Texture2D
	alchemyTex       rl.Texture2D
	chaosTex         rl.Texture2D
	exaltedTex       rl.Texture2D
	divineTex        rl.Texture2D
	// Radial Texture Background for Materials
	radialTex rl.Texture2D
)

func materialAmountToText(amount int32) string {
	return strconv.Itoa(int(amount))
}

func getColorFromRarity(rarity int8) color.RGBA {
	if rarity == NORMAL {
		return rl.White
	}
	if rarity == MAGIC {
		return rl.Blue
	}
	if rarity == RARE {
		return rl.Yellow
	}
	if rarity == UNIQUE {
		return rl.Orange
	}
	return rl.Gray
}

func drawMaterialBox(posX, posY int32, tex rl.Texture2D, material Material) {
	imageOffset := (materialBoxSize - materialImageSize) / 2
	rarityColor := getColorFromRarity(material.Rarity)

	// Border (active, inactive, etc.)
	rl.DrawRectangle(posX, posY, materialBoxSize, materialBoxSize, rarityColor)

	// Image of the Material
	imagePosX := posX + imageOffset
	imagePosY := posY + imageOffset
	rl.DrawRectangle(imagePosX, imagePosY, materialImageSize, materialImageSize, rl.Black)

	// Radial Background based on rarity
	rl.DrawTexture(radialTex, imagePosX, imagePosY, rarityColor)

	rl.DrawTexture(tex, imagePosX, imagePosY, rl.White)

	// Material Amount
	amountInText := materialAmountToText(material.Amount)
	amountPosX := posX + materialBoxSize - imageOffset - rl.MeasureText(amountInText, materialFontSize) - 2 // -2 pixels padding to not touch border
	amountPosY := posY + materialBoxSize - materialFontSize - 2                                             // -2 pixels padding to not touch border
	rl.DrawText(amountInText, amountPosX, amountPosY, materialFontSize, rl.White)
}

func renderMaterialMenu(gs *GameState) {
	var currentXOffset int32
	var currentYOffset int32

	// row 0
	currentYOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 0)

	currentXOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 0)
	drawMaterialBox(materialMenuOffsetX+currentXOffset, materialMenuOffsetY, transmutationTex, gs.Materials[ORB_TRANS])

	currentXOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 1)
	drawMaterialBox(materialMenuOffsetX+currentXOffset, materialMenuOffsetY, augmentationTex, gs.Materials[ORB_AUG])

	currentXOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 2)
	drawMaterialBox(materialMenuOffsetX+currentXOffset, materialMenuOffsetY, alchemyTex, gs.Materials[ORB_ALCH])

	currentXOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 3)
	drawMaterialBox(materialMenuOffsetX+currentXOffset, materialMenuOffsetY+currentYOffset, regalTex, gs.Materials[ORB_REGAL])

	// row 1
	currentYOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 1)

	currentXOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 0)
	drawMaterialBox(materialMenuOffsetX+currentXOffset, materialMenuOffsetY+currentYOffset, exaltedTex, gs.Materials[ORB_EXALT])

	currentXOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 1)
	drawMaterialBox(materialMenuOffsetX+currentXOffset, materialMenuOffsetY+currentYOffset, chaosTex, gs.Materials[ORB_CHAOS])

	currentXOffset = getSizeGapPerElementOffset(materialBoxSize, materialMenuGap, 2)
	drawMaterialBox(materialMenuOffsetX+currentXOffset, materialMenuOffsetY+currentYOffset, divineTex, gs.Materials[ORB_DIVINE])
}

func getSelectionTextureFromOrbType(orbType string) rl.Texture2D {
	return divineTex
}

func renderOrbSelectionUnderCursor(selectedOrb string) {
	mousePos := rl.GetMousePosition()

	selectedOrbTex := getSelectionTextureFromOrbType(selectedOrb)
	selectedOrbTexHalfHeight := selectedOrbTex.Height / 2

	rl.DrawTexture(
		selectedOrbTex,
		int32(mousePos.X)-selectedOrbTexHalfHeight,
		int32(mousePos.Y)-selectedOrbTexHalfHeight,
		rl.White,
	)
}

func handleInput(gs *GameState) {
	if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		// Hardcoded, because there is no system in place to detect which orb was selected with mouse position
		gs.SelectedOrb = ORB_DIVINE
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		gs.SelectedOrb = ""
	}
}

func renderFrame(gs *GameState) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.DarkGray)

	// [1] Input handling
	handleInput(gs)

	// [2] Rendering
	renderMaterialMenu(gs)

	// If player has active selection then display smaller/equal version of the orb under cursor
	if gs.SelectedOrb != "" {
		renderOrbSelectionUnderCursor(gs.SelectedOrb)
	}

	rl.EndDrawing()
}

func loadTextureFromFile(fileName string, size int32) rl.Texture2D {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("Image file %s does not exist", fileName))
	}

	img := rl.LoadImage(fileName)
	defer rl.UnloadImage(img)

	rl.ImageResize(img, size, size)

	return rl.LoadTextureFromImage(img)
}

func main() {
	rl.InitWindow(800, 460, "PoE2 Item Gambling")
	rl.SetTargetFPS(144)

	augmentationTex = loadTextureFromFile("./assets/augmentation-orb.png", materialImageSize)
	transmutationTex = loadTextureFromFile("./assets/transmutation-orb.png", materialImageSize)
	regalTex = loadTextureFromFile("./assets/regal-orb.png", materialImageSize)
	alchemyTex = loadTextureFromFile("./assets/alchemy-orb.png", materialImageSize)
	chaosTex = loadTextureFromFile("./assets/chaos-orb.png", materialImageSize)
	exaltedTex = loadTextureFromFile("./assets/exalted-orb.png", materialImageSize)
	divineTex = loadTextureFromFile("./assets/divine-orb.png", materialImageSize)

	radialImg := rl.GenImageGradientRadial(int(materialImageSize), int(materialImageSize), 0.01, rl.White, rl.Black)
	radialTex = rl.LoadTextureFromImage(radialImg)

	// Before main shuts down, do some cleanup
	defer func() {
		// General Cleanup
		rl.CloseWindow()

		// Image Unloading
		rl.UnloadImage(radialImg)

		// Texture Unloading
		rl.UnloadTexture(chaosTex)
		rl.UnloadTexture(augmentationTex)
		rl.UnloadTexture(regalTex)
		rl.UnloadTexture(alchemyTex)
		rl.UnloadTexture(exaltedTex)
		rl.UnloadTexture(transmutationTex)
		rl.UnloadTexture(radialTex)
	}()

	gs := GameState{
		Materials: map[string]Material{
			ORB_TRANS: {
				Name:        "Orb of Transmutation",
				Type:        ORB_TRANS,
				Description: "Converts NORMAL item to NORMAL one and adds an affix",
				Rarity:      MAGIC,
				Amount:      16,
			},
			ORB_AUG: {
				Name:        "Orb of Augmentation",
				Type:        ORB_AUG,
				Description: "Adds a new Random Modifier to MAGIC item",
				Rarity:      MAGIC,
				Amount:      12,
			},
			ORB_ALCH: {
				Name:        "Orb of Alchemy",
				Type:        ORB_ALCH,
				Description: "Converts normal item to RARE one with 4 affixes",
				Rarity:      RARE,
				Amount:      9,
			},
			ORB_REGAL: {
				Name:        "Regal Orb",
				Type:        ORB_AUG,
				Description: "Converts the item to RARE one and adds a Random Modifier",
				Rarity:      RARE,
				Amount:      4,
			},
			ORB_EXALT: {
				Name:        "Exalted Orb",
				Type:        ORB_EXALT,
				Description: "Adds a Random Modifier to RARE item",
				Rarity:      RARE,
				Amount:      111,
			},
			ORB_CHAOS: {
				Name:        "Chaos Orb",
				Type:        ORB_CHAOS,
				Description: "Removes random affix and replaces it with another Random Modifier",
				Rarity:      RARE,
				Amount:      3,
			},
			ORB_DIVINE: {
				Name:        "Divine Orb",
				Type:        ORB_DIVINE,
				Description: "Re-rolls all the affix values to a new ones within the range",
				Rarity:      UNIQUE,
				Amount:      0,
			},
		},
		Item: &Item{
			Name:    "Ring",
			Affixes: make([]ItemAffix, 6),
		},
		SelectedOrb: "",
	}

	for !rl.WindowShouldClose() {
		renderFrame(&gs)
	}
}
