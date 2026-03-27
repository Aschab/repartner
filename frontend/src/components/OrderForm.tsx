import { useState, FormEvent } from 'react';

type OrderFormProps = {
  onSubmit: (quantity: number) => void;
  isLoading: boolean;
};

export function OrderForm({ onSubmit, isLoading }: OrderFormProps) {
  const [inputValue, setInputValue] = useState('');
  const [validationError, setValidationError] = useState('');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();

    const trimmed = inputValue.trim();
    if (!trimmed) {
      setValidationError('Please enter a quantity');
      return;
    }

    const quantity = parseInt(trimmed, 10);
    if (isNaN(quantity) || quantity <= 0) {
      setValidationError('Please enter a positive integer');
      return;
    }

    setValidationError('');
    onSubmit(quantity);
  };

  const handleChange = (value: string) => {
    setInputValue(value);
    if (validationError) {
      setValidationError('');
    }
  };

  return (
    <form className="order-form" onSubmit={handleSubmit} data-testid="order-form">
      <div className="form-group">
        <label htmlFor="quantity">Order Quantity</label>
        <input
          id="quantity"
          type="number"
          min="1"
          step="1"
          value={inputValue}
          onChange={(e) => handleChange(e.target.value)}
          placeholder="Enter quantity..."
          className={validationError ? 'invalid' : ''}
          disabled={isLoading}
        />
        {validationError && (
          <div className="validation-message">{validationError}</div>
        )}
      </div>
      <button type="submit" className="submit-button" disabled={isLoading}>
        {isLoading ? 'Calculating...' : 'Calculate Packs'}
      </button>
    </form>
  );
}
