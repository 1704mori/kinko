export async function api(request: RequestInfo, init?: RequestInit) {
    const res = await fetch(request, init);
    if (!res.ok) {
        const response = await res.json();
        throw new Error(response.message ?? 'Failed to fetch, check console for more information');
    }

    return res.json();
}