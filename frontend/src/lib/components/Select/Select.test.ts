// @ts-nocheck
import { afterAll, beforeEach, describe, expect, it } from 'vitest';
import { render, type RenderResult } from '@testing-library/svelte';
import Select from './Select.svelte';
import userEvent from '@testing-library/user-event';

const PROPS = {
	props: {
		options: [
			{
				label: 'Option 1',
				value: 'option_1'
			},
			{
				label: 'Option 2',
				value: 'option_2'
			},
			{
				label: 'Option 3',
				value: 'option_3'
			}
		]
	}
};

describe('Select.svelte', () => {
	let select: RenderResult<Select>;

	beforeEach((context) => {
		if (context.task.name != 'should allow multiple selection when isMultiple is true') {
			select = render(Select, PROPS);
		}
	});

	afterAll(() => {
		select.unmount();
	});

	it('should render a select element', () => {
		expect(select.getByText('Select')).toBeInTheDocument();
	});

	it('should render options', async () => {
		const user = userEvent.setup();
		expect(select.getByText('Select')).toBeInTheDocument();

		await user.click(select.getByText('Select'));
		expect(select.getByText('Option 1')).toBeInTheDocument();
		expect(select.getByText('Option 2')).toBeInTheDocument();
	});

	it('should select an option', async () => {
		const user = userEvent.setup();
		expect(select.getByText('Select')).toBeInTheDocument();

		await user.click(select.getByText('Select'));
		expect(select.getByText('Option 1')).toBeInTheDocument();
		expect(select.getByText('Option 2')).toBeInTheDocument();

		await user.click(select.getByText('Option 1'));
		expect(select.getByTestId('select')).toHaveTextContent('Option 1');

		expect(select.getByTestId('select').getAttribute('data-state')).toEqual('hidden');
	});
	
	it('should allow multiple selection when isMultiple is true', async () => {
		select = render(Select, {
			props: {
				...PROPS.props,
				isMultiple: true,
			}
		});
		const user = userEvent.setup();
		expect(select.getByText('Select')).toBeInTheDocument();

		await user.click(select.getByText('Select'));
		expect(select.getByText('Option 1')).toBeInTheDocument();
		expect(select.getByText('Option 2')).toBeInTheDocument();
		expect(select.getByText('Option 3')).toBeInTheDocument();

		await user.click(select.getByText('Option 1'));
		await user.click(select.getByText('Option 2'));

		expect(select.getByTestId('select')).toHaveTextContent('2 selected');
		expect(select.getByText('Option 3')).toBeInTheDocument();

		await user.click(select.getByText('Option 3'));

		expect(select.getByTestId('select')).toHaveTextContent('3 selected');

		await user.click(select.getByText('×'));
		
		expect(select.getByText('Select')).toBeInTheDocument();
	});
});
