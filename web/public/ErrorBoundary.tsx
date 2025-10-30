import React from 'react';

interface State {
  hasError: boolean;
  error: Error | null;
  info: React.ErrorInfo | null;
}

class ErrorBoundary extends React.Component<React.PropsWithChildren<{}>, State> {
  constructor(props: React.PropsWithChildren<{}>) {
    super(props);
    this.state = { hasError: false, error: null, info: null };
  }

  static getDerivedStateFromError(error: Error) {
    // Этот метод должен вернуть объект для обновления состояния.
    return { hasError: true };
  }

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    // You can log error to an external service here
    // eslint-disable-next-line no-console
    console.error('Uncaught error in component tree:', error, info);
    this.setState({ error, info });
  }

  render() {
    if (this.state.hasError) {
      return (
        <div style={{ padding: 24 }}>
          <h1>Произошла ошибка</h1>
          <p>Пожалуйста, откройте консоль разработчика, чтобы посмотреть стек ошибки.</p>
          {this.state.error && (
            <details style={{ whiteSpace: 'pre-wrap' }}>
              <summary>Показать детали</summary>
              {this.state.error.stack}
              {this.state.info && this.state.info.componentStack}
            </details>
          )}
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
