import { render, screen } from '@testing-library/react'
import App from '../../App'

describe('routing smoke test', () => {
  it('renders the public resume route content on the home page', () => {
    render(<App />)

    expect(screen.getByText('温庆京')).toBeInTheDocument()
  })
})
