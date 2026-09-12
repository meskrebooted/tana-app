import { Container, getRandom } from "@cloudflare/containers";
import { env } from "cloudflare:workers";

const INSTANCE_COUNT = 1; // alza questo numero (e max_instances in wrangler.jsonc) se ti serve più capacità

// Il container esegue esattamente l'immagine definita in backend/Dockerfile:
// lo stesso backend Go usato in locale e con docker-compose, nessuna riscrittura.
export class Backend extends Container {
  defaultPort = 8080;
  sleepAfter = "10m"; // si "addormenta" (e non costa nulla) dopo 10 minuti di inattività

  envVars = {
    PORT: "8080",
    DB_HOST: env.DB_HOST,
    DB_PORT: env.DB_PORT,
    DB_USER: env.DB_USER,
    DB_PASSWORD: env.DB_PASSWORD,
    DB_NAME: env.DB_NAME,
    DB_SSLMODE: env.DB_SSLMODE,
    JWT_SECRET: env.JWT_SECRET,
    // Il frontend è servito dallo stesso Worker (stessa origine), quindi il CORS
    // del backend non entra nemmeno in gioco per le richieste del browser.
    FRONTEND_ORIGIN: "*"
  };
}

export default {
  async fetch(request, workerEnv) {
    const url = new URL(request.url);

    if (url.pathname.startsWith("/api") || url.pathname === "/health") {
      const instance = await getRandom(workerEnv.BACKEND, INSTANCE_COUNT);
      return instance.fetch(request);
    }

    return workerEnv.ASSETS.fetch(request);
  }
};
