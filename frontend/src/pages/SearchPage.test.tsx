import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import SearchPage from './SearchPage'

vi.mock('../services/api', () => ({
  searchBooks: vi.fn().mockResolvedValue({
    data: {
      items: [
        { id: 1, name: 'Test Book', author: 'Author', price: 9.99, genre: 'Fiction', condition: 'Good', isInStock: true, isLowInStock: false, coverImageUrl: '' },
      ],
      pageIndex: 1,
      totalPages: 1,
      totalCount: 1,
      hasNext: false,
      hasPrevious: false,
    }
  }),
  addToCart: vi.fn().mockResolvedValue({}),
  addToWishlist: vi.fn().mockResolvedValue({}),
  getMe: vi.fn().mockRejectedValue(new Error('not authenticated')),
}))

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>{ui}</MemoryRouter>
    </QueryClientProvider>
  )
}

describe('SearchPage', () => {
  it('renders search heading', async () => {
    renderWithProviders(<SearchPage />)
    expect(screen.getByText('Search Books')).toBeInTheDocument()
  })

  it('displays search results', async () => {
    renderWithProviders(<SearchPage />)
    expect(await screen.findByText('Test Book')).toBeInTheDocument()
  })
})
