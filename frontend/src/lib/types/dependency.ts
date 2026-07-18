export interface DependencyReason {
	type: string;
	id: number;
	name: string;
	description: string;
}

export interface DependencyCheck {
	canDelete: boolean;
	reasons: DependencyReason[];
}
