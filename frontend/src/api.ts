import { CalculationResponse, PackSizesResponse, ErrorResponse } from './types';

function getApiBaseUrl(): string {
  // 1. Check for explicit env var (set at build time)
  if (import.meta.env.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }

  // 2. In production, try same origin (if backend is proxied or same domain)
  if (import.meta.env.PROD) {
    return window.location.origin;
  }

  // 3. Default for local development
  return 'http://localhost:8080';
}

const API_BASE_URL = getApiBaseUrl();

export async function getPackSizes(): Promise<PackSizesResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/packs`);

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch pack sizes');
  }

  return response.json();
}

export async function calculateOrder(orderQuantity: number, packSizes: number[]): Promise<CalculationResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/calculate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      order_quantity: orderQuantity,
      pack_sizes: packSizes,
    }),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Calculation failed');
  }

  return response.json();
}
