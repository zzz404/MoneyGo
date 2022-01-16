package deposit

import (
	"fmt"

	bk "github.com/zzz404/MoneyGo/internal/bank"
	"github.com/zzz404/MoneyGo/internal/coin"
	mb "github.com/zzz404/MoneyGo/internal/member"
	ut "github.com/zzz404/MoneyGo/internal/utils"
)

type timeDepositController struct {
}

var TimeDepController = new(timeDepositController)

func (c *timeDepositController) List(r *ut.HttpRequest, w *ut.HttpResponse) {
	form := new(QueryForm)
	err := form.ReadParameters(r)
	if w.ResponseForError(err) {
		return
	}

	tds, err := TimeDepService.Query(form)
	if w.ResponseForError(err) {
		return
	}

	totalYearIncome := 0.0
	for _, td := range tds {
		income := td.EspectedYearIncome()
		if income != nil {
			totalYearIncome += *td.EspectedYearIncome()
		}
	}

	tpl, err := ut.GetTemplate("/timeDepositList.html")
	if w.ResponseForError(err) {
		return
	}
	err = tpl.Execute(w, map[string]interface{}{
		"form":            form,
		"members":         mb.Members,
		"banks":           bk.Banks,
		"coinTypes":       coin.CoinTypes,
		"tds":             tds,
		"count":           len(tds),
		"totalYearIncome": fmt.Sprintf("%.2f", totalYearIncome),
	})
	w.ResponseForError(err)
}
