package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	BudgetCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "budget_created_total",
			Help: "Number of budgets created",
		},
	)
)

func Init() {
	prometheus.MustRegister(BudgetCreated)
}
