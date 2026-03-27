import { CalculationResponse, PackSizesResponse, ErrorResponse } from './types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export async function getPackSizes(): Promise<PackSizesResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/packs`);

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch pack sizes');
  }

  return response.json();
}

export async function calculateOrder(orderQuantity: number): Promise<CalculationResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/calculate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ order_quantity: orderQuantity }),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Calculation failed');
  }

  return response.json();
}
