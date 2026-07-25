export function cookieSnapshot(label: string) {
	const rawCookie = document.cookie;
	const cookies = rawCookie
		.split("; ")
		.filter(Boolean);

	console.group(`COOKIE SNAPSHOT -> ${label}`);

	console.log("document.cookie (raw):", rawCookie);
	console.log("auth_token via includes:", rawCookie.includes("auth_token="));
	console.log("auth_token via split:", cookies.some(c => c.startsWith("auth_token=")));

	console.log("cookies separados:");
	console.table(cookies);

	console.log("auth_token encontrado:",
		cookies.some(c => c.startsWith("auth_token=")));

	console.groupEnd();
}

export async function checkCookieStoreAPI() {
	if ("cookieStore" in window) {
		const cookies = await (window as any).cookieStore.getAll();

		console.group("COOKIE STORE");
		console.table(cookies);
		console.groupEnd();
	} else {
		console.log("CookieStore API não disponível");
	}
}

export function cookieAPICheck() {
	console.group("COOKIE API");

	console.log(
		"auth_token:",
		document.cookie.includes("auth_token=")
	);

	console.log(
		"teste_cookie:",
		document.cookie.includes("teste_cookie=")
	);

	console.groupEnd();
}
