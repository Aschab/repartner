import { useState, useEffect } from 'react';
import { getPackSizes, calculateOrder } from './api';
import { CalculationResponse } from './types';
import { OrderForm } from './components/OrderForm';
import { ResultCard } from './components/ResultCard';
import { ErrorBanner } from './components/ErrorBanner';

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
    setIsLoading(true);
    setError('');
    setResult(null);

    try {
      const data = await calculateOrder(quantity);
      setResult(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Calculation failed');
    } finally {
      setIsLoading(false);
    }
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
          <div className="pack-sizes">
            <h3>Available Pack Sizes</h3>
            <div className="pack-sizes-list">
              {packSizes.map((size) => (
                <span key={size} className="pack-size-badge">
                  {size.toLocaleString()}
                </span>
              ))}
            </div>
          </div>

          <OrderForm onSubmit={handleSubmit} isLoading={isLoading} />

          {error && <ErrorBanner message={error} />}

          {result && <ResultCard result={result} />}
        </>
      )}
    </div>
  );
}

export default App;
