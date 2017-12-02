package crossword

import "sort"

const EMPTY_CHAR = ' '
const EMPTY_INDEX = 0

type Cell struct {
	char  rune
	index int
	value rune
}

type Grid [][]Cell

type Word struct {
	Word string
	Clue string
}

type Words []Word

func (w Words) Len() int           { return len(w) }
func (w Words) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w Words) Less(i, j int) bool { return len(w[i].Word) > len(w[j].Word) }

type ActiveWord struct {
	word     Word
	x        int
	y        int
	vertical bool
	number   int
}

type Crossword struct {
	cols           int
	rows           int
	words          Words
	grid           Grid
	activeWordList []ActiveWord
	downCount      int
	acrossCount    int
}

func New(cols, rows int, words Words) *Crossword {
	var grid Grid
	grid = make([][]Cell, rows)
	for x := 0; x < rows; x++ {
		grid[x] = make([]Cell, cols)
		for y := 0; y < cols; y++ {
			grid[x][y] = Cell{
				EMPTY_CHAR,
				EMPTY_INDEX,
				EMPTY_CHAR,
			}
		}
	}

	sort.Sort(words)

	return &Crossword{
		cols,
		rows,
		words,
		grid,
		make([]ActiveWord, 0),
		0,
		0,
	}
}

func (c *Crossword) Generate(seed int, loops int) {
	//manually place the longest word horizontally at 0,0, try others if the generated board is too weak
	c.placeWord(c.words[seed], c.rows/2, 0, false)
	c.generate(seed, loops)
}

func (c *Crossword) generate(seed int, loops int) {

	//attempt to fill the rest of the board
	for iy := 0; iy < loops; iy++ { //usually 2 times is enough for max fill potential
		for ix := 1; ix < len(c.words); ix++ {
			if !c.isActiveWord(c.words[ix].Word) { //only add if not already in the active word list
				topScore := 0
				bestScoreIndex := 0
				fitScore := 0

				coordList := c.suggestCoords([]rune(c.words[ix].Word)) //fills coordList and coordCount

				if len(coordList) > 0 {
					//coordList = shuffleArray(coordList)     //adds some randomization
					for cl := 0; cl < len(coordList); cl++ { //get the best fit score from the list of possible valid coordinates
						fitScore = c.checkFitScore([]rune(c.words[ix].Word), coordList[cl].x, coordList[cl].y, coordList[cl].vertical)
						if fitScore > topScore {
							topScore = fitScore
							bestScoreIndex = cl
						}
					}
				}

				if topScore > 1 { //only place a word if it has a fitscore of 2 or higher
					c.placeWord(c.words[ix], coordList[bestScoreIndex].x, coordList[bestScoreIndex].y, coordList[bestScoreIndex].vertical)
				}
			}
		}
	}
}

func (c *Crossword) placeWord(w Word, x, y int, vertical bool) bool { //places a new active word on the board

	wordPlaced := false

	word := []rune(w.Word)
	l := len(word)

	if vertical {
		if l+x < c.rows {
			for i := 0; i < l; i++ {
				c.grid[x+i][y].char = word[i]
			}
			wordPlaced = true
		}
	} else {
		if l+y < c.cols {
			for i := 0; i < l; i++ {
				c.grid[x][y+i].char = word[i]
			}
			wordPlaced = true
		}
	}

	if wordPlaced {
		number := 0
		if vertical {
			c.downCount++
			number = c.downCount
		} else {
			c.acrossCount++
			number = c.acrossCount
		}

		aw := ActiveWord{
			w,
			x,
			y,
			vertical,
			number,
		}

		c.activeWordList = append(c.activeWordList, aw)
	}
	return wordPlaced
}

func (c *Crossword) isActiveWord(word string) bool {
	l := len(c.activeWordList)
	for w := 0; w < l; w++ {
		if word == c.activeWordList[w].word.Word {
			return true
		}
	}
	return false
}

type coord struct {
	x        int
	y        int
	score    int
	vertical bool
}

func (c *Crossword) suggestCoords(word []rune) []coord { //search for potential cross placement locations
	coordList := make([]coord, 0)
	coordCount := 0
	for i := 0; i < len(word); i++ { //cycle through each character of the word
		ch := word[i]
		for x := 0; x < c.rows; x++ {
			for y := 0; y < c.cols; y++ {
				if c.grid[x][y].char == ch { //check for letter match in cell
					if x-i+1 > 0 && x-i+len(word)-1 < c.rows { //would fit vertically?
						coordList = append(coordList, coord{
							x - i,
							y,
							0,
							true,
						})
						coordCount++
					}

					if y-i+1 > 0 && y-i+len(word)-1 < c.cols { //would fit horizontally?
						coordList = append(coordList, coord{
							x,
							y - i,
							0,
							false,
						})
						coordCount++
					}
				}
			}
		}
	}
	return coordList
}

func (c *Crossword) checkFitScore(word []rune, x, y int, vertical bool) int {
	fitScore := 1 //default is 1, 2+ has crosses, 0 is invalid due to collision

	if vertical { //vertical checking
		for i := 0; i < len(word); i++ {
			if i == 0 && x > 0 { //check for empty space preceeding first character of word if not on edge
				if c.grid[x-1][y].char != EMPTY_CHAR { //adjacent letter collision
					fitScore = 0
					break
				}
			} else if i == len(word) && x < c.rows { //check for empty space after last character of word if not on edge
				if c.grid[x+i+1][y].char != EMPTY_CHAR { //adjacent letter collision
					fitScore = 0
					break
				}
			}
			if x+i < c.rows {
				if c.grid[x+i][y].char == word[i] { //letter match - aka cross point
					fitScore += 1
				} else if c.grid[x+i][y].char != EMPTY_CHAR { //letter doesn't match and it isn't empty so there is a collision
					fitScore = 0
					break
				} else { //verify that there aren't letters on either side of placement if it isn't a crosspoint
					if y < c.cols-1 { //check right side if it isn't on the edge
						if c.grid[x+i][y+1].char != EMPTY_CHAR { //adjacent letter collision
							fitScore = 0
							break
						}
					}
					if y > 0 { //check left side if it isn't on the edge
						if c.grid[x+i][y-1].char != EMPTY_CHAR { //adjacent letter collision
							fitScore = 0
							break
						}
					}
				}
			}

		}

	} else { //horizontal checking
		for i := 0; i < len(word); i++ {
			if i == 0 && y > 0 { //check for empty space preceeding first character of word if not on edge
				if c.grid[x][y-1].char != EMPTY_CHAR { //adjacent letter collision
					fitScore = 0
					break
				}
			} else if i == len(word)-1 && y+i < c.cols-1 { //check for empty space after last character of word if not on edge
				if c.grid[x][y+i+1].char != EMPTY_CHAR { //adjacent letter collision
					fitScore = 0
					break
				}
			}
			if y+i < c.cols {
				if c.grid[x][y+i].char == word[i] { //letter match - aka cross point
					fitScore++
				} else if c.grid[x][y+i].char != EMPTY_CHAR { //letter doesn't match and it isn't empty so there is a collision
					fitScore = 0
					break
				} else { //verify that there aren't letters on either side of placement if it isn't a crosspoint
					if x < c.rows-1 { //check top side if it isn't on the edge
						if c.grid[x+1][y+i].char != EMPTY_CHAR { //adjacent letter collision
							fitScore = 0
							break
						}
					}
					if x > 0 { //check bottom side if it isn't on the edge
						if c.grid[x-1][y+i].char != EMPTY_CHAR { //adjacent letter collision
							fitScore = 0
							break
						}
					}
				}
			}
		}
	}

	return fitScore
}

func (c *Crossword) String() string {
	result := ""
	for i := 0; i < len(c.grid); i++ {
		for j := 0; j < len(c.grid[i]); j++ {
			result = result + string(c.grid[i][j].char)
		}
		result = result + "\n"
	}
	return result
}
