import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { vi, describe, it, expect } from 'vitest';
import SearchPage from './SearchPage';

vi.mock('../services/api', () => ({
  default: {
    get: vi.fn().mockResolvedValue({
      data: {
        items: [
          { id: 1, name: 'Test Book', author: 'Test Author', price: 9.99, quantity: 5, publisher: { text: 'Pub' }, bookType: { text: 'HC' }, genre: { text: 'Fiction' }, condition: { text: 'Good' } },
        ],
        totalCount: 1,
        pageIndex: 1,
        pageSize: 10,
        totalPages: 1,
      },
    }),
    post: vi.fn().mockResolvedValue({}),
  },
}));

describe('SearchPage', () => {
  it('renders search page with books', async () => {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
          <SearchPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('Test Book')).toBeInTheDocument();
    expect(screen.getByText('by Test Author')).toBeInTheDocument();
    expect(screen.getByText('$9.99')).toBeInTheDocument();
  });

  it('renders search input', () => {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
          <SearchPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByPlaceholderText('Search by name, author, or ISBN...')).toBeInTheDocument();
  });
});
