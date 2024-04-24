<script lang="ts">
	import type { Secret } from '$lib/typings';
	import { to } from '$lib/utils';
	import { toast } from 'svelte-sonner';
	import Icon from '../icons/Icon.svelte';
	import Trash from '../icons/Trash.svelte';
	import Dialog from '../Dialog.svelte';
	import Button from '../Button.svelte';
	import { invalidateAll } from '$app/navigation';
	import { api } from '$lib/api';
	import { PUBLIC_API_TOKEN, PUBLIC_API_URL } from '$env/static/public';

	export let data: Secret;

	let showDeleteDialog = false;

	async function handleDelete() {
		const [res, err] = await to(
			api(`${PUBLIC_API_URL}/api/v1/secret/${data.secret_name}asdasd?key=${data.key}`, {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json',
					Authorization: PUBLIC_API_TOKEN
				}
			})
		);

		if (err) {
			toast.error(err.message);
			return;
		}

		toast.success('Secret deleted successfully');
		await invalidateAll();
	}
</script>

<Dialog bind:showModal={showDeleteDialog}>
	<div slot="header">Delete Secret</div>

	<div class="flex flex-col gap-1">
		<p>Are you sure you want to delete this secret?</p>
	</div>

	<div slot="footer" class="flex justify-end gap-2">
		<Button on:click={handleDelete}>Delete</Button>
	</div>
</Dialog>
<div class="flex items-center justify-center gap-2">
	<button type="button" on:click={() => (showDeleteDialog = true)}>
		<Icon color="red" size="20">
			<Trash />
		</Icon>
	</button>
</div>
