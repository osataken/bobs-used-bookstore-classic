import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { vi, describe, it, expect } from 'vitest';
import HomePage from './HomePage';

vi.mock('../services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({
      data: {
        bestSellers: [
          { id: 1, name: '2020: The Apocalypse', author: 'Li Juan', price: 10.95, quantity: 25 },
          { id: 2, name: 'Children Of Iron', author: 'Nikki Wolf', price: 13.95, quantity: 3 },
        ],
      },
    }),
  },
}));

describe('HomePage', () => {
  it('renders best sellers heading', async () => {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
          <HomePage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('Best Sellers')).toBeInTheDocument();
  });

  it('renders book list', async () => {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
          <HomePage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('2020: The Apocalypse')).toBeInTheDocument();
    expect(await screen.findByText('Children Of Iron')).toBeInTheDocument();
  });
});
