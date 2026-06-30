import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Home from '../Home'

vi.mock('../../hooks/useBooks', () => ({
  useBestSellers: () => ({
    data: [
      { id: 1, name: '2020: The Apocalypse', author: 'Li Juan', price: 10.95, quantity: 25, coverImageUrl: '' },
      { id: 2, name: 'Children Of Iron', author: 'Nikki Wolf', price: 13.95, quantity: 3, coverImageUrl: '' },
    ],
    isLoading: false,
  }),
}))

vi.mock('../../hooks/useCart', () => ({
  useAddToCart: () => ({ mutate: vi.fn() }),
}))

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>{ui}</BrowserRouter>
    </QueryClientProvider>
  )
}

describe('Home', () => {
  it('renders best sellers heading', () => {
    renderWithProviders(<Home />)
    expect(screen.getByText('Best Sellers')).toBeInTheDocument()
  })

  it('renders book titles', () => {
    renderWithProviders(<Home />)
    expect(screen.getByText('2020: The Apocalypse')).toBeInTheDocument()
    expect(screen.getByText('Children Of Iron')).toBeInTheDocument()
  })

  it('renders add to cart buttons', () => {
    renderWithProviders(<Home />)
    const buttons = screen.getAllByText('Add to Cart')
    expect(buttons.length).toBe(2)
  })
})
