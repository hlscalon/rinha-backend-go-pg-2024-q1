package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const QUERY_OBTER_TRANSACOES = "SELECT valor, tipo, descricao, realizada_em, limite_atual, saldo_atual " +
	"FROM transacao " +
	"WHERE cliente_id = $1 " +
	"ORDER BY id DESC " +
	"LIMIT 10 "

func handleExtrato(clienteId int, conn *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rows, err := conn.Query(context.Background(), QUERY_OBTER_TRANSACOES, clienteId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro extrato: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ultimasTransacoes, err := pgx.CollectRows(rows, pgx.RowToStructByName[TransacaoExtratoResponse])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro extrato transacoes: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if len(ultimasTransacoes) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	response := &ExtratoResponse{
		Saldo: SaldoExtratoResponse{
			DataExtrato: time.Now(),
			Limite:      ultimasTransacoes[0].LimiteAtual,
			Total:       ultimasTransacoes[0].SaldoAtual,
		},
		UltimasTransacoes: ultimasTransacoes,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro encoding: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}
