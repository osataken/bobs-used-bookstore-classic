import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Search from '../Search'

vi.mock('../../hooks/useBooks', () => ({
  useSearchBooks: () => ({
    data: {
      items: [
        { id: 1, name: '2020: The Apocalypse', author: 'Li Juan', price: 10.95, quantity: 25, coverImageUrl: '' },
        { id: 3, name: 'Gold In The Dark', author: 'Richard Roe', price: 6.50, quantity: 10, coverImageUrl: '' },
      ],
      pageIndex: 1,
      pageSize: 10,
      totalCount: 2,
      totalPages: 1,
    },
    isLoading: false,
  }),
}))

vi.mock('../../hooks/useCart', () => ({
  useAddToCart: () => ({ mutate: vi.fn() }),
  useAddToWishlist: () => ({ mutate: vi.fn() }),
}))

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>{ui}</BrowserRouter>
    </QueryClientProvider>
  )
}

describe('Search', () => {
  it('renders browse books heading', () => {
    renderWithProviders(<Search />)
    expect(screen.getByText('Browse Books')).toBeInTheDocument()
  })

  it('renders search results', () => {
    renderWithProviders(<Search />)
    expect(screen.getByText('2020: The Apocalypse')).toBeInTheDocument()
    expect(screen.getByText('Gold In The Dark')).toBeInTheDocument()
  })

  it('renders search input', () => {
    renderWithProviders(<Search />)
    expect(screen.getByPlaceholderText('Search by name, author, or ISBN...')).toBeInTheDocument()
  })

  it('renders wishlist buttons', () => {
    renderWithProviders(<Search />)
    const wishlistButtons = screen.getAllByText('Wishlist')
    expect(wishlistButtons.length).toBe(2)
  })
})
