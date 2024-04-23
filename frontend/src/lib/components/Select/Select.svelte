<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { cn } from '$lib/utils';
	import { clickOutside } from '$lib/clickOutside';
	import { fade } from 'svelte/transition';

	const dispatch = createEventDispatcher<{
		select: typeof options;
	}>();

	export let options: {
		label: string;
		value: string;
	}[] = [];
	export let isMultiple = false;
	export let value: string | string[] | null = null;
	export let placeholder = 'Select';

	let selectedOptions: typeof options = [];

	let showOptions = false;

	function toggleOptions() {
		showOptions = !showOptions;
	}

	function selectOption(option: (typeof options)[0], selected: boolean) {
		if (isMultiple) {
			if (selected) {
				selectedOptions = selectedOptions.filter((opt) => opt.value != option.value);
			} else {
				selectedOptions = [...selectedOptions, option];
			}
		} else {
			selectedOptions = [option];
			toggleOptions();
		}

		dispatch('select', selectedOptions);
	}

	function clear() {
		selectedOptions = [];
		dispatch('select', selectedOptions);
	}

	$: isSelected = (option: (typeof options)[0]) =>
		!!selectedOptions.find((opt) => opt.value == option.value)!;
	$: if (value) {
		selectedOptions = options.filter((opt) =>
			Array.isArray(value) ? value.includes(opt.value) : opt.value === value
		);
	} else {
		selectedOptions = [];
	}

	$: dataValue = () => {
		if (value) return value;

		if (!!selectedOptions.length) {
			if (isMultiple) {
				return JSON.stringify(selectedOptions.map((opt) => opt.value));
			} else {
				return selectedOptions[0].value;
			}
		}

		return '';
	};

	let className: string | undefined = undefined;
	export { className as class };
</script>

<div class="flex flex-col gap-1 ww-80 relative">
	<button
		type="button"
		class={cn(
			'select bg-white dark:bg-black hover:bg-neutral-200 border border-neutral-800 dark:border-neutral-800 hover:border-neutral-500 dark:hover:border-neutral-600 rounded-md h-10 relative disabled:opacity-80 disabled:cursor-not-allowed px-2 transition-colors',
			className
		)}
		data-testid="select"
		data-state={showOptions ? 'showing' : 'hidden'}
		data-value={dataValue()}
		on:click={toggleOptions}
		use:clickOutside={['button']}
		on:clickOutside={() => (showOptions = false)}
		{...$$restProps}
	>
		{#if selectedOptions.length === 0}
			<span>{placeholder}</span>
		{:else}
			<span>
				{#if isMultiple}
					{#if selectedOptions.length === 1}
						{selectedOptions[0].label}
					{:else}
						{selectedOptions.length} selected
					{/if}
				{:else}
					{selectedOptions[0].label}
				{/if}
			</span>
		{/if}

	</button>
	{#if !!selectedOptions.length}
		<button
			type="button"
			class="absolute right-2 top-1/2 -translate-y-1/2  px-1.5"
			on:click={clear}
		>
			&times;
		</button>
	{/if}
	{#if showOptions}
		<div
			class="flex flex-col gap-1 bg-neutral-50 dark:bg-neutral-950 border border-neutral-200 dark:border-neutral-900 p-1 rounded-md absolute top-[calc(100%+1.5rem/2)] w-full z-50 shadow-lg"
			transition:fade={{
				delay: 50,
				duration: 100
			}}
		>
			{#if !options.length}
				<span>
					No options available
				</span>
			{/if}
			{#each options as option}
				<button
					type="button"
					class={cn(
						'flex items-center gap-1 px-2 py-1 hover:bg-neutral-200 dark:hover:bg-neutral-800 group rounded-md',
						selectedOptions.find((o) => o.value === option.value) ? 'bg-neutral-50 dark:bg-neutral-950' : 'bg-neutral-50 dark:bg-neutral-950',
						'h-8'
					)}
					on:click={() => selectOption(option, isSelected(option))}
				>
					{#if isMultiple}
						<input
							type="checkbox"
							checked={isSelected(option)}
							on:change={() => selectOption(option, isSelected(option))}
						/>
					{/if}
					<span
						class={cn(
							// 'group-hover:text-white',
							// selectedOptions.find((o) => o.value === option.value) ? 'text-white' : 'text-tvt-red'
						)}
					>
						{option.label}
					</span>
				</button>
			{/each}
		</div>
	{/if}
</div>
