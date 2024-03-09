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

// var conn *pgx.Conn
var conn *pgxpool.Pool

func handleCliente(w http.ResponseWriter, r *http.Request) {
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

	if parts[3] == "transacoes" {
		handleTransacao(id, conn, w, r)
		return
	} else if parts[3] == "extrato" {
		handleExtrato(id, conn, w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	var err error
	for {
		retries := 1

		fmt.Println("Tentativa conectar banco: ", retries)

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
		// conn, err = pgx.Connect(context.Background(), dsn)
		conn, err = pgxpool.New(context.Background(), dsn)

		if err == nil || retries > 5 {
			break
		}

		// Espera 3 segs e tenta novamente
		time.Sleep(time.Second * 3)
		retries += 1
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao conectar com o banco de dados: %v\n", err)
		return
	}
	// defer conn.Close(context.Background())
	defer conn.Close()

	http.HandleFunc("/clientes/", handleCliente)

	fmt.Println("Servidor rodando na porta 9000...")

	err = http.ListenAndServe(":9000", nil)

	if err != nil {
		fmt.Println("Erro ao iniciar servidor")
	}
}
