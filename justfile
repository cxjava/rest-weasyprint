set shell := ["sh", "-c"]
set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]
set allow-duplicate-recipes
set positional-arguments
set dotenv-load
set export

alias b := build
alias f := format

_default:
    @just --list

# Define variables
out_dir := "dist"
backend_out := out_dir + "/rest-weasyprint"

# Backend recipe
build:
    go build -ldflags="-w -s" -o {{backend_out}} ./cmd

format:
    go fmt ./...