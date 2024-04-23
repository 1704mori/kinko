export function submit<T = any>(handler: (data: T) => void): (event: SubmitEvent) => void {
	return function factory(event) {
		// @ts-ignore
		const form = event.target!.closest('form') as HTMLFormElement;

		let valid = form.checkValidity();

		const inputs = form.querySelectorAll('input, textarea, .select');

		const payload: Record<string, any> = {};

		for (const input of inputs) {
			// parse 'Select' values
			if (!['INPUT', 'TEXTAREA'].includes(input.nodeName)) {
				let value = input.getAttribute('data-value');
				if (value) {
					if (value.startsWith('[')) {
						value = JSON.parse(value);
					}

					if (!input.getAttribute('name')) {
						continue;
					}

					payload[input.getAttribute('name')!] = value;
				}
			} else {
				if (import.meta.env.MODE != 'test') {
					validateForm(input as HTMLElement);
				}
			}
		}

		if (valid) {
			const data = new FormData(form);

			for (const [key, value] of data.entries()) {
				payload[key] = value;
			}

			handler(payload as T);
		} else {
			event.preventDefault();
			if (import.meta.env.MODE != 'test') {
				// @ts-ignore
				[...form].find((element) => element.validity.valid == false)?.focus();
			}
			form.checkValidity();
		}
	};
}

// finalizar validador caso venhamos a usar
// labels em cima dos inputs para mostrar errors
// e ajustar validator para aceitar Select
function validateForm(node: HTMLElement): { destroy: () => void } {
	const output = node.offsetParent?.previousElementSibling;

	node.addEventListener('invalid', onInvalid);
	node.addEventListener('input', onInput);

	function onInvalid() {
		// ...
	}

	function onInput() {}

	return {
		destroy() {
			node.removeEventListener('invalid', onInvalid);
			node.removeEventListener('input', onInput);
		}
	};
}
