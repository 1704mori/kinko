export async function api(request: RequestInfo, init?: RequestInit) {
    const res = await fetch(request, init);
    if (!res.ok) {
        console.error('[api] Failed to fetch', res);
        const response = await res.json();
        console.error('[api] Response', response);
        throw new Error(response.message ?? 'Failed to fetch, check console for more information');
    }

    return res.json();
}