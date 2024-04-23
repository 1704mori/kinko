<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiToken, apiUrl } from '$lib';
	import Button from '$lib/components/Button.svelte';
	import Dialog from '$lib/components/Dialog.svelte';
	import Input from '$lib/components/Input.svelte';
	import Select from '$lib/components/Select/Select.svelte';
	import SelectSecret from '$lib/components/SelectSecret.svelte';
	import Table from '$lib/components/Table.svelte';
	import Actions from '$lib/components/Table/Actions.svelte';
	import Icon from '$lib/components/icons/Icon.svelte';
	import Plus from '$lib/components/icons/Plus.svelte';
	import { submit } from '$lib/form';
	import type { Secret } from '$lib/typings';
	import { onMount } from 'svelte';

	export let data: { props: { secrets: Secret[]; apiUrl: string; apiToken: string } };
	let selectedSecret: string | null = null;

	let createNewSecret = false;
	let showCreateDialog = false;

	$: {
		if (!showCreateDialog) {
			resetForm();
		}
	}
	
	function resetForm() {
		selectedSecret = null;
		createNewSecret = false;
		showCreateDialog = false;
	}

	async function selectSecret(name: string) {
		if (!name) {
			selectedSecret = null;
			goto(`/`, {
				replaceState: true,
				keepFocus: true,
				state: {
					secret_name: null
				}
			});
			return;
		}

		selectedSecret = name;

		const searchParams = new URLSearchParams();
		if (name) searchParams.set('secret_name', name);
		const secret_name = searchParams.toString();
		goto(`?${secret_name}`, {
			replaceState: true,
			keepFocus: true,
			state: {
				secret_name
			}
		});
	}

	const handleSubmit = submit<{
		secret_name: string;
		key: string;
		value: string;
	}>(async (data) => {
		console.log(data);
	});

	onMount(() => {
		const searchParams = new URLSearchParams(location.search);
		const secret_name = searchParams.get('secret_name');
		selectedSecret = secret_name;

		if (data.props.apiUrl) {
			apiUrl.set(data.props.apiUrl);
		}

		if (data.props.apiToken) {
			apiToken.set(data.props.apiToken);
		}
	});
</script>

<Dialog bind:showModal={showCreateDialog} class="w-96" on:submit={handleSubmit}>
	<div slot="header">Secrets</div>

	<div class="flex flex-col gap-1">
		<span class="text-xs uppercase">Select Secret</span>
		<SelectSecret
			name={createNewSecret ? '' : 'secret_name'}
			placeholder="Select Secret"
			secrets={data.props.secrets
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
				} else {
					selectSecret(e.detail[0].value);
				}
			}}
		/>
	</div>

	{#if createNewSecret}
		<div class="flex flex-col gap-1">
			<span class="text-xs uppercase">Secret Name</span>
			<Input name="secret_name" type="text" placeholder="Secret Name" />
		</div>
	{/if}

	<div class="flex flex-col gap-1">
		<span class="text-xs uppercase">Key</span>
		<Input name="key" type="text" placeholder="Key" />
	</div>
	<div class="flex flex-col gap-1">
		<span class="text-xs uppercase">Value</span>
		<Input name="value" type="text" placeholder="Value" />
	</div>

	<div slot="footer">
		<Button type="submit">Create Secret</Button>
	</div>
</Dialog>

<div class="flex items-center justify-center w-full">
	{#if data.props}
		<div class="flex flex-col gap-2">
			<div class="flex items-center">
				<SelectSecret
					class="w-40"
					secrets={data.props.secrets}
					{selectedSecret}
					on:select={(e) => selectSecret(!!e.detail.length ? e.detail[0].value : '')}
				/>
				<Button class="ml-auto" on:click={() => (showCreateDialog = true)}>
					<Icon>
						<Plus />
					</Icon>
					Create Secret
				</Button>
			</div>
			<Table
				label
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
				data={data.props.secrets.map((secret) => {
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
	{/if}
</div>
