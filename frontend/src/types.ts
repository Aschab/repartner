export type PackSelection = {
  pack_size: number;
  count: number;
};

export type CalculationRequest = {
  order_quantity: number;
  pack_sizes: number[];
};

export type CalculationResponse = {
  order_quantity: number;
  total_shipped: number;
  total_packs: number;
  packs: PackSelection[];
};

export type PackSizesResponse = {
  pack_sizes: number[];
};

export type ErrorResponse = {
  error: string;
};
