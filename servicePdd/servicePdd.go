package servicePdd

import (
	"math/rand"
)

const lenBlock = 200

const (
	ticket            = iota
	addition_question = iota
	result            = iota
)

type DataTicket struct {
	Mode          int64   `json:"mode"`
	Ticket        int64   `json:"ticket"`
	Question      int64   `json:"question"`
	Answer        int64   `json:"answer"`
	CorrectAnswer int64   `json:"correct_answer"`
	BlockBad      []int64 `json:"block_bad"`
	BlockAdd      []int64 `json:"block_add"`
	AdditionQues  []int64 `json:"addition_ques"`
}

func _modB(res *DataTicket) {
	if res.AdditionQues == nil {
		res.Mode = result
	} else {
		lenlistQ := len(res.AdditionQues)
		res.Ticket = res.AdditionQues[lenlistQ-1]/20 + 1
		res.Question = res.AdditionQues[lenlistQ-1] - (res.Ticket-1)*20
		res.AdditionQues = res.AdditionQues[:lenlistQ-1]

	}
}

func NextQuestion(res DataTicket) DataTicket {
	if res.Mode == ticket {
		if res.Answer != res.CorrectAnswer {
			res.BlockBad = append(res.BlockBad, (res.Ticket*20 + res.Question))
			randomIndex := rand.Intn(lenBlock)
			switch res.Question {
			case 1, 2, 3, 4, 5:
				res.AdditionQues = append(res.AdditionQues, block1[randomIndex])
			case 6, 7, 8, 9, 10:
				res.AdditionQues = append(res.AdditionQues, block2[randomIndex])
			case 11, 12, 13, 14, 15:
				res.AdditionQues = append(res.AdditionQues, block3[randomIndex])
			case 16, 17, 18, 19, 20:
				res.AdditionQues = append(res.AdditionQues, block4[randomIndex])
			}
		}
		if res.Question < 20 {
			res.Question += 1
		} else if res.Question == 20 {
			res.Mode = addition_question
			_modB(&res)
		}
	} else if res.Mode == addition_question {
		if res.Answer != res.CorrectAnswer {
			res.BlockBad = append(res.BlockAdd, (res.Ticket*20 + res.Question))
		}
		_modB(&res)
	}
	return res
}

// func Test() ResultQuestion {
// 	return ResultQuestion{Mode: 0, Ticket: 1, Question: 13, Answer: 2, CorrectAnswer: 1, BlockBad: nil, BlockAdd: nil, AdditionQues: nil}
// }
