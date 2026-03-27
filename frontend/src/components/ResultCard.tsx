import { CalculationResponse } from '../types';

type ResultCardProps = {
  result: CalculationResponse;
};

export function ResultCard({ result }: ResultCardProps) {
  return (
    <div className="result-card" data-testid="result-card">
      <h2>Calculation Result</h2>

      <div className="result-summary">
        <div className="summary-item">
          <div className="label">Ordered</div>
          <div className="value">{result.order_quantity.toLocaleString()}</div>
        </div>
        <div className="summary-item">
          <div className="label">Shipped</div>
          <div className="value">{result.total_shipped.toLocaleString()}</div>
        </div>
        <div className="summary-item">
          <div className="label">Packs</div>
          <div className="value">{result.total_packs.toLocaleString()}</div>
        </div>
      </div>

      <div className="pack-breakdown">
        <h3>Pack Breakdown</h3>
        <div className="pack-list">
          {result.packs.map((pack) => (
            <div key={pack.pack_size} className="pack-item">
              <span className="pack-size">{pack.pack_size.toLocaleString()} items/pack</span>
              <span className="pack-count">&times; {pack.count.toLocaleString()}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
