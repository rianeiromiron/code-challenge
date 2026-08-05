package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rianeiromiron/code-challenge/internal/api"
	"github.com/rianeiromiron/code-challenge/web"
)

func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	web.Tmpl.Execute(w, nil)
}

func ClimaHandler(w http.ResponseWriter, r *http.Request) {
	pais := r.URL.Query().Get("pais")
	if pais == "" {
		http.Error(w, `{"error":"parámetro 'pais' requerido"}`, http.StatusBadRequest)
		return
	}

	limit := 1
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	resultado, err := api.ObtenerClima(ctx, pais, limit)
	if err != nil {
		code := http.StatusBadGateway
		if strings.Contains(err.Error(), "país no encontrado") {
			code = http.StatusNotFound
		}
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}
