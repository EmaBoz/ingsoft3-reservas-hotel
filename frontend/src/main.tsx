import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import * as Sentry from '@sentry/react'
import './index.css'
import App from './App.tsx'

// El DSN del front no es un secreto: sólo sirve para mandar datos a Sentry,
// no para leerlos, y viaja igual en el JavaScript que baja el navegador.
Sentry.init({
  dsn: 'https://25e6c1f4705651d4f5732d52e3ae02b8@o4512317152624640.ingest.us.sentry.io/4512317185654784',
  // QA y PROD corren los dos en localhost (fallback local): el puerto dice cuál es
  environment: window.location.host,
  integrations: [Sentry.browserTracingIntegration()],
  tracesSampleRate: 1.0,
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
