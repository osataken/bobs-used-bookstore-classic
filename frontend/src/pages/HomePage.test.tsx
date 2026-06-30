import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import HomePage from './HomePage'

vi.mock('../services/api', () => ({
  getHome: vi.fn().mockResolvedValue({
    data: {
      bestSellingBooks: [
        { id: 1, name: '2020: The Apocalypse', author: 'Li Juan', price: 10.95, isInStock: true, coverImageUrl: '' },
        { id: 2, name: 'Children Of Iron', author: 'Nikki Wolf', price: 13.95, isInStock: true, coverImageUrl: '' },
      ]
    }
  }),
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

describe('HomePage', () => {
  it('renders welcome heading', async () => {
    renderWithProviders(<HomePage />)
    expect(await screen.findByText("Welcome to Bob's Used Bookstore")).toBeInTheDocument()
  })

  it('renders best selling books', async () => {
    renderWithProviders(<HomePage />)
    expect(await screen.findByText('2020: The Apocalypse')).toBeInTheDocument()
    expect(await screen.findByText('Children Of Iron')).toBeInTheDocument()
  })
})
