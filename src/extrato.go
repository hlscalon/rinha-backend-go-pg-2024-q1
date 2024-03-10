package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

const QUERY_OBTER_TRANSACOES = "SELECT valor, tipo, descricao, realizada_em, limite_atual, saldo_atual " +
	"FROM transacao " +
	"WHERE cliente_id = $1 " +
	"ORDER BY id DESC " +
	"LIMIT 11 " // Deve pegar uma a mais para ignorar a inicial depois, se necessário

func (app *App) handleExtrato(clienteId int, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rows, err := app.conn.Query(context.Background(), QUERY_OBTER_TRANSACOES, clienteId)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ultimasTransacoes, err := pgx.CollectRows(rows, pgx.RowToStructByName[TransacaoExtratoResponse])

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if len(ultimasTransacoes) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	limiteAtual := ultimasTransacoes[0].LimiteAtual
	saldoAtual := ultimasTransacoes[0].SaldoAtual

	// Sempre remove a última transação
	// Se tem 11 transações, remove e fica com 10 (nenhuma sendo a inicial)
	// Se tem menos, remove a última, que será sempre o saldo inicial
	ultimasTransacoes = ultimasTransacoes[:len(ultimasTransacoes)-1]

	response := &ExtratoResponse{
		Saldo: SaldoExtratoResponse{
			DataExtrato: time.Now(),
			Limite:      limiteAtual,
			Total:       saldoAtual,
		},
		UltimasTransacoes: ultimasTransacoes,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}
