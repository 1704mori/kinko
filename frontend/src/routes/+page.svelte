<script lang="ts">
	import { onMount } from 'svelte';
	import SelectSecret from '$lib/components/SelectSecret.svelte';
	import Table from '$lib/components/Table.svelte';
	import Button from '$lib/components/Button.svelte';
	import Icon from '$lib/components/icons/Icon.svelte';
	import Plus from '$lib/components/icons/Plus.svelte';
	import { fetchSecrets, secrets } from '$lib/stores/secrets';
	import Actions from '$lib/components/Table/Actions.svelte';
	import Dialog from '$lib/components/Dialog.svelte';
	import { submit } from '$lib/form';
	import Input from '$lib/components/Input.svelte';
	import { cn, to } from '$lib/utils';
	import { api } from '$lib/api';
	import { toast } from 'svelte-sonner';

	let selectedSecret: string | null = null;
	let showCreateDialog = false;
	let createNewSecret = false;

	let selectedModalSecret: string | null = null;

	onMount(() => {
		fetchSecrets();
	});

	function handleSelectSecret(secret: string): void {
		selectedSecret = secret ?? null;
		fetchSecrets(secret);
	}

	const handleSubmit = submit<{
		secret_name: string;
		key: string;
		value: string;
	}>(async ({ secret_name, key, value }) => {
		const [res, err] = await to(
			api(`/secret/${secret_name}`, {
				method: 'PUT',
				body: JSON.stringify({
					[key]: value,
				})
			})
		);

		if (err) {
			toast.error(err.message);
			return;
		}

		toast.success('Secret created/updated successfully');
		showCreateDialog = false;
		fetchSecrets();
	});
</script>

{#if showCreateDialog}
	<Dialog bind:showModal={showCreateDialog} class="w-96" on:submit={handleSubmit}>
		<div slot="header">Secrets</div>

		<div class="flex flex-col gap-1">
			<span class={cn('text-xs font-semibold uppercase')}>
				Select Secret
				<span class="font-normal italic text-red-500"></span>
			</span>
			<SelectSecret
				required={!createNewSecret}
				name={createNewSecret ? '' : 'secret_name'}
				placeholder="Select Secret"
				secrets={$secrets
					.concat([
						{
							secret_name: 'New Secret',
							key: 'New Secret',
							value: 'new_secret',
							created_at: '',
							updated_at: ''
						}
					])
					.sort((a, b) => a.secret_name.localeCompare(b.secret_name))}
				{selectedSecret}
				on:select={(e) => {
					if (e.detail[0].value === 'New Secret') {
						createNewSecret = true;
						selectedModalSecret = null;
					} else {
						createNewSecret = false;
						selectedModalSecret = e.detail[0].value;
					}
				}}
			/>
		</div>

		{#if createNewSecret}
			<div class="flex flex-col gap-1">
				<span class={cn('text-xs font-semibold uppercase')}>
					Secret Name
					<span class="font-normal italic text-red-500"></span>
				</span>
				<Input name="secret_name" type="text" placeholder="Secret Name" />
			</div>
		{/if}

		<div class="flex flex-col gap-1">
			<span class={cn('text-xs font-semibold uppercase')}>
				Key
				<span class="font-normal italic text-red-500"></span>
			</span>
			<Input required name="key" type="text" placeholder="Key" />
		</div>
		<div class="flex flex-col gap-1">
			<span class={cn('text-xs font-semibold uppercase')}>
				Value
				<span class="font-normal italic text-red-500"></span>
			</span>
			<Input required name="value" type="text" placeholder="Value" />
		</div>

		<div slot="footer">
			<Button type="submit">Create Secret</Button>
		</div>
	</Dialog>
{/if}

<div class="flex items-center justify-center w-full">
	<div class="flex flex-col gap-2">
		<div class="flex items-center">
			<SelectSecret
				class="w-40"
				secrets={$secrets}
				{selectedSecret}
				on:select={(e) => handleSelectSecret(e.detail[0].value)}
			/>
			<Button class="ml-auto" on:click={() => (showCreateDialog = true)}>
				<Icon>
					<Plus />
				</Icon>
				Create Secret
			</Button>
		</div>
		<Table
			columns={[
				{
					key: 'secret_name',
					label: 'Secret Name'
				},
				{
					key: 'key',
					label: 'Key'
				},
				{
					key: 'value',
					label: 'Value'
				},
				{
					key: 'created_at',
					label: 'Created At',
					modifyValue: (value) =>
						new Date(value).toLocaleDateString(undefined, {
							year: 'numeric',
							month: 'long',
							day: 'numeric',
							hour: '2-digit',
							minute: '2-digit'
						})
				},
				{
					key: 'updated_at',
					label: 'Updated At',
					modifyValue: (value) =>
						new Date(value).toLocaleDateString(undefined, {
							year: 'numeric',
							month: 'long',
							day: 'numeric',
							hour: '2-digit',
							minute: '2-digit'
						})
				},
				{
					key: 'actions',
					label: 'Actions'
				}
			]}
			data={$secrets.map((secret) => {
				return {
					...secret,
					actions: {
						component: Actions,
						props: {
							data: secret
						}
					}
				};
			})}
		/>
	</div>
</div>
