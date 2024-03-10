package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const QUERY_CREDITAR = "UPDATE cliente SET saldo = saldo + $1 WHERE id = $2 RETURNING nome, saldo, limite"
const QUERY_DEBITAR = "UPDATE cliente SET saldo = saldo - $1 WHERE id = $2 AND saldo - $3 >= -ABS(limite) RETURNING nome, saldo, limite"

const QUERY_INSERIR_TRANSACAO_CREDITAR = "WITH cliente_atualizado AS (%s) " +
	"INSERT INTO transacao (cliente_id, valor, tipo, descricao, limite_atual, saldo_atual) " +
	"SELECT $3, $4, $5, $6, cliente_atualizado.limite, cliente_atualizado.saldo " +
	"FROM cliente_atualizado " +
	"RETURNING limite_atual, saldo_atual"

const QUERY_INSERIR_TRANSACAO_DEBITAR = "WITH cliente_atualizado AS (%s) " +
	"INSERT INTO transacao (cliente_id, valor, tipo, descricao, limite_atual, saldo_atual) " +
	"SELECT $4, $5, $6, $7, cliente_atualizado.limite, cliente_atualizado.saldo " +
	"FROM cliente_atualizado " +
	"RETURNING limite_atual, saldo_atual"

func (app *App) handleTransacao(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	transacao := &TransacaoRequest{}
	err := json.NewDecoder(r.Body).Decode(transacao)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if len(transacao.Descricao) < 1 || len(transacao.Descricao) > 10 || (transacao.Tipo != "c" && transacao.Tipo != "d") {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var updateQuery string
	var argsQuery []any
	if transacao.Tipo == "c" {
		updateQuery = QUERY_CREDITAR
		argsQuery = append(argsQuery, transacao.Valor, app.clienteId)
		updateQuery = fmt.Sprintf(QUERY_INSERIR_TRANSACAO_CREDITAR, updateQuery)
	} else {
		updateQuery = QUERY_DEBITAR
		argsQuery = append(argsQuery, transacao.Valor, app.clienteId, transacao.Valor)
		updateQuery = fmt.Sprintf(QUERY_INSERIR_TRANSACAO_DEBITAR, updateQuery)
	}

	argsQuery = append(argsQuery, app.clienteId, transacao.Valor, transacao.Tipo, transacao.Descricao)

	var novoLimite int
	var novoSaldo int
	err = app.conn.QueryRow(context.Background(), updateQuery, argsQuery...).Scan(&novoLimite, &novoSaldo)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro transacao: %v\n", err)
		fmt.Fprintf(os.Stderr, "args: %v\n", argsQuery)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	response := &TransacaoResponse{
		Saldo:  novoSaldo,
		Limite: novoLimite,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro encoding: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}
