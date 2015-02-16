package lang

////////// TokenLocation

type TokenLocation struct {
	Pos      int // position within the file
	Line     int
	Filename string
}
