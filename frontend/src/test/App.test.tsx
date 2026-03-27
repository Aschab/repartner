import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import App from '../App';
import * as api from '../api';

vi.mock('../api');

const mockPackSizes = { pack_sizes: [250, 500, 1000, 2000, 5000] };
const mockCalculationResult = {
  order_quantity: 501,
  total_shipped: 750,
  total_packs: 2,
  packs: [
    { pack_size: 500, count: 1 },
    { pack_size: 250, count: 1 },
  ],
};

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders pack sizes on load', async () => {
    vi.mocked(api.getPackSizes).mockResolvedValue(mockPackSizes);

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText('5,000')).toBeInTheDocument();
    });

    expect(screen.getByText('2,000')).toBeInTheDocument();
    expect(screen.getByText('1,000')).toBeInTheDocument();
    expect(screen.getByText('500')).toBeInTheDocument();
    expect(screen.getByText('250')).toBeInTheDocument();
  });

  it('shows error when pack sizes fail to load', async () => {
    vi.mocked(api.getPackSizes).mockRejectedValue(new Error('Network error'));

    render(<App />);

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('Network error');
    });
  });

  it('shows result after successful calculation', async () => {
    vi.mocked(api.getPackSizes).mockResolvedValue(mockPackSizes);
    vi.mocked(api.calculateOrder).mockResolvedValue(mockCalculationResult);

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText('5,000')).toBeInTheDocument();
    });

    const input = screen.getByLabelText('Order Quantity');
    fireEvent.change(input, { target: { value: '501' } });

    const button = screen.getByRole('button', { name: 'Calculate Packs' });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByTestId('result-card')).toBeInTheDocument();
    });

    expect(screen.getByText('501')).toBeInTheDocument();
    expect(screen.getByText('750')).toBeInTheDocument();
  });

  it('shows error when calculation fails', async () => {
    vi.mocked(api.getPackSizes).mockResolvedValue(mockPackSizes);
    vi.mocked(api.calculateOrder).mockRejectedValue(
      new Error('order_quantity must be greater than 0')
    );

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText('5,000')).toBeInTheDocument();
    });

    const input = screen.getByLabelText('Order Quantity');
    fireEvent.change(input, { target: { value: '501' } });

    const button = screen.getByRole('button', { name: 'Calculate Packs' });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(
        'order_quantity must be greater than 0'
      );
    });
  });
});

describe('OrderForm validation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(api.getPackSizes).mockResolvedValue(mockPackSizes);
  });

  it('shows validation error for empty input', async () => {
    render(<App />);

    await waitFor(() => {
      expect(screen.getByText('5,000')).toBeInTheDocument();
    });

    const button = screen.getByRole('button', { name: 'Calculate Packs' });
    fireEvent.click(button);

    expect(screen.getByText('Please enter a quantity')).toBeInTheDocument();
  });

  it('shows validation error for invalid input', async () => {
    render(<App />);

    await waitFor(() => {
      expect(screen.getByText('5,000')).toBeInTheDocument();
    });

    const input = screen.getByLabelText('Order Quantity');
    fireEvent.change(input, { target: { value: '0' } });

    const form = screen.getByTestId('order-form');
    fireEvent.submit(form);

    await waitFor(() => {
      expect(screen.getByText('Please enter a positive integer')).toBeInTheDocument();
    });
  });
});
