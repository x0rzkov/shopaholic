package tg_bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"strings"
)

type HistoryCommand struct {
	Options
}

func (c *HistoryCommand) Execute(m *tb.Message) error {
	log.Printf("[INFO] bot is preparing history")

	user, err := c.Store.Details(strconv.Itoa(m.Sender.ID))
	if err != nil {
		return err
	}

	transactions, err := c.Store.List(user)
	if err != nil {
		return err
	}

	results := []string{}
	for _, transaction := range transactions {
		amount := float64(transaction.Amount.Amount / 100)
		balanceWas := float64(transaction.BalanceWas.Amount / 100)
		balanceNow := float64(transaction.BalanceNow.Amount / 100)
		time := transaction.CreatedAt.Format("02.01.2006 15:04")
		result := fmt.Sprintf("🕰 %s. %.2f💲 at %s. Balance %.2f$ → %.2f$",
			time, amount, transaction.Category.Title, balanceWas, balanceNow)
		results = append(results, result)
	}

	if len(transactions) == 0 {
		results = append(results, "There no transactions for you yet🤷‍‍", user.Name)
	}

	_, err = c.Bot.Send(m.Sender, strings.Join(results, "\n"))
	return err
}
