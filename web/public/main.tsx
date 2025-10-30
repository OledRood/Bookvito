import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import { AuthProviderWithManager } from './AuthProviderWithManager';
import { ThemeProvider } from './ThemeContext';
import { NotificationProvider } from './NotificationContext';
import ErrorBoundary from './ErrorBoundary';
import ServerCheck from './ServerCheck';

// Global image fallback: if any <img> on the page fails to load,
// replace its src with the default book image. We use a capturing
// error listener so it catches <img> load errors on the document.
const DEFAULT_IMAGE = '/images/default-book.png';

function globalImageErrorHandler(e: Event) {
  const target = e.target as HTMLElement | null;
  if (!target) return;

  // Only handle <img> elements
  if (target.tagName === 'IMG') {
    const img = target as HTMLImageElement;
    // Avoid infinite loops: set a flag if we've already applied fallback
    if (img.dataset.fallbackApplied) return;
    img.dataset.fallbackApplied = '1';
    // remove handler on this element to be safe
    img.onerror = null;
    try {
      img.src = DEFAULT_IMAGE;
    } catch (err) {
      // ignore
    }
  }
}

// Use capture so image errors (which don't bubble) are caught here
window.addEventListener('error', globalImageErrorHandler, true);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <ThemeProvider>
        <AuthProviderWithManager>
          <NotificationProvider>
            <ErrorBoundary>
              <ServerCheck>
                <App />
              </ServerCheck>
            </ErrorBoundary>
          </NotificationProvider>
        </AuthProviderWithManager>
      </ThemeProvider>
    </BrowserRouter>
  </React.StrictMode>,
);