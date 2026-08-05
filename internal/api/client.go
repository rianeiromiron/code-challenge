package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rianeiromiron/code-challenge/internal/models"
)

func ObtenerClima(ctx context.Context, pais string, limit int) (*models.ClimaResponse, error) {
	query := url.Values{}
	query.Set("name", pais)
	query.Set("language", "es")

	apiURL := "https://geocoding-api.open-meteo.com/v1/search?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error geocoding: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var localizacion models.GeocodingResponse
	if err := json.Unmarshal(body, &localizacion); err != nil {
		return nil, fmt.Errorf("error parseando geocoding: %w", err)
	}
	if len(localizacion.Results) == 0 {
		return nil, fmt.Errorf("país no encontrado: %s", pais)
	}

	lat := localizacion.Results[0].Latitude
	lon := localizacion.Results[0].Longitude

	query = url.Values{}
	query.Set("lat", fmt.Sprint(lat))
	query.Set("lon", fmt.Sprint(lon))

	apiURL = "https://api.met.no/weatherapi/locationforecast/2.0/compact.json?" + query.Encode()
	req, err = http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "PruebasClima/1.0 rianeiromiron@gmail.com")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error API clima: %w", err)
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var datosclima models.Clima
	if err := json.Unmarshal(body, &datosclima); err != nil {
		return nil, fmt.Errorf("error parseando clima: %w", err)
	}

	if limit > len(datosclima.Properties.Timeseries) || limit <= 0 {
		limit = len(datosclima.Properties.Timeseries)
	}

	resp := &models.ClimaResponse{
		Pais:     pais,
		Latitud:  lat,
		Longitud: lon,
	}

	for _, t := range datosclima.Properties.Timeseries[:limit] {
		r := models.Resultado{
			FechaHora:   t.Time,
			Temperatura: t.Data.Instant.Details.AirTemperature,
			Presion:     t.Data.Instant.Details.AirPressureAtSeaLevel,
			Humedad:     t.Data.Instant.Details.RelativeHumidity,
			Viento:      t.Data.Instant.Details.WindSpeed,
			VientoDir:   t.Data.Instant.Details.WindFromDirection,
			Nubes:       t.Data.Instant.Details.CloudAreaFraction,
		}
		if t.Data.Next12Hours.Summary.SymbolCode != "" {
			r.Simbolo12h = t.Data.Next12Hours.Summary.SymbolCode
			r.Precipitacion12h = t.Data.Next12Hours.Details.PrecipitationAmount
		}
		if t.Data.Next1Hours.Summary.SymbolCode != "" {
			r.Simbolo1h = t.Data.Next1Hours.Summary.SymbolCode
			r.Precipitacion1h = t.Data.Next1Hours.Details.PrecipitationAmount
		}
		if t.Data.Next6Hours.Summary.SymbolCode != "" {
			r.Simbolo6h = t.Data.Next6Hours.Summary.SymbolCode
			r.Precipitacion6h = t.Data.Next6Hours.Details.PrecipitationAmount
		}
		resp.Resultados = append(resp.Resultados, r)
	}

	return resp, nil
}
