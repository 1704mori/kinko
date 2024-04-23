<script lang="ts">
	import type { Secret } from '$lib/typings';
	import Select from './Select/Select.svelte';

	export let secrets: Secret[];
	export let selectedSecret: string | null;
</script>

<Select
	{...$$props}
	options={secrets
		.filter(
			(secret, index, self) => self.findIndex((s) => s.secret_name === secret.secret_name) === index
		)
		.map((secret) => ({
			label: secret.secret_name,
			value: secret.secret_name
		}))}
	value={selectedSecret}
	on:select
/>
