package Structs

type Status struct {
	Status          byte
	BeatmapUpdate   bool
	StatusText      string
	BeatmapChecksum string
	CurrentMods     uint16
}
