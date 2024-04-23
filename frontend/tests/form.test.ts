import { submit } from '$lib/form';
import { describe, it, expect, beforeEach, vi, afterAll } from 'vitest';

function createMockForm() {
	document.body.innerHTML = `
    <form id="test-form">
      <input name="name" value="John Doe" required />
      <input name="age" type="number" value="30" required />
      <textarea name="bio">Developer</textarea>
      <button type="submit">Submit</button>
    </form>
  `;
	return document.getElementById('test-form');
}

describe('submit function', () => {
	let form: Element;

	beforeEach(() => {
		createMockForm();
		form = document.querySelector('#test-form')!;
	});

	afterAll(() => {
		form = {} as Element;
	});

	it('calls the handler with correct payload when the form is valid', async () => {
		const mockHandler = vi.fn();
		const submitFunction = submit(mockHandler);

		// @ts-ignore
		form.addEventListener('submit', submitFunction);

		form.dispatchEvent(new Event('submit', { cancelable: true, bubbles: true }));

		expect(mockHandler).toHaveBeenCalledWith({
			name: 'John Doe',
			age: '30',
			bio: 'Developer'
		});
	});

	it('does not call the handler when the form is invalid', async () => {
		const mockHandler = vi.fn();

		// torna o formulario invalido
		// @ts-ignore
		form.querySelector('input[name="name"]').value = ''; // Assuming 'name' is required

		const submitFunction = submit(mockHandler);
		// @ts-ignore
		form.addEventListener('submit', submitFunction);

		form.dispatchEvent(new Event('submit', { cancelable: true, bubbles: true }));

		expect(mockHandler).not.toHaveBeenCalled();
	});
});
