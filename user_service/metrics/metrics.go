package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    RegisterSuccess = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "register_success_total",
            Help: "Number of successful user registrations",
        },
    )
    LoginSuccess = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "login_success_total",
            Help: "Number of successful logins",
        },
    )
    LoginFailed = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "login_failed_total",
            Help: "Number of failed login attempts",
        },
    )
)

func Init() {
    prometheus.MustRegister(RegisterSuccess, LoginSuccess, LoginFailed)
}
