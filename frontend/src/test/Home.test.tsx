import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Home from '../pages/Home';

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

function renderWithProviders(ui: React.ReactElement) {
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {ui}
      </BrowserRouter>
    </QueryClientProvider>
  );
}

describe('Home Page', () => {
  it('renders the page title', () => {
    renderWithProviders(<Home />);
    expect(screen.getByText("Bob's Used Bookstore")).toBeInTheDocument();
  });

  it('shows loading state', () => {
    renderWithProviders(<Home />);
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });
});
