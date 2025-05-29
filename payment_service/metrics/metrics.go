package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	PaymentCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "payment_created_total",
			Help: "Number of payments created",
		},
	)
)

func Init() {
	prometheus.MustRegister(PaymentCreated)
}
