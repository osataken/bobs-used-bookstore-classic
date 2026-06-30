import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect, vi } from 'vitest'
import Navbar from './Navbar'

vi.mock('../hooks/useAuth', () => ({
  useAuth: () => ({ user: null, loading: false }),
}))

vi.mock('../services/api', () => ({
  login: vi.fn(),
  logout: vi.fn(),
}))

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>{ui}</MemoryRouter>
    </QueryClientProvider>
  )
}

describe('Navbar', () => {
  it('renders store name', () => {
    renderWithProviders(<Navbar />)
    expect(screen.getByText("Bob's Used Bookstore")).toBeInTheDocument()
  })

  it('shows login button when not authenticated', () => {
    renderWithProviders(<Navbar />)
    expect(screen.getByText('Login')).toBeInTheDocument()
  })

  it('shows navigation links', () => {
    renderWithProviders(<Navbar />)
    expect(screen.getByText('Search')).toBeInTheDocument()
    expect(screen.getByText('Cart')).toBeInTheDocument()
    expect(screen.getByText('Wishlist')).toBeInTheDocument()
  })
})
