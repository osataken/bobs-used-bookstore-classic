import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import CartPage from './CartPage'

vi.mock('../services/api', () => ({
  getCart: vi.fn().mockResolvedValue({
    data: {
      items: [
        { id: 1, shoppingCartId: 1, bookId: 1, quantity: 2, wantToBuy: true, book: { id: 1, name: 'Test Book', price: 10.00, author: 'Author' } },
      ],
      subTotal: 10.00,
    }
  }),
  deleteCartItem: vi.fn().mockResolvedValue({}),
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

describe('CartPage', () => {
  it('renders cart heading', async () => {
    renderWithProviders(<CartPage />)
    expect(await screen.findByText('Shopping Cart')).toBeInTheDocument()
  })

  it('displays cart items', async () => {
    renderWithProviders(<CartPage />)
    expect(await screen.findByText('Test Book')).toBeInTheDocument()
  })
})
