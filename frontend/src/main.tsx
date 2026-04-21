import { ClerkProvider } from '@clerk/react'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { ROUTES } from './lib/routes.ts'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ClerkProvider
      publishableKey={import.meta.env.VITE_CLERK_PUBLISHABLE_KEY}
      signInUrl={ROUTES.MAIN.AUTH.SIGN_IN}
      signUpUrl={ROUTES.MAIN.AUTH.SIGN_UP}
    >
      <App />
    </ClerkProvider>
  </StrictMode>,
)
