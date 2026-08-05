package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestObtenerClima_Success verifica que parsea correctamente 1 resultado.
func TestObtenerClima_Success(t *testing.T) {
	// Servidor falso de geocoding
	geo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"results":[{"latitude":19.43,"longitude":-99.13}]}`)
	}))
	defer geo.Close()

	// Servidor falso de clima (devuelve 2 registros para probar el límite)
	weather := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"geometry":{"type":"Point","coordinates":[-99.13,19.43]},
			"properties":{
				"timeseries":[
					{
						"time":"2026-08-04T12:00:00Z",
						"data":{
							"instant":{
								"details":{
									"air_pressure_at_sea_level":1013.0,
									"air_temperature":25.0,
									"cloud_area_fraction":20.0,
									"relative_humidity":60.0,
									"wind_from_direction":180.0,
									"wind_speed":3.5
								}
							},
							"next_1_hours":{"summary":{"symbol_code":"clearsky_day"},"details":{"precipitation_amount":0.0}},
							"next_6_hours":{"summary":{"symbol_code":"fair_day"},"details":{"precipitation_amount":0.1}},
							"next_12_hours":{"summary":{"symbol_code":"partlycloudy_day"},"details":{"precipitation_amount":0.5}}
						}
					},
					{
						"time":"2026-08-04T13:00:00Z",
						"data":{
							"instant":{
								"details":{
									"air_pressure_at_sea_level":1012.0,
									"air_temperature":26.0,
									"cloud_area_fraction":30.0,
									"relative_humidity":55.0,
									"wind_from_direction":190.0,
									"wind_speed":4.0
								}
							}
						}
					}
				]
			}
		}`)
	}))
	defer weather.Close()

	// Cliente con URLs falsas
	client := &Client{
		GeocodingBaseURL: geo.URL,
		WeatherBaseURL:   weather.URL,
		HTTPClient:       &http.Client{},
	}

	ctx := context.Background()
	result, err := client.ObtenerClima(ctx, "Mexico", 1)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if result.Pais != "Mexico" {
		t.Errorf("pais = %s; want Mexico", result.Pais)
	}
	if result.Latitud != 19.43 {
		t.Errorf("latitud = %f; want 19.43", result.Latitud)
	}
	if len(result.Resultados) != 1 {
		t.Errorf("resultados = %d; want 1", len(result.Resultados))
	}
	if result.Resultados[0].Temperatura != 25.0 {
		t.Errorf("temperatura = %f; want 25.0", result.Resultados[0].Temperatura)
	}
	if result.Resultados[0].Simbolo1h != "clearsky_day" {
		t.Errorf("simbolo1h = %s; want clearsky_day", result.Resultados[0].Simbolo1h)
	}
}

// TestObtenerClima_Limit verifica que respeta el límite solicitado.
func TestObtenerClima_Limit(t *testing.T) {
	geo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"results":[{"latitude":10.0,"longitude":20.0}]}`)
	}))
	defer geo.Close()

	weather := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Devuelve 3 registros idénticos simplificados
		fmt.Fprintln(w, `{
			"geometry":{"type":"Point","coordinates":[20.0,10.0]},
			"properties":{
				"timeseries":[
					{"time":"T1","data":{"instant":{"details":{"air_temperature":1,"air_pressure_at_sea_level":1000,"relative_humidity":50,"wind_speed":1,"wind_from_direction":0,"cloud_area_fraction":0}}}},
					{"time":"T2","data":{"instant":{"details":{"air_temperature":2,"air_pressure_at_sea_level":1000,"relative_humidity":50,"wind_speed":1,"wind_from_direction":0,"cloud_area_fraction":0}}}},
					{"time":"T3","data":{"instant":{"details":{"air_temperature":3,"air_pressure_at_sea_level":1000,"relative_humidity":50,"wind_speed":1,"wind_from_direction":0,"cloud_area_fraction":0}}}}
				]
			}
		}`)
	}))
	defer weather.Close()

	client := &Client{
		GeocodingBaseURL: geo.URL,
		WeatherBaseURL:   weather.URL,
		HTTPClient:       &http.Client{},
	}

	result, err := client.ObtenerClima(context.Background(), "Test", 2)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(result.Resultados) != 2 {
		t.Errorf("limit 2 devolvió %d resultados", len(result.Resultados))
	}
	if result.Resultados[1].Temperatura != 2 {
		t.Errorf("segundo resultado temperatura = %f; want 2", result.Resultados[1].Temperatura)
	}
}

// TestObtenerClima_NotFound verifica el error cuando el país no existe.
func TestObtenerClima_NotFound(t *testing.T) {
	geo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"results":[]}`)
	}))
	defer geo.Close()

	client := &Client{
		GeocodingBaseURL: geo.URL,
		WeatherBaseURL:   "http://localhost", // no se usará
		HTTPClient:       &http.Client{},
	}

	_, err := client.ObtenerClima(context.Background(), "Atlantis", 1)
	if err == nil {
		t.Fatal("se esperaba error por país no encontrado")
	}
	if err.Error() != "país no encontrado: Atlantis" {
		t.Errorf("mensaje de error inesperado: %v", err)
	}
}
