package game

import (
	"image/color"

	"github.com/z46-dev/game-dev-project/util"
	"golang.org/x/image/colornames"
)

var factionColors []color.RGBA = []color.RGBA{
	colornames.Red, colornames.Blue, colornames.Green, colornames.Yellow,
	colornames.Purple, colornames.Orange, colornames.Cyan, colornames.Magenta,
	colornames.Lime, colornames.Pink, colornames.Teal, colornames.Violet,
	colornames.Gold, colornames.Silver, colornames.Brown, colornames.Maroon,
}

func NewFaction(g *Game, name string) (f *Faction) {
	g.nextFactionID++
	f = &Faction{
		ID:               g.nextFactionID,
		Name:             name,
		Color:            factionColors[int(g.nextFactionID)%len(factionColors)],
		ShipsSpatialHash: util.NewSpatialHash[*Ship](),
	}

	g.FactionsMu.Lock()
	g.Factions[f.ID] = f
	g.FactionsMu.Unlock()
	return
}

func (f *Faction) Update() {
	f.ShipsSpatialHash.Clear()
}
