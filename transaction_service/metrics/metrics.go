package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TransactionCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "transaction_created_total",
			Help: "Number of transactions created",
		},
	)
)

func Init() {
	prometheus.MustRegister(TransactionCreated)
}
