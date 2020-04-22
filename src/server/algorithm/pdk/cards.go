package pdk

type Cards []byte

type Card byte

func (c Card) getCardIndex() int{
	cardValue := c & 0x0f
	return int(cardValue)
}



