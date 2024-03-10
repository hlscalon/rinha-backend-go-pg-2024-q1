package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	conn      *pgxpool.Pool
	clienteId int
}

func (app *App) handleCliente(w http.ResponseWriter, r *http.Request) {
	// Tratar os seguintes casos:
	// 	- /clientes/{id}/transacoes
	// 	- /clientes/{id}/extrato

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro parsing: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Limitação do projeto: somente tem 5 clientes
	// Performance: paramos aqui
	if id < 1 || id > 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	app.clienteId = id

	if parts[3] == "transacoes" {
		app.handleTransacao(w, r)
		return
	} else if parts[3] == "extrato" {
		app.handleExtrato(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func createDB() (conn *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	conn, err = pgxpool.New(context.Background(), dsn)

	if err != nil {
		return nil, err
	}

	retries := 1
	for {
		fmt.Println("Tentativa conectar banco: ", retries)

		err = conn.Ping(context.Background())

		if err == nil || retries >= 20 {
			break
		}

		// Espera 5 segs e tenta novamente
		time.Sleep(time.Second * 5)
		retries++
	}

	return
}

func main() {
	conn, err := createDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao conectar com o banco de dados: %v\n", err)
		return
	}

	defer conn.Close()

	app := &App{
		conn: conn,
	}
	http.HandleFunc("/clientes/", app.handleCliente)

	serverPort := os.Getenv("SERVER_PORT")
	fmt.Printf("Servidor rodando na porta %s...\n", serverPort)

	err = http.ListenAndServe(":"+serverPort, nil)

	if err != nil {
		fmt.Println("Erro ao iniciar servidor")
	}
}
