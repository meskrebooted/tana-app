# Tana — coworking app

Progetto full-stack sviluppato per la Tech Journey "Dev Run" (Golang + Vue + Postgres):
un'app di co-working dove prenoti postazioni, sale riunioni e phone booth. Sotto il tono
giocoso (nomi degli spazi, "vibe", streak di presenza, pulsante "Sorprendimi") il
funzionamento è quello di un vero sistema di prenotazione: autenticazione JWT, controllo
dei conflitti di orario a livello di database, ruoli utente/admin.

## Stack

- **Backend**: Go, Gin, GORM, Postgres, JWT, bcrypt
- **Frontend**: Vue 3 (Composition API), Vite, Pinia, Vue Router, Axios
- **DevOps**: Docker multi-stage per backend e frontend, Docker Compose, GitHub Actions CI

## Struttura

```
coworking-app/
├── backend/          # API Go
├── frontend/         # SPA Vue
├── docker-compose.yml
└── .github/workflows/ci.yml
```

## Avvio rapido con Docker (consigliato)

Richiede solo Docker e Docker Compose installati.

```bash
docker compose up --build
```

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api
- Postgres: localhost:5432 (utente/password/db: `coworking`)

Le tabelle vengono create automaticamente (auto-migration) e vengono seminati alcuni
spazi di esempio al primo avvio.

## Avvio in locale senza Docker

### 1. Postgres

Avvia un Postgres locale e crea un database `coworking` con utente `coworking` /
password `coworking` (oppure modifica le variabili d'ambiente).

### 2. Backend

```bash
cd backend
cp .env.example .env   # personalizza se necessario
go mod tidy
go run main.go
```

Il server Go legge le variabili con `os.Getenv`, quindi se non usi un tool tipo
`direnv`/`godotenv` esportale a mano oppure lancia con:

```bash
export $(cat .env | xargs) && go run main.go
```

L'API risponde su `http://localhost:8080`.

### 3. Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

L'app è su `http://localhost:5173`.

## Come funziona

1. **Registrazione/Login** → `/api/auth/register`, `/api/auth/login` restituiscono un
   JWT salvato nel browser (localStorage) e allegato automaticamente alle richieste
   successive (interceptor Axios).
2. **Spazi** → lista di postazioni/sale/booth, ciascuna con una breve "vibe". Il
   pulsante **🎲 Sorprendimi** chiama `/api/spaces/random`, che sceglie a caso uno
   spazio libero nella prossima ora.
3. **Prenotazioni** → creare una prenotazione (`POST /api/bookings`) controlla lato
   database che non esista già una prenotazione confermata che si sovrappone allo
   stesso spazio e fascia oraria; in caso di conflitto risponde `409`.
4. **Streak** → `GET /api/bookings/streak` calcola per quanti giorni consecutivi
   l'utente ha prenotato almeno uno spazio, mostrato nella sidebar come 🔥.

## Riferimento API

| Metodo | Endpoint                     | Auth   | Descrizione                                  |
|--------|-------------------------------|--------|-----------------------------------------------|
| POST   | `/api/auth/register`          | -      | Crea un utente, restituisce JWT               |
| POST   | `/api/auth/login`             | -      | Login, restituisce JWT                        |
| GET    | `/api/auth/me`                | utente | Profilo dell'utente autenticato               |
| GET    | `/api/spaces`                 | utente | Elenco spazi                                  |
| GET    | `/api/spaces/random`          | utente | Uno spazio libero scelto a caso                |
| POST   | `/api/spaces`                 | admin  | Crea uno spazio                               |
| DELETE | `/api/spaces/:id`             | admin  | Elimina uno spazio                            |
| POST   | `/api/bookings`                | utente | Crea una prenotazione (con controllo conflitti)|
| GET    | `/api/bookings/me`             | utente | Le prenotazioni dell'utente                   |
| GET    | `/api/bookings/streak`         | utente | Giorni consecutivi di prenotazione            |
| PATCH  | `/api/bookings/:id/cancel`     | utente | Annulla una prenotazione propria              |
| GET    | `/health`                      | -      | Health check                                  |

Per rendere un utente admin (es. per creare/eliminare spazi), aggiorna manualmente
la colonna `role` a `admin` sulla riga corrispondente nella tabella `users`.

## Deploy su Cloudflare (Workers + Containers)

Il repo include già `wrangler.jsonc` e `worker/index.js`: un unico Worker Cloudflare
serve il frontend Vue come asset statico e instrada le richieste `/api` a un
Container che esegue **la stessa immagine Docker del backend Go**, senza riscrivere
nulla. Frontend e backend finiscono sulla stessa origine, quindi niente problemi di
CORS in produzione. Il database resta esterno (Cloudflare non offre Postgres
gestito): usa un Postgres gratuito come [Neon](https://neon.tech) o
[Supabase](https://supabase.com).

### Passo passo

1. **Carica il progetto su GitHub** (se non l'hai già fatto):
   ```bash
   git init && git add . && git commit -m "prima versione"
   git branch -M main
   git remote add origin <url-del-tuo-repo>
   git push -u origin main
   ```

2. **Crea il database** su [neon.tech](https://neon.tech): un progetto gratuito ti
   dà host, utente, password e nome del database da usare nei passi successivi.

3. **Crea un account Cloudflare** (free tier) su [dash.cloudflare.com](https://dash.cloudflare.com)
   e assicurati che il piano **Workers Paid** sia attivo: i Container richiedono
   questo piano (ha comunque una soglia gratuita mensile generosa).

4. **Collega il repo** in Workers & Pages → Create → Import a repository, scegliendo
   il tuo repo GitHub. Cloudflare userà "Workers Builds" per buildare e deployare
   automaticamente ad ogni push.

5. **Imposta il comando di build** su `npm run deploy` (root directory del progetto,
   dove si trova `wrangler.jsonc`). Questo comando fa: build del frontend Vue,
   poi `wrangler deploy`, che pubblica il Worker e builda/pubblica l'immagine
   Docker del container — tutto lato Cloudflare, senza bisogno di Docker sulla tua
   macchina.

6. **Configura i secrets** del Worker (Settings → Variables and Secrets, oppure via
   CLI `wrangler secret put NOME`): `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`,
   `DB_NAME`, `DB_SSLMODE` (metti `require` per Neon), `JWT_SECRET`. Il file
   `worker/index.js` li inoltra al container come variabili d'ambiente.

7. **Fai partire il deploy**: push su `main`, oppure trigger manuale dalla dashboard.
   Al primo avvio il backend farà comunque l'auto-migration e il seed degli spazi,
   proprio come in locale.

8. **Apri l'URL** che Cloudflare assegna al Worker (tipo
   `https://tana-coworking.<tuo-account>.workers.dev`): frontend e API sono lì,
   sulla stessa origine.

Per iterare: modifichi il codice, fai `git push`, Cloudflare rifà build e deploy da
solo. Per un dominio personalizzato, aggiungilo dalle impostazioni del Worker.



`.github/workflows/ci.yml` esegue ad ogni push/PR su `main`:

1. build + `go vet` del backend
2. `npm install` + build del frontend
3. build delle immagini Docker di entrambi i servizi

`docker-compose.yml` orchestra Postgres, backend e frontend con healthcheck sul
database, così il backend attende che Postgres sia pronto prima di partire.
