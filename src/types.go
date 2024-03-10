package main

import "time"

type TransacaoRequest struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"` // char, byte, ?
	Descricao string `json:"descricao"`
}

type TransacaoResponse struct {
	Limite int `json:"limite"` // unsigned?
	Saldo  int `json:"saldo"`
}

type SaldoExtratoResponse struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type TransacaoExtratoResponse struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em" db:"realizada_em"`
	LimiteAtual int       `json:"-" db:"limite_atual"`
	SaldoAtual  int       `json:"-" db:"saldo_atual"`
}

type ExtratoResponse struct {
	Saldo             SaldoExtratoResponse       `json:"saldo"`
	UltimasTransacoes []TransacaoExtratoResponse `json:"ultimas_transacoes"`
}
