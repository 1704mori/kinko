import { PUBLIC_API_TOKEN, PUBLIC_API_URL } from "$env/static/public";

export async function api(request: RequestInfo, init?: RequestInit) {
    const baseUrl = PUBLIC_API_URL || "__API_URL__";
    const token = PUBLIC_API_TOKEN || "__API_TOKEN__";

    if (typeof request === 'string') {
        request = `${baseUrl}/api/v1${request}`;
    }

    if (init?.headers && !('Authorization' in init.headers)) {
        init.headers = {
            ...init.headers,
            Authorization: token,
        };
    }

    if (!init?.headers) {
        init = {
            ...init,
            headers: {
                Authorization: token,
                'Content-Type': 'application/json',
            },
        };
    }

    const res = await fetch(request, init);
    if (!res.ok) {
        const response = await res.json();
        throw new Error(response.message ?? 'Failed to fetch, check console for more information');
    }

    return res.json();
}