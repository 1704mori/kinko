import type { Load } from "@sveltejs/kit";
import { KINKO_URL, KINKO_TOKEN } from '$env/static/private'
import type { Secret } from "$lib/typings";

export const load: Load = async ({url}) => {
    const offset = url.searchParams.get('page') ?? 0;
    const limit = 10;

    if (!KINKO_URL || !KINKO_TOKEN) {
        return {
            status: 401,
            error: new Error('Unauthorized')
        };
    }

    let _url = `${KINKO_URL}/api/v1/secrets?offset=${offset}&limit=${limit}`;

    if (url.searchParams.has('secret_name')) {
        _url += `&secret_name=${url.searchParams.get('secret_name')}`;
    }

    const res = await fetch(_url, {
        headers: {
            'Authorization': `${KINKO_TOKEN}`,
            'Content-Type': 'application/json'
        },
    });


    if (!res.ok) {
        console.error('[%2Bpage.server.ts] Failed to fetch', res);
        return {
            status: res.status,
            error: new Error('Failed to fetch, check console for more information')
        };
    }

    const secrets: Secret[] = await res.json();
    console.log('[%2Bpage.server.ts] Fetched secrets', secrets);

    return {
        props: {
            secrets,
            apiUrl: KINKO_URL,
            apiToken: KINKO_TOKEN
        }
    };
}