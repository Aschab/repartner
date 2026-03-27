import { useState } from 'react';

type PackSizesManagerProps = {
  packSizes: number[];
  onAdd: (size: number) => void;
  onRemove: (size: number) => void;
  onUpdate: (oldSize: number, newSize: number) => void;
};

export function PackSizesManager({ packSizes, onAdd, onRemove, onUpdate }: PackSizesManagerProps) {
  const [newSize, setNewSize] = useState('');
  const [editingSize, setEditingSize] = useState<number | null>(null);
  const [editValue, setEditValue] = useState('');
  const [error, setError] = useState('');

  const handleAdd = () => {
    const size = parseInt(newSize, 10);
    if (isNaN(size) || size <= 0) {
      setError('Please enter a positive integer');
      return;
    }
    if (packSizes.includes(size)) {
      setError('Pack size already exists');
      return;
    }
    setError('');
    onAdd(size);
    setNewSize('');
  };

  const handleStartEdit = (size: number) => {
    setEditingSize(size);
    setEditValue(size.toString());
    setError('');
  };

  const handleSaveEdit = () => {
    if (editingSize === null) return;

    const newSizeValue = parseInt(editValue, 10);
    if (isNaN(newSizeValue) || newSizeValue <= 0) {
      setError('Please enter a positive integer');
      return;
    }
    if (newSizeValue !== editingSize && packSizes.includes(newSizeValue)) {
      setError('Pack size already exists');
      return;
    }
    setError('');
    onUpdate(editingSize, newSizeValue);
    setEditingSize(null);
    setEditValue('');
  };

  const handleCancelEdit = () => {
    setEditingSize(null);
    setEditValue('');
    setError('');
  };

  const handleKeyDown = (e: React.KeyboardEvent, action: 'add' | 'edit') => {
    if (e.key === 'Enter') {
      if (action === 'add') handleAdd();
      else handleSaveEdit();
    }
    if (e.key === 'Escape' && action === 'edit') {
      handleCancelEdit();
    }
  };

  return (
    <div className="pack-sizes-manager" data-testid="pack-sizes-manager">
      <h3>Pack Sizes</h3>

      <div className="pack-sizes-list">
        {packSizes.map((size) => (
          <div key={size} className="pack-size-item">
            {editingSize === size ? (
              <div className="pack-size-edit">
                <input
                  type="number"
                  value={editValue}
                  onChange={(e) => setEditValue(e.target.value)}
                  onKeyDown={(e) => handleKeyDown(e, 'edit')}
                  autoFocus
                />
                <button onClick={handleSaveEdit} className="btn-save" title="Save">
                  ✓
                </button>
                <button onClick={handleCancelEdit} className="btn-cancel" title="Cancel">
                  ✕
                </button>
              </div>
            ) : (
              <>
                <span className="pack-size-value">{size.toLocaleString()}</span>
                <div className="pack-size-actions">
                  <button
                    onClick={() => handleStartEdit(size)}
                    className="btn-edit"
                    title="Edit"
                  >
                    ✎
                  </button>
                  <button
                    onClick={() => onRemove(size)}
                    className="btn-remove"
                    title="Remove"
                    disabled={packSizes.length <= 1}
                  >
                    ✕
                  </button>
                </div>
              </>
            )}
          </div>
        ))}
      </div>

      <div className="pack-size-add">
        <input
          type="number"
          value={newSize}
          onChange={(e) => setNewSize(e.target.value)}
          onKeyDown={(e) => handleKeyDown(e, 'add')}
          placeholder="New pack size..."
          min="1"
        />
        <button onClick={handleAdd} className="btn-add">
          Add
        </button>
      </div>

      {error && <div className="pack-sizes-error">{error}</div>}
    </div>
  );
}
