package pdk

type Cards []Card

type Card byte

type CardType struct {
	Type uint8
	OriginCards []Card
	MainCards []Card `主牌型比如555`
	ExtraCards []Card `带的牌  比如三代二 二`
	Count int
}

var CardTypeError = CardType{Type: TYPE_ERROR}

var StraightMask = []uint16{15872, 7936, 3968, 1984, 992, 496, 248, 124, 62, 31}

func (c Card) getCardIndex() int{
	cardValue := c & 0x0f
	//if cardValue == 0x01 {
	//	return 14
	/*}else*/ if cardValue == 0x02{
		return 16
	}else if cardValue == 0x0e{ //小王
		return 17
	}else if cardValue == 0x0f{ //大王
		return 18
	}
	return int(cardValue)
}

func (c *Cards) getType() CardType{
	tempCards := Cards{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}  //-a-k  大小王  15张序列
	for _,v := range *c {
		vIndex := v.getCardIndex()
		tempCards[vIndex] = tempCards[vIndex] + 1 //
	}

	counts := []int{0,0,0,0,0} //--1,2,3,4组合的数目
	cards := []Cards{{},{},{},{},{}}         //数目分别为1,2,3,4的牌的序列
	for i,v := range tempCards{
		if v > 0 {
			counts[v] = counts[v] +1
			cards[v] = append(cards[v], Card(i))
		}
	}

	if counts[4] > 0 {
		return getType4(counts,cards)
	}else if counts[3] > 0 {
		return getType3(counts,cards)
	}else if counts[2] > 0 {
		return getType2(counts,cards)
	}else if counts[1] > 0{
		return getType1(counts,cards)
	}
	return CardTypeError
}

func getType4(counts []int,cards []Cards) CardType{
	if counts[4] >1{
		return CardTypeError
	}

	//飞机
	if counts[3] > 1{
		if !isContinue(cards[3]) {
			return CardTypeError
		}
		if counts[3] == 4 {
			if counts[2] + counts[1] != 0 {
				return CardTypeError
			}
			return CardType{Type: TYPE_FEI_JI_1,MainCards: cards[3],Count: counts[3]}
		}else if counts[3] == 2 + counts[2] && counts[1] == 0 {
			return CardType{Type: TYPE_FEI_JI_2,MainCards: cards[3],Count: counts[3]}
		}
		return CardTypeError
	}

	sum := counts[3] + counts[2] + counts[1]
	if sum > 1 {
		return CardTypeError
	}

	//card := cards[4][0]
	if sum == 0 {
		return CardType{Type: TYPE_SI_GE,MainCards: cards[4]}
	}
	if counts[3] ==1 {
		return CardType{Type: TYPE_SI_GE_3,MainCards: cards[4]}
	}else if counts[2] == 1 {
		return CardType{Type: TYPE_SI_GE_2,MainCards: cards[4]}
	}else if counts[1] == 1 {
		return CardType{Type: TYPE_SI_GE_1,MainCards: cards[4]}
	}

	return CardTypeError
}

func getType3(counts []int,cards []Cards) CardType{
	//card := cards[3][len(cards[3])-1] //结尾的牌是多少
	count3 := counts[3]
	count2 := counts[2]
	count1 := counts[1]
	sum := count3 + count2 + count1
	//三带一，或三带二，三个不带
	if count3 == 1 {
		if count2 ==1 && count1 == 0 {
			return CardType{Type: TYPE_SAN_GE_2,MainCards: cards[3]}
		}else if count2 == 0 && count1 == 1 {
			return CardType{Type: TYPE_SAN_GE_1,MainCards: cards[3]}
		}else if count2 == 0 && count1 == 0 {
			return CardType{Type: TYPE_SAN_GE,MainCards: cards[3]}
		}
	}
	if isContinue(cards[3]) {
		if sum%4 == 0 {
			return CardType{Type: TYPE_FEI_JI_1,MainCards: cards[3]}
		}else if sum % 5 == 0 && count1 == 0 {
			return CardType{Type: TYPE_FEI_JI_2,MainCards: cards[3]}
		}else if count2 == 0 && count1 == 0 {
			return CardType{Type: TYPE_FEI_JI,MainCards: cards[3]}
		}
	}
	//--如 444555666999   三带一飞机
	if count3 == 4 {
		cardsNotFirst := cards[3][1:4]
		cardsNotLast := cards[3][:3]
		if isContinue(cardsNotFirst) || isContinue(cardsNotLast) {
			return CardType{Type: TYPE_FEI_JI_1,MainCards: cards[3]}
		}
	}
	return CardTypeError
}

func getType2(counts []int,cards []Cards) CardType{
	if counts[1] > 0 {
		return CardTypeError
	}
	//card := cards[2][0]
	if counts[2] == 1 {
		return CardType{Type: TYPE_DUI_Zi,MainCards: cards[2]}
	}
	if !isContinue(cards[2]){
		return CardTypeError
	}

	return CardType{Type: TYPE_LIAN_DUI,MainCards: cards[2]}
}

func getType1(counts []int,cards []Cards) CardType{
	count := counts[1]
	if count < 5 && count != 2{
		return CardTypeError
	}
	if !isContinue(cards[1]){
		return CardTypeError
	}

	if count == 2 {
		//王炸
		if cards[1][0] == 0x5e && cards[1][1] == 0x5f {
			return CardType{Type: TYPE_DOUBLE_KING,MainCards: cards[1]}
		}
	}
	return CardType{Type: TYPE_SHUN_Zi,MainCards: cards[1]}
}

func isContinue(cards Cards) bool{
	var card Card
	for _,v := range cards{
		if card > 0 && card + 1 != v {
			return  false
		}
		card = v
	}
	return true
}

//func (this *Cards) straight() uint32 {
//	var handvalue uint16
//	for _, v := range (*this) {
//		value := v & 0xF
//		if value == 0xE {
//			handvalue |= 1
//		}
//		handvalue |= (1 << (value - 1 ) )
//	}
//
//	for i := uint8(0); i < 10; i++ {
//		if handvalue&StraightMask[i] == StraightMask[i] {
//			return En(5, uint32(10-i+4))
//		}
//	}
//	return 0
//}
//
//func De(v uint32) (uint8, uint32) {
//	return uint8(v >> 24), v & 0xFFFFFF
//}
//
//func En(t uint8, v uint32) uint32 {
//	v1 := v | ( uint32(t) << 24)
//	return v1
//}


