package brot

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/teadove/awesome-fractals/internal/palette"

	"github.com/stretchr/testify/assert"
	"github.com/teadove/awesome-fractals/internal/utils"
)

func TestUnit_Brot_Render_Ok(t *testing.T) {
	t.Parallel()

	file, err := os.ReadFile("test_image.png")
	utils.Check(err)

	image, err := png.Decode(bytes.NewReader(file))
	utils.Check(err)

	input, _ := NewInputFromBuilderInput(&BuilderInput{
		Colors:       palette.ColorPalettes["Plan9"],
		ColorStep:    6000,
		XPos:         -0.00275,
		YPos:         0.78912,
		Width:        1000,
		Height:       1000,
		MaxIteration: 800,
		EscapeRadius: .125689,
	})
	imageCompile := Render(input)

	assert.Equal(t, image, imageCompile)
}
