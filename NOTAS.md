# Checkpoint del Proyecto Clima API

## Estado actual (2026-08-04)
Aplicación Go dual: puede correr como CLI o como servidor REST con frontend HTML embebido.

## Cómo ejecutar
```bash
# Modo CLI interactivo
.\code-challenge.exe

# Modo CLI con flags
.\code-challenge.exe -pais Mexico -limit 3

# Modo servidor + web
.\code-challenge.exe -server :8080

Arquitectura actual (TODO: separar en paquetes)
main.go: Todo en un solo archivo (structs, lógica, handlers, template).
APIs externas: Open-Meteo (geocoding) + MET Norway (clima).
Template HTML embebido en constante indexTemplate.
Flags: -server, -pais, -limit.
Decisiones clave tomadas
Usar flag en vez de os.Args posicionales para CLI.
runCLI() recibe parámetros; runCLIInteractive() lee del teclado y llama a runCLI().
init() compila el template HTML al arrancar.
enableCORS es middleware que envuelve el handler /clima.
go.mod usa module github.com/rianeiromiron/code-challenge.
Siguiente paso pendiente
Separar en paquetes según estructura estándar de Go:
cmd/cli/
cmd/server/
internal/api/
internal/handlers/
internal/models/
web/

Comandos git útiles
git status
git add .
git commit -m "mensaje"
git push origin master

