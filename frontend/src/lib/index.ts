import { writable } from "svelte/store";

export const apiUrl = writable<string>('');
export const apiToken = writable<string>('');