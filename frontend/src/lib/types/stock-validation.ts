export interface StockValidationRequest {
	items: {
		productId: number;
		quantity: number;
	}[];
}

export interface InsufficientIngredient {
	ingredientId: number;
	ingredientName: string;
	required: number;
	available: number;
	shortage: number;
	unit: string;
}

export interface StockValidationResponse {
	valid: boolean;
	insufficientStock: InsufficientIngredient[];
}
