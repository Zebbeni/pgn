package pgn

import (
	"errors"
	"fmt"
	"strings"
)

type FEN struct {
	FOR                 string
	ToMove              Color
	WhiteCastleStatus   CastleStatus
	BlackCastleStatus   CastleStatus
	EnPassantVulnerable Position
	HalfmoveClock       int
	Fullmove            int
}

func ParseFEN(fenstr string) (*FEN, error) {
	fen := FEN{}
	colorStr := ""
	castleStr := ""
	enPassant := ""
	_, err := fmt.Sscanf(fenstr, "%s %s %s %s %d %d",
		&fen.FOR,
		&colorStr,
		&castleStr,
		&enPassant,
		&fen.HalfmoveClock,
		&fen.Fullmove,
	)
	if err != nil {
		return nil, err
	}
	switch colorStr {
	case "w":
		fen.ToMove = White
	case "b":
		fen.ToMove = Black
	default:
		return nil, errors.New("pgn: invalid color")
	}

	if strings.Contains(castleStr, "k") {
		fen.BlackCastleStatus = Kingside
	}
	if strings.Contains(castleStr, "q") {
		if fen.BlackCastleStatus == Kingside {
			fen.BlackCastleStatus = Both
		} else {
			fen.BlackCastleStatus = Queenside
		}
	}

	if strings.Contains(castleStr, "K") {
		fen.WhiteCastleStatus = Kingside
	}
	if strings.Contains(castleStr, "Q") {
		if fen.WhiteCastleStatus == Kingside {
			fen.WhiteCastleStatus = Both
		} else {
			fen.WhiteCastleStatus = Queenside
		}
	}

	if enPassant == "-" {
		fen.EnPassantVulnerable = NoPosition
	} else {
		fen.EnPassantVulnerable, err = ParsePosition(enPassant)
		if err != nil {
			return nil, err
		}
	}
	return &fen, nil
}

func FORFromBoard(b *Board) string {
	f := ""
	for y := '8'; y > '0'; y-- {
		count := 0
		for x := 'a'; x <= 'h'; x++ {
			pos, _ := ParsePosition(fmt.Sprintf("%c%c", x, y))
			p := b.GetPiece(pos)
			if p == Empty {
				count++
			} else {
				if count > 0 {
					f += fmt.Sprintf("%d", count)
					count = 0
				}
				f += string(p)
			}
		}
		if count > 0 {
			f += fmt.Sprintf("%d", count)
		}
		if y != '1' {
			f += "/"
		}
	}
	return f
}

func FENFromBoard(b *Board) FEN {
	f := FEN{}
	f.FOR = FORFromBoard(b)
	f.ToMove = b.toMove
	f.WhiteCastleStatus = b.wCastle
	f.BlackCastleStatus = b.bCastle
	f.HalfmoveClock = b.halfmoveClock
	f.Fullmove = b.fullmove
	return f
}

func (fen FEN) String() string {
	castleStatus := fen.WhiteCastleStatus.String(White) + fen.BlackCastleStatus.String(Black)
	if castleStatus == "--" {
		castleStatus = "-"
	}
	return fmt.Sprintf("%s %v %s %s %d %d",
		fen.FOR,
		fen.ToMove,
		castleStatus,
		fen.EnPassantVulnerable.String(),
		fen.HalfmoveClock,
		fen.Fullmove,
	)
}
