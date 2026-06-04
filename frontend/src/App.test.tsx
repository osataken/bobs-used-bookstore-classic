import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi } from 'vitest';
import App from './App';

vi.mock('./context/AuthContext', () => ({
  AuthProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  useAuth: () => ({
    auth: { authenticated: false },
    loading: false,
    login: vi.fn(),
    logout: vi.fn(),
  }),
}));

vi.mock('./services/api', () => ({
  getBestSellingBooks: vi.fn().mockResolvedValue([]),
  getAuthStatus: vi.fn().mockResolvedValue({ authenticated: false }),
}));

function renderApp() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </QueryClientProvider>
  );
}

describe('App', () => {
  it('renders the navigation with bookstore name', () => {
    renderApp();
    expect(screen.getByText("Bob's Used Bookstore")).toBeInTheDocument();
  });

  it('renders Browse link', () => {
    renderApp();
    expect(screen.getByText('Browse')).toBeInTheDocument();
  });

  it('renders Cart link', () => {
    renderApp();
    expect(screen.getByText('Cart')).toBeInTheDocument();
  });

  it('renders Login button when not authenticated', () => {
    renderApp();
    expect(screen.getByText('Login')).toBeInTheDocument();
  });
});
