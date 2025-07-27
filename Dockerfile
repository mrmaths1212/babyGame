# Utilise une image Go officielle
FROM golang:1.24.5 as builder

# Répertoire de travail dans le conteneur
WORKDIR /app

# Copie les fichiers
COPY . .

# Compilation du binaire
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server
# Étape finale - image légère
FROM debian:bullseye-slim

# Installer uniquement ce qu'il faut
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copier le binaire depuis l'étape précédente
COPY --from=builder /app/server /app/server
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/data /app/data

WORKDIR /app

EXPOSE 8080

CMD ["./server"]
