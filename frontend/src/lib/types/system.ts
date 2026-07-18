export interface Health {
	status: string;
	database: string;
	storage: string;
	version: string;
	uptime: string;
}

export interface Version {
	version: string;
	commit: string;
	build: string;
	environment: string;
}

export interface Capabilities {
	upload: boolean;
	seo: boolean;
	marketplace: boolean;
	ifood: boolean;
	pix: boolean;
	fiscal: boolean;
	delivery: boolean;
	cardapioDigital: boolean;
}
