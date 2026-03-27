import { useState, useEffect } from 'react';
import { getPackSizes, calculateOrder } from './api';
import { CalculationResponse } from './types';
import { OrderForm } from './components/OrderForm';
import { ResultCard } from './components/ResultCard';
import { ErrorBanner } from './components/ErrorBanner';
import { PackSizesManager } from './components/PackSizesManager';

function App() {
  const [packSizes, setPackSizes] = useState<number[]>([]);
  const [result, setResult] = useState<CalculationResponse | null>(null);
  const [error, setError] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingPacks, setIsLoadingPacks] = useState(true);

  useEffect(() => {
    getPackSizes()
      .then((data) => {
        // Sort descending for display
        setPackSizes([...data.pack_sizes].sort((a, b) => b - a));
      })
      .catch((err) => {
        setError(err.message || 'Failed to load pack sizes');
      })
      .finally(() => {
        setIsLoadingPacks(false);
      });
  }, []);

  const handleSubmit = async (quantity: number) => {
    if (packSizes.length === 0) {
      setError('Please add at least one pack size');
      return;
    }

    setIsLoading(true);
    setError('');
    setResult(null);

    try {
      const data = await calculateOrder(quantity, packSizes);
      setResult(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Calculation failed');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddPackSize = (size: number) => {
    setPackSizes((prev) => [...prev, size].sort((a, b) => b - a));
    setResult(null);
  };

  const handleRemovePackSize = (size: number) => {
    setPackSizes((prev) => prev.filter((s) => s !== size));
    setResult(null);
  };

  const handleUpdatePackSize = (oldSize: number, newSize: number) => {
    setPackSizes((prev) =>
      prev.map((s) => (s === oldSize ? newSize : s)).sort((a, b) => b - a)
    );
    setResult(null);
  };

  return (
    <div className="app">
      <header className="header">
        <h1>Pack Calculator</h1>
        <p>Calculate the optimal pack combination for your order</p>
      </header>

      {isLoadingPacks ? (
        <div className="loading">Loading pack sizes...</div>
      ) : (
        <>
          <PackSizesManager
            packSizes={packSizes}
            onAdd={handleAddPackSize}
            onRemove={handleRemovePackSize}
            onUpdate={handleUpdatePackSize}
          />

          <OrderForm onSubmit={handleSubmit} isLoading={isLoading} />

          {error && <ErrorBanner message={error} />}

          {result && <ResultCard result={result} />}
        </>
      )}
    </div>
  );
}

export default App;
