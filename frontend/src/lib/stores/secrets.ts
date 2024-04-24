import { writable, type Writable } from 'svelte/store';
import type { Secret } from '$lib/typings';
import { PUBLIC_API_TOKEN, PUBLIC_API_URL } from '$env/static/public';

interface SecretsStore {
  secrets: Writable<Secret[]>;
  fetchSecrets: (selectedSecret?: string) => Promise<void>;
}

function createSecretsStore(): SecretsStore {
  const secrets = writable<Secret[]>([]);

  const fetchSecrets = async (selectedSecret?: string): Promise<void> => {
    const baseUrl = PUBLIC_API_URL || "__API_URL__";
    const token = PUBLIC_API_TOKEN || "__API_TOKEN";
    if (!baseUrl || !token) {
      console.error('API URL or token is missing');
      return;
    }

    let url = `${baseUrl}/api/v1/secrets?offset=0&limit=10`;
    if (selectedSecret) {
      url += `&secret_name=${encodeURIComponent(selectedSecret)}`;
    }

    try {
      const response = await fetch(url, {
        headers: {
          'Authorization': token,
          'Content-Type': 'application/json'
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch secrets');
      }

      const data: Secret[] = await response.json();
      secrets.set(data);
    } catch (error) {
      console.error(error);
    }
  };

  return { secrets, fetchSecrets };
}

export const { secrets, fetchSecrets } = createSecretsStore();
